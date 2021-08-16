package main

import (
	"context"
	"log"

	"github.com/MangioneAndrea/airhockey/client/geometry/vectors"
	"github.com/MangioneAndrea/airhockey/gamepb"
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

func (g *MainMenu) Draw() {
	button.Draw()
}

func (g *MainMenu) OnConstruction(screenWidth int, screenHeight int, gui *GUI) error {
	buttonImage, err := GetImageFromFilePath("client/graphics/button/idle.png")
	if err != nil {
		log.Fatal(err)
	}
	button = &Button{
		Position: vectors.Vector2D{X: float64(screenWidth / 2), Y: float64(screenHeight / 2)}, Image: &buttonImage, OnClick: func() {
			_, err := connection.RequestGame(context.Background(), &gamepb.GameRequest{})
			if err != nil {
				log.Fatal(err)
			}

			//gui.ChangeScene(&Game{token})
		},
	}
	return nil
}
