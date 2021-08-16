package game

import (
	"syscall/js"

	"github.com/MangioneAndrea/airhockey/client/entities"
	"github.com/MangioneAndrea/airhockey/client/entities/actors"
	"github.com/MangioneAndrea/airhockey/client/geometry/figures"
)

var (
	button *actors.Button
)

type MainMenu struct {
	actors []entities.Actor
}

func (g *MainMenu) GetActors() *[]entities.Actor {
	return &g.actors
}

func (g *MainMenu) Tick(delta int) {
}

func (g *MainMenu) Draw(ctx js.Value) {
	for _, actor := range g.actors {
		actor.Draw(ctx)
	}
}

func (g *MainMenu) OnConstruction(c entities.SceneController) {
	button = actors.NewButton(
		figures.NewCircle(figures.NewPoint(75, 75), 50),
		func() { println("clicked") },
	)
	button.OnConstruction(c)
	g.actors = append(g.actors, button)
}
