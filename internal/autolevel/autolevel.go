package autolevel

import (
	"errors"

	"github.com/rafaelmartins/pcb-gcode-sender/internal/gcode"
	"github.com/rafaelmartins/pcb-gcode-sender/internal/grbl"
	"github.com/rafaelmartins/pcb-gcode-sender/internal/interp2d"
	"github.com/rafaelmartins/pcb-gcode-sender/internal/point"
)

func AutoLevel(gc gcode.Job, spline *interp2d.Spline, gcs grbl.GCodeStates) (gcode.Job, error) {
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
			z := state.Z
			if spline != nil {
				dz, err := spline.At(state.X, state.Y)
				if err != nil {
					return nil, err
				}
				z += dz
			}
			l.SetPosition(&point.Point{
				X: state.X,
				Y: state.Y,
				Z: z,
			})
		}

		gcs.ProcessLine(l)
		rv = append(rv, l)
	}

	return rv, nil
}
