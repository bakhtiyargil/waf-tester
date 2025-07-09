package utility

import (
	"sync"
	"time"
	"waf-tester/logger"
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
	logger     *logger.AppLogger
	start      time.Time
}

func NewWorkerPoolExecutor(id string, workers int, logger *logger.AppLogger) *WorkerPoolExecutor {
	return &WorkerPoolExecutor{
		id:         id,
		numWorkers: workers,
		stopChan:   make(chan struct{}),
		logger:     logger,
	}
}

func (wp *WorkerPoolExecutor) Start() {
	wp.logger.Infof("starting worker pool executor [ID]: %s", wp.id)
	wp.start = time.Now()
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go func() {
			for {
				if wp.taskQ.IsEmpty() {
					break
				}
				task := wp.taskQ.Dequeue()
				task.routine(task.staticParam, task.param)
			}
			defer wp.wg.Done()
		}()
	}
}

func (wp *WorkerPoolExecutor) Finish() {
	wp.wg.Wait()
	wp.logger.Infof("stopped worker pool executor [ID]: %s", wp.id)
}

func (wp *WorkerPoolExecutor) Submit(t *Task) {
	wp.taskQ.Enqueue(t)
}
