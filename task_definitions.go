package dogpool

import (
	"context"

	"github.com/chinmayrelkar/dogpool/dao"
)

type TaskName string
type TaskExecutionFunc func(appContext context.Context, task dao.Task) error

var taskDefinitions = map[TaskName]TaskExecutionFunc{}
