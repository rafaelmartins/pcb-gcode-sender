package commands

import (
	"context"
	"fmt"

	"github.com/eiannone/keyboard"
	"github.com/rafaelmartins/pcb-gcode-sender/internal/actions"
)

type jogCommand struct{}

func (*jogCommand) GetName() string {
	return "jog"
}

func (*jogCommand) GetCompletions(args []string) []string {
	return nil
}

func (*jogCommand) Run(ctx context.Context, a *actions.Actions, args []string) error {
	if err := keyboard.Open(); err != nil {
		return err
	}
	defer keyboard.Close()

	step := 1.

	fmt.Println("Press 'q' to quit jogging. F1/F2 to change step.")
	fmt.Printf("step: %.3f\n", step)

	for {
		ch, key, err := keyboard.GetKey()
		if err != nil {
			return err
		}

		if ch == 'q' {
			return nil
		}

		if err := func() error {
			switch key {
			case keyboard.KeyArrowLeft:
				return a.Jog(ctx, -step, 0, 0)
			case keyboard.KeyArrowRight:
				return a.Jog(ctx, step, 0, 0)
			case keyboard.KeyArrowDown:
				return a.Jog(ctx, 0, -step, 0)
			case keyboard.KeyArrowUp:
				return a.Jog(ctx, 0, step, 0)
			case keyboard.KeyPgdn:
				return a.Jog(ctx, 0, 0, -step)
			case keyboard.KeyPgup:
				return a.Jog(ctx, 0, 0, step)
			case keyboard.KeyF1:
				if step <= 10. {
					step /= 10
					fmt.Printf("step: %.3f\n", step)
				}
			case keyboard.KeyF2:
				if step <= 1. {
					step *= 10
					fmt.Printf("step: %.3f\n", step)
				}
			}

			return nil
		}(); err != nil {
			return err
		}
	}
}
