package commands

import (
	"context"
	"sort"
	"strings"

	"github.com/google/shlex"
	"github.com/rafaelmartins/pcb-gcode-sender/internal/actions"
)

type Command interface {
	GetName() string
	GetCompletions(args []string) []string
	Run(ctx context.Context, a *actions.Actions, args []string) error
}

var (
	commands = []Command{
		&autolevelCommand{},
		&autolevelLoadCommand{},
		&gotoOriginCommand{},
		&homeCommand{},
		&jogCommand{},
		&loadCommand{},
		&resetCommand{},
		&startCommand{},
		&unlockCommand{},
		&xyZeroCommand{},
		&zProbeCommand{},
	}
)

func Lookup(cmd string) Command {
	if len(cmd) == 0 {
		return nil
	}

	for _, command := range commands {
		if cmd == command.GetName() {
			return command
		}
	}

	return nil
}

func Completer(line string) []string {
	cmds := []string{}
	for _, cmd := range commands {
		cmds = append(cmds, cmd.GetName())
	}
	cmds = append(cmds, "quit") // FIXME
	sort.Strings(cmds)

	parts, err := shlex.Split(line)
	if err != nil {
		return nil
	}

	rv := []string{}
	for _, n := range cmds {
		if len(parts) > 0 && n == parts[0] {
			rv = []string{}
			for _, c := range commands {
				if c.GetName() == n {
					for _, cp := range c.GetCompletions(parts[1:]) {
						rv = append(rv, n+" "+cp)
					}
					break
				}
			}
			return rv
		}

		if strings.HasPrefix(n, strings.ToLower(line)) {
			rv = append(rv, n)
		}
	}
	return rv
}
