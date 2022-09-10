package main

import (
	"context"
	"github.com/MangioneAndrea/airhockey/client/geometry/vectors"
	"github.com/MangioneAndrea/airhockey/gamepb"
	"github.com/MangioneAndrea/gonsole"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	button *Button
)

type MainMenu struct {
}

func (g *MainMenu) Tick() error {
	button.CheckClicked()
	return nil
}

func (g *MainMenu) Draw(screen *ebiten.Image) {
	button.Draw(screen)
}

func (g *MainMenu) OnConstruction(screenWidth int, screenHeight int, gui *GUI) error {
	buttonImage, err := GetImageFromFilePath("client/graphics/button/idle.png")
	if err != nil {
		gonsole.Error(err)
	}
	button = &Button{
		Position: vectors.Vector2D{X: float64(screenWidth / 2), Y: float64(screenHeight) / 1.3}, Image: buttonImage, OnClick: func() {
			token, err := connection.RequestGame(context.Background(), &gamepb.GameRequest{})
			if err != nil {
				gonsole.Error(err)
			}

			gui.ChangeStage(&Game{token})
		},
	}
	return nil
}
