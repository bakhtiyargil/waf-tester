package utility

import (
	"sync"
	"waf-tester/logger"
)

type Executor interface {
	Submit(task *Task)
	Start() (key string, err error)
	Finish()
	Terminate() error
}

type WorkerPoolExecutor struct {
	id         string
	numWorkers int
	taskQ      TaskQueue
	wg         sync.WaitGroup
	logger     logger.Logger
	terminate  chan struct{}
}

func NewWorkerPoolExecutor(id string, workers int, logger logger.Logger) Executor {
	return &WorkerPoolExecutor{
		id:         id,
		numWorkers: workers,
		logger:     logger,
		terminate:  make(chan struct{}),
	}
}

func (wp *WorkerPoolExecutor) Start() (poolKey string, err error) {
	poolKey, err = PlContext.add(wp)
	if err != nil {
		return "", err
	}
	wp.id = poolKey
	wp.logger.Infof("starting worker pool executor [ID]: %s", poolKey)
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
						break
					}
					task := wp.taskQ.Dequeue()
					task.routine(task.staticParam, task.param)

				}
			}
			defer wp.wg.Done()
		}()
	}
	return wp.id, nil
}

func (wp *WorkerPoolExecutor) Finish() {
	wp.wg.Wait()
	PlContext.remove(wp.id)
	wp.logger.Infof("stopped worker pool executor [ID]: %s", wp.id)
}

func (wp *WorkerPoolExecutor) Submit(t *Task) {
	wp.taskQ.Enqueue(t)
}

func (wp *WorkerPoolExecutor) Terminate() error {
	close(wp.terminate)
	return nil
}
