package outbound

import (
	"github.com/lazyironf4ur/hproxy/proxy/http"
)

type Config struct {
	Mode          string      `yaml:"mode"`
	Workers       int         `yaml:"workers"`
	ServerAddress string      `yaml:"server_address"`
	Port          int         `yaml:"port"`
	AppSettings   []*http.App `yaml:"proxyApp"`
}
