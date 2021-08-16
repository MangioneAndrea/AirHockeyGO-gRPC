package figures

import (
	"math"
	"syscall/js"
)

type Circle struct {
	Center *Point
	Radius float64
}

func NewCircle(center *Point, radius float64) *Circle {
	return &Circle{
		Center: center,
		Radius: radius,
	}
}

func (circle *Circle) Intersects(elem Figure) bool {
	switch other := (elem).(type) {
	case *Circle:
		// Get the distance between the center of the 2 circles
		d := circle.Center.Minus(other.Center).Abs()
		// The distance is larger than the sum of both radiuses, too far
		if d.X > circle.Radius+other.Radius || d.Y > circle.Radius+other.Radius {
			return false
		}
		// Pitagoras
		return d.SquaredSize() < math.Pow(float64(circle.Radius+other.Radius), 2)
	case *Point:
		// Get the distance between the center and the point
		d := circle.Center.Minus(other).Abs()
		// The distance is larger than a square wrapping the circle
		if d.X > circle.Radius || d.Y > circle.Radius {
			return false
		}
		// The distance is less than a square contained by the circle or pitagora's
		return d.X+d.Y <= circle.Radius || d.SquaredSize() <= math.Pow(circle.Radius, 2)
	case *Rectangle:
		// Get the distance between the center of the circle and the center of the rectangle
		d := other.Start.Plus(other.End).Times(0.5).Minus(circle.Center).Abs()
		// The distance is bigger than the radius and half the length/height of the rectangle
		if d.X > other.Width/2+circle.Radius || d.Y > other.Height/2+circle.Radius {
			return false
		}
		// The center of the circle is inside the rectangle
		if other.Intersects(circle.Center) {
			return true
		}
		// The distance may overlap the corners of the rectangle (pitagoras)
		return (math.Pow(d.X-other.Width/2, 2)+math.Pow(d.Y-other.Height/2, 2) <= math.Pow(circle.Radius, 2))
	case *Segment:
		p := other.ToLine().NearestPointTo(circle.Center)
		// If the point is not in the segment, take the nearest end
		if !p.Intersects(other) {
			if other.Start.DistanceTo(p) < other.End.DistanceTo(p) {
				p = other.Start
			} else {
				p = other.End
			}
		}
		// if the nearest point is in the circle, the segment intersects it
		return other.Intersects(p)
	case *Line:
		// If the distance of the nearest point of the line is smaller than the radius --> intersection
		return circle.Center.DistanceToLine(other) < circle.Radius
	}

	return false
}
func (circle *Circle) GetCenter() *Point {
	return circle.Center
}
func (circle *Circle) MoveTo(where *Point) {
	circle.Center.MoveTo(where)
}

func (circle *Circle) Draw(ctx js.Value) {
	ctx.Call("beginPath")
	ctx.Call("arc", circle.Center.X, circle.Center.Y, circle.Radius, 0, 2*math.Pi)
	ctx.Call("stroke")
}
