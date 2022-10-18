package blockchain

import (
	"github.com/rockjoon/nomadcoin/db"
	"github.com/rockjoon/nomadcoin/utils"
	"sync"
)

type blockchain struct {
	NewestHash        string `json:"newestHash"`
	Height            int    `json:"height"`
	CurrentDifficulty int    `json:"current_difficulty"`
}

var b *blockchain
var once sync.Once

const (
	defaultDifficulty       int = 2
	difficultyInterval      int = 3
	allowedMiningSecondsMin int = 30
	allowedMiningSecondsMax int = 40
)

func GetBlockChain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{Height: 0}
			checkpoint := db.GetCheckPoint()
			if checkpoint == nil {
				b.AddBlock()
			} else {
				b.restore(checkpoint)
			}
		})
	}
	return b
}

func (b *blockchain) restore(data []byte) {
	utils.FromBytes(data, b)
}

func (b *blockchain) difficulty() int {
	if b.Height == 0 {
		return defaultDifficulty
	}
	if b.Height%difficultyInterval == 0 {
		return b.recalculateDifficulty()
	}
	return b.CurrentDifficulty
}

func (b *blockchain) recalculateDifficulty() int {
	blocks := b.AllBlocks()
	newestBlock := blocks[0]
	lastRecalculated := blocks[difficultyInterval-1]
	actualTime := newestBlock.Timestamp - lastRecalculated.Timestamp
	if actualTime > allowedMiningSecondsMax {
		if b.CurrentDifficulty == 1 {
			return b.CurrentDifficulty
		}
		return b.CurrentDifficulty - 1
	}
	if actualTime < allowedMiningSecondsMin {
		return b.CurrentDifficulty + 1
	}
	return b.CurrentDifficulty
}

func (b *blockchain) AddBlock() {
	block := createBlock(b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
	b.CurrentDifficulty = block.Difficulty
	b.persist()
}

func (b *blockchain) persist() {
	db.SaveBlockChain(b)
}

func (b *blockchain) AllBlocks() []*Block {
	var blocks []*Block
	blockCursor := b.NewestHash
	for {
		block, _ := FindBlock(blockCursor)
		blocks = append(blocks, block)
		if block.PrevHash != "" {
			blockCursor = block.PrevHash
		} else {
			return blocks
		}
	}

}

func (b *blockchain) txOuts() []*TxOut {
	blocks := b.AllBlocks()
	var txOuts []*TxOut
	for _, block := range blocks {
		for _, tx := range block.Transactions {
			txOuts = append(txOuts, tx.TxOuts...)
		}
	}
	return txOuts
}

func (b *blockchain) BalanceByAddress(address string) int {
	var balance int
	for _, txOut := range b.UTxOutsByAddress(address) {
		balance += txOut.Amount
	}
	return balance
}

func (b *blockchain) UTxOutsByAddress(address string) []*UTxOut {
	var utxOuts []*UTxOut
	creatorTxs := make(map[string]bool)

	for _, block := range b.AllBlocks() {
		for _, tx := range block.Transactions {
			for _, txIn := range tx.TxIns {
				if txIn.Owner == address {
					creatorTxs[txIn.TxID] = true
				}
			}
			for i, txOut := range tx.TxOuts {
				if txOut.Owner == address {
					_, ok := creatorTxs[tx.Id]
					if !ok {
						utxOut := &UTxOut{tx.Id, i, txOut.Amount}
						if !isOnMempool(utxOut) {
							utxOuts = append(utxOuts, utxOut)
						}
					}
				}
			}
		}
	}
	return utxOuts
}
