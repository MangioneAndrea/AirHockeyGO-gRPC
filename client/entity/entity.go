package entity

import "syscall/js"

type Entity interface {
	Tick(delta int)
	Draw(ctx js.Value)
	OnConstruction(s *SceneController)
}
