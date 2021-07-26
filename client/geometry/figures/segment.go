package figures

type Segment struct {
	Start *Point
	End   *Point
}

func (segment *Segment) Intersects(elem Figure) bool {
	switch other := (elem).(type) {
	case *Segment:
		return false
	case *Rectangle:
	case *Point:
	case *Circle:
	case *Line:
	}
	return false
}

func NewSegment(start *Point, end *Point) *Segment {
	return &Segment{Start: start, End: end}
}
