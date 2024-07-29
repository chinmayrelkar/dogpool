package dogpool

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chinmayrelkar/dogpool/dao"
	"gorm.io/gorm"
)

type WorkerError string

func (we WorkerError) Error() error {
	return errors.New(string(we))
}

func (we WorkerError) String() string {
	return string(we)
}

const (
	FailedToMarkTaskAsRunning   WorkerError = "failed to mark task as running"
	FailedToFetchTaskToRun      WorkerError = "failed to fetch task to run"
	FailedToExecuteTask         WorkerError = "failed to execute task"
	FailedToMarkTaskAsSucceeded WorkerError = "failed to mark task as succeeded"
	FailedToMarkTaskAsFailed    WorkerError = "failed to mark task as failed"
	ExecuteTaskFunctionNotFound WorkerError = "execute task function not found"
)

type Worker interface {
	Run() error
	Register(name string, executionFunc TaskExecutionFunc)
	Exit()
}

func NewWorker(taskRunnerContext context.Context, db *gorm.DB) Worker {
	return &worker{
		taskRunnerContext: taskRunnerContext,
		dao:               dao.NewTaskDao(db),
		db:                db,
		logger:            NewLogger(),
		exitChannel:       make(chan struct{}),
	}
}

type worker struct {
	taskRunnerContext context.Context
	dao               dao.WorkerTaskDao
	db                *gorm.DB
	logger            Logger
	exitChannel       chan struct{}
}

func (w *worker) Register(name string, executionFunc TaskExecutionFunc) {
	taskDefinitions[TaskName(name)] = executionFunc
}

func (w *worker) Run() error {
	for {
		select {
		case <-w.exitChannel:
			return nil
		default:
			w.runSingleTask()
		}
	}
}

func (w *worker) Exit() {
	w.exitChannel <- struct{}{}
}

func (w *worker) runSingleTask() {
	task, err := w.getSingleTask()
	if err != nil {
		w.logger.Error(FailedToFetchTaskToRun.String(), err)
		w.logger.Warn("Sleeping for 10 seconds")
		time.Sleep(10 * time.Second)
		return
	}
	if task == nil {
		w.logger.Warn("No Task found, sleeping for 10 seconds")
		time.Sleep(10 * time.Second)
		return
	}

	logIdentifier := fmt.Sprintf("[%s][%s]", task.Name, task.ID)

	executionFunc := taskDefinitions[TaskName(task.Name)]
	if executionFunc == nil {
		w.logger.Error(logIdentifier, ExecuteTaskFunctionNotFound)
		if _, err := w.dao.MarkTaskAsFailed(w.db, task.ID, err); err != nil {
			w.logger.Error(logIdentifier, FailedToMarkTaskAsFailed.String(), err)
			return
		}
	}
	err = executionFunc(w.taskRunnerContext, *task)
	if err != nil {
		w.logger.Error(logIdentifier, FailedToExecuteTask.String(), err)
		return
	}

	if _, err := w.dao.MarkTaskAsSucceeded(w.db, task.ID); err != nil {
		w.logger.Error(logIdentifier, FailedToMarkTaskAsSucceeded.String(), err)
		return
	}
}

func (w *worker) getSingleTask() (*dao.Task, error) {
	tx := w.db.Begin()
	defer tx.Commit()
	task, err := w.dao.GetTaskToBeRun(tx)
	if err != nil {
		return nil, fmt.Errorf(FailedToFetchTaskToRun.String()+"Error: %v", err)
	}
	if task == nil {
		return nil, nil
	}
	if _, err := w.dao.MarkTaskAsRunning(tx, task.ID); err != nil {
		return nil, fmt.Errorf(FailedToMarkTaskAsRunning.String()+". TaskId:%s. Error: %v", task.ID, err)
	}
	tx.Commit()
	return task, nil
}
