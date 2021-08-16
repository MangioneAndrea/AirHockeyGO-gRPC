package entities

import (
	"syscall/js"

	"github.com/MangioneAndrea/airhockey/client/geometry/figures"
)

type Actor interface {
	HasCollisions() bool
	GetHitbox() figures.Figure
	Tick(delta int)
	Draw(ctx js.Value)
	OnConstruction(s SceneController)
}
