package figures

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"

	"github.com/MangioneAndrea/airhockey/client/geometry/vectors"
)

type Point struct {
	Vector *vectors.Vector2D
	X      float64
	Y      float64
}

func NewPoint(x float64, y float64) *Point {
	return &Point{Vector: &vectors.Vector2D{X: x, Y: y}, X: x, Y: y}
}
func NewPoint2(vector *vectors.Vector2D) *Point {
	return &Point{Vector: vector, X: vector.X, Y: vector.Y}
}

func (point *Point) Intersects(elem Figure) bool {
	switch other := (elem).(type) {
	case *Point:
		return point.Vector.X == other.Vector.X && point.Vector.Y == other.Vector.Y
	case *Rectangle:
		return other.Start.X < point.Vector.X &&
			other.End.X > point.Vector.X &&
			other.Start.Y < point.Vector.Y &&
			other.End.Y > point.Vector.Y
	case *Circle:
		// Get the distance between the center and the point
		dx := math.Abs(point.Vector.X - other.Center.X)
		dy := math.Abs(point.Vector.Y - other.Center.Y)
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
	return point.Vector.DistanceTo(other.Vector)
}
func (point *Point) DistanceToLine(line *Line) float64 {
	return line.NearestPointTo(point).DistanceToPoint(point)
}

func (point *Point) LineTo(other *Point) *Line {
	return NewLine(point, other)
}
func (point *Point) SegmentTo(other *Point) *Segment {
	return NewSegment(point, other, "")
}

func (point *Point) Draw(screen *ebiten.Image) {

}

func (point *Point) GetAnchor() *Point { return point }

func (point *Point) SetAnchor(p *Point) {
	point = p
}
