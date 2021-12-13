package commands

import (
	"context"

	"github.com/rafaelmartins/pcb-gcode-sender/internal/actions"
)

type xyZeroCommand struct{}

func (*xyZeroCommand) GetName() string {
	return "xy-zero"
}

func (*xyZeroCommand) GetCompletions(args []string) []string {
	return nil
}

func (z *xyZeroCommand) Run(ctx context.Context, a *actions.Actions, args []string) error {
	return a.SetZeroXY(ctx)
}
