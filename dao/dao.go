package dao

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskDao interface {
	SchedulerTaskDao
	WorkerTaskDao
}

var taskDaoInstance TaskDao

func NewTaskDao(db *gorm.DB) TaskDao {
	if taskDaoInstance != nil {
		return taskDaoInstance
	}
	db.AutoMigrate(&Task{})
	taskDaoInstance = &taskDaoImpl{db: db}
	return taskDaoInstance
}

type SchedulerTaskDao interface {
	ScheduleTask(name string, args interface{}) (*Task, error)
}

type WorkerTaskDao interface {
	GetTaskToBeRun(db *gorm.DB) (*Task, error)
	MarkTaskAsRunning(db *gorm.DB, id string) (*Task, error)
	MarkTaskAsSucceeded(db *gorm.DB, id string) (*Task, error)
	MarkTaskAsFailed(db *gorm.DB, id string, err error) (*Task, error)
}

type taskDaoImpl struct {
	db *gorm.DB
}

func (t *taskDaoImpl) ScheduleTask(name string, args interface{}) (*Task, error) {
	task := Task{
		ID:          generateID(),
		Name:        name,
		ScheduledAt: time.Now().Format(time.RFC3339),
		State:       TaskStateScheduled,
	}
	task.WriteArgs(args)
	if err := t.db.Create(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (t *taskDaoImpl) GetTaskToBeRun(db *gorm.DB) (*Task, error) {
	var task Task
	if err := db.Where("state = ?", TaskStateScheduled).
		Order("scheduled_at asc").
		First(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (t *taskDaoImpl) MarkTaskAsRunning(db *gorm.DB, id string) (*Task, error) {
	var task Task
	if err := db.First(&task, "id = ? and state = ?", id, TaskStateScheduled).Error; err != nil {
		return nil, err
	}
	task.State = TaskStateRunning
	task.StartedAt = time.Now().Format(time.RFC3339)
	if err := db.Save(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (t *taskDaoImpl) MarkTaskAsSucceeded(db *gorm.DB, id string) (*Task, error) {
	var task Task
	if err := db.First(&task, "id = ? and state = ?", id, TaskStateRunning).Error; err != nil {
		return nil, err
	}
	task.State = TaskStateSucceeded
	task.CompletedAt = time.Now().Format(time.RFC3339)
	if err := db.Save(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func (t *taskDaoImpl) MarkTaskAsFailed(db *gorm.DB, id string, err error) (*Task, error) {
	var task Task
	if err := db.First(&task, "id = ? and state = ?", id, TaskStateRunning).Error; err != nil {
		return nil, err
	}
	task.State = TaskStateFailed
	task.CompletedAt = time.Now().Format(time.RFC3339)
	task.Error = err.Error()
	if err := db.Save(&task).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

func generateID() string {
	return uuid.New().String()
}
