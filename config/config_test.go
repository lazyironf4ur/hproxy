package config

import "testing"

func TestLoadConfig(t *testing.T) {
	global := GlobalConf
	if global == nil {
		t.Fail()
	}

}
