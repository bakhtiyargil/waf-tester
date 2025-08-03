package utility

import (
	"sync"
)

type Task struct {
	routine func()
	next    *Task
	prev    *Task
}

func NewTask(routine func()) *Task {
	return &Task{
		routine: routine,
	}
}

type TaskQueue struct {
	size  int
	head  *Task
	tail  *Task
	mutex sync.Mutex
}

func (taskQ *TaskQueue) Enqueue(newTask *Task) {
	if taskQ.head == nil && taskQ.tail == nil {
		taskQ.tail = newTask
		taskQ.head = taskQ.tail
	} else if taskQ.tail != nil && taskQ.tail == taskQ.head {
		taskQ.tail = newTask
		taskQ.head.prev = taskQ.tail
		taskQ.tail.next = taskQ.head
	} else {
		taskQ.tail.prev = newTask
		newTask.next = taskQ.tail
		taskQ.tail = newTask
	}
	taskQ.size++
}

func (taskQ *TaskQueue) Dequeue() (eldest *Task) {
	taskQ.mutex.Lock()
	defer taskQ.mutex.Unlock()
	if taskQ.head == nil {
		return nil
	}

	eldest = taskQ.head
	taskQ.head = taskQ.head.prev

	if taskQ.head != nil {
		taskQ.head.next = nil
	} else {
		taskQ.tail = nil
	}

	eldest.prev = nil
	eldest.next = nil
	taskQ.size--
	return eldest
}

func (taskQ *TaskQueue) GetSize() int {
	return taskQ.size
}

func (taskQ *TaskQueue) IsEmpty() bool {
	return taskQ.size == 0
}
