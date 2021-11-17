package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"

	"github.com/ethereum/go-ethereum/common"
)

const TxGas = 21
const TxGasPriceDefault = 1
const TxFee = uint(50)

type State struct {
	Balances      map[common.Address]uint
	Account2Nonce map[common.Address]uint

	dbFile *os.File

	latestBlock     Block
	latestBlockHash Hash
	hasGenesisBlock bool

	miningDifficulty uint

	forkTIP1 uint64
}

func NewStateFromDisk(dataDir string, miningDifficulty uint) (*State, error) {
	err := InitDataDirIfNotExist(dataDir, []byte(genesisJson))
	if err != nil {
		return nil, err
	}

	gen, err := loadGenesis(getGenesisJsonFilePath(dataDir))
	if err != nil {
		return nil, err
	}

	balances := make(map[common.Address]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	account2nonce := make(map[common.Address]uint)

	dbFilePath := getBlocksDbFilePath(dataDir)
	f, err := os.OpenFile(dbFilePath, os.O_APPEND|os.O_RDWR, 0600)

	scanner := bufio.NewScanner(f)

	state := &State{
		balances, account2nonce, f, Block{}, Hash{}, false, miningDifficulty, gen.ForkTIP1,
	}

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		blockFsJson := scanner.Bytes()

		if len(blockFsJson) == 0 {
			break
		}

		var blockFs BlockFS
		err = json.Unmarshal(blockFsJson, &blockFs)
		if err != nil {
			return nil, err
		}

		err = applyBlock(blockFs.Value, state)
		if err != nil {
			return nil, err
		}

		state.latestBlock = blockFs.Value
		state.latestBlockHash = blockFs.Key
		state.hasGenesisBlock = true

	}

	return state, nil

}

func (s *State) LatestBlockHash() Hash {
	return s.latestBlockHash
}

func (s *State) Close() error {
	return s.dbFile.Close()
}

func (s *State) NextBlockNumber() uint64 {
	if !s.hasGenesisBlock {
		return uint64(0)
	}

	return s.LatestBlock().Header.Number
}

func (s *State) LatestBlock() Block {
	return s.latestBlock
}

func (s *State) GetNextAccountNonce(account common.Address) uint {
	return s.Account2Nonce[account] + 1
}

func (s *State) IsTIP1Fork() bool {
	return s.NextBlockNumber() >= s.forkTIP1
}

func applyBlock(b Block, s *State) error {
	nextExpectedBlockNumber := s.latestBlock.Header.Number + 1

	if s.hasGenesisBlock && b.Header.Number != nextExpectedBlockNumber {
		return fmt.Errorf("Next Expected BlockNumber must be '%d' no '%d'", nextExpectedBlockNumber, b.Header.Number)
	}

	if s.hasGenesisBlock && s.latestBlock.Header.Number > 0 && !reflect.DeepEqual(b.Header.Parent, s.latestBlockHash) {
		return fmt.Errorf("next block parent hash must '%x'not '%x' ", s.latestBlockHash, b.Header.Parent)
	}

	hash, err := b.Hash()
	if err != nil {
		return err
	}

	if !IsBlockHashValid(hash, s.miningDifficulty) {
		return fmt.Errorf("Invalid block hash %x", hash)
	}

	err = applyTXs(b.TXs, s)
	if err != nil {
		return err
	}

	return nil
}

func applyTXs(txs []SignedTx, s *State) error {
	sort.Slice(txs, func(i, j int) bool {
		return txs[i].Time < txs[j].Time

	})

	for _, tx := range txs {
		err := ApplyTx(tx, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func ApplyTx(tx SignedTx, s *State) error {
	err := ValidateTx(tx, s)
	if err != nil {
		return err
	}

	return nil
}

func ValidateTx(tx SignedTx, s *State) error {
	ok, err := tx.IsAuthentic()
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("wrong Tx. Sender %s is forged", tx.From.String())
	}

	expectedNonce := s.GetNextAccountNonce(tx.From)
	if tx.Nonce != expectedNonce {
		return fmt.Errorf("wrong TX. Sender '%s' next nonce must be '%d', not '%d'", tx.From.String(), expectedNonce, tx.Nonce)
	}

	if s.IsTIP1Fork() {
		// all tx must pay 21 gas
		if tx.Gas != TxGas {
			return fmt.Errorf("insufficient TX gas %v. required: %v", tx.Gas, TxGas)

		}

		if tx.GasPrice < TxGasPriceDefault {
			return fmt.Errorf("insufficient TX gasPrice %v. required at least: %v", tx.GasPrice, TxGasPriceDefault)
		}

	} else {
		if tx.Gas != 0 || tx.GasPrice != 0 {
			return fmt.Errorf("invalid TX. `Gas` and `GasPrice` can't be populated before TIP1 fork is active")
		}
	}

	if tx.Cost(s.IsTIP1Fork()) > s.Balances[tx.From] {
		return fmt.Errorf("wrong TX. sender '%s' balance is '%d' OBBOY. Tx cost is '%d' OBBOY", tx.From.String(), s.Balances[tx.From], tx.Cost(s.IsTIP1Fork()))
	}

	return nil
}
