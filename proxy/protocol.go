package proxy

import (
	"unsafe"
)

type Protocol string

var (
	count uint64
)

const (
	HTTP = "http"
)

type AppSetting interface {
	GetAppName() string
	DAddress() string
	DPort() int
}

// HProtocol origin client len(8bit):ip(32bit):port(16bit):appName(len*8bit)
type HProtocol struct {
	Identifier uint64
	//AppName string
	//RAddr   string
	//RPort   int
	Payload []byte
}

func (h *HProtocol) DeSerialize(buf []byte) *HProtocol {
	identifierBytes := buf[:8]
	identifier := *(*uint64)(unsafe.Pointer(&identifierBytes))
	h.Identifier = identifier
	payload := buf[8:]
	h.Payload = payload
	return h
}

func (h *HProtocol) Serialize() []byte {
	idb := *(*[8]byte)(unsafe.Pointer(&h.Identifier))
	newBuf := idb[:]
	newBuf = append(newBuf, h.Payload...)
	return newBuf
}
