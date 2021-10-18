package main

import (
	"fmt"
	"log"
	"myblockchain/blockchain"
	"os"
	"runtime"
)

type CommandLine struct{}

func main() {
	defer os.Exit(0)
	cli := CommandLine{}
	cli.run()
}

func (cli *CommandLine) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("1: Init Blockchain")
	fmt.Println("2: Print Blockchain")
	fmt.Println("3: Send coin by addresses")
	fmt.Println("4: View balance of an addresses")
	fmt.Println("0: Exit Program")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) printChain() {
	chain := blockchain.LoadBlockChain()
	defer chain.Database.Close()
	iter := chain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Prev Hash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Println()

		if block.PrevHash == nil {
			break
		}
	}
}

func (cli *CommandLine) createBlockChain(address string) {
	chain := blockchain.InitBlockChain(address)
	chain.Database.Close()
	fmt.Println("Finished!")
}

func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.LoadBlockChain()
	defer chain.Database.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)
	chain.AddBlock([]*blockchain.Transaction{tx})
	fmt.Println("Success!")
}

func (cli *CommandLine) viewBalance(address string) {
	chain := blockchain.LoadBlockChain()
	defer chain.Database.Close()

	trans := chain.FindUnspentTxOutputs(address)
	balance := 0
	for _, tx := range trans {
		for _, output := range tx.Outputs {
			if output.Address == address {
				balance += output.Value
			}
		}
	}

	fmt.Printf("Balance of %s is %d\n", address, balance)
}

func (cli *CommandLine) run() {
	userChoice := 1

	cli.printUsage()

	for userChoice != 0 {
		fmt.Printf("--------------------------------------------------------------------------------\n")
		fmt.Printf("Please enter your choice: ")
		_, err := fmt.Scan(&userChoice)
		if err != nil || userChoice > 4 || userChoice < 0 {
			cli.handleErrors(err)
			continue
		}

		switch userChoice {
		case 1:
			fmt.Printf("-------------Create blockchain and sends genesis reward to address-------------\n")
			var address string
			fmt.Printf("Please enter first address in blockchain: ")
			_, err := fmt.Scan(&address)
			if err != nil {
				cli.handleErrors(err)
				continue
			}
			cli.createBlockChain(address)
		case 2:
			fmt.Printf("-------------------------------------Print Chain--------------------------------\n")
			cli.printChain()
		case 3:
			fmt.Printf("-------------------------------Send coin over addresses-------------------------\n")
			fmt.Printf("Please enter FromAddress ToAddress Amount: ")
			var fromAddress, toAddress string
			var amount int
			_, err := fmt.Scanf("%s %s %d", &fromAddress, &toAddress, &amount)
			if err != nil {
				cli.handleErrors(err)
				continue
			} else {
				cli.send(fromAddress, toAddress, amount)
			}
		case 4:
			var address string
			fmt.Printf("-------------------------------view balance of an addresses---------------------\n")
			fmt.Printf("Please enter address: ")
			_, err := fmt.Scan(&address)
			if err != nil {
				cli.handleErrors(err)
				continue
			}
			cli.viewBalance(address)
		}
	}
}

func (cli *CommandLine) handleErrors(err error) {
	fmt.Printf("Sorry, there some errors happen. Please make sure that input's type is correct or you don't forget anything\n")
	if err != nil {
		log.Print(err)
	}
}
