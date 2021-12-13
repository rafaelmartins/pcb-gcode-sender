package point

type ByDistanceXY struct {
	From   *Point
	Points []*Point
}

func (d ByDistanceXY) Len() int {
	return len(d.Points)
}

func (d ByDistanceXY) Swap(i int, j int) {
	d.Points[i], d.Points[j] = d.Points[j], d.Points[i]
}

func (d ByDistanceXY) Less(i int, j int) bool {
	return d.From.DistanceXY(d.Points[i]) < d.From.DistanceXY(d.Points[j])
}
