package autolevel

import (
	"github.com/rafaelmartins/pcb-gcode-sender/internal/gcode"
	"github.com/rafaelmartins/pcb-gcode-sender/internal/point"
)

func AutoLevel(gc gcode.Job, probe []*point.Point) (gcode.Job, error) {
	rv := []gcode.Line{}

	var (
		state  = &point.Point{}
		foundX bool
		foundY bool
		foundZ bool
	)

	for _, l := range gc {
		if !l.HasPosition() {
			rv = append(rv, l)
			continue
		}

		if x := l.Get('X'); x != nil {
			foundX = true
			state.X = x.Value
		}

		if y := l.Get('Y'); y != nil {
			foundY = true
			state.Y = y.Value
		}

		if z := l.Get('Z'); z != nil {
			foundZ = true
			state.Z = z.Value
		}

		if foundX && foundY && foundZ {
			np, err := state.InterpolateZ(probe)
			if err != nil {
				return nil, err
			}
			l.SetPosition(np)
		}

		rv = append(rv, l)
	}

	return rv, nil
}
