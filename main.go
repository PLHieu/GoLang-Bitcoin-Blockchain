package main

import (
	"fmt"
	"myblockchain/blockchain"
)

func main() {
	chain := blockchain.InitBlockChain()

	fmt.Printf("%x", chain.LastHash)

	//chain.AddBlock("First Block after Genesis")
	//chain.AddBlock("Second Block after Genesis")
	//chain.AddBlock("Third Block after Genesis")

	//for _, block := range chain.Blocks {
	//	fmt.Printf("Prev: %x\n", block.PrevHash)
	//	fmt.Printf("Data: %s\n", block.Data)
	//	fmt.Printf("Hash: %x\n", block.Hash)
	//	fmt.Printf("Nonce: %d\n", block.Nonce)
	//	fmt.Printf("---------------------------------\n")
	//}
}
