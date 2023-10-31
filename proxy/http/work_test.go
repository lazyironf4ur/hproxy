package http

import (
	"fmt"
	"github.com/lazyironf4ur/hproxy"
	"github.com/lazyironf4ur/hproxy/proxy"
	"testing"
)

func TestWorker(t *testing.T) {
	hp := &proxy.HProtocol{
		Identifier: uint64(16516184948849),
		Payload:    []byte("GET /api/asdljjTTASASDACasdgfga/WEQWTQWDASdasdasgg?param1=dasdagfasg HTTP/1.1\r\nContentType=application/json\r\n\r\nGASGASGWRQWEQWEQWDASDASGGASGASGGGGGGGGGASDASDCXCVVSDFSDGSDGSDGSDHSDHDH"),
	}
	task := &hproxy.Task{Hp: hp}
	appMap := make(map[string]*App)
	appMap["api"] = &App{
		AppName: "api",
		Address: "127.0.0.1",
		Port:    8088,
	}
	work := Work{
		appMap: appMap,
		task:   task,
	}
	resTask, err := work.DoTask()
	if err != nil {
		t.Logf("do task err: %v", err)
	}
	fmt.Printf("resTask content: %s", string(resTask.Hp.Serialize()))
}
