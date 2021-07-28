package figures

import (
	"testing"
)

func TestNewCircle(t *testing.T) {
	circle := NewCircle(NewPoint(5, 9), 15)

	if circle.Center.X != 5 && circle.Center.Y != 9 {
		t.Errorf("The center of the circle was incorrect, got: %f %f, want: %f %f.", circle.Center.X, circle.Center.Y, 5., 9.)
	}
}

func TestIntersectsPoint(t *testing.T) {
	circle := NewCircle(NewPoint(15, 15), 15)
	pointInsideNear := NewPoint(7, 7)
	pointInsideFar := NewPoint(5, 5)
	pointOutsideNear := NewPoint(4, 4)
	pointOutsideFar := NewPoint(-1, 7)

	if !circle.Intersects(pointInsideNear) {
		t.Errorf("The circle O:(%f %f) R:%f should contain (%f %f)", circle.Center.X, circle.Center.Y, circle.Radius, pointInsideNear.X, pointInsideNear.Y)
	}
	if !circle.Intersects(pointInsideFar) {
		t.Errorf("The circle O:(%f %f) R:%f should contain (%f %f)", circle.Center.X, circle.Center.Y, circle.Radius, pointInsideFar.X, pointInsideFar.Y)
	}
	if circle.Intersects(pointOutsideNear) {
		t.Errorf("The circle O:(%f %f) R:%f should not contain (%f %f)", circle.Center.X, circle.Center.Y, circle.Radius, pointOutsideNear.X, pointOutsideNear.Y)
	}
	if circle.Intersects(pointOutsideFar) {
		t.Errorf("The circle O:(%f %f) R:%f should not contain (%f %f)", circle.Center.X, circle.Center.Y, circle.Radius, pointOutsideFar.X, pointOutsideFar.Y)
	}
}
func TestIntersectsCircle(t *testing.T) {
	circle := NewCircle(NewPoint(15, 15), 15)
	circleOutside := NewCircle(NewPoint(2, 2), 2)
	circleOverlap := NewCircle(NewPoint(2, 2), 4)
	circleInside := NewCircle(NewPoint(10, 10), 5)

	if circle.Intersects(circleOutside) {
		t.Errorf("The circle O:(%f %f) R:%f should intersect O:(%f %f) R:%f", circle.Center.X, circle.Center.Y, circle.Radius, circleOutside.Center.X, circleOutside.Center.Y, circleOutside.Radius)
	}
	if !circle.Intersects(circleOverlap) {
		t.Errorf("The circle O:(%f %f) R:%f should intersect O:(%f %f) R:%f", circle.Center.X, circle.Center.Y, circle.Radius, circleOverlap.Center.X, circleOverlap.Center.Y, circleOverlap.Radius)
	}
	if !circle.Intersects(circleInside) {
		t.Errorf("The circle O:(%f %f) R:%f should intersect O:(%f %f) R:%f", circle.Center.X, circle.Center.Y, circle.Radius, circleInside.Center.X, circleInside.Center.Y, circleInside.Radius)
	}
}
