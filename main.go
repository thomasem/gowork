package main

import (
	"fmt"
	"time"
)

type WorkFn func() string

type SchuylerSister struct {
	name     string
	incoming <-chan WorkFn
	log      chan<- string
	done     chan<- struct{}
}

func (w *SchuylerSister) Work() {
	for {
		work, ok := <-w.incoming
		if !ok {
			break
		}
		w.Log(work())
		time.Sleep(100 * time.Millisecond)
	}
	w.done <- struct{}{}
}

func (w *SchuylerSister) Log(log string) {
	w.log <- fmt.Sprintf("%s: %s", w.name, log)
}

func NewSchuylerSister(name string, incoming <-chan WorkFn, log chan<- string, done chan<- struct{}) *SchuylerSister {
	return &SchuylerSister{
		name:     name,
		incoming: incoming,
		log:      log,
		done:     done,
	}
}

type SchuylerSisters struct {
	workers  []*SchuylerSister
	outgoing chan<- WorkFn
	done     <-chan struct{}
	log      <-chan string
}

func (wp *SchuylerSisters) DoWork(fns []WorkFn) {
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

func (wp *SchuylerSisters) Log() {
	for {
		log, ok := <-wp.log
		if !ok {
			break
		}
		fmt.Println(log)
	}
}

func NewSchuylerSisters() *SchuylerSisters {
	workCh := make(chan WorkFn)
	logCh := make(chan string)
	doneCh := make(chan struct{})
	sisters := []string{"Angelica", "Peggy", "Eliza"}
	workers := make([]*SchuylerSister, len(sisters))
	for i, name := range sisters {
		workers[i] = NewSchuylerSister(name, workCh, logCh, doneCh)
	}
	return &SchuylerSisters{
		workers:  workers,
		outgoing: workCh,
		log:      logCh,
		done:     doneCh,
	}
}

func main() {
	workFn := func() string {
		return "Work! Work!"
	}
	work := make([]WorkFn, 20)
	for i := 0; i < len(work); i++ {
		work[i] = workFn
	}
	wp := NewSchuylerSisters()
	wp.DoWork(work)
}
