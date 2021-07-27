package figures

import (
	"image/color"
	"math"
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
		// Get the distance between the center of the circle and the center of the rectangle
		d := rectangle.Start.Vector.Plus(rectangle.End.Vector).Times(0.5).Minus(other.Center.Vector).Abs()
		// The distance is bigger than the radius and half the length/height of the rectangle
		if d.X > rectangle.Width/2+other.Radius || d.Y > rectangle.Height/2+other.Radius {
			return false
		}
		// The center of the circle is inside the rectangle
		if rectangle.Intersects(other.Center) {
			return true
		}
		// The distance may overlap the corners of the rectangle (pitagoras)
		return (math.Pow(d.X-rectangle.Width/2, 2)+math.Pow(d.Y-rectangle.Height/2, 2) <= math.Pow(other.Radius, 2))
	}

	return false
}

func (rectangle *Rectangle) Sides() (bot *Segment, right *Segment, top *Segment, left *Segment) {
	return NewSegment(NewPoint(rectangle.Start.X, rectangle.End.Y), rectangle.End),
		NewSegment(NewPoint(rectangle.End.X, rectangle.Start.Y), rectangle.End),
		NewSegment(rectangle.Start, NewPoint(rectangle.End.X, rectangle.Start.Y)),
		NewSegment(rectangle.Start, NewPoint(rectangle.Start.X, rectangle.End.Y))
}
