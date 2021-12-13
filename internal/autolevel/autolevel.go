package autolevel

import (
	"errors"

	"github.com/rafaelmartins/pcb-gcode-sender/internal/gcode"
	"github.com/rafaelmartins/pcb-gcode-sender/internal/grbl"
	"github.com/rafaelmartins/pcb-gcode-sender/internal/point"
)

func AutoLevel(gc gcode.Job, probe []*point.Point, gcs grbl.GCodeStates) (gcode.Job, error) {
	rv := []gcode.Line{}

	var (
		state  = &point.Point{}
		foundX bool
		foundY bool
		foundZ bool
	)

	for _, l := range gc {
		if gcs.Distance == "G91" {
			return nil, errors.New("autolevel: can't autolevel incremental positions") // FIXME
		}

		if !l.HasPosition() {
			gcs.ProcessLine(l)
			rv = append(rv, l)
			continue
		}

		if x := l.Get('X'); x != nil {
			foundX = true
			state.X = gcs.ToMM(x.Value)
		}

		if y := l.Get('Y'); y != nil {
			foundY = true
			state.Y = gcs.ToMM(y.Value)
		}

		if z := l.Get('Z'); z != nil {
			foundZ = true
			state.Z = gcs.ToMM(z.Value)
		}

		if foundX && foundY && foundZ {
			np, err := state.InterpolateZ(probe)
			if err != nil {
				return nil, err
			}
			l.SetPosition(np)
		}

		gcs.ProcessLine(l)
		rv = append(rv, l)
	}

	return rv, nil
}
