package utility

import (
	"sync"
)

type Worker interface {
	Submit(task *Task)
	Start()
	Stop()
}

type WorkerPool struct {
	numWorkers int
	taskQ      TaskQueue
	stopChan   chan struct{}
	wg         sync.WaitGroup
	running    bool
	cond       *sync.Cond
}

func NewWorkerPool(workers int) *WorkerPool {
	return &WorkerPool{
		numWorkers: workers,
		cond:       sync.NewCond(&sync.Mutex{}),
	}
}

func (wp *WorkerPool) Start() {
	if wp.running {
		return
	}
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			for {
				select {
				case <-wp.stopChan:
					return
				default:
					wp.cond.L.Lock()
					for wp.taskQ.IsEmpty() {
						wp.cond.Wait()
					}
					wp.cond.L.Unlock()
					task := wp.taskQ.Dequeue()
					task.routine(task.staticParam, task.param)
				}
			}
		}()
	}
	wp.running = true
}

func (wp *WorkerPool) Stop() {
	wp.cond.Broadcast()
	close(wp.stopChan)
	wp.wg.Wait()
}

func (wp *WorkerPool) Submit(t *Task) {
	wp.taskQ.Enqueue(t)
	wp.cond.Signal()
}
