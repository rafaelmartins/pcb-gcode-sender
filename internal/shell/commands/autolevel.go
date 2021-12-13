package commands

import (
	"context"

	"github.com/rafaelmartins/pcb-gcode-sender/internal/actions"
)

type autolevelCommand struct{}

func (*autolevelCommand) GetName() string {
	return "autolevel"
}

func (*autolevelCommand) GetCompletions(args []string) []string {
	return nil
}

func (*autolevelCommand) Run(ctx context.Context, a *actions.Actions, args []string) error {
	return a.AutoLevel(ctx)
}
