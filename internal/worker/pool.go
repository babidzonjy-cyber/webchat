package worker

import (
	"sync"
)

type Pool struct {
	tasks chan func()
	wg    sync.WaitGroup
}

func NewPool(workers int) *Pool {
	p := &Pool{
		tasks: make(chan func(), 256),
	}

	for range workers {
		p.wg.Add(1)
		go p.worker()
	}

	return p
}

func (p *Pool) worker() {
	defer p.wg.Done()

	for task := range p.tasks {
		task()
	}
}

func (p *Pool) Submit(task func()) {
	p.tasks <- task
}

func (p *Pool) Shutdown() {
	close(p.tasks)
	p.wg.Wait()
}
