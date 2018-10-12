package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

func IntToByte(num int64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer,binary.BigEndian,num)
	if err != nil{
		fmt.Println("int 转换数组错误:")
		CheckErr(err)
	}

	return buffer.Bytes()

}

func CheckErr(err error)  {
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
}