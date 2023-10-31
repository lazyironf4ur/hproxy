package inbound

import "context"

type InBound interface {
	Init(ctx context.Context) InBound
	Start(ctx context.Context) error
	Close(ctx context.Context) error
}

//type InBound struct {
//	Ctx          context.Context
//	Conn         net.Conn
//	ReadableChan chan []byte
//	WritableChan chan []byte
//}
