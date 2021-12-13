package commands

import (
	"context"

	"github.com/rafaelmartins/pcb-gcode-sender/internal/actions"
)

type autolevelLoadCommand struct{}

func (*autolevelLoadCommand) GetName() string {
	return "autolevel-load"
}

func (*autolevelLoadCommand) GetCompletions(args []string) []string {
	return nil
}

func (*autolevelLoadCommand) Run(ctx context.Context, a *actions.Actions, args []string) error {
	return a.AutoLevelLoad(ctx)
}
