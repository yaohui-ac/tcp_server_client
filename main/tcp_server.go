package main

import (
	"fmt"
	"tcp_server_client/util"
	"net"
	"os"
	"strconv"
)

var enableNagle = true

func HandleTcpConnection(conn *net.TCPConn) {
	tcpReceiver := util.NewTcpReader(conn)
	if conn == nil {
		return
	}
	defer func () {
		fmt.Printf("%s closed.\n",conn.RemoteAddr().String())
		_ = conn.Close()
	}()
	for { //循环读取数据并发送
		bytes, err := tcpReceiver.GetBytes()
		if err != nil {
			fmt.Printf("connection %s: status: %v .\n", conn.RemoteAddr().String(), err)
			break
		}
		if bytes == nil {
			continue
		}
		num1, num2, err := util.SplitExpression(string(bytes))
		if err != nil {
			num1, num2 = 0, 0
			//非法 置为0
		}
		sum := num1 + num2

		fmt.Printf("%s: %d + %d = %d\n", conn.RemoteAddr().String(), num1, num2, sum) //log信息

		_, err = conn.Write(util.GenWriteMessage(strconv.FormatInt(sum, 10)))
		if err != nil {
			break
		}

	}

}
func main() {
	fmt.Println("tcp server start.")
	if len(os.Args) < 2 {
		fmt.Println("参数数量不正确")
		return
	}
	
	Addr := os.Args[1]
	if len(os.Args) < 3 || os.Args[2] == "nodelay" {
		// 禁用nagle
		enableNagle = false

	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", Addr)
	if err != nil {
		fmt.Println("ResolveTCPAddr error")
		return
	}
	listenFd, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("ListenTCP error")
		return
	}
	fmt.Println("begin to listen.")
	for {
		connFd, err := listenFd.AcceptTCP()
		if enableNagle {
			err = connFd.SetNoDelay(!enableNagle)
			if err != nil {
				fmt.Println("connFd.SetNoDelay error")
				continue
			}
		}
		if err != nil {
			fmt.Println("accept error")
			continue
		}
		fmt.Printf("connect to %s\n", connFd.RemoteAddr().String())
		go HandleTcpConnection(connFd)
	}

}
