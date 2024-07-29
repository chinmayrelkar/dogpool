package dao

import "encoding/json"

type TaskState string

const (
	TaskStateScheduled TaskState = "TaskStateScheduled"
	TaskStateRunning   TaskState = "TaskStateRunning"
	TaskStateSucceeded TaskState = "TaskStateSucceeded"
	TaskStateFailed    TaskState = "TaskStateFailed"
)

type Task struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Args        string    `json:"args"`
	ScheduledAt string    `json:"scheduled_at"`
	State       TaskState `json:"state"`
	StartedAt   string    `json:"executed_at"`
	CompletedAt string    `json:"execution_completed_at"`
	Error       string    `json:"execution_result"`
}

func (t Task) TableName() string {
	return "dogpool_tasks"
}

func (t Task) ReadArgs(object any) error {
	return json.Unmarshal([]byte(t.Args), object)
}

func (t *Task) WriteArgs(object any) error {
	bytes, err := json.Marshal(object)
	if err != nil {
		return err
	}
	t.Args = string(bytes)
	return nil
}
