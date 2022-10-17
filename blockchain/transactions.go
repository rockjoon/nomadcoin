package blockchain

import (
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

func (*mempool) makeTxs(from, to string, amount int) (*Tx, error) {
	return nil, nil
}

func (m *mempool) AddTxs(to string, amount int) error {
	tx, err := m.makeTxs("joon", to, amount)
	if err != nil {
		return err
	}
	m.Txs = append(m.Txs, tx)
	return nil
}

func (m *mempool) TxToConfirm() []*Tx {
	coinbase := makeCoinbaseTx("nico")
	txs := m.Txs
	txs = append(txs, coinbase)
	m.Txs = nil
	return txs
}
