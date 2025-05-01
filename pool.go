package gowork

import (
	"runtime"
	"sync"
)

type basePool interface {
	Start()
	Stop()
	registerCleanup(func())
	setSize(int)
	setLogger(Logger)
}

type Pool interface {
	basePool
	Submit(func())
}

type ResultingPool[T any] interface {
	basePool
	Submit(func() T)
}

type pool[T any] struct {
	jobs    chan T
	logs    chan Log
	logger  Logger
	size    int
	workers []worker
	workWg  *sync.WaitGroup
	logWg   *sync.WaitGroup
	cleanup []func()
}

func (p *pool[T]) Start() {
	p.workWg.Add(len(p.workers))
	p.logWg.Add(1)
	for _, w := range p.workers {
		go w.run(p.workWg)
	}
	go p.handleLog()
}

func (p *pool[T]) Stop() {
	close(p.jobs)
	p.workWg.Wait() // wait for all work to finish after closing jobs channel
	close(p.logs)
	p.logWg.Wait()
	for _, f := range p.cleanup {
		f()
	}
}

func (p *pool[T]) registerCleanup(f func()) {
	p.cleanup = append(p.cleanup, f)
}

func (p *pool[T]) setSize(size int) {
	p.size = size
}

func (p *pool[T]) setLogger(logger Logger) {
	p.logger = logger
}

func (p *pool[T]) handleLog() {
	defer p.logWg.Done()
	for l := range p.logs {
		p.logger.Log(l)
	}
}

func (p *pool[T]) Submit(j T) {
	p.jobs <- j
}

type Option func(basePool)

func WithSize(size int) Option {
	return func(p basePool) {
		p.setSize(size)
	}
}

func WithLogger(logger Logger) Option {
	return func(p basePool) {
		p.setLogger(logger)
	}
}

func newWorkerPool[F any](opts ...Option) *pool[F] {
	p := &pool[F]{
		logs:    make(chan Log),
		logger:  NewDefaultLogger(LogInfo),
		workWg:  &sync.WaitGroup{},
		logWg:   &sync.WaitGroup{},
		size:    runtime.NumCPU() / 2,
		cleanup: []func(){},
	}

	for _, opt := range opts {
		opt(p)
	}

	p.workers = make([]worker, p.size)
	return p
}

func NewPool(opts ...Option) Pool {
	p := newWorkerPool[func()](opts...)
	p.jobs = make(chan func(), p.size)
	for i := range p.workers {
		p.workers[i] = newSimpleWorker(i+1, p.jobs, p.logs)
	}

	return p
}

func NewResultingPool[T any](opts ...Option) (<-chan T, ResultingPool[T]) {
	p := newWorkerPool[func() T](opts...)
	results := make(chan T)
	p.jobs = make(chan func() T, p.size)

	for i := range p.workers {
		p.workers[i] = newResultingWorker(i+1, p.jobs, results, p.logs)
	}

	p.registerCleanup(func() { close(results) })

	return results, p
}
