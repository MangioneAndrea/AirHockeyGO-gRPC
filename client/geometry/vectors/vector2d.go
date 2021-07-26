package vectors

import (
	"math"
)

type Vector2D struct {
	X float64
	Y float64
}

func (vector *Vector2D) SquaredSize() float64 {
	return math.Pow(vector.X, 2) + math.Pow(vector.Y, 2)
}

func (vector *Vector2D) Size() float64 {
	return math.Sqrt(vector.SquaredSize())
}

func (vector *Vector2D) Plus(other *Vector2D) *Vector2D {
	return &Vector2D{X: vector.X + other.X, Y: vector.Y + other.Y}
}
func (vector *Vector2D) Minus(other *Vector2D) *Vector2D {
	return &Vector2D{X: vector.X - other.X, Y: vector.Y - other.Y}
}
func (vector *Vector2D) Abs() *Vector2D {
	return &Vector2D{X: math.Abs(vector.X), Y: math.Abs(vector.Y)}
}
func (vector Vector2D) Dot(other *Vector2D) float64 {
	return vector.X*other.X + vector.Y*other.Y
}
func (vector *Vector2D) Times(other float64) *Vector2D {
	return &Vector2D{X: vector.X * other, Y: vector.Y * other}
}

func (vector *Vector2D) DistanceTo(other *Vector2D) float64 {
	return vector.Minus(other).Size()
}
