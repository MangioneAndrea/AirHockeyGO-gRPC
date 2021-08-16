package game

import (
	"context"
	"log"
	"syscall/js"

	"github.com/MangioneAndrea/airhockey/client/entity"
	"github.com/MangioneAndrea/airhockey/client/geometry/vectors"
	"github.com/MangioneAndrea/airhockey/gamepb"
)

var (
	button *Button
)

type MainMenu struct {
	actors *[]*entity.Actor
}

func (g *MainMenu) GetActors() *[]*entity.Actor {
	return g.actors
}

func (g *MainMenu) Tick(delta int) {
	button.CheckClicked()
}

func (g *MainMenu) Draw(ctx js.Value) {
	for _, actor := range *g.actors {
		(*actor).Draw(canvas)
	}
}

func (g *MainMenu) OnConstruction(c entity.SceneController) {
	/*
		buttonImage, err := GetImageFromFilePath("client/graphics/button/idle.png")
		if err != nil {
			log.Fatal(err)
		}
	*/
	button = &Button{
		Position: vectors.Vector2D{X: float64(c.GetWidth() / 2), Y: float64(c.GetHeight() / 2)}, OnClick: func() {
			_, err := c.GetConnection().RequestGame(context.Background(), &gamepb.GameRequest{})
			if err != nil {
				log.Fatal(err)
			}

			//c.ChangeScene(&Game{token})
		},
	}
}
