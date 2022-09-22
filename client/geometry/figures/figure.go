package figures

import "github.com/hajimehoshi/ebiten/v2"

type Figure interface {
	Intersects(other Figure) bool
	Draw(screen *ebiten.Image)
	GetAnchor() *Point
	SetAnchor(point *Point)
}
