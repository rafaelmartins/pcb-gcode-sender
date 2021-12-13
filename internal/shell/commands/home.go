package commands

import (
	"context"

	"github.com/rafaelmartins/pcb-gcode-sender/internal/actions"
)

type homeCommand struct{}

func (*homeCommand) GetName() string {
	return "home"
}

func (*homeCommand) GetCompletions(args []string) []string {
	return nil
}

func (*homeCommand) Run(ctx context.Context, a *actions.Actions, args []string) error {
	return a.Home(ctx)
}
