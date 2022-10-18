package blockchain

import (
	"errors"
	"github.com/rockjoon/nomadcoin/utils"
	"time"
)

const (
	mineReward = 50
)

var Mempool *mempool = &mempool{}

type Tx struct {
	Id        string   `json:"id"`
	Timestamp int      `json:"timestamp"`
	TxIns     []*TxIn  `json:"tx_ins"`
	TxOuts    []*TxOut `json:"tx_outs"`
}

type TxIn struct {
	TxID  string `json:"tx_id"`
	Index int    `json:"index"`
	Owner string `json:"owner"`
}

type TxOut struct {
	Owner  string `json:"owner"`
	Amount int    `json:"amount"`
}

type UTxOut struct {
	TxID   string `json:"tx_id"`
	Index  int    `json:"index"`
	Amount int    `json:"amount"`
}

type mempool struct {
	Txs []*Tx
}

func (t *Tx) setId() {
	t.Id = utils.Hash(t)
}

func makeCoinbaseTx(address string) *Tx {
	txIns := []*TxIn{
		{"", -1, "COINBASE"},
	}
	txOut := []*TxOut{
		{address, mineReward},
	}
	tx := Tx{
		Id:        "",
		Timestamp: int(time.Now().Unix()),
		TxIns:     txIns,
		TxOuts:    txOut,
	}
	tx.setId()
	return &tx
}

func makeTxs(from, to string, amount int) (*Tx, error) {
	if BalanceByAddress(from, GetBlockChain()) < amount {
		return nil, errors.New("not enough balance")
	}
	var txIns []*TxIn
	var txOuts []*TxOut
	total := 0
	uTxOuts := UTxOutsByAddress(from, GetBlockChain())
	for _, utxOut := range uTxOuts {
		if total >= amount {
			break
		}
		txIn := &TxIn{utxOut.TxID, utxOut.Index, from}
		txIns = append(txIns, txIn)
		total = total + utxOut.Amount
	}
	if change := total - amount; change > 0 {
		txOuts = append(txOuts, &TxOut{from, change})
	}
	txOuts = append(txOuts, &TxOut{to, amount})
	var tx = &Tx{
		"",
		int(time.Now().Unix()),
		txIns,
		txOuts,
	}
	tx.setId()
	return tx, nil
}

func (m *mempool) AddTxs(to string, amount int) error {
	tx, err := makeTxs("joon", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx("joon")
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}

func isOnMempool(utxOut *UTxOut) bool {
	for _, tx := range Mempool.Txs {
		for _, txIn := range tx.TxIns {
			if isSameTx(utxOut, txIn) {
				return true
			}
		}
	}
	return false
}

func isSameTx(utxOut *UTxOut, txIn *TxIn) bool {
	return utxOut.TxID == txIn.TxID && utxOut.Index == txIn.Index
}
