package main

import (
	"github.com/rockjoon/nomadcoin/db"
	"github.com/rockjoon/nomadcoin/rest"
)

func main() {
	rest.Start(4000)
	defer db.Close()
	//blockchain.GetBlockChain()

}
