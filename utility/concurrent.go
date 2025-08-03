package utility

import (
	"sync"
	"waf-tester/logger"
)

type Executor interface {
	GetId() string
	Submit(task *Task)
	Start() error
	Finish()
	TerminateGracefully() error
}

type WorkerPoolExecutor struct {
	ID         string
	name       string
	numWorkers int
	taskQ      TaskQueue
	wg         sync.WaitGroup
	logger     logger.Logger
	terminate  chan struct{}
}

func NewWorkerPoolExecutor(name string, workers int, logger logger.Logger) Executor {
	wrk := &WorkerPoolExecutor{
		name:       name,
		numWorkers: workers,
		logger:     logger,
		terminate:  make(chan struct{}),
	}
	wrk.ID = PlContext.generateWorkerKey(wrk)
	return wrk
}

func (wp *WorkerPoolExecutor) GetId() string {
	return wp.ID
}

func (wp *WorkerPoolExecutor) Start() error {
	err := PlContext.add(wp)
	if err != nil {
		return err
	}
	wp.logger.Infof("starting worker pool executor [ID]: %s", wp.ID)
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go func() {
		outer:
			for {
				select {
				case <-wp.terminate:
					break outer
				default:
					if wp.taskQ.IsEmpty() {
						break outer
					}
					task := wp.taskQ.Dequeue()
					task.routine()
				}
			}
			defer wp.wg.Done()
		}()
	}
	return nil
}

func (wp *WorkerPoolExecutor) Finish() {
	wp.wg.Wait()
	PlContext.remove(wp.ID)
	wp.logger.Infof("stopped worker pool executor [ID]: %s", wp.ID)
}

func (wp *WorkerPoolExecutor) Submit(t *Task) {
	wp.taskQ.Enqueue(t)
}

func (wp *WorkerPoolExecutor) TerminateGracefully() error {
	close(wp.terminate)
	return nil
}
