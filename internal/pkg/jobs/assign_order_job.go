package jobs

import (
	"context"
	"delivery/internal/core/application/usecasese/commands"
	"delivery/internal/pkg/errs"

	"github.com/labstack/gommon/log"
	"github.com/robfig/cron/v3"
)

var _ cron.Job = &AssignOrderJob{}

type AssignOrderJob struct {
	assignOrderHandler commands.AssignOrderToCourierHandler
}

func NewAssignOrderJob(
	assignOrderHandler commands.AssignOrderToCourierHandler) (cron.Job, error) {
	if assignOrderHandler == nil {
		return nil, errs.NewValueIsRequiredError("assignOrderHandler")
	}

	return &AssignOrderJob{
		assignOrderHandler: assignOrderHandler,
	}, nil
}

func (j *AssignOrderJob) Run() {
	ctx := context.Background()
	command, err := commands.NewAssignOrderToCourierCommand()
	if err != nil {
		log.Error(err)
	}
	err = j.assignOrderHandler.Handle(ctx, command)
	if err != nil {
		//log.Error(err)
	}
}
