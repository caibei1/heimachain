package main

import (
	"os"
	"github.com/boltdb/bolt"
	"fmt"
)

const dbfile  = "blockChainDB.db"

const blockBucket  = "block"

const lasthash  = "lasthash"

const genesisBlockInfo = "The Time 13/July/2018"


type BlockChain struct {
	//block []*Block
	db *bolt.DB
	lastHash []byte
}


//创建含有创世块的区块链
func NewBlockChain(addr string) *BlockChain {

	if IsBlockChainExist(){
		fmt.Println("区块链已存在")
		os.Exit(1)
	}

	db,err := bolt.Open(dbfile,0600,nil)
	CheckErr(err)

	var lastHash []byte

	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))

		coinbase := NewCoinbaseTX(addr,genesisBlockInfo)
		genesis := NewGenesisBlock(coinbase)

		bucket,err := tx.CreateBucket([]byte(blockBucket))
		CheckErr(err)

		err = bucket.Put(genesis.Hash,genesis.Serialize())
		CheckErr(err)
		//err = bucket.Put([]byte("wsj"),[]byte(genesis.Transactions[0].TXOutputs[0].LockScript))
		//CheckErr(err)

		err = bucket.Put([]byte(lasthash),genesis.Hash)
		CheckErr(err)
		lastHash = genesis.Hash

		return nil
	})
	CheckErr(err)
	return &BlockChain{db,lastHash}
}


//获取BlockChain
func GetBlockChainHandler() *BlockChain {
	if !IsBlockChainExist(){
		fmt.Println("区块链不存在  请先创建")
		os.Exit(1)
	}
	db,err := bolt.Open(dbfile,0600,nil)
	CheckErr(err)

	var lastHash []byte

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket != nil{
			lastHash = bucket.Get([]byte(lasthash))
		}
		return nil
	})

	return &BlockChain{db,lastHash}
}




//添加区块
func (bc *BlockChain)AddBlock(transaction []*Transaction)  {
	var prevBlockHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		lastHash := bucket.Get([]byte(lasthash))
		prevBlockHash = lastHash
		return nil
	})

	CheckErr(err)

	block := NewBlock(transaction,prevBlockHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		err = bucket.Put(block.Hash,block.Serialize())
		CheckErr(err)

		err = bucket.Put([]byte(lasthash),block.Hash)
		CheckErr(err)
		return nil
	})
	CheckErr(err)

}

type BlockChainIterator struct {
	db *bolt.DB
	currentHash []byte
}

func (bc *BlockChain) Iterater()*BlockChainIterator  {
	return &BlockChainIterator{bc.db,bc.lastHash}
}

//
//func (it *BlockChainIterator) Test() *[]byte {
//	var b []byte
//	it.db.View(func(tx *bolt.Tx) error {
//		bucket := tx.Bucket([]byte(blockBucket))
//		b = bucket.Get([]byte("wsj"))
//		return nil
//	})
//	return &b
//}

//迭代器
func (it *BlockChainIterator) Next() *Block {
	var block *Block
	err := it.db.View(func(tx *bolt.Tx) error {
		bucket :=tx.Bucket([]byte(blockBucket))
		if bucket == nil{
			os.Exit(1)
		}else {
			byteBlock :=bucket.Get(it.currentHash)
			block = Deserialize(byteBlock)
			it.currentHash = block.PrevBlockHash
		}

		return nil
	})
	CheckErr(err)
	return block
}

//func (it *BlockChainIterator) HasNext() bool {
//	var block *Block
//	err := it.db.View(func(tx *bolt.Tx) error {
//		bucket :=tx.Bucket([]byte(blockBucket))
//		byteBlock :=bucket.Get(it.currentHash)
//		block = Deserialize(byteBlock)
//		//hash := block.PrevBlockHash
//		return nil
//	})
//	CheckErr(err)
//	return string(block.PrevBlockHash) != string([]byte{})
//}


//检查文件是否存在
func IsBlockChainExist() bool {
	_,err := os.Stat(dbfile)
	if os.IsNotExist(err){
		return false
	}
	return true
}


//查找所需要的utxo
func (bc *BlockChain)FindSuitableUTXOs(addr string,amount float64) (float64,map[string][]int64) {
	txs := bc.FoundUnspandTransaction(addr)
	var countTotal float64
	var container = make(map[string][]int64)
	for _,tx := range txs{
		for index,output := range tx.TXOutputs{
			if output.CanBeUnlockedByAddr(addr){
				countTotal += output.Value
				container[string(tx.TXID)] = append(container[string(tx.TXID)],int64(index))
				if countTotal >= amount{
					break
				}
			}
		}
	}
	//if countTotal < amount{
	//	//余额不足
	//}
	return countTotal,container
}









