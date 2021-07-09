package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Vector interface {
	Plus(*Vector) *Vector
	Minus(*Vector) *Vector
	Dot(*Vector) *Vector
	Times(*float64) *Vector

	Size() float64
}

type Vector2D struct {
	X float64
	Y float64
}

type Vector3D struct {
	X float64
	Y float64
	Z float64
}

func (vector *Vector2D) Size() float64 {
	return math.Sqrt(math.Pow(vector.X, 2) + math.Pow(vector.Y, 2))
}

func (vector *Vector2D) Plus(other *Vector2D) *Vector2D {
	return &Vector2D{X: vector.X + other.X, Y: vector.Y + other.Y}
}
func (vector *Vector2D) Minus(other *Vector2D) *Vector2D {
	return &Vector2D{X: vector.X - other.X, Y: vector.Y - other.Y}
}
func (vector Vector2D) Dot(other *Vector2D) *Vector2D {
	return &Vector2D{X: vector.X - other.X, Y: vector.Y - other.Y}
}
func (vector *Vector2D) Times(other float64) *Vector2D {
	return &Vector2D{X: vector.X * other, Y: vector.Y * other}
}

func (vector *Vector2D) DistanceTo(other *Vector2D) float64 {
	return vector.Minus(other).Size()
}
func (vector *Vector2D) To(other *Vector2D) *Line2D {
	return &Line2D{Start: vector, Direction: other}
}

type Line2D struct {
	Start     *Vector2D
	Direction *Vector2D
}

func (line *Line2D) Draw(screen *ebiten.Image) {
	ebitenutil.DrawLine(screen, line.Start.X, line.Start.Y, line.Direction.X, line.Direction.Y, color.White)
}

func (line *Line2D) DrawAxis(screen *ebiten.Image) {
	ebitenutil.DrawLine(screen,
		-line.YIntercept()/line.Slope(),
		0,
		(screenHeight-line.YIntercept())/line.Slope(),
		screenHeight,
		color.White)
}

func (line *Line2D) Slope() float64 {
	return (line.Direction.Y - line.Start.Y) / (line.Direction.X - line.Start.X)
}

func (line *Line2D) YIntercept() float64 {
	return -line.Slope()*line.Start.X + line.Start.Y
}
func (line *Line2D) Intersection(other *Line2D) *Vector2D {
	// Division by 0 (0 or infinite vectors)r
	if line.Slope() == other.Slope() {
		return nil
	}
	X := (other.YIntercept() - line.YIntercept()) / (other.Slope() - line.Slope())

	return &Vector2D{
		X: X,
		Y: line.Slope()*X + line.YIntercept(),
	}
}

/*
func (line *Line2D) Mirror(other *Line2D) *Line2D {

}
*/
