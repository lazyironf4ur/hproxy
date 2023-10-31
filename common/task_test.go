package common

import (
	"testing"
	"time"
)

func TestDoRecoveredAsyncTask(t *testing.T) {
	a := make([]int, 5)
	b := func(a []int) {
		_ = a[6]
	}
	go RecoveredTask(b, a)()
	time.Sleep(10 * time.Second)
}

func TestDoRecoveredAsyncTaskArgOfPoint(t *testing.T) {
	type A struct {
		arr []byte
		len int
	}

	arr := make([]byte, 0)
	a := &A{
		arr: arr,
		len: len(arr),
	}

	fn := func(ap *A) {
		_ = ap.arr[ap.len]
	}
	go RecoveredTask(fn, a)()
	time.Sleep(10 * time.Second)
}
