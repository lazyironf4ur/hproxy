package inbound

import (
	"context"
	"github.com/lazyironf4ur/hproxy"
	"github.com/lazyironf4ur/hproxy/common"
	"github.com/lazyironf4ur/hproxy/common/network"
	"github.com/lazyironf4ur/hproxy/proxy"
	"net"
	"strconv"
	"sync"
)

const (
	InBoundConfigKey = "inbound_config"
)

// TCPInbound
// Supply some abilities like monitor client requests, manage client conn
type TCPInbound struct {
	sync.Mutex
	InBoundConfig *Config
	listener      net.Listener
	// client map
	cliMap         map[uint64]net.Conn
	resNotifierMap map[uint64]*ResNotifier

	shutdown    chan struct{}
	taskManager *hproxy.TaskManager
}

type ResNotifier struct {
	Id      uint64
	Done    chan struct{}
	Payload []byte
}

func (t *TCPInbound) Init(ctx context.Context) InBound {
	val := ctx.Value(InBoundConfigKey)
	if conf, ok := val.(*Config); ok {
		t.InBoundConfig = conf
	}
	t.resNotifierMap = make(map[uint64]*ResNotifier)
	t.shutdown = make(chan struct{}, 1)
	t.taskManager = hproxy.NewDefaultTaskManager()
	return t
}

func (t *TCPInbound) Start(ctx context.Context) error {

	listener, err := net.Listen("tcp", t.InBoundConfig.Address+strconv.Itoa(t.InBoundConfig.Port))
	if err != nil {
		return err
	}
	t.listener = listener
	common.DoAsyncRecoveredTask(t.listenerHandler, listener)
	t.startAsyncConsumer()
	return nil
}

func (t *TCPInbound) Close(ctx context.Context) error {
	var err error
	t.Lock()
	defer t.Unlock()
	err = t.listener.Close()
	for _, c := range t.cliMap {
		err = c.Close()
	}
	return err
}

func (t *TCPInbound) listenerHandler(listener net.Listener) {
	handler := func() error {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		id := common.EncodeCliId(conn.RemoteAddr().String())
		t.Lock()
		t.cliMap[id] = conn
		rn := &ResNotifier{
			Id:   id,
			Done: make(chan struct{}, 1),
		}
		t.resNotifierMap[id] = rn

		// Async submit task
		go func() {
			fromConn := network.ReadFromConn(conn)
			task := &hproxy.Task{
				Hp: &proxy.HProtocol{
					Identifier: id,
					Payload:    fromConn,
				},
			}
			t.taskManager.AddRelayTask(task)
		}()
		return nil
	}
	for {
		select {
		case <-t.shutdown:
			err2 := t.Close(context.Background())
			if err2 != nil {
				common.GetLogger().Printf("TCPInBound close error: %v", err2)
			}
			return
		default:
			err := handler()
			common.GetLogger().Printf("listen accept error: %v", err)
		}
	}
}

func (t *TCPInbound) startAsyncConsumer() {
	resHandler := func() {
		for _, n := range t.resNotifierMap {
			if _, ok := <-n.Done; ok {
				conn := t.cliMap[n.Id]
				go response(conn, n.Payload)
			}
		}
	}
	go func() {
		for {
			select {
			case <-t.shutdown:
				return
			default:
				resHandler()
			}
		}
	}()
}

func response(conn net.Conn, payload []byte) {
	conn.Write(payload)
}

func (t *TCPInbound) NotifyRes(id uint64, payload []byte) {
	rn := t.resNotifierMap[id]
	rn.Payload = payload
	rn.Done <- struct{}{}
}
