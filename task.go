package hproxy

import (
	"errors"
	"github.com/lazyironf4ur/hproxy/proxy"
	"sync"
)

const (
	DefaultTaskQueueSize = 100
	DefaultWorkerSize    = 20
)

type TaskManager struct {
	workLocker          sync.Mutex
	TaskQueueLocker     sync.Mutex
	ProxyReqQueueLocker sync.Mutex
	ProxyResQueueLocker sync.Mutex
	Workers             []Worker
	RelayTaskQueue      []*Task
	ProxyReqQueue       []*Task
	ProxyResQueue       []*Task
}

type Task struct {
	Hp *proxy.HProtocol
}

func NewDefaultTaskManager() *TaskManager {
	return &TaskManager{
		Workers:        make([]Worker, DefaultWorkerSize),
		RelayTaskQueue: make([]*Task, DefaultTaskQueueSize),
		ProxyReqQueue:  make([]*Task, DefaultTaskQueueSize),
	}
}

func (t *TaskManager) AddRelayTask(task *Task) {
	t.TaskQueueLocker.Lock()
	defer t.TaskQueueLocker.Unlock()
	t.RelayTaskQueue = append(t.RelayTaskQueue, task)
}

func (t *TaskManager) AddProxyReqTask(task *Task) {
	t.ProxyReqQueueLocker.Lock()
	defer t.ProxyReqQueueLocker.Unlock()
	t.ProxyReqQueue = append(t.ProxyReqQueue, task)
}

func (t *TaskManager) AddProxyResTask(task *Task) {
	t.ProxyResQueueLocker.Lock()
	defer t.ProxyResQueueLocker.Unlock()
	t.ProxyResQueue = append(t.ProxyResQueue, task)
}

//func (t *TaskManager) DoTask() error {
//
//	task, err := t.GetRelayNextTaskSafely()
//	if err != nil {
//		return err
//	}
//	for _, worker := range t.Workers {
//		if worker.GetState() == Free {
//			err = worker.DoTask(context.Background(), task)
//			break
//		}
//	}
//	return err
//}

func (t *TaskManager) GetRelayNextTaskSafely() (*Task, error) {
	var task *Task
	t.TaskQueueLocker.Lock()
	defer t.TaskQueueLocker.Unlock()
	if len(t.RelayTaskQueue) > 0 {
		task = t.RelayTaskQueue[0]
		t.RelayTaskQueue = t.RelayTaskQueue[1:]
		return task, nil
	}
	return nil, errors.New("no more task")
}

func (t *TaskManager) GetProxyReqNextTaskSafely() (*Task, error) {
	var task *Task
	t.ProxyReqQueueLocker.Lock()
	defer t.ProxyReqQueueLocker.Unlock()
	if len(t.ProxyReqQueue) > 0 {
		task = t.ProxyReqQueue[0]
		t.ProxyReqQueue = t.ProxyReqQueue[1:]
		return task, nil
	}
	return nil, errors.New("no more task")
}

func (t *TaskManager) GetProxyResNextTaskSafely() (*Task, error) {
	var task *Task
	t.ProxyReqQueueLocker.Lock()
	defer t.ProxyReqQueueLocker.Unlock()
	if len(t.ProxyResQueue) > 0 {
		task = t.ProxyResQueue[0]
		t.ProxyResQueue = t.ProxyResQueue[1:]
		return task, nil
	}
	return nil, errors.New("no more task")
}

//func (t *TaskManager) AddWorker(worker Worker) {
//	t.Workers = append(t.Workers, worker)
//}
//
//func (t *TaskManager) ParseAndAddTask(protocol *proxy.HProtocol) error {
//
//	return nil
//}
