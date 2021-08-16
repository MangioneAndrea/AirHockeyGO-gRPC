package actors

import (
	"syscall/js"

	"github.com/MangioneAndrea/airhockey/client/entities"
	"github.com/MangioneAndrea/airhockey/client/geometry/figures"
)

type Button struct {
	Figure    figures.Figure
	OnClick   func()
	isClicked bool
}

func NewButton(fig figures.Figure, onClick func()) *Button {
	return &Button{
		Figure:  fig,
		OnClick: onClick,
	}
}

func (button *Button) OnConstruction(s entities.SceneController) {
	s.GetCanvas().Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if button.Figure.Intersects(figures.NewPoint(float64(args[0].Get("layerX").Int()), float64(args[0].Get("layerY").Int()))) {
			button.OnClick()
		}
		return nil
	}))
}
func (button *Button) Tick(delta int) {
	return
}
func (button *Button) Draw(ctx js.Value) {
	button.Figure.Draw(ctx)
}

func (button *Button) HasCollisions() bool {
	return false
}
func (button *Button) GetHitbox(delta int) {
	return
}
