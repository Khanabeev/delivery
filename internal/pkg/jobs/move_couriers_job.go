package jobs

import (
	"context"
	"delivery/internal/core/application/usecasese/commands"
	"delivery/internal/pkg/errs"

	"github.com/labstack/gommon/log"
	"github.com/robfig/cron/v3"
)

var _ cron.Job = &MoveCouriersJob{}

type MoveCouriersJob struct {
	moveCouriersCommandHandler commands.MoveCourierHandler
}

func NewMoveCouriersJob(
	moveCouriersCommandHandler commands.MoveCourierHandler) (cron.Job, error) {
	if moveCouriersCommandHandler == nil {
		return nil, errs.NewValueIsRequiredError("moveCouriersCommandHandler")
	}

	return &MoveCouriersJob{
		moveCouriersCommandHandler: moveCouriersCommandHandler}, nil
}

func (j *MoveCouriersJob) Run() {
	ctx := context.Background()
	command, err := commands.NewMoveCouriersCommand()
	if err != nil {
		log.Error(err)
	}
	err = j.moveCouriersCommandHandler.Handler(ctx, command)
	if err != nil {
		log.Error(err)
	}
}
