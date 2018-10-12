package main

import (
	"math/big"
	"bytes"
	"math"
	"crypto/sha256"
	"fmt"
)

const targetBits  = 24


type ProofOfWork struct {
	block *Block
	targetBit *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	var IntTarget  = big.NewInt(1)

	//左移256-24位
	IntTarget.Lsh(IntTarget,uint(256 - targetBits))

	return &ProofOfWork{block,IntTarget}
}


//连接数据成byte切片
func (pow *ProofOfWork)PrepareRawData(nonce int64)[]byte {

	block := pow.block

	tmp := [][]byte{
		IntToByte(block.Version),
		block.PrevBlockHash,
		IntToByte(block.TimeStamp),
		block.MerKelRoot,
		IntToByte(nonce),
		IntToByte(targetBits),
		//block.Transactions,//TO DO
	}
	//通过空切片进行数据连接
	data := bytes.Join(tmp, []byte{})

	return data
}


//hash碰撞
func (pow *ProofOfWork)Run() (int64,[]byte) {
	var nonce  int64
	var hash  [32]byte
	var hashInt big.Int

	fmt.Println("开始挖矿...")
	//fmt.Printf("target hash:0000%x\n",pow.targetBit.Bytes())
	for nonce < math.MaxInt64{
		data := pow.PrepareRawData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(pow.targetBit) == 1{
			fmt.Printf("found hash %x\n",hash)
			break
		}else {
			nonce++
		}
	}
	return nonce,hash[:]
}

//验证hash是否正确
func (pow *ProofOfWork) IsValid() bool {
	data := pow.PrepareRawData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	var intHash big.Int
	intHash.SetBytes(hash[:])
	//(string(hash[:]) == string(pow.block.Hash))
	fmt.Printf("Hash:%x\n",hash)
	return intHash.Cmp(pow.targetBit) == -1 && (string(hash[:]) == string(pow.block.Hash))
}