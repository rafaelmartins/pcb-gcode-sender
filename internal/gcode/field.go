package gcode

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

type Field struct {
	Letter rune
	Value  float64
}

func NewField(field string) (*Field, error) {
	field = strings.TrimSpace(field)
	if field == "" || !unicode.IsLetter(rune(field[0])) {
		return nil, fmt.Errorf("gcode: bad field: %s", field)
	}

	v, err := strconv.ParseFloat(field[1:], 64)
	if err != nil {
		return nil, err
	}

	return &Field{
		Letter: rune(field[0]),
		Value:  v,
	}, nil
}

func (f *Field) String() string {
	if f.Value == math.Trunc(f.Value) {
		return fmt.Sprintf("%c%.0f", f.Letter, f.Value)
	}
	return fmt.Sprintf("%c%.4f", f.Letter, f.Value)
}
