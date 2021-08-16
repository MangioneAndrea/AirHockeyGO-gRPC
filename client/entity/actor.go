package entity

import "github.com/MangioneAndrea/airhockey/client/geometry/figures"

type Actor struct {
	Collisions bool
	Hitbox     figures.Figure
}

func (s *Actor) OnConstruction(interface{}) {

}
func (s *Actor) Tick(delta int) {

}
