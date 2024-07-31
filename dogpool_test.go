package dogpool_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/chinmayrelkar/dogpool"
	"github.com/chinmayrelkar/dogpool/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, error) {
	// Set up a test database connection
	// You might want to use an in-memory SQLite database for testing
	return gorm.Open(mysql.New(mysql.Config{
		DriverName: "mysql",
		DSN:        os.Getenv("TEST_DB_URL"),
	}), &gorm.Config{})
}

func TestTaskExecution(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to set up test database: %v", err)
	}

	worker := dogpool.NewWorker(context.Background(), db)
	defer worker.Exit()

	scheduler := dogpool.NewScheduler(db)

	t.Run("TestSuccessfulTaskExecution", func(t *testing.T) {
		worker.Register("TestTaskExecution", func(ctx context.Context, task dao.Task) error {
			return nil
		})

		task, err := scheduler.ScheduleTask(
			"TestTaskExecution",
			map[string]string{},
		)
		if err != nil {
			t.Fatalf("Failed to schedule task: %v", err)
		}
		if task.State != dao.TaskStateScheduled {
			t.Fatalf("Task is not marked as scheduled")
		}
		go worker.Run()
		time.Sleep(10 * time.Second)

		// check if task is marked as completed
		db.Raw("SELECT * FROM tasks WHERE id =?", task.ID).Scan(&task)
		if task.State != dao.TaskStateSucceeded {
			t.Fatalf("Task is not marked as completed")
		}
	})

	t.Run("TestFailedTaskExecution", func(t *testing.T) {
		worker.Register("TestFailedTaskExecution", func(ctx context.Context, task dao.Task) error {
			return errors.New("Failed task")
		})
		go worker.Run()

		task, err := scheduler.ScheduleTask(
			"TestFailedTaskExecution",
			map[string]string{},
		)

		if err != nil {
			t.Fatalf("Failed to schedule task: %v", err)
		}
		time.Sleep(10 * time.Second)
		// check if task is marked as failed
		db.Raw("SELECT * FROM tasks WHERE id =?", task.ID).Scan(&task)
		if task.State != dao.TaskStateFailed {
			t.Fatalf("Task is not marked as failed")
		}
		if task.Error != "Failed task" {
			t.Fatalf("Task error is not set correctly")
		}
	})
}
