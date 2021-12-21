package point

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

func NewFromInches(x float64, y float64, z float64) *Point {
	return &Point{
		X: x / 25.4,
		Y: y / 25.4,
		Z: z / 25.4,
	}
}

func NewFromStringMM(str string) (*Point, error) {
	parts := strings.Split(str, ",")
	if len(parts) != 3 {
		return nil, fmt.Errorf("grbl: status: failed to convert to axis: %s", str)
	}

	rv := make([]float64, 3)

	var err error
	for i, part := range parts {
		rv[i], err = strconv.ParseFloat(strings.TrimSpace(part), 64)
		if err != nil {
			return nil, err
		}
	}

	return &Point{
		X: rv[0],
		Y: rv[1],
		Z: rv[2],
	}, nil
}

func (p *Point) ToInches() (float64, float64, float64) {
	return p.X * 25.4, p.Y * 25.4, p.Z * 25.4
}

func (p *Point) String() string {
	return fmt.Sprintf("X=%.3f,Y=%.3f,Z=%.3f", p.X, p.Y, p.Z)
}

func (p *Point) Copy() *Point {
	return &Point{
		X: p.X,
		Y: p.Y,
		Z: p.Z,
	}
}

func (p *Point) Add(pp *Point) *Point {
	return &Point{
		X: p.X + pp.X,
		Y: p.Y + pp.Y,
		Z: p.Z + pp.Z,
	}
}

func (p *Point) Sub(pp *Point) *Point {
	return &Point{
		X: p.X - pp.X,
		Y: p.Y - pp.Y,
		Z: p.Z - pp.Z,
	}
}

func (p *Point) Equals(pp *Point) bool {
	// for grbl purposes, 4 digits of precision are enough
	if math.Abs(p.X-pp.X) >= 0.001 {
		return false
	}
	if math.Abs(p.Y-pp.Y) >= 0.001 {
		return false
	}
	if math.Abs(p.Z-pp.Z) >= 0.001 {
		return false
	}

	return true
}
