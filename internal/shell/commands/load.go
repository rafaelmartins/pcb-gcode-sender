package commands

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/rafaelmartins/pcb-gcode-sender/internal/actions"
)

type loadCommand struct{}

func (*loadCommand) GetName() string {
	return "load"
}

func (*loadCommand) GetCompletions(args []string) []string {
	if len(args) > 1 {
		return nil
	}

	fname := ""
	if len(args) == 1 {
		fname = args[0]
	}

	if st, err := os.Stat(fname); err == nil && !st.IsDir() {
		return []string{fname}
	}

	fs, err := ioutil.ReadDir(".")
	if err != nil {
		return nil
	}

	rv := []string{}
	for _, f := range fs {
		if !f.IsDir() && strings.HasPrefix(f.Name(), strings.ToLower(fname)) {
			rv = append(rv, f.Name())
		}
	}
	return rv
}

func (*loadCommand) Run(ctx context.Context, a *actions.Actions, args []string) error {
	if len(args) == 0 {
		return errors.New("load: g-code file not defined")
	}

	return a.LoadGCode(ctx, args[0])
}
