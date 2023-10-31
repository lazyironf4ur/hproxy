package common

import (
	"log"
	"os"
)

var glog *log.Logger

func init() {
	glog = log.New(os.Stdout, "", log.LstdFlags|log.Llongfile)
}

func GetLogger() *log.Logger {
	return glog
}
