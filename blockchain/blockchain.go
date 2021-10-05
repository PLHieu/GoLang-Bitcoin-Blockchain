package blockchain

import (
	"github.com/dgraph-io/badger"
	"log"
)

const (
	dbPath = "./tmp/blocks"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

func InitBlockChain() *BlockChain {
	db, err := badger.Open(badger.DefaultOptions(dbPath))
	handle(err)

	// Create new Blockchain
	newBlockChain := &BlockChain{[]byte{}, db}

	// Create Genesis Block
	newBlock := CreateGenesis()

	// Insert block into Badge DB & Update LastHash
	err = newBlockChain.Database.Update(func(txn *badger.Txn) error {
		lastHash, err := txn.Get([]byte("lasthash"))

		// if blockchain haven't ever exist
		if err == badger.ErrKeyNotFound {
			err := txn.Set(newBlock.Hash, newBlock.Serialize())
			err = txn.Set([]byte("lasthash"), newBlock.Hash)

			newBlockChain.LastHash = newBlock.Data

			return err
		} else { // if it already existed
			err := lastHash.Value(func(val []byte) error {
				newBlockChain.LastHash = val
				return nil
			})
			handle(err)
		}

		return err
	})
	handle(err)

	return newBlockChain
}

func (chain *BlockChain) AddBlock(data string) {
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
	handle(err)
	newBlock := CreateBlock(data, lastHash)

	// them block vao trong Badger DB &  Update lastHash cua Badger DB & blockchain
	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		err = txn.Set([]byte("lasthash"), newBlock.Hash)

		chain.LastHash = newBlock.Data

		return err
	})
	handle(err)
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
