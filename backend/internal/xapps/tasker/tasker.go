package tasker

import (
	"context"

	"github.com/hay-kot/hookfeed/backend/internal/core/tasks"
)

type Tasker struct{}

var _ tasks.Queue = &Tasker{}

func New() *Tasker {
	return &Tasker{}
}

func (t *Tasker) Start(ctx context.Context) error {
	return nil
}

func (t *Tasker) Enqueue(task tasks.Task) error {
	return nil
}
