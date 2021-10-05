package main

import (
	"fmt"
	"myblockchain/blockchain"
)

func main() {
	chain := blockchain.InitBlockChain()

	//chain.AddBlock("First Block after Genesis")
	//chain.AddBlock("Second Block after Genesis")
	//chain.AddBlock("Third Block after Genesis")

	iter := chain.Iterator()
	block := iter.Next()
	for block != nil {
		fmt.Printf("Prev: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		block = iter.Next()
	}
}
