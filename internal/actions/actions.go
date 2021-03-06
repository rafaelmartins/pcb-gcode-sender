package actions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/rafaelmartins/pcb-gcode-sender/internal/autolevel"
	"github.com/rafaelmartins/pcb-gcode-sender/internal/gcode"
	"github.com/rafaelmartins/pcb-gcode-sender/internal/grbl"
	"github.com/rafaelmartins/pcb-gcode-sender/internal/interp2d"
	"github.com/rafaelmartins/pcb-gcode-sender/internal/point"
)

var (
	ErrGrblNotSet = errors.New("actions: grbl not set")
)

type Actions struct {
	Grbl           *grbl.Grbl
	CurrentJob     gcode.Job
	CurrentJobFile string
	Probe          [][]*point.Point
	ProbeSpline    *interp2d.Spline
}

func (a *Actions) Home(ctx context.Context) error {
	if a == nil || a.Grbl == nil {
		return ErrGrblNotSet
	}

	return a.Grbl.SendCommands(ctx, "$H")
}

func (a *Actions) Reset(ctx context.Context) error {
	if a == nil || a.Grbl == nil {
		return ErrGrblNotSet
	}

	return a.Grbl.SendCommands(ctx, string([]byte{0x18}))
}

func (a *Actions) Jog(ctx context.Context, x float64, y float64, z float64) error {
	if a == nil || a.Grbl == nil {
		return ErrGrblNotSet
	}

	cmd := fmt.Sprintf("$J=G91 X%.3f Y%.3f Z%.3f F10000", x, y, z)
	return a.Grbl.SendCommands(ctx, cmd)
}

func (a *Actions) Unlock(ctx context.Context) error {
	if a == nil || a.Grbl == nil {
		return ErrGrblNotSet
	}

	return a.Grbl.SendCommands(ctx, "$X")
}

func (a *Actions) GotoOrigin(ctx context.Context) error {
	if a == nil || a.Grbl == nil {
		return ErrGrblNotSet
	}

	return a.Grbl.SendGCodeInline(ctx, `
G90
G01 Z2 F10000
G01 X0 Y0
G04 P0.001
F100
`)
}

func (a *Actions) SetZeroXY(ctx context.Context) error {
	if a == nil || a.Grbl == nil {
		return ErrGrblNotSet
	}

	return a.Grbl.SendGCodeInline(ctx, "G10 L20 P1 X0 Y0")
}

func (a *Actions) ProbeZ(ctx context.Context) error {
	if a == nil || a.Grbl == nil {
		return ErrGrblNotSet
	}

	if err := a.Grbl.SendGCodeInline(ctx, `
G91
G38.2 Z-100 F50
G01 Z1 F100
G38.2 Z-2 F10
G90
G04 P0.001`); err != nil {
		return err
	}

	if err := a.Grbl.SendCommands(ctx, "?"); err != nil {
		return err
	}

	if a.Grbl.LastProbe == nil || a.Grbl.MPos == nil {
		return errors.New("actions: probe-z: probe failed")
	}

	return a.Grbl.SendGCodeInline(ctx, fmt.Sprintf(`
G10 L20 P1 Z%.3f
G01 Z2 F100
G04 P0.001`, a.Grbl.MPos.Z-a.Grbl.LastProbe.Z))
}

func (a *Actions) LoadGCode(ctx context.Context, file string) error {
	if a == nil || a.Grbl == nil {
		return ErrGrblNotSet
	}

	j, err := gcode.NewJobFromFile(file)
	if err != nil {
		return err
	}
	a.CurrentJob = j
	a.CurrentJobFile = file

	return nil
}

func (a *Actions) autoLevelProbe(ctx context.Context, x float64, y float64) error {
	if err := a.Grbl.SendGCodeInline(ctx, fmt.Sprintf(`
G90
G01 Z2 F10000
G01 X%.3f Y%.3f
G91
G38.2 Z-100 F50
G01 Z1 F100
G38.2 Z-2 F10
G90
G01 Z2 F100
G04 P0.001`, x, y)); err != nil {
		return err
	}

	if a.Grbl.LastProbe == nil || a.Grbl.MPos == nil {
		return errors.New("actions: autolevel: probe failed")
	}

	return nil
}

