package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type Account string

type Tx struct {
	From Account `json:"from"`
	To Account `json:"to"`
	Value uint `json:"value"`
	Data string `json:"data"`
}

func (t Tx) IsReward() bool {
	return t.Data == "reward"
}

type State struct {
	Balances map[Account]uint `json:"balances"`
	txMempool []Tx

	dbFile *os.File
}

type Genesis struct {
	Time time.Time `json:"genesis_time"`
	ChainId string `json:"chain-id"`
	Balances map[Account]uint `json:"balances"`
}

func loadGenesis(path string) (Genesis, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return Genesis{}, err
	}

	var loadedGenesis Genesis

	err = json.Unmarshal(content, &loadedGenesis)
	if err != nil {
		return Genesis{}, err
	}

	return loadedGenesis, nil

}

func (s *State) apply(tx Tx) error {
	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}

	if tx.Value > s.Balances[tx.From] {
		return fmt.Errorf("insufficient balance")
	}

	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value

	return nil
}

func (s *State) Add(tx Tx) error {
	if err := s.apply(tx); err != nil {
		return err
	}

	s.txMempool = append(s.txMempool, tx)

	return nil
}

func (s *State) Persist() error {
	// make a copy of the mempool as it will be modified
	mempool := make([]Tx, len(s.txMempool))
	copy(mempool, s.txMempool)

	for i := 0; i < len(mempool); i++ {
		txJson, err := json.Marshal(mempool[i])
		if err != nil {
			return err
		}

		if _, err = s.dbFile.Write(append(txJson, '\n')); err != nil {
			return err
		}

		s.txMempool = s.txMempool[1:]
	}

	return nil
}



func NewStateFromDisk() (*State, error) {
	// get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	genFilePath := filepath.Join(cwd, "database", "genesis.json")

	gen, err := loadGenesis(genFilePath)
	if err != nil {
		return nil, err
	}

	balances := make(map[Account]uint)

	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	txDbFilePath := filepath.Join(cwd, "database", "tx.db")
	f, err := os.OpenFile(txDbFilePath, os.O_APPEND | os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	state := &State{balances, make([]Tx, 0), f}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		// Convert JSON encoded TX into an object (struct)
		var tx Tx 
		json.Unmarshal(scanner.Bytes(), &tx)

		// Rebuild the state (user balances) as a series of events
		if err := state.apply(tx); err != nil {
			return nil, err
		}

	}

	return state, nil

}

func ReadGenesis() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	genFilePath := filepath.Join(cwd, "database", "genesis.json")
	
	genesis, err := loadGenesis(genFilePath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Successfully opend genesis file")
	fmt.Println(genesis)

	return nil
}
