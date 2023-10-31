package main

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var blockChan chan os.Signal

func main() {
	blockChan = make(chan os.Signal, 1)
	signal.Notify(blockChan, os.Interrupt, syscall.SIGUSR1)
	listen, err := net.Listen("tcp", "0.0.0.0:9998")
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Printf("accept error: %v\n", err)
				continue
			}
			go handleConn(conn)
		}
	}()

	<-blockChan
}

func handleConn(conn net.Conn) {
	defer func() {
		if exception := recover(); exception != nil {
			log.Printf("panic recovered: %v", exception)
		}
	}()
	//var err error
	var lastbuf []byte
	for {
		buf := make([]byte, 1024)
		if n, err := conn.Read(buf); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
		} else {
			log.Printf("read: %d", n)
			lastbuf = append(lastbuf, buf[:n]...)
			if n < len(buf) {
				break
			}
		}

	}
	method, url, version := parseBuf(lastbuf)
	log.Printf("method: %s, url: %s, version: %s", method, url, version)
	conn.Write([]byte{'o', 'k'})
	conn.Close()
}

func parseBuf(b []byte) (method, url, protoc_version string) {
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
