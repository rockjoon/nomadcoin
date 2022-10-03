package blockchain

import (
	"github.com/rockjoon/nomadcoin/db"
	"github.com/rockjoon/nomadcoin/utils"
	"sync"
)

type blockchain struct {
	NewestHash string `json:"newestHash"`
	Height     int    `json:"height"`
}

var b *blockchain
var once sync.Once

func GetBlockChain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{"", 0}
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

func (b *blockchain) AddBlock(data string) {
	block := createBlock(data, b.NewestHash, b.Height+1)
	b.NewestHash = block.Hash
	b.Height = block.Height
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
