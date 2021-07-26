package figures

import (
	"image/color"

	"github.com/MangioneAndrea/airhockey/client/geometry/vectors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Line struct {
	Start      *Point
	Direction  *Point
	slope      float64
	yIntercept float64
}

func NewLine(start *Point, direction *Point) *Line {
	res := &Line{
		Start:     start,
		Direction: direction,
		slope:     (direction.Y - start.Y) / (direction.X - start.X),
	}
	res.yIntercept = res.slope*start.X + start.Y
	return res
}

func (line *Line) Intersects(elem Figure) bool {
	switch other := (elem).(type) {
	case *Line:
		// Either a different slope (1 point), or it's identical (infinite points)
		return line.Slope() != other.Slope() || line.YIntercept() == other.YIntercept()
	case *Segment:
		// TODO
	case *Point:
		return line.Slope() == NewLine(line.Start, other).Slope()
	case *Rectangle:
		bot, right, top, left := other.Sides()
		return line.Intersects(bot) || line.Intersects(right) || line.Intersects(top) || line.Intersects(left)
	case *Circle:
		// If the distance of the nearest point of the line is smaller than the radius --> intersection
		return other.Center.DistanceToLine(line) < other.Radius
	}
	return false
}

func (line *Line) Draw(screen *ebiten.Image) {
	ebitenutil.DrawLine(screen, line.Start.X, line.Start.Y, line.Direction.X, line.Direction.Y, color.White)
}

func (line *Line) DrawAxis(screen *ebiten.Image, bounds *Rectangle) {
	ebitenutil.DrawLine(screen,
		-line.YIntercept()/line.Slope()+bounds.Start.X,
		bounds.Start.Y,
		(bounds.End.X-line.YIntercept())/line.Slope(),
		bounds.End.Y,
		color.White)
}

func (line *Line) Slope() float64 {
	return line.slope
}

func (line *Line) YIntercept() float64 {
	return line.yIntercept
}
func (line *Line) Intersection(other *Line) *vectors.Vector2D {
	// Division by 0 (0 or infinite vectors)r
	if line.Slope() == other.Slope() {
		return nil
	}
	X := (other.YIntercept() - line.YIntercept()) / (other.Slope() - line.Slope())

	return &vectors.Vector2D{
		X: X,
		Y: line.Slope()*X + line.YIntercept(),
	}
}

func (line *Line) NearestPointTo(point *Point) *Point {
	t := ((point.Vector.X-line.Start.X)*line.Direction.X + (point.Vector.Y-line.Start.Y)*line.Direction.Y) / (line.Direction.Y*line.Direction.Y + line.Direction.X*line.Direction.X)
	Fx := line.Start.X + line.Direction.X*t
	Fy := line.Start.Y + line.Direction.Y*t
	return NewPoint(Fx, Fy)
}