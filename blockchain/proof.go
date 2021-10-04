package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math"
	"math/big"
)

// Take the data from the block

// create a counter (nonce) which starts at 0

// create a hash of the data plus the counter

// check the hash to see if it meets a set of requirements

// Requirements:
// The First few bytes must contain 0s

const Difficulty = 12

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

// CreatePow constructor
func CreatePow(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))
	// target now will have binary form: 0..(nums 0 = Difficulty)....1....0

	pow := &ProofOfWork{b, target}
	return pow
}

// ToHex convert int into []byte
func toHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func (pow *ProofOfWork) combineDataForPow(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.Block.Data,
		pow.Block.PrevHash,
		toHex(int64(nonce)),
		toHex(int64(Difficulty)),
	}, []byte{})

	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var nonce = 0
	var hash [32]byte
	var intHash big.Int

	for nonce < math.MaxInt64 {
		data := pow.combineDataForPow(nonce)
		hash = sha256.Sum256(data)
		//fmt.Printf("Try nonce=%d hash=%x\n", nonce, hash)
		//fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		// if intHash < Target which mean first few byte of intHash contain number of 0 == Difficulty
		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.combineDataForPow(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}
