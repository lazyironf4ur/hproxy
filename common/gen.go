package common

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const Symbol = "hp"

func EncodeCliId(addr string) uint64 {
	a := strings.Split(addr, ":")
	if len(a) < 2 {
		panic("invalid addr")
	}
	ip := a[0]
	port := a[1]
	var (
		cliId uint64
		t     uint64
		o     uint64
	)
	for i, c := range strings.Split(ip, ".") {
		n, _ := strconv.Atoi(c)
		t = t | (uint64(n) << ((5 - i) * 8))
	}

	p, _ := strconv.Atoi(port)
	t = t | uint64(p)

	for i, w := range Symbol {
		o = o | uint64(w)<<((7-i)*8)
	}
	cliId = o | t
	return cliId
}

func ValidateAndDecodeCliId(cliId uint64) (addr string, err error) {
	symbol := cliId >> 48
	symbol1 := byte(symbol)
	symbol2 := byte(symbol >> 8)
	if !strings.EqualFold(Symbol, fmt.Sprintf("%s%s", string(symbol2), string(symbol1))) {
		return "", errors.New("bad cliId")
	}
	port := uint16(cliId)
	o := uint32(cliId >> 16)
	var ip string
	for i := 3; i >= 0; i-- {
		p := strconv.Itoa(int(uint8(o >> (i * 8))))
		if ip == "" {
			ip = fmt.Sprintf("%s", p)
			continue
		}
		ip = fmt.Sprintf("%s.%s", ip, p)
	}
	addr = ip + ":" + strconv.Itoa(int(port))
	return
}
