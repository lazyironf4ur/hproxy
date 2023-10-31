package http

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/lazyironf4ur/hproxy"
	"github.com/lazyironf4ur/hproxy/proxy"

	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	regx, _ = regexp.Compile("/(\\w)+")
)

const SchemeHead = "http://"

type Work struct {
	appMap map[string]*App
	task   *hproxy.Task
}

func NewHttpWork(task *hproxy.Task, appMap map[string]*App) *Work {
	return &Work{
		appMap: appMap,
		task:   task,
	}
}

func (w Work) DoTask() (*hproxy.Task, error) {
	if w.task == nil {
		return nil, errors.New("nil task to do")
	}
	url, method, origin, body, err := parseTask(w.task)
	loc := regx.FindStringIndex(url)
	appName := url[loc[0]+1 : loc[1]]
	realUrl := url[loc[1]:]
	app := w.appMap[appName]
	if err != nil {
		return nil, err
	}
	fullUrl := SchemeHead + app.Address + ":" + strconv.Itoa(app.Port) + realUrl
	c := http.Client{}
	req, err := http.NewRequest(method, fullUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	b := bytes.NewBuffer(nil)
	res.Write(b)
	resContent := b.Bytes()
	hpb := append(origin, resContent...)
	hp := &proxy.HProtocol{}
	hp.DeSerialize(hpb)
	return &hproxy.Task{Hp: hp}, nil
}

func parseTask(task *hproxy.Task) (url, method string, origin, body []byte, err error) {
	buf := task.Hp.Serialize()
	origin = buf[:8]
	content := buf[8:]
	proto_idx := bytes.IndexAny(content, "\r\n")
	proto := content[:proto_idx]
	arr := strings.Split(string(proto), " ")
	if len(arr) < 3 {
		err = errors.New("bad http request")
		return
	}
	method = arr[0]
	url = arr[1]
	reader := bytes.NewReader(content)
	lineReader := bufio.NewReader(reader)
	var (
		eof   error
		line  []byte
		count int
	)

	for eof == nil {
		line, eof = lineReader.ReadBytes('\n')
		if eof != nil {
			break
		}

		count = count + len(line)
		if bytes.Equal(line, []byte{'\r', '\n'}) {
			break
		}
	}
	body = content[count:]
	return
}
