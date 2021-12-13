package commands

import (
	"context"

	"github.com/rafaelmartins/pcb-gcode-sender/internal/actions"
)

type zProbeCommand struct{}

func (*zProbeCommand) GetName() string {
	return "z-probe"
}

func (*zProbeCommand) GetCompletions(args []string) []string {
	return nil
}

func (z *zProbeCommand) Run(ctx context.Context, a *actions.Actions, args []string) error {
	return a.ProbeZ(ctx)
}
