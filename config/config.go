package config

import (
	"fmt"
	"github.com/lazyironf4ur/hproxy/feature/inbound"
	"github.com/lazyironf4ur/hproxy/feature/outbound"
	"gopkg.in/yaml.v3"
	"os"
)

var GlobalConf = new(GlobalConfig)

func init() {
	//conf.ReadFileConf(defaultConfigPath, GlobalConf)
	fmt.Println(os.Getwd())
	file, err := os.ReadFile("./config_example.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(file, GlobalConf)
	if err != nil {
		panic(err)
	}
}

type GlobalConfig struct {
	InboundConfig  *inbound.Config  `yaml:"inbound"`
	OutBoundConfig *outbound.Config `yaml:"outbound"`
}
