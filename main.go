package main

import (
	"fmt"
	"time"
)

type WorkFn func() string

type Worker struct {
	id       int
	incoming <-chan WorkFn
	log      chan<- string
	done     chan<- struct{}
}

func (w *Worker) Work() {
	w.Log("coming online...")
	for {
		work, ok := <-w.incoming
		if !ok {
			break
		}
		w.Log(work())
		w.Log("resting for a little now...")
		time.Sleep(100 * time.Millisecond)
	}
	w.Log("all done!")
	w.done <- struct{}{}
}

func (w *Worker) Log(log string) {
	w.log <- fmt.Sprintf("Worker %v: %s", w.id, log)
}

func NewWorker(id int, incoming <-chan WorkFn, log chan<- string, done chan<- struct{}) *Worker {
	return &Worker{
		id:       id,
		incoming: incoming,
		log:      log,
		done:     done,
	}
}

type WorkerPool struct {
	workers  []*Worker
	outgoing chan<- WorkFn
	done     <-chan struct{}
	log      <-chan string
}

func (wp *WorkerPool) DoWork(fns []WorkFn) {
	go wp.Log()
	for _, w := range wp.workers {
		go w.Work()
	}
	for _, fn := range fns {
		wp.outgoing <- fn
	}
	close(wp.outgoing)
	count := 0
	for count < len(wp.workers) {
		<-wp.done
		count++
	}
}

func (wp *WorkerPool) Log() {
	for {
		log, ok := <-wp.log
		if !ok {
			break
		}
		fmt.Println(log)
	}
}

func NewWorkerPool(count int) *WorkerPool {
	workCh := make(chan WorkFn)
	logCh := make(chan string)
	doneCh := make(chan struct{})
	workers := make([]*Worker, count)
	for i := 0; i < len(workers); i++ {
		workers[i] = NewWorker(i, workCh, logCh, doneCh)
	}
	return &WorkerPool{
		workers:  workers,
		outgoing: workCh,
		log:      logCh,
		done:     doneCh,
	}
}

func main() {
	workFn := func() string {
		return "All work and no play makes Jack a dull boy"
	}
	work := make([]WorkFn, 100)
	for i := 0; i < len(work); i++ {
		work[i] = workFn
	}
	wp := NewWorkerPool(5)
	wp.DoWork(work)
}
