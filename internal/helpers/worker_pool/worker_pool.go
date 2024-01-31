package workerpool

import (
	"github.com/pramod-janardhana/dir-watcher/internal/zlog"
)

type Task struct {
	id      int
	details any
}

func NewTask(id int, details any) Task {
	return Task{id, details}
}

func (t *Task) Id() int { return t.id }

func (t *Task) Details() any { return t.details }

type TaskFeedback struct {
	TaskId  int
	Success bool
	Data    any
	Err     error
}

type WorkerPool struct {
	numWorkers          int
	taskChannel         chan Task
	taskFeedbackChannel chan TaskFeedback
	needNewWorker       chan struct{}
	workerQuitAck       chan struct{}
	handler             func(task Task) (any, error)
}

func (w *WorkerPool) TaskChan() chan Task { return w.taskChannel }

func (w *WorkerPool) TaskFeedbackChan() chan TaskFeedback { return w.taskFeedbackChannel }

func (w *WorkerPool) workerQuitAckChan() chan struct{} { return w.workerQuitAck }

// Start starts the worker and makes sure that the worker are always active until the flushed.
func (w *WorkerPool) Start() {
	i := 1
	for ; i <= w.numWorkers; i++ {
		go worker(i, w)
		zlog.Debugf("started worker %d", i)
	}

	for {
		_, ok := <-w.needNewWorker
		if !ok {
			return
		}

		i++
		go worker(i, w)
		zlog.Infof("started new worker %d", i)
	}
}

// Flush closes all the channels and waits for the workers to finish.
// The same worker pool can be used after flush.
func (w *WorkerPool) Flush() {

	// close only the task channels and keep the other channels alive until the workers finish
	close(w.taskChannel)

	// waiting for workers to terminate
	for i := 0; i < w.numWorkers; i++ {
		<-w.workerQuitAck
	}

	close(w.workerQuitAck)
	close(w.needNewWorker)
	close(w.taskFeedbackChannel)
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(numWorkers int, handler func(task Task) (any, error)) *WorkerPool {
	return &WorkerPool{
		numWorkers:          numWorkers,
		taskChannel:         make(chan Task, numWorkers),
		taskFeedbackChannel: make(chan TaskFeedback, numWorkers),
		needNewWorker:       make(chan struct{}),
		workerQuitAck:       make(chan struct{}, numWorkers),
		handler:             handler,
	}
}

func worker(workerId int, w *WorkerPool) {
	defer func() {
		if r := recover(); r != nil {
			zlog.Errorf("work-%d failed with panic: %v", workerId, r)
			w.needNewWorker <- struct{}{}
		} else {
			w.workerQuitAckChan() <- struct{}{}
			zlog.Debugf("work-%d completed and sent ack signal", workerId)
		}
	}()

	for task := range w.taskChannel {
		data, err := w.handler(task)
		if err != nil {
			w.TaskFeedbackChan() <- TaskFeedback{
				TaskId:  task.id,
				Success: false,
				Err:     err,
			}
		} else {
			w.TaskFeedbackChan() <- TaskFeedback{
				TaskId:  task.id,
				Success: true,
				Data:    data,
			}
		}

		zlog.Debugf("worker %d completed task %d", workerId, task.Id())
	}
}
