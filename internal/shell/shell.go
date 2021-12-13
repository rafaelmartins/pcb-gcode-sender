package shell

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"

	"github.com/google/shlex"
	"github.com/peterh/liner"
	"github.com/rafaelmartins/pcb-gcode-sender/internal/actions"
	"github.com/rafaelmartins/pcb-gcode-sender/internal/point"
	"github.com/rafaelmartins/pcb-gcode-sender/internal/shell/commands"
	"golang.org/x/sys/unix"
)

func formatAxis(v string, p *point.Point) string {
	if p == nil {
		return fmt.Sprintf(" | %s:undefined", v)
	}

	return fmt.Sprintf(" | %s:%s", v, p)
}

func formatFile(file string) string {
	if file != "" {
		return " | G:" + file
	}
	return " | G:none"
}

func Run(a *actions.Actions) error {
	if a.Grbl == nil {
		return errors.New("shell: grbl undefined")
	}

	line := liner.NewLiner()
	defer line.Close()

	line.SetCompleter(commands.Completer)

	sig := make(chan os.Signal)
	signal.Notify(sig, unix.SIGINT, unix.SIGKILL, unix.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			<-sig
			if _, err := line.Prompt("Cancel? (ctrl-d cancel, anything else to continue): "); err == io.EOF {
				fmt.Println()
				break
			}
		}

		log.Print("cancelling")
		cancel()
	}()

	first := true

	for {
		if err := a.Grbl.SendCommands(ctx, "?"); err != nil {
			if first {
				return fmt.Errorf("shell: %w", err)
			}

			log.Printf("error: shell: %s", err)
		}

		l, err := line.Prompt("pcb-gcode-sender | " + a.Grbl.StateName + formatAxis("M", a.Grbl.MPos) + formatAxis("W", a.Grbl.WPos) + formatFile(a.CurrentJobFile) + "> ")
		if err != nil {
			if err == io.EOF {
				fmt.Println()
				return nil
			}

			if first {
				return fmt.Errorf("shell: %w", err)
			}

			log.Printf("error: shell: %s", err)
			continue
		}

		first = false

		parts, err := shlex.Split(l)
		if err != nil {
			log.Printf("error: shell: %s", err)
			continue
		}

		if len(parts) == 0 {
			continue
		}

		if parts[0] == "quit" {
			line.AppendHistory(l)
			break
		}

		c := commands.Lookup(parts[0])
		if c == nil {
			log.Printf("error: shell: command not found: %s", parts[0])
			continue
		}

		line.AppendHistory(l)

		if err := c.Run(ctx, a, parts[1:]); err != nil {
			log.Printf("error: shell: %s", err)
		}

		select {
		case <-ctx.Done():
			return nil
		default:
		}
	}

	a.Grbl.SendCommands(context.Background(), "G04 P0.001\nM5")
	return nil
}
