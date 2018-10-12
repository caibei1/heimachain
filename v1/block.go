package main

import (
	"time"
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
)

type Block struct {
	Version int64
	PrevBlockHash []byte
	Hash []byte
	TimeStamp int64
	TargetBits int64  //难度
	Nonce int64
	MerKelRoot []byte
	//Data []byte
	Transactions []*Transaction
}

func NewBlock(transactions []*Transaction ,prevBlockHash []byte) *Block {

	block := &Block{
		Version:1,
		PrevBlockHash:prevBlockHash,
		TimeStamp:time.Now().Unix(),
		TargetBits:targetBits,
		Nonce:0,
		MerKelRoot:[]byte{},
		//Data:[]byte(data),
		Transactions:transactions,
	}
	//fmt.Println(transactions[0].TXOutputs[0].lockScript)
	pow := NewProofOfWork(block)
	nonce,hash := pow.Run()
	block.Hash = hash
	block.Nonce = nonce
	//block.SetHash()
	return block
}



//创世块
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase},[]byte{})
}


//序列化
func (block *Block)Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(block)
	CheckErr(err)
	return buffer.Bytes()
}

//反序列化
func Deserialize(data []byte) *Block {

	if len(data) == 0{
		fmt.Println("数据为空")
		os.Exit(1)
	}

	var block *Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	CheckErr(err)
	return block
}

















