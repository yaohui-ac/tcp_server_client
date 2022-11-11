package util

import (
	"encoding/binary"
	"fmt"
	"testing"
)

func TestGenWriteMessage(t *testing.T) {
	str := "1+24"
	message := GenWriteMessage(str)
	messageLen := binary.LittleEndian.Uint32(message[0:4])
	fmt.Println(messageLen)
	fmt.Println(string(message[4:]))
}
