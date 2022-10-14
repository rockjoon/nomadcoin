package blockchain

import (
	"github.com/rockjoon/nomadcoin/db"
	"github.com/rockjoon/nomadcoin/utils"
	"strings"
	"time"
)

type Block struct {
	Transactions []*Tx  `json:"transactions"`
	Hash         string `json:"hash"`
	PrevHash     string `json:"prev_hash"`
	Height       int    `json:"height"`
	Difficulty   int    `json:"difficulty"`
	Nonce        int    `json:"nonce"`
	Timestamp    int    `json:"timestamp"`
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, b)
}

func (b *Block) restore(blockBytes []byte) {
	utils.FromBytes(blockBytes, b)
}

func createBlock(prevhash string, height int) *Block {
	block := Block{
		Transactions: []*Tx{makeCoinbaseTx("joon")},
		Hash:         "",
		PrevHash:     prevhash,
		Height:       height,
		Difficulty:   GetBlockChain().difficulty(),
		Nonce:        0,
		Timestamp:    int(time.Now().Unix()),
	}
	block.Hash, block.Timestamp = block.mine()
	block.persist()
	return &block
}

func (b *Block) mine() (string, int) {
	target := strings.Repeat("0", b.Difficulty)

	for {
		hash := utils.Hash(b)
		if strings.HasPrefix(hash, target) {
			return hash, int(time.Now().Unix())
		}
		b.Nonce++
	}
}

func FindBlock(hash string) (*Block, error) {
	blockBytes := db.Block(hash)
	if blockBytes == nil {
		return nil, BlockNotFoundError
	} else {
		var block = &Block{}
		block.restore(blockBytes)
		return block, nil
	}
}
