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
	NumWorkers int
	taskQ      TaskQueue
	stopChan   chan struct{}
	wg         sync.WaitGroup
	running    bool
}

func (wp *WorkerPool) Start() {
	if wp.running {
		return
	}
	for i := 0; i < wp.NumWorkers; i++ {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			for {
				select {
				case <-wp.stopChan:
					return
				default:
					task := wp.taskQ.Dequeue()
					if task != nil && task.routine != nil {
						task.routine(task.staticParam, task.param)
					}
				}
			}
		}()
	}
	wp.running = true
}

// when to stop this, refactor all
func (wp *WorkerPool) Stop() {
	close(wp.stopChan)
	wp.wg.Wait()
}

func (wp *WorkerPool) Submit(t *Task) {
	wp.taskQ.Enqueue(t)
}
