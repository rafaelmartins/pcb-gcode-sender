package commands

import (
	"context"

	"github.com/rafaelmartins/pcb-gcode-sender/internal/actions"
)

type resetCommand struct{}

func (*resetCommand) GetName() string {
	return "reset"
}

func (*resetCommand) GetCompletions(args []string) []string {
	return nil
}

func (*resetCommand) Run(ctx context.Context, a *actions.Actions, args []string) error {
	return a.Reset(ctx)
}
