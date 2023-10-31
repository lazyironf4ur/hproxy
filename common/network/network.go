package network

import (
	"errors"
	"io"
	"log"
	"net"
	"strings"
)

const (
	defaultBuffSize = 2048
)

func ReadFromConn(conn net.Conn) []byte {
	var lastBuf []byte
	for {
		buf := make([]byte, defaultBuffSize)
		// will block until conn is closed; TODO: Conn.SetReadDeadline(t time.Time)
		if n, err := conn.Read(buf); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
		} else {
			log.Printf("read: %d", n)
			lastBuf = append(lastBuf, buf[:n]...)
			// read over
			if n < len(buf) {
				break
			}
		}
	}
	return lastBuf
}

func ParseBuf(b []byte) (method, url, protoc_version string) {
	content := string(b)
	protoc := strings.Split(content, "\r\n")
	protoc_line := strings.Split(protoc[0], " ")
	if len(protoc_line) < 3 {
		return
	}
	method = protoc_line[0]
	url = protoc_line[1]
	protoc_version = protoc_line[2]
	return
}
