package main

import (
	"crypto/sha256"
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
)

const reward  = 12.5

type Transaction struct {
	TXID []byte //交易id
	TXInputs []Input
	TXOutputs []Output
}

type Input struct {
	Txid []byte //交易id
	ReferOutputIndex int64  //索引
	UnlockScript string    //解锁脚本
}

type Output struct {
	Value float64
	LockScript string    //解定脚本
}


func (input *Input) CanUnlockUTXOByAddress(unlockdata string) bool {
	return input.UnlockScript == unlockdata
}

func (output *Output)CanBeUnlockedByAddr(unlockdata string) bool {
	return output.LockScript == unlockdata
}

func (tx *Transaction) IsCoinbase () bool {
	if len(tx.TXInputs) == 1{
		if tx.TXInputs[0].Txid == nil && tx.TXInputs[0].ReferOutputIndex ==-1{
			return true
		}
	}
	return false
}


//设置交易id
func (tx *Transaction)SetTXID()  {
	//data := bytes.Join([][]byte{},[]byte{})
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(tx)
	hash := sha256.Sum256(buffer.Bytes())
	tx.TXID = hash[:]
}


//创建交易
func NewCoinbaseTX(addr string,data string) *Transaction {

	//if data == ""{
	//	fmt.Sprintf(data,"current reward is : %f\n",reward)
	//}

	input := Input{nil,-1,data}
	inputs := []Input{input}

	output := Output{reward,addr}
	outputs := []Output{output}

	tx := Transaction{nil,inputs,outputs}
	//fmt.Println(tx.TXOutputs[0].lockScript)
	tx.SetTXID()
	return &tx
}


//找到地址相关的交易
func (bc *BlockChain)FoundUnspandTransaction(addr string) []Transaction {
	var transactions  []Transaction
	var spentUTXOs = make(map[string/*交易的txid*/][]int64)
	bci := bc.Iterater()
	//倒序遍历
	for{
		block := bci.Next()

		for _,tx := range block.Transactions{

			//遍历当前交易的inputs  找到当前地址消耗的utxo
			for _,input := range tx.TXInputs{
				if tx.IsCoinbase() == false{
					if input.CanUnlockUTXOByAddress(addr){
						spentUTXOs[string(input.Txid)] = append(spentUTXOs[string(input.Txid)],input.ReferOutputIndex)
					}
				}
			}


		//for _,tx := range block.Transactions{

			LABLE1:
			//遍历当前交易的outputs  通过output解锁条件  确定满足条件的交易
			for outputIndex,output := range tx.TXOutputs{

				//判断是否被消耗
				if spentUTXOs[string(tx.TXID)] != nil{
					for _,usedIndex := range spentUTXOs[string(tx.TXID)]{
						if int64(outputIndex) == usedIndex{
							continue LABLE1
						}
					}
				}

				if output.CanBeUnlockedByAddr(addr){

					//????????????是否多次添加
					transactions = append(transactions,*tx)
				}
			}
		}

		if len(block.PrevBlockHash) == 0{
			break
		}
	}

	return transactions
}


//查找未消耗的outputs
func (bc *BlockChain)FindUTXOs(addr string) []Output  {
	var outputs []Output
	txs := bc.FoundUnspandTransaction(addr)
	for _,tx := range txs{
		for _,output := range tx.TXOutputs {
			if output.CanBeUnlockedByAddr(addr){
				outputs = append(outputs,output)
			}
		}
	}
	return outputs
}


//创建一个交易
func NewTransaction(from,to string,amount float64,bc *BlockChain) *Transaction {

	counted,container := bc.FindSuitableUTXOs(from,amount)
	if counted < amount{
		fmt.Println("余额不足")
		os.Exit(1)
	}

	var inputs []Input
	//遍历  构造input
	for txid,outputIndexs := range container{
		for _,outputIndex := range outputIndexs{
			input := Input{[]byte(txid),outputIndex,from}
			inputs = append(inputs,input)
		}
	}

	var outputs []Output
	//构造outputs
	if amount == counted{
		output := Output{amount,to}
		outputs = append(outputs,output)
	}else {
		output1 := Output{amount,to}
		output2 := Output{counted-amount,from}
		outputs = append(outputs,output1)
		outputs = append(outputs,output2)
	}

	tx := Transaction{nil,inputs,outputs}
	tx.SetTXID()
	return &tx
}

