package http

import "github.com/lazyironf4ur/hproxy/proxy"

type Setting struct {
	Address string
	Port    int
	Url     string
}

func (s Setting) Protocol() proxy.Protocol {
	return proxy.HTTP
}
