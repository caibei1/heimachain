package main

import (
	"os"
	"fmt"
	"flag"
)

const Usage  = `
	creatchain -address ADDRESS  "create block chain"  
	send -from FROM -to TO -amount AMOUNT    "make a transaction"
	go build ./ && v1.exe printchain		  "print all blocks"
	getbalance -address ADDRESS   "get balance"
`

type CLI struct {
	//bc *BlockChain
}



func (cli *CLI)run()  {

	//fmt.Println("进入客户端命令")

	//设置命令
	addBlockCmd := flag.NewFlagSet("addBlock",flag.ExitOnError)
	printCmd := flag.NewFlagSet("printchain",flag.ExitOnError)
	creatChainCmd := flag.NewFlagSet("creatchain",flag.ExitOnError)
	getbalanceCmd := flag.NewFlagSet("getbalance",flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send",flag.ExitOnError)

	//获取参数
	addBlockCmdPara :=addBlockCmd.String("data","", "block info")
	creatChainCmdPara :=creatChainCmd.String("address","", "address data")
	getbalanceCmdPara :=getbalanceCmd.String("address","", "address data")

	//创建交易的3个参数
	fromPara :=sendCmd.String("from","", "from address data")
	toPara :=sendCmd.String("to","", "to address data")
	amountPara :=sendCmd.Float64("amount",0, "amount value")


	if len(os.Args)>=2{
		//fmt.Println("判断命令:", os.Args[1])
		switch os.Args[1] {

		case "creatchain":
			err := creatChainCmd.Parse(os.Args[2:])
			CheckErr(err)
			if  creatChainCmd.Parsed(){
				if *creatChainCmdPara != ""{
					cli.CreateChain(*creatChainCmdPara)
				}else {
					fmt.Println("数据为空")
				}
			}
		
		case "addBlock":
			err := addBlockCmd.Parse(os.Args[2:])
			CheckErr(err)
			if addBlockCmd.Parsed(){
				if *addBlockCmdPara != ""{
					fmt.Println("插入数据：",*addBlockCmdPara)
					//cli.AddBlock(*addBlockCmdPara)
					//cli.AddBlock("aa")
				}else {
					fmt.Println("数据为空")
				}
			}

		case "printchain":
			//fmt.Println("查询")
			err := printCmd.Parse(os.Args[2:])
			CheckErr(err)
			if printCmd.Parsed(){
				cli.PrintChain()
			}

		case "getbalance":
			err := getbalanceCmd.Parse(os.Args[2:])
			CheckErr(err)
			if  getbalanceCmd.Parsed(){
				if *getbalanceCmdPara != ""{
					cli.getBalance(*getbalanceCmdPara)
				}else {
					fmt.Println("数据为空")
				}
			}

		case "send":
			err := sendCmd.Parse(os.Args[2:])
			CheckErr(err)
			if sendCmd.Parsed(){
				if *fromPara == "" || *toPara =="" || *amountPara ==0{
					fmt.Println(Usage)
					os.Exit(1)
				}else {
					cli.Send(*fromPara,*toPara,*amountPara)
				}
			}

		default:
			fmt.Printf("invalid cmd: %v\n",Usage)

		}
	}


}