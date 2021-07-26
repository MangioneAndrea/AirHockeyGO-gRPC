package figures

import (
	"image/color"
)

type Rectangle struct {
	Start  *Point
	End    *Point
	Width  float64
	Height float64
	Color  color.Color
}

func NewRectangle(start *Point, width float64, height float64) *Rectangle {
	return &Rectangle{
		Start:  start,
		End:    NewPoint(start.X+width, start.Y+height),
		Width:  width,
		Height: height,
		Color:  color.White,
	}
}

func NewRectangle2(start *Point, end *Point) *Rectangle {
	return &Rectangle{
		Start:  start,
		End:    end,
		Width:  end.X - start.X,
		Height: end.Y - start.Y,
		Color:  color.White,
	}
}

func (rectangle *Rectangle) Intersects(elem Figure) bool {
	switch other := (elem).(type) {
	case *Rectangle:
		return rectangle.Start.X < other.End.X &&
			rectangle.End.X > other.Start.X &&
			rectangle.Start.Y < other.End.Y &&
			rectangle.End.Y > other.Start.Y

	case *Point:
		return rectangle.Start.X < other.Vector.X &&
			rectangle.End.X > other.Vector.X &&
			rectangle.Start.Y < other.Vector.Y &&
			rectangle.End.Y > other.Vector.Y
	case *Line, *Segment:
		bot, right, top, left := rectangle.Sides()
		return bot.Intersects(other) || right.Intersects(other) || top.Intersects(other) || left.Intersects(other)
	case *Circle:
	}

	return false
}

func (rectangle *Rectangle) Sides() (bot *Segment, right *Segment, top *Segment, left *Segment) {
	return NewSegment(NewPoint(rectangle.Start.X, rectangle.End.Y), rectangle.End),
		NewSegment(NewPoint(rectangle.End.X, rectangle.Start.Y), rectangle.End),
		NewSegment(rectangle.Start, NewPoint(rectangle.End.X, rectangle.Start.Y)),
		NewSegment(rectangle.Start, NewPoint(rectangle.Start.X, rectangle.End.Y))
}
