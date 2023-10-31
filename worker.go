package hproxy

import "context"

type State int

const (
	Processing State = iota
	Free
	Initialization
	Uninit
)

type Worker interface {
	GetState() State
	DoTask(ctx context.Context, task *Task) error
	Init() error
}
