package figures

import (
	"math"
	"syscall/js"
)

type Point struct {
	X float64
	Y float64
}

func NewPoint(x float64, y float64) *Point {
	return &Point{X: x, Y: y}
}

func (point *Point) Intersects(elem Figure) bool {
	switch other := (elem).(type) {
	case *Point:
		return point.X == other.X && point.Y == other.Y
	case *Rectangle:
		return other.Start.X < point.X &&
			other.End.X > point.X &&
			other.Start.Y < point.Y &&
			other.End.Y > point.Y
	case *Circle:
		// Get the distance between the center and the point
		dx := math.Abs(point.X - other.Center.X)
		dy := math.Abs(point.Y - other.Center.Y)
		// The distance is larger than a square wrapping the circle
		if dx > other.Radius || dy > other.Radius {
			return false
		}
		// The distance is less than a square contained by the circle or pitagora's
		if dx+dy <= other.Radius || math.Pow(dx, 2)+math.Pow(dy, 2) <= math.Pow(other.Radius, 2) {
			return true
		}
		return false
	case *Line:
		return other.Slope() == NewLine(other.Start, point).Slope()
	case *Segment:
		// The point is between the 2 x's and the slope is the same (same immaginary line)
		return point.X <= math.Max(other.Start.X, other.End.X) && point.X >= math.Min(other.Start.X, other.End.X) &&
			point.Y <= math.Max(other.Start.Y, other.End.Y) && point.Y >= math.Min(other.Start.Y, other.End.Y) &&
			other.Slope() == NewLine(other.Start, point).Slope()

	}
	// The point cannot check if it is contained, as it doesn't know the other figure
	return elem.Intersects(point)
}

func (point *Point) DistanceToPoint(other *Point) float64 {
	return point.DistanceTo(other)
}
func (point *Point) DistanceToLine(line *Line) float64 {
	return line.NearestPointTo(point).DistanceToPoint(point)
}

func (point *Point) LineTo(other *Point) *Line {
	return NewLine(point, other)
}
func (point *Point) SegmentTo(other *Point) *Segment {
	return NewSegment(point, other)
}

func (point *Point) GetCenter() *Point {
	return point
}
func (point *Point) MoveTo(where *Point) {
	point.X, point.Y = where.Values()
}

func (point *Point) Draw(ctx js.Value) {
	ctx.Call("beginPath")
	ctx.Call("arc", point.X, point.Y, 1, 0, 2*math.Pi)
	ctx.Call("stroke")
}

func (point *Point) Values() (float64, float64) {
	return point.X, point.Y
}

func (point *Point) SquaredSize() float64 {
	return math.Pow(point.X, 2) + math.Pow(point.Y, 2)
}

func (point *Point) Size() float64 {
	return math.Sqrt(point.SquaredSize())
}

func (point *Point) Plus(other *Point) *Point {
	return &Point{X: point.X + other.X, Y: point.Y + other.Y}
}
func (point *Point) Minus(other *Point) *Point {
	return &Point{X: point.X - other.X, Y: point.Y - other.Y}
}
func (point *Point) Abs() *Point {
	return &Point{X: math.Abs(point.X), Y: math.Abs(point.Y)}
}
func (point *Point) Dot(other *Point) float64 {
	return point.X*other.X + point.Y*other.Y
}
func (point *Point) Times(other float64) *Point {
	return &Point{X: point.X * other, Y: point.Y * other}
}
func (point *Point) Avg(other *Point) *Point {
	return other.Minus(point).Times(0.5).Plus(point)
}

func (point *Point) DistanceTo(other *Point) float64 {
	return point.Minus(other).Size()
}
