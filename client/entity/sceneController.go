package entity

import "github.com/MangioneAndrea/airhockey/gamepb"

type SceneController interface {
	ChangeScene(scene *Scene)
	GetConnection() gamepb.PositionServiceClient
	GetWidth() float32
	GetHeight() float32
}