func (a *Actions) autoLevelLoadProbe(pts [][]*point.Point, wco *point.Point) error {
	a.Probe = make([][]*point.Point, len(pts))
	for j, lp := range pts {
		a.Probe[j] = make([]*point.Point, len(lp))
		for i, p := range lp {
			a.Probe[j][i] = p.Sub(wco)
		}
	}

	sp, err := interp2d.NewSpline(a.Probe)
	if err != nil {
		return err
	}
	a.ProbeSpline.Close()
	a.ProbeSpline = sp

	return nil
}

func (a *Actions) AutoLevel(ctx context.Context) error {
	if a == nil || a.Grbl == nil {
		return ErrGrblNotSet
	}

	if a.CurrentJob == nil {
		return errors.New("actions: autolevel: no g-code loaded")
	}

	minx, miny, maxx, maxy, err := a.CurrentJob.GetBoundingBox()
	if err != nil {
		return err
	}
	minx -= 0.2
	miny -= 0.2
	maxx += 0.2
	maxy += 0.2

	distx := maxx - minx
	numx := int(distx / 10)

	disty := maxy - miny
	numy := int(disty / 10)

	xgap := distx / (float64(numx) - 1)
	ygap := disty / (float64(numy) - 1)

	fmt.Println(numx, numy, maxx, maxy)

	pts := make([][]*point.Point, numy)
	y := miny

	for j := 0; j < numy; j++ {
		pts[j] = make([]*point.Point, numx)

		if j%2 == 0 {
			x := minx
			for i := 0; i < numx; i++ {
				select {
				case <-ctx.Done():
					return nil
				default:
				}

				if err := a.autoLevelProbe(ctx, x, y); err != nil {
					return err
				}
				pts[j][i] = a.Grbl.LastProbe

				x += xgap
			}
		} else {
			x := maxx
			for i := numx - 1; i >= 0; i-- {
				select {
				case <-ctx.Done():
					return nil
				default:
				}

				if err := a.autoLevelProbe(ctx, x, y); err != nil {
					return err
				}
				pts[j][i] = a.Grbl.LastProbe

				x -= xgap
			}
		}

		y += ygap
	}

	fp, err := os.OpenFile(a.CurrentJobFile+".json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer fp.Close()

	if err := json.NewEncoder(fp).Encode(map[string]interface{}{
		"wco":    a.Grbl.WCO,
		"points": pts,
	}); err != nil {
		return err
	}

	return a.autoLevelLoadProbe(pts, a.Grbl.WCO)
}

func (a *Actions) AutoLevelLoad(ctx context.Context) error {
	if a == nil || a.Grbl == nil {
		return ErrGrblNotSet
	}

	if a.CurrentJob == nil || a.CurrentJobFile == "" {
		return errors.New("actions: autolevel-load: no g-code loaded")
	}

	fp, err := os.Open(a.CurrentJobFile + ".json")
	if err != nil {
		return err
	}
	defer fp.Close()

	data := struct {
		WCO    *point.Point     `json:"wco"`
		Points [][]*point.Point `json:"points"`
	}{}

	if err := json.NewDecoder(fp).Decode(&data); err != nil {
		return err
	}

	if a.Grbl.WCO != nil && !a.Grbl.WCO.Equals(data.WCO) {
		return errors.New("actions: autolevel-load: stored WCO differs from current WCO")
	}

	return a.autoLevelLoadProbe(data.Points, data.WCO)
}

func (a *Actions) Start(ctx context.Context) error {
	if a == nil || a.Grbl == nil {
		return ErrGrblNotSet
	}

	if a.CurrentJob == nil || a.CurrentJobFile == "" {
		return errors.New("actions: start: no g-code loaded")
	}

	if len(a.Probe) > 0 {
		g, err := autolevel.AutoLevel(a.CurrentJob, a.ProbeSpline, *a.Grbl.GCodeState)
		if err != nil {
			return err
		}

		log.Print("autolevel enabled")
		return a.Grbl.SendJob(ctx, g)
	}

	return a.Grbl.SendJob(ctx, a.CurrentJob)
}
