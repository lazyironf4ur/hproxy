package common

import (
	"fmt"
	"strings"
	"testing"
)

func TestGen(t *testing.T) {
	p := uint64(255)
	c := uint64(256)
	fmt.Println(p | c<<8)
	fmt.Printf("%b\n", p|c<<8)

	r := (p | c<<8) ^ 0xff
	fmt.Printf("%b", r)
}

func TestGenCliId(t *testing.T) {
	addr := "127.0.0.1:34523"
	id := EncodeCliId(addr)
	fmt.Printf("%b\n", id)
}

func TestEncodeAndDecode(t *testing.T) {
	ipaddr := "127.0.0.1:6657"
	id := EncodeCliId(ipaddr)
	addr, err := ValidateAndDecodeCliId(id)
	if err != nil {
		t.Errorf("validate cliid failed, cause: %v", err)
		t.FailNow()
	}

	if !strings.EqualFold(ipaddr, addr) {
		t.Errorf("origin ipaddr does not equal the decode addr")
	}

}

func Test1(t *testing.T) {
	n := 127
	s := fmt.Sprintf("%b", n)
	println(s)
}
