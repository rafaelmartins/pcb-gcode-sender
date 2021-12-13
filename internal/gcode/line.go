package gcode

import (
	"bufio"
	"strings"

	"github.com/rafaelmartins/pcb-gcode-sender/internal/point"
)

type Line []*Field

func NewLine(line string) (Line, error) {
	scanner := bufio.NewScanner(strings.NewReader(line))
	scanner.Split(bufio.ScanWords)

	rv := []*Field{}
	for scanner.Scan() {
		f, err := NewField(scanner.Text())
		if err != nil {
			return nil, err
		}
		rv = append(rv, f)
	}

	return Line(rv), nil
}

func (l Line) String() string {
	rv := ""
	for _, f := range l {
		rv += f.String()
	}
	return rv
}

func (l Line) Get(letter rune) *Field {
	for _, f := range l {
		if f.Letter == letter {
			return f
		}
	}

	return nil
}

func (l Line) IsMotion() bool {
	g := l.Get('G')
	return g != nil && (g.Value == 0. || g.Value == 1.) // there's more supported by flatcam?
}

func (l Line) HasPosition() bool {
	return l.Get('X') != nil || l.Get('Y') != nil || l.Get('Z') != nil
}

func (l *Line) SetPosition(p *point.Point) {
	if p == nil {
		return
	}

	if x := l.Get('X'); x != nil {
		x.Value = p.X
	} else {
		*l = append(*l, &Field{
			Letter: 'X',
			Value:  p.X,
		})
	}

	if y := l.Get('Y'); y != nil {
		y.Value = p.Y
	} else {
		*l = append(*l, &Field{
			Letter: 'Y',
			Value:  p.Y,
		})
	}

	if z := l.Get('Z'); z != nil {
		z.Value = p.Z
	} else {
		*l = append(*l, &Field{
			Letter: 'Z',
			Value:  p.Z,
		})
	}
}
