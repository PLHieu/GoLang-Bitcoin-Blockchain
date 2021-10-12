package blockchain

import (
	"encoding/hex"
	"fmt"
	"github.com/dgraph-io/badger"
	"myblockchain/utils"
	"os"
	"runtime"
)

const (
	dbPath = "./tmp/blocks"
	dbFile      = "./tmp/blocks/MANIFEST"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func DBexists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

func LoadBlockChain() *BlockChain {
	if DBexists() == false {
		fmt.Println("No existing blockchain found, create one!")
		runtime.Goexit()
	}

	var lastHash []byte

	db, err := badger.Open(badger.DefaultOptions(dbPath))
	utils.Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lasthash"))
		utils.Handle(err)

		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		utils.Handle(err)

		return err
	})
	utils.Handle(err)

	chain := BlockChain{lastHash, db}

	return &chain
}

func InitBlockChain(firstAddress string) *BlockChain {
	if DBexists() {
		fmt.Println("Blockchain already exists")
		runtime.Goexit()
	}

	db, err := badger.Open(badger.DefaultOptions(dbPath))
	utils.Handle(err)

	// Create new Blockchain
	newBlockChain := &BlockChain{[]byte{}, db}

	// Create Genesis Block
	genesisCoinbaseTx := CoinbaseTx(firstAddress, "First Transaction from Genesis")
	newBlock := CreateGenesis(genesisCoinbaseTx)

	// Insert block into Badge DB & Update LastHash
	err = newBlockChain.Database.Update(func(txn *badger.Txn) error {
		err = txn.Set(newBlock.Hash, newBlock.Serialize())
		utils.Handle(err)
		err = txn.Set([]byte("lasthash"), newBlock.Hash)
		utils.Handle(err)
		newBlockChain.LastHash = newBlock.Hash

		return err
	})
	utils.Handle(err)

	return newBlockChain
}

func (chain *BlockChain) AddBlock(txs []*Transaction) {
	// tao block tu data + lastHash trong Badger DB
	var lastHash []byte
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lasthash"))

		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})

		return err
	})
	utils.Handle(err)
	newBlock := CreateBlock(txs, lastHash)

	// them block vao trong Badger DB &  Update lastHash cua Badger DB & blockchain
	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		err = txn.Set([]byte("lasthash"), newBlock.Hash)

		chain.LastHash = newBlock.Hash

		return err
	})
	utils.Handle(err)
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}

	return iter
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	// get Value from hash key
	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		if err != nil {
			block = nil
		} else {
			err = item.Value(func(encodedBlock []byte) error {
				block = Deserialize(encodedBlock)
				return nil
			})
			utils.Handle(err)
		}

		return err
	})
	utils.Handle(err)

	// Move current to next
	iter.CurrentHash = block.PrevHash

	return block
}

// FindUnspentTxOutputs
// Find all address's transaction that have unspent outputs
// Input: address that want to transfer coin
// Output: list of transactions
///*
func (chain *BlockChain) FindUnspentTxOutputs(address string) []Transaction {
	var unspentTxs []Transaction

	// list of outputs that were spent, map: transactionID -> index of spent Output
	spentTXOs := make(map[string][]int)

	iter := chain.Iterator()
	for {
		block := iter.Next()
		if block == nil {
			break
		}

		// Iterate each transaction in a block
		for _, tx := range block.Transactions {
			// Get string version of transaction's ID
			strTxId := hex.EncodeToString(tx.ID)

		Outputs:
			// iterate each of outputs of transaction
			for outIdx, output := range tx.Outputs {
				// If output ot this transaction has spent -> continue
				if spentTXOs[strTxId] != nil {
					for _, spentOut := range spentTXOs[strTxId] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				// If output of this transaction is not spent and is this address
				if output.Address == address {
					unspentTxs = append(unspentTxs, *tx)
				}
			}

			// if this transaction is not coinbase transaction
			if !tx.IsCoinBase() {
				for _, input := range tx.Inputs {
					// Todo: something with signature
					if input.Sig == address {
						// Add spent output into spentTXOs
						spentTXOs[strTxId] = append(spentTXOs[strTxId], input.OutIndex)
					}
				}
			}
		}
	}

	return unspentTxs
}

// FindSpendableOutputs
// Input: address that want to transfer coin, amount of coin
// Output: Sum of spendable outputs
// Output: Map: transactionID -> array of Output's index in that transaction
///*
func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	accumulated := 0

	// Get all address's transaction that have unspent outputs
	unspentTxs := chain.FindUnspentTxOutputs(address)

	for _, tx := range unspentTxs {
		// Get string version of transaction's ID
		txID := hex.EncodeToString(tx.ID)

		// Iterate each of outputs in that transaction
		for outIdx, output := range tx.Outputs {
			// If that Output belong to address and acc still < amount
			if output.Address == address {
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)
				accumulated += output.Value
			}

			if accumulated > amount {
				return accumulated, unspentOuts
			}
		}
	}

	return accumulated, unspentOuts
}
