package commands

import (
	"context"

	"github.com/rafaelmartins/pcb-gcode-sender/internal/actions"
)

type startCommand struct{}

func (*startCommand) GetName() string {
	return "start"
}

func (*startCommand) GetCompletions(args []string) []string {
	return nil
}

func (*startCommand) Run(ctx context.Context, a *actions.Actions, args []string) error {
	return a.Start(ctx)
}
