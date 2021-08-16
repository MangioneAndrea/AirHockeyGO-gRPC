package entity

type Entity interface {
	Tick(delta int)
	OnConstruction(interface{})
}
