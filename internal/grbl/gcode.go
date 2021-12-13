package grbl

import (
	"fmt"
	"strings"

	"github.com/rafaelmartins/pcb-gcode-sender/internal/gcode"
)

type GCodeStates struct {
	Distance string
	Units    string
}

func NewGCodeStates(line string) (*GCodeStates, error) {
	rv := &GCodeStates{}
	for _, part := range strings.Split(line, " ") {
		rv.set(part)
	}

	if rv.Distance == "" {
		return nil, fmt.Errorf("grbl: failed to detect gcode distance state: %s", line)
	}

	if rv.Units == "" {
		return nil, fmt.Errorf("grbl: failed to detect gcode units state: %s", line)
	}

	return rv, nil
}

func (g *GCodeStates) set(f string) {
	if f == "G90" || f == "G91" {
		g.Distance = f
	}
	if f == "G20" || f == "G21" {
		g.Units = f
	}
}

func (g *GCodeStates) ProcessLine(l gcode.Line) {
	for _, field := range l {
		g.set(field.String())
	}
}
