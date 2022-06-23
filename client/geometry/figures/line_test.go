package figures

import (
	"github.com/google/go-cmp/cmp"
	"math"
	"testing"
)

var compareFloats = cmp.Comparer(func(x, y float64) bool {
	diff := math.Abs(x - y)
	mean := math.Abs(x+y) / 2.0
	if math.IsNaN(diff / mean) {
		return true
	}
	return (diff / mean) < 0.00001
})

func TestLine_NearestPointToVerticalX(t *testing.T) {
	line := NewLine(NewPoint(10, 10), NewPoint(10, 100))
	point := NewPoint(234, 75)

	np := line.NearestPointTo(point)

	if !cmp.Equal(line.Start.X, np.X, compareFloats) {
		t.Errorf("The nearest point of a vertical line should have the same X as the line. Wanted: %f Got: %f", line.Start.X, np.X)
	}
}

func TestLine_NearestPointToVerticalY(t *testing.T) {
	line := NewLine(NewPoint(10, 10), NewPoint(10, 100))
	point := NewPoint(234, 75)

	np := line.NearestPointTo(point)

	if !cmp.Equal(point.Y, np.Y, compareFloats) {
		t.Errorf("The nearest point of a vertical line should have the same Y as the point. Wanted: %f Got: %f", point.Y, np.Y)
	}
}

func TestLine_NearestPointToHorizontalY(t *testing.T) {
	line := NewLine(NewPoint(10, 10), NewPoint(100, 10))
	point := NewPoint(75, 234)

	np := line.NearestPointTo(point)

	if !cmp.Equal(line.Start.Y, np.Y, compareFloats) {
		t.Errorf("The nearest point of a vertical line should have the same Y as the line. Wanted: %f Got: %f", line.Start.Y, np.Y)
	}
}

func TestLine_NearestPointToHorizontalX(t *testing.T) {
	line := NewLine(NewPoint(10, 10), NewPoint(100, 10))
	point := NewPoint(75, 234)

	np := line.NearestPointTo(point)

	if !cmp.Equal(point.X, np.X, compareFloats) {
		t.Errorf("The nearest point of a vertical line should have the same X as the point. Wanted: %f Got: %f", point.X, np.X)
	}
}
