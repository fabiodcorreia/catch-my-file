package worker

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/fabiodcorreia/catch-my-file/pkg/clog"
)

// WorkerTask represents a block of work that should be executed by the worker.
//
// The context provided will allow cancelation, timeout or deadline of the
// task allowing the called to stop the work execution.
type WorkerTask func(ctx context.Context)

// WorkerPool is the interface that wraps the basic methods
// to implement an WorkerPool.
type WorkerPool interface {
	// Run receives a context and start the WorkerPool. This context will be
	// injected on each WorkerTask.
	Run(ctx context.Context)
	// Stop will prevent new tasks be queued and wait for the workers to finish.
	Stop()
	// AddTask will add a new WorkerTask to the WorkerPool.
	AddTask(task WorkerTask) error
	// IsFull return if all the workers are busy or not.
	IsFull() bool
}

// workerPool is the private implementation of the WorkerPool.
type workerPool struct {
	maxWorkers  int
	wg          sync.WaitGroup
	queuedTasks chan WorkerTask
	running     bool
	busyWorkers int32
}

// NewPool will create a new instance of WorkerPool and return it.
//
// The maxWorkers will set the maximum workers running on the pool.
func NewPool(maxWorkers int) WorkerPool {
	return &workerPool{
		maxWorkers:  maxWorkers,
		queuedTasks: make(chan WorkerTask),
	}
}

// Run will start the WorkerPool and start the workers to execute.
func (p *workerPool) Run(ctx context.Context) {
	for i := 0; i < p.maxWorkers; i++ {
		p.wg.Add(1)
		go func(workerId int) {
			for task := range p.queuedTasks {
				atomic.AddInt32(&p.busyWorkers, 1)
				task(ctx)
				atomic.AddInt32(&p.busyWorkers, -1)
			}
			p.wg.Done()
		}(i + 1)
	}
	p.running = true
}

// Stop will block the WorkePool to prevent new task to arrive and wait
// for all the workers to finish the running task before closing the pool.
func (p *workerPool) Stop() {
	p.running = false
	close(p.queuedTasks) // This will prevent the workers to start new jobs
	p.wg.Wait()
	clog.Info("Worker pool is closed")
}

// IsFull will return true if all the workers are busy running tasks or
// false otherwise.
func (p *workerPool) IsFull() bool {
	return int(p.busyWorkers) >= p.maxWorkers
}

// AddTask will add a new task to the workers queue, if the WorkerPool is
// not running or if all the workers are busy an error is returned.
func (p *workerPool) AddTask(task WorkerTask) error {
	if !p.running {
		return fmt.Errorf("worker is not running")
	}

	if p.IsFull() {
		return fmt.Errorf("worker doesn't have free workers")
	}

	p.queuedTasks <- task
	return nil
}
