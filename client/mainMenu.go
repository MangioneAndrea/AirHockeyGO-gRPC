package main

import (
	"log"

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
		log.Fatal(err)
	}
	button = &Button{
		X: screenWidth / 2, Y: screenHeight / 2, Image: buttonImage, OnClick: func() { gui.ChangeStage(&Game{}) },
	}
	return nil
}
