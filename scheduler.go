package dogpool

import (
	"github.com/chinmayrelkar/dogpool/dao"
	"gorm.io/gorm"
)

type Scheduler interface {
	ScheduleTask(taskName TaskName, args any) (*dao.Task, error)
}

func NewScheduler(db *gorm.DB) Scheduler {
	return &scheduler{
		dao: dao.NewTaskDao(db),
	}
}

type scheduler struct {
	dao dao.SchedulerTaskDao
}

func (s *scheduler) ScheduleTask(name TaskName, args any) (*dao.Task, error) {
	return s.dao.ScheduleTask(string(name), args)
}
