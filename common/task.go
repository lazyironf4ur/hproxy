package common

import (
	"log"
	"reflect"
)

type Args interface{}

func RecoveredTask(f interface{}, args ...Args) func() {
	// TODO: consider to specific general func type for f
	var v reflect.Value
	if v = reflect.ValueOf(f); v.Kind() != reflect.Func {
		return nil
	}
	return func() {
		defer func() {
			if exception := recover(); exception != nil {
				log.Printf("do recovered async task error: %v\n", exception)
			}
		}()
		var args_v []reflect.Value
		for _, va := range args {
			args_v = append(args_v, reflect.ValueOf(va))
		}
		v.Call(args_v)
	}

}

func DoAsyncRecoveredTask(f interface{}, args ...Args) {
	go RecoveredTask(f, args)
}
