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

type TcpReader struct {
	conn        *net.TCPConn
	buf         []byte
	haveReadIdx int64
	haveHandIdx int
}

func NewTcpReader(conn *net.TCPConn) *TcpReader {
	t := &TcpReader{}
	t.conn = conn
	t.buf = make([]byte, 1<<20)
	t.haveReadIdx = 0
	t.haveHandIdx = -1
	return t
}
func (t *TcpReader) GetBytes() ([]byte, error) {
	cnt, err := t.conn.Read(t.buf[t.haveHandIdx+1:])
	var messageLen uint32 = 0
	if err != nil || cnt == 0 {
		return nil, errors.New("xnx")
	}
	t.haveReadIdx += int64(cnt)
	if int64(t.haveHandIdx+MessageHeadLen) > t.haveReadIdx {
		fmt.Println("need next read from remote 1")
		return nil, nil //代表还需要下一次读取
	}
	messageLen = binary.LittleEndian.Uint32(t.buf[t.haveHandIdx+1 : t.haveHandIdx+1+MessageHeadLen])
	if int64(t.haveHandIdx)+int64(messageLen)+MessageHeadLen > t.haveReadIdx {
		fmt.Println("need next read from remote 2")
		return nil, nil
	}

	defer func() {
		t.haveHandIdx += int(MessageHeadLen + messageLen)
	}()
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
