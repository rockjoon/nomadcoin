package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type block struct {
	Data     string
	Hash     string
	PrevHash string
}

type blockchain struct {
	blocks []*block
}

func (b *block) setHash() {
	hash := sha256.Sum256([]byte(b.Data))
	b.Hash = fmt.Sprintf("%x", hash)
}

func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, createBlock(data))
}

func (b *blockchain) AllBlocks() []*block {
	return GetBlockChain().blocks
}

func createBlock(data string) *block {
	newBlock := block{data, "", getLastHash()}
	newBlock.setHash()
	return &newBlock
}

func getLastHash() string {
	totalLength := len(GetBlockChain().blocks)
	if totalLength == 0 {
		return ""
	}
	return GetBlockChain().blocks[totalLength-1].Hash
}

var b *blockchain
var once sync.Once

func GetBlockChain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{}
			b.AddBlock("genesisBlock")
		})
	}
	return b
}
