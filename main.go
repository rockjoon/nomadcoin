package main

import (
	"fmt"
	"github.com/rockjoon/nomadcoin/blockchain"
)

func main() {
	chain := blockchain.GetBlockChain()
	chain.AddBlock("second")
	chain.AddBlock("thrid")
	for _, block := range chain.AllBlocks() {
		fmt.Println(block)
	}
}
