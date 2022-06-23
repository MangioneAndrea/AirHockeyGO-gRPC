package figures

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Line struct {
	Start      *Point
	Direction  *Point
	slope      float64
	yIntercept float64
}

func NewLine(start *Point, secondPoint *Point) *Line {
	res := &Line{
		Start:     start,
		Direction: NewPoint2(secondPoint.Vector.Minus(start.Vector).Normalize()),
	}
	if res.Direction.X == start.X {
		res.slope = math.Inf(0)
	} else {
		res.slope = (res.Direction.Y - start.Y) / (res.Direction.X - start.X)
	}
	res.yIntercept = -res.slope*start.X + start.Y
	return res
}

func (line *Line) Intersects(elem Figure) bool {
	switch other := (elem).(type) {
	case *Line:
		// Either a different slope (1 point), or it's identical (infinite points)
		return line.Slope() != other.Slope() || line.YIntercept() == other.YIntercept()
	case *Segment:
		// If slope and yintercepts are the same, the segment is contained in the line
		if other.Slope() == line.Slope() && other.YIntercept() == line.YIntercept() {
			return true
		}
		// If the intersection of the line and a virtual line along the segment is contained in the segment, it intersects
		return other.Intersects(line.LineIntersection(other.ToLine()))

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

func (line *Line) SnapSegment(screen *ebiten.Image, bounds *Rectangle) *Segment {
	bot, right, top, left := bounds.Sides()
	sides := []*Segment{bot, right, top, left}

	var points []*Point

	for _, side := range sides {
		p := line.SegmentIntersection(side)
		if p != nil {
			points = append(points[:], p)
		}
	}

	if len(points) != 2 {
		fmt.Println("ui ui")
		return nil
	}

	return NewSegment(points[0], points[1], "")
}

func (line *Line) Slope() float64 {
	return line.slope
}

func (line *Line) YIntercept() float64 {
	return line.yIntercept
}
func (line *Line) LineIntersection(other *Line) *Point {
	// Division by 0 (0 or infinite vectors)
	if line.Slope() == other.Slope() {
		return nil
	}
	var X float64
	var Y float64
	// The first line is not functional (vertical)
	if math.IsInf(line.Slope(), 0) {
		X = line.Direction.X
		Y = other.Slope()*X + other.YIntercept()
		// The second line is not functional (vertical)
	} else if math.IsInf(other.Slope(), 0) {
		X = other.Direction.X
		Y = line.Slope()*X + line.YIntercept()
	} else {
		X = (line.YIntercept() - other.YIntercept()) / (other.Slope() - line.Slope())
		Y = line.Slope()*X + line.YIntercept()
	}
	return NewPoint(X, Y)
}
func (line *Line) SegmentIntersection(other *Segment) *Point {
	p := line.LineIntersection(other.ToLine())
	if p == nil || other == nil || !other.Intersects(p) {
		return nil
	}
	return p
}

func (line *Line) NearestPointTo(point *Point) *Point {
	normDir := line.Direction.Vector.Normalize()
	lhs := point.Vector.Minus(line.Start.Vector)
	dotp := lhs.Dot(normDir)
	p := line.Start.Vector.Plus(normDir.Times(dotp))
	return NewPoint(p.X, p.Y)
}
