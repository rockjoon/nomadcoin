package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
)

type Block struct {
	Data     string
	Hash     string
	PrevHash string
}

type blockchain struct {
	blocks []*Block
}

var b *blockchain
var once sync.Once

func GetBlockChain() *blockchain {
	if b == nil {
		once.Do(func() {
			b = &blockchain{}
			b.AddBlock("genesis")
		})
	}
	return b
}

func (b *blockchain) AddBlock(data string) {
	b.blocks = append(b.blocks, CreateBlock(data))
}

func CreateBlock(data string) *Block {
	newBlock := Block{data, "", getLastHash()}
	newBlock.Hash = calculateHash(newBlock.Data, newBlock.PrevHash)
	return &newBlock
}

func calculateHash(data string, prevHash string) string {
	hash := sha256.Sum256([]byte(data + prevHash))
	return fmt.Sprintf("%x", hash)
}

func getLastHash() string {
	totalBlockCount := len(GetBlockChain().blocks)
	if totalBlockCount == 0 {
		return ""
	}
	return GetBlockChain().blocks[totalBlockCount-1].Hash
}

func (b *blockchain) AllBlocks() []*Block {
	return b.blocks
}
