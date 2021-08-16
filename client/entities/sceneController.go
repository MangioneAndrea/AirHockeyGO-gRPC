package entities

import (
	"syscall/js"

	"github.com/MangioneAndrea/airhockey/gamepb"
)

type SceneController interface {
	ChangeScene(scene Scene)
	GetConnection() gamepb.PositionServiceClient
	GetWidth() float32
	GetHeight() float32
	GetCanvas() js.Value
	GetCtx() js.Value
}
