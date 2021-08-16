package entity

type Scene struct {
	Actors *[]*Actor
}

func (s *Scene) OnConstruction(interface{}) {

}
func (s *Scene) Tick(delta int) {

}
