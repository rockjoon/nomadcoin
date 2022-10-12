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
				b.AddBlock("genesis")
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

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
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
