package common

import "sync"

// FuncWorker 任务
type FuncWorker interface {
	Task()
}

// GoroutinePool 协程池
type GoroutinePool struct {
	work chan FuncWorker
	wg   sync.WaitGroup
}

// NewGoroutinePool 新建协程池
func NewGoroutinePool(maxGoroutines int) *GoroutinePool {
	p := GoroutinePool{
		work: make(chan FuncWorker),
	}

	p.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		go func() {
			for w := range p.work {
				w.Task()
			}
			p.wg.Done()
		}()
	}

	return &p
}

// Run 协程池执行任务
func (p *GoroutinePool) Run(w FuncWorker) {
	p.work <- w
}

// Shutdown 关闭协程池
func (p *GoroutinePool) Shutdown() {
	close(p.work)
	p.wg.Wait()
}
