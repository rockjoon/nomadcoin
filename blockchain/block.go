package blockchain

import (
	"crypto/sha256"
	"fmt"
	"github.com/rockjoon/nomadcoin/db"
)

type Block struct {
	Data     string
	Hash     string
	PrevHash string
	Height   int
}

func (b *Block) persist() {
	db.SaveBlock(b.Hash, b)
	//db.Save(db.BlocksBucket, b.Hash, b)

}

func createBlock(data string, prevhash string, height int) *Block {
	block := Block{
		Data:     data,
		PrevHash: prevhash,
		Height:   height,
	}
	payload := data + prevhash + fmt.Sprint(height)
	block.Hash = fmt.Sprintf("%x", sha256.Sum256([]byte(payload)))
	fmt.Println(block)
	block.persist()
	return &block
}
