package main

import (
	"fmt"
	"tcp_server_client/util"
	"math/rand"
	"net"
	"strconv"
	"time"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("参数数量不正确")
		return
	}
	remoteAddr := os.Args[1]
	tcpAddr, err := net.ResolveTCPAddr("tcp", remoteAddr)
	if err != nil {

	}
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {

	}
	tcpReceiver := util.NewTcpReader(tcpConn)
	rand.NewSource(time.Now().Unix())
	var bytes []byte
	for i := 0; i <= 20; i++ {
		num1 := rand.Int() % 1000
		num2 := rand.Int() % 1000
		message := strconv.FormatInt(int64(num1), 10) + "+" + strconv.FormatInt(int64(num2), 10)
		fmt.Println("send: "+ message)
		_, err = tcpConn.Write(util.GenWriteMessage(message))
		if err != nil {
			break
		}

		for len(bytes) == 0 { //持续读取 直到返回数据
			bytes, err = tcpReceiver.GetBytes()
			if err != nil {
				fmt.Println("get bytes error")
				_ = tcpConn.Close()
				return
			}
		}
		fmt.Printf("%d + %d = %s\n", num1, num2, string(bytes))
		bytes = make([]byte, 0)
	}
	_ = tcpConn.Close()
}
