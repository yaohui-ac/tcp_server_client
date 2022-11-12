package util

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	MessageHeadLen = 4
)
type TcpReadStatus int

const (
	Read_Finish TcpReadStatus = 1
	On_Read		TcpReadStatus = 2
)
type TcpReader struct {
	conn        *net.TCPConn
	buf         []byte
	haveReadIdx int64
	haveHandIdx int
	readStatus TcpReadStatus
	
}

func NewTcpReader(conn *net.TCPConn) *TcpReader {
	t := &TcpReader{}
	t.conn = conn
	t.buf = make([]byte, 1<<20)
	t.haveReadIdx = 0
	t.haveHandIdx = -1
	t.readStatus = Read_Finish
	return t
}
func (t *TcpReader) GetBytes() ([]byte, error) {
	if (t.readStatus == Read_Finish) {
		t.haveReadIdx = 0
		t.haveHandIdx = -1
	}

	cnt, err := t.conn.Read(t.buf[t.haveHandIdx+1:])
	var messageLen uint32 = 0
	if err != nil || cnt == 0 {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("remote closed")
	}
	t.haveReadIdx += int64(cnt)
	if int64(t.haveHandIdx+MessageHeadLen) > t.haveReadIdx {
		fmt.Println("need next read from remote 1")
		t.readStatus = On_Read
		return nil, nil //代表还需要下一次读取
	}
	messageLen = binary.LittleEndian.Uint32(t.buf[t.haveHandIdx+1 : t.haveHandIdx+1+MessageHeadLen])
	if int64(t.haveHandIdx)+int64(messageLen)+MessageHeadLen > t.haveReadIdx {
		fmt.Println("need next read from remote 2")
		t.readStatus = On_Read
		return nil, nil
	}

	defer func() {
		t.haveHandIdx += int(MessageHeadLen + messageLen)
	}()
	t.readStatus = Read_Finish //复位
	return t.buf[t.haveHandIdx+1+MessageHeadLen : t.haveHandIdx+1+MessageHeadLen+int(messageLen)], nil

}
func SplitExpression(str string) (num1 int64, num2 int64, err error) {
	//解析表达式
	num1, num2 = 0, 0
	nums := strings.Split(str, "+")
	if len(nums) != 2 {
		return 0, 0, errors.New("str expression error")
	}
	num1, err = strconv.ParseInt(nums[0], 10, 64)
	if err != nil {
		return 0, 0, errors.New("str expression error")
	}
	num2, err = strconv.ParseInt(nums[1], 10, 64)
	if err != nil {
		return 0, 0, errors.New("str expression error")
	}

	return num1, num2, nil
}
func GenWriteMessage(message string) []byte {
	bytebuff := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytebuff, binary.LittleEndian, int32(len(message)))
	_ = binary.Write(bytebuff, binary.LittleEndian, []byte(message))
	return bytebuff.Bytes()
}
