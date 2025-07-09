package utility

import (
	"sync"
)

type Executor interface {
	Submit(task *Task)
	Start()
	Stop()
}

type WorkerPoolExecutor struct {
	id         string
	numWorkers int
	taskQ      TaskQueue
	stopChan   chan struct{}
	wg         sync.WaitGroup
	cond       *sync.Cond
}

func NewWorkerPoolExecutor(id string, workers int) *WorkerPoolExecutor {
	return &WorkerPoolExecutor{
		id:         id,
		numWorkers: workers,
		cond:       sync.NewCond(&sync.Mutex{}),
		stopChan:   make(chan struct{}),
	}
}

func (wp *WorkerPoolExecutor) Start() {
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
						select {
						case <-wp.stopChan:
							return
						default:
						}
						wp.cond.Wait()
					}
					wp.cond.L.Unlock()
					task := wp.taskQ.Dequeue()
					task.routine(task.staticParam, task.param)
				}
			}
		}()
	}
}

func (wp *WorkerPoolExecutor) Stop() {
	close(wp.stopChan)
	wp.cond.Broadcast()
	wp.wg.Wait()
}

func (wp *WorkerPoolExecutor) Submit(t *Task) {
	wp.taskQ.Enqueue(t)
	wp.cond.Signal()
}
