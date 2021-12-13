package commands

import (
	"context"

	"github.com/rafaelmartins/pcb-gcode-sender/internal/actions"
)

type gotoOriginCommand struct{}

func (*gotoOriginCommand) GetName() string {
	return "goto-origin"
}

func (*gotoOriginCommand) GetCompletions(args []string) []string {
	return nil
}

func (*gotoOriginCommand) Run(ctx context.Context, a *actions.Actions, args []string) error {
	return a.GotoOrigin(ctx)
}
