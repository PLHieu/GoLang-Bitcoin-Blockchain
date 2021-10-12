package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
	"myblockchain/utils"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

type TxInput struct {
	OutTxId    []byte
	OutIndex int
	Sig      string
}

type TxOutput struct {
	Value   int
	Address string
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	utils.Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

func CoinbaseTx(outputAddress, data string) *Transaction {
	txOut := TxOutput{100, outputAddress}
	// Todo: Generate signature for transaction
	txIn := TxInput{nil, -1, data}

	tx := &Transaction{nil, []TxInput{txIn}, []TxOutput{txOut}}
	tx.SetID()

	return tx
}

func NewTransaction(fromAddress string, toAddress string, amount int, chain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	// Get all unspent outputs of fromAddress
	acc, unspentTxOutputs := chain.FindSpendableOutputs(fromAddress, amount)
	if acc < amount {
		log.Panic("Error: Not enough coins")
	}

	// Create Inputs
	for strTxId, outIdxs := range unspentTxOutputs {
		// Get []byte version of transaction ID
		txId, err := hex.DecodeString(strTxId)
		utils.Handle(err)

		// Append inputs
		for _, outIdx := range outIdxs {
			// Todo: Generate Signature
			input := TxInput{txId,outIdx, fromAddress}
			inputs = append(inputs, input )
		}
	}

	// Create Outputs
	// Firstly, one part to destination
	outputs = append(outputs, TxOutput{amount, toAddress})

	// If there is any left over -> send to fromAddress
	if acc > amount {
		outputs = append(outputs, TxOutput{acc - amount, fromAddress })
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}

func (tx *Transaction) IsCoinBase() bool {
	return tx.Inputs[0].OutTxId == nil  && tx.Inputs[0].OutIndex == -1
}
