package entities

import (
	"syscall/js"
)

type Actor interface {
	HasCollisions() bool
	GetHitbox(delta int)
	Tick(delta int)
	Draw(ctx js.Value)
	OnConstruction(s SceneController)
}
