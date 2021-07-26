package figures

import (
	"image/color"
	"math"

	"github.com/MangioneAndrea/airhockey/client/geometry/vectors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Circle struct {
	Center         *Point
	Radius         float64
	memoizedPoints []*vectors.Vector2D
}

func (circle *Circle) Intersects(elem Figure) bool {
	switch other := (elem).(type) {
	case *Circle:
		// Get the distance between the center of the 2 circles
		d := circle.Center.Vector.Minus(other.Center.Vector).Abs()
		// The distance is larger than the sum of both radiuses, too far
		if d.X > circle.Radius+other.Radius || d.Y > circle.Radius+other.Radius {
			return false
		}
		// Pitagoras
		return d.SquaredSize() < math.Pow(float64(circle.Radius+other.Radius), 2)
	case *Point:
		// Get the distance between the center and the point
		d := circle.Center.Vector.Minus(other.Vector).Abs()
		// The distance is larger than a square wrapping the circle
		if d.X > circle.Radius || d.Y > circle.Radius {
			return false
		}
		// The distance is less than a square contained by the circle or pitagora's
		return d.X+d.Y <= circle.Radius || d.SquaredSize() <= math.Pow(circle.Radius, 2)
	case *Rectangle:
		// Get the distance between the center and the starting point of the rectangle
		d := circle.Center.Vector.Minus(other.Start.Vector).Abs()
		// The distance is too far for an interception
		if d.X > other.Width/2+circle.Radius || d.Y > other.Height/2+circle.Radius {
			return false
		}
		// The distance is very little, meaning there must be an interception
		if d.X <= other.Width/2 || d.X <= other.Height/2 {
			return true
		}
		// The distance may overlap the corners of the rectangle (pitagoras)
		return (math.Pow(d.X-other.Width/2, 2)+math.Pow(d.Y-other.Height/2, 2) <= math.Pow(circle.Radius, 2))
	case *Segment:

	case *Line:
		// If the distance of the nearest point of the line is smaller than the radius --> intersection
		return circle.Center.DistanceToLine(other) < circle.Radius
	}

	return false
}

func (circle *Circle) Draw(screen *ebiten.Image) {
	// Memoize calc of the circle to speed up the process
	if circle.memoizedPoints == nil || len(circle.memoizedPoints) == 0 {
		for theta := float64(0); theta < 2*math.Pi; theta += math.Pi * 0.1 {
			x := +float64(circle.Radius) * math.Cos(theta)
			y := -float64(circle.Radius) * math.Sin(theta)
			circle.memoizedPoints = append(circle.memoizedPoints, &vectors.Vector2D{X: x, Y: y})
		}
	}
	for index, vector := range circle.memoizedPoints {
		var other *vectors.Vector2D
		if index == 0 {
			other = circle.memoizedPoints[len(circle.memoizedPoints)-1]
		} else {
			other = circle.memoizedPoints[index-1]
		}
		ebitenutil.DrawLine(
			screen,
			float64(circle.Center.X)+other.X,
			float64(circle.Center.Y)-other.Y,
			float64(circle.Center.X)+vector.X,
			float64(circle.Center.Y)-vector.Y,
			color.White)
	}
}
