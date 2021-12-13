package point

import (
	"errors"
	"fmt"
	"math"
	"sort"
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

func (p *Point) distance(pd *Point, withz bool) float64 {
	dx := pd.X - p.X
	dy := pd.Y - p.Y
	d := math.Pow(dx, 2) + math.Pow(dy, 2)
	if withz {
		dz := pd.X - p.X
		d += math.Pow(dz, 2)
	}
	return math.Sqrt(d)
}

func (p *Point) Distance(pd *Point) float64 {
	return p.distance(pd, true)
}

func (p *Point) DistanceXY(pd *Point) float64 {
	return p.distance(pd, false)
}

func areColinear(p1 *Point, p2 *Point, p3 *Point) bool {
	return math.Abs(p1.X*(p2.Y-p3.Y)+p2.X*(p3.Y-p1.Y)+p3.X*(p1.Y-p2.Y)) < 0.001
}

func (p *Point) ThreeClosest(points []*Point) []*Point {
	if len(points) < 3 {
		return nil
	}

	sort.Sort(ByDistanceXY{
		From:   p,
		Points: points,
	})

	rv := points[:2]

	for _, n := range points[2:] {
		if !areColinear(rv[0], rv[1], n) {
			rv = append(rv, n)
			break
		}
	}

	if len(rv) != 3 {
		return nil
	}

	return rv
}

func (p *Point) InterpolateZ(points []*Point) (*Point, error) {
	pts := p.ThreeClosest(append([]*Point{}, points...))
	if len(pts) != 3 {
		return nil, errors.New("point: failed to find 3 closest points")
	}

	s1 := pts[1].Sub(pts[0])
	s2 := pts[2].Sub(pts[0])

	normal := &Point{
		X: s1.Y*s2.Z - s1.Z*s2.Y,
		Y: -(s1.X*s2.Z - s1.Z*s2.X),
		Z: s1.X*s2.Y - s1.Y*s2.X,
	}

	z := 0.
	if normal.Z != 0 {
		z = pts[0].Z - (normal.X*(p.X-pts[0].X)+normal.Y*(p.Y-pts[0].Y))/normal.Z
	}

	return &Point{
		X: p.X,
		Y: p.Y,
		Z: p.Z + z,
	}, nil
}
