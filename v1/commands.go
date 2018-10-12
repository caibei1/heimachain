package main

import "fmt"

func (cli *CLI)CreateChain(addr string)  {
	bc := NewBlockChain(addr)
	defer bc.db.Close()
	fmt.Println("成功创建区块链")
}


//func (cli *CLI)AddBlock(data string)  {
//	cli.bc.AddBlock(data)
//	fmt.Println("success")
//}

func (cli *CLI)PrintChain()  {

	bc := GetBlockChainHandler()
	it := bc.Iterater()


	for  {
		block := it.Next()

		for _,o := range block.Transactions{
			fmt.Println("input: \n",o.TXInputs)
		}

		for _,o := range block.Transactions{
			fmt.Println("output: \n",o.TXOutputs)
		}

		fmt.Println("Version:",block.Version)
		fmt.Printf("PrevBlockHash:%x\n",block.PrevBlockHash)
		fmt.Printf("Hash:%x\n",block.Hash)
		fmt.Printf("TimeStamp:%d\n",block.TimeStamp)
		fmt.Printf("Nonce:%d\n",block.Nonce)
		pow := NewProofOfWork(block)
		fmt.Printf("IsValid: %v \n",pow.IsValid())
		fmt.Println()
		fmt.Println("===============================")
		fmt.Println()
		if len(block.PrevBlockHash) == 0{
			break
		}
	}
}

func (cli *CLI)getBalance(addr string)  {
	bc := GetBlockChainHandler()
	defer bc.db.Close()
	utxos := bc.FindUTXOs(addr)
	var total float64
	for _,utxo := range utxos{
		total += utxo.Value
	}
	fmt.Printf("%v的余额：%f",addr,total)
}

func (cli *CLI)Send(from,to string,amount float64)  {
	bc := GetBlockChainHandler()
	tx := NewTransaction(from,to,amount,bc)

	//for _,o := range tx{
	//	fmt.Println("input: \n",o.TXInputs)
	//}
	//
	//for _,o := range tx{
	//	fmt.Println("output: \n",o.TXOutputs)
	//}

	bc.AddBlock([]*Transaction{tx})
	fmt.Println("转账成功")
}