package gowork

import (
	"fmt"
	"sync"
)

type worker interface {
	run(wg *sync.WaitGroup)
}

type simpleWorker struct {
	id       int
	incoming <-chan func()
	logs     chan<- Log
}

func (w *simpleWorker) run(wg *sync.WaitGroup) {
	w.logs <- Log{Level: LogInfo, Message: fmt.Sprintf("worker %d started", w.id)}
	defer wg.Done()
	for job := range w.incoming {
		job()
	}
	w.logs <- Log{Level: LogInfo, Message: fmt.Sprintf("worker %d stopped", w.id)}
}

func newSimpleWorker(id int, jobs <-chan func(), logs chan<- Log) worker {
	return &simpleWorker{
		id:       id,
		incoming: jobs,
		logs:     logs,
	}
}

type resultingWorker[T any] struct {
	id       int
	incoming <-chan func() T
	results  chan<- T
	logs     chan<- Log
}

func (w *resultingWorker[T]) run(wg *sync.WaitGroup) {
	w.logs <- Log{Level: LogInfo, Message: fmt.Sprintf("resulting worker %d started", w.id)}
	defer wg.Done()
	for job := range w.incoming {
		result := job()
		w.results <- result
	}
	w.logs <- Log{Level: LogInfo, Message: fmt.Sprintf("resulting worker %d stopped", w.id)}
}

func newResultingWorker[T any](id int, jobs <-chan func() T, results chan<- T, logs chan<- Log) worker {
	return &resultingWorker[T]{
		id:       id,
		incoming: jobs,
		results:  results,
		logs:     logs,
	}
}
