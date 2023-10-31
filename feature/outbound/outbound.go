package outbound

import (
	"context"
	"fmt"
	"github.com/lazyironf4ur/hproxy"
	"github.com/lazyironf4ur/hproxy/common"
	"github.com/lazyironf4ur/hproxy/common/network"
	"github.com/lazyironf4ur/hproxy/feature/inbound"
	"github.com/lazyironf4ur/hproxy/proxy"
	"github.com/lazyironf4ur/hproxy/proxy/http"
	"net"
	"strconv"
)

const (
	MODE_RELAY        = "relay"
	MODE_PROXY        = "proxy"
	OutBoundConfigKey = "outbound_config"
)

// OutBound
// consumer task in taskManager
type OutBound struct {
	Ctx         context.Context
	Config      *Config
	taskManager *hproxy.TaskManager
	inbound     inbound.InBound
	tunnels     []net.Conn
	appMap      map[string]*http.App
	shutdown    chan struct{}
}

func (o *OutBound) Init(ctx context.Context, inbound inbound.InBound, manager *hproxy.TaskManager) error {
	val := ctx.Value(OutBoundConfigKey)
	if conf, ok := val.(*Config); ok {
		o.Config = conf
	}
	o.appMap = make(map[string]*http.App)
	for _, a := range o.Config.AppSettings {
		o.appMap[a.GetAppName()] = a
	}
	o.tunnels = make([]net.Conn, 10)
	o.inbound = inbound
	o.taskManager = manager
	return nil
}

func (o *OutBound) Start() error {
	var err error
	switch o.Config.Mode {
	case MODE_RELAY:
		err = startRelayOutBound(o)

	case MODE_PROXY:
		err = startProxyOutBound(o)
	default:
	}
	startAsyncConsumer(o)
	return err
}

// start to listen incoming connection, add it to tunnel list
func startRelayOutBound(bound *OutBound) error {
	config := bound.Config
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", config.ServerAddress, strconv.Itoa(config.Port)))
	if err != nil {
		return err
	}

	common.DoAsyncRecoveredTask(func(listener2 net.Listener) {
		for {
			conn, err2 := listener2.Accept()
			if err2 != nil {
				common.GetLogger().Println("accept error, cause: %v", err2)
			}
			bound.tunnels = append(bound.tunnels, conn)
		}

	}, listener)
	return err
}

// Actively connect the relay server to establish tunnel
func startProxyOutBound(bound *OutBound) error {
	config := bound.Config
	common.DoAsyncRecoveredTask(func() {
		for i := 1; i < config.Workers; i++ {
			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", config.ServerAddress, strconv.Itoa(config.Port)))
			if err != nil {
				panic(err)
			}
			bound.tunnels = append(bound.tunnels, conn)
		}

		for _, t := range bound.tunnels {
			conn := t
			// async read
			go func() {
				for {
					select {
					case <-bound.shutdown:
						return
					default:
						buf := network.ReadFromConn(conn)
						hp := &proxy.HProtocol{}
						hp.DeSerialize(buf)
						task := &hproxy.Task{
							Hp: hp,
						}
						bound.taskManager.AddProxyReqTask(task)
					}
				}
			}()
		}

		tm := bound.taskManager
		// async write
		go func() {
			for {
				select {
				case <-bound.shutdown:
					return
				default:
					for task, err := tm.GetProxyResNextTaskSafely(); err == nil; {
						idx := int(task.Hp.Identifier) % len(bound.tunnels)
						tunnel := bound.tunnels[idx]
						go tunnel.Write(task.Hp.Serialize())
					}
				}
			}
		}()

	})

	return nil
}

func startAsyncConsumer(outBound *OutBound) {
	switch outBound.Config.Mode {
	case MODE_RELAY:
		startRelayAsyncSendConsumer(outBound)
		startRelayAsyncRevConsumer(outBound)
	case MODE_PROXY:
		startProxyAsyncSendConsumer(outBound)
	default:
	}
}

func startRelayAsyncSendConsumer(outBound *OutBound) {
	tm := outBound.taskManager
	common.DoAsyncRecoveredTask(func(tm *hproxy.TaskManager) {
		for {
			select {
			case <-outBound.shutdown:
				return
			default:
				for task, err := tm.GetRelayNextTaskSafely(); err == nil; {
					idx := int(task.Hp.Identifier) % len(outBound.tunnels)
					tunnel := outBound.tunnels[idx]
					go tunnel.Write(task.Hp.Serialize())
				}
			}
		}
	}, tm)
}

func startRelayAsyncRevConsumer(outBound *OutBound) {
	for i := 0; i < len(outBound.tunnels); i++ {
		conn := outBound.tunnels[i]
		if ti, ok := (outBound.inbound).(*inbound.TCPInbound); ok {
			go func() {
				for {
					select {
					case <-outBound.shutdown:
						return
					default:
						relayResponse(conn, ti)
					}
				}
			}()
		} else {
			panic("wrong type for outbound.inbound to turn")
		}
	}
}

// monitor
func startProxyAsyncSendConsumer(outBound *OutBound) {
	tm := outBound.taskManager
	common.DoAsyncRecoveredTask(func(tm *hproxy.TaskManager) {
		for {
			select {
			case <-outBound.shutdown:
				return
			default:
				for task, err := tm.GetProxyReqNextTaskSafely(); err == nil; {
					worker := http.NewHttpWork(task, outBound.appMap)
					returnTask, err := worker.DoTask()
					if err != nil {
						common.GetLogger().Println(err)
						continue
					}
					tm.AddProxyResTask(returnTask)
				}
			}
		}
	}, tm)
}

//func startProxyAsyncRevConsumer(outBound *OutBound) {
//	for i := 0; i < len(outBound.tunnels); i++ {
//		conn := outBound.tunnels[i]
//		if ti, ok := (outBound.inbound).(*inbound.TCPInbound); ok {
//			go func() {
//				for {
//					select {
//					case <-outBound.shutdown:
//						return
//					default:
//						relayResponse(conn, ti)
//					}
//				}
//			}()
//		} else {
//			panic("wrong type for outbound.inbound to turn")
//		}
//	}
//}

func relayResponse(conn net.Conn, tcpInbound *inbound.TCPInbound) {
	buf := network.ReadFromConn(conn)
	hp := &proxy.HProtocol{}
	hp.DeSerialize(buf)
	tcpInbound.NotifyRes(hp.Identifier, hp.Payload)
}

func proxyResponse(conn net.Conn) {

}
