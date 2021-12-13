package commands

import (
	"context"

	"github.com/rafaelmartins/pcb-gcode-sender/internal/actions"
)

type unlockCommand struct{}

func (*unlockCommand) GetName() string {
	return "unlock"
}

func (*unlockCommand) GetCompletions(args []string) []string {
	return nil
}

func (*unlockCommand) Run(ctx context.Context, a *actions.Actions, args []string) error {
	return a.Unlock(ctx)
}
