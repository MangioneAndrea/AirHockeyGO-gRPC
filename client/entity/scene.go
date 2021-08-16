package entity

import "syscall/js"

type Scene interface {
	GetActors() *[]*Actor
	Tick(delta int)
	Draw(ctx js.Value)
	OnConstruction(s SceneController)
}
