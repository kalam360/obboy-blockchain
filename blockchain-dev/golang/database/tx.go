package database

import (
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)


func NewAccount(value string) common.Address {
	return common.HexToAddress(value)
}

type Tx struct {
	From     common.Address `json:"from"`
	To       common.Address `json:"to"`
	Gas      uint           `json:"gas"`
	GasPrice uint           `json:"gasPrice"`
	Value    uint           `json:"value"`
	Nonce    uint           `json:"nonce"`
	Data     string         `json:"data"`
	Time     uint64         `json:"time"`
}

type SignedTx struct {
	Tx
	Sig []byte `json:"signature"`
}

func (t Tx) Encode() ([]byte, error) {

	return json.Marshal(t)
}

func (t Tx) Cost(isTip1Fork bool) uint {
	if isTip1Fork {
		return t.Value + t.GasCost() 
	}

	return t.Value + TxFee
}

func (t Tx) GasCost() uint {
	return t.Gas * t.GasPrice
}

func (t Tx) Hash() (Hash, error) {
	txJson, err := t.Encode()
	if err != nil {
		return Hash{}, err
	}

	return sha256.Sum256(txJson), nil
}

func (t SignedTx) IsAuthentic() (bool, error) {
	txHash, err := t.Tx.Hash()
	if err != nil {
		return false, err
	}

	recoveredPubKey, err := crypto.SigToPub(txHash[:], t.Sig)
	if err != nil {
		return false, err
	}

	recoveredPubKeyBytes := elliptic.Marshal(crypto.S256(), recoveredPubKey.X, recoveredPubKey.Y)
	recoveredPubKeyBytesHash := crypto.Keccak256(recoveredPubKeyBytes[1:])
	recoveredAccount := common.BytesToAddress(recoveredPubKeyBytesHash[12:])

	return recoveredAccount.Hex() == t.From.Hex(), nil

}
