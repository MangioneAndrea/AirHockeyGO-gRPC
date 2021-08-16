package figures

import (
	"math"
	"syscall/js"
)

type Segment struct {
	Start      *Point
	End        *Point
	slope      float64
	yIntercept float64
}

func NewSegment(start *Point, end *Point) *Segment {
	res := &Segment{Start: start, End: end}
	if end.X == start.X {
		res.slope = math.Inf(1)
	} else {
		res.slope = (end.Y - start.Y) / (end.X - start.X)
	}
	res.yIntercept = res.slope*start.X + start.Y
	return res
}

func (segment *Segment) Intersects(elem Figure) bool {
	switch other := (elem).(type) {
	case *Segment:
		// The rightmost point of A is more left than the leftmost of B --> no intersection
		if math.Max(segment.Start.X, segment.End.X) < math.Min(other.Start.X, other.End.X) {
			return false
		}
		interceptionSlope := (other.yIntercept - segment.yIntercept) / (other.slope - segment.slope)

		return interceptionSlope >= math.Max(math.Min(segment.Start.X, segment.End.X), math.Min(other.Start.X, other.End.X)) &&
			interceptionSlope <= math.Min(math.Max(segment.Start.X, segment.End.X), math.Max(other.Start.X, other.End.X))
	case *Rectangle:
		bot, right, top, left := other.Sides()
		return bot.Intersects(segment) || right.Intersects(segment) || top.Intersects(segment) || left.Intersects(segment)
	case *Point:
		// The point is between the 2 x's and the slope is the same (same immaginary line)
		return other.Y <= math.Max(segment.Start.Y, segment.End.Y) && other.Y >= math.Min(segment.Start.Y, segment.End.Y) &&
			other.X <= math.Max(segment.Start.X, segment.End.X) && other.X >= math.Min(segment.Start.X, segment.End.X) &&
			segment.Slope() == NewLine(segment.Start, other).Slope()
	case *Circle:
		p := segment.ToLine().NearestPointTo(other.Center)
		// If the point is not in the segment, take the nearest end
		if !p.Intersects(segment) {
			if segment.Start.Vector.DistanceTo(p.Vector) < segment.End.Vector.DistanceTo(p.Vector) {
				p = segment.Start
			} else {
				p = segment.End
			}
		}
		// if the nearest point is in the circle, the segment intersects it
		return other.Intersects(p)
	case *Line:
		// If slope and yintercepts are the same, the segment is contained in the line
		if segment.Slope() == other.Slope() && segment.YIntercept() == other.YIntercept() {
			return true
		}
		// If the intersection of the line and a virtual line along the segment is contained in the segment, it intersects
		return segment.Intersects(other.LineIntersection(segment.ToLine()))
	}
	return false
}

func (segment *Segment) Slope() float64 {
	return segment.slope
}

func (segment *Segment) YIntercept() float64 {
	return segment.yIntercept
}

func (segment *Segment) ToLine() *Line {
	return &Line{Start: segment.Start, Direction: segment.End, slope: segment.slope, yIntercept: segment.yIntercept}
}

func (segment *Segment) Draw(ctx js.Value) {
	ctx.Call("beginPath")
	ctx.Call("moveTo", segment.Start.X, segment.Start.Y)
	ctx.Call("lineTo", segment.End.X, segment.End.Y)
	ctx.Call("stroke")
}
