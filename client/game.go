package main

import (
	"context"
	"fmt"
	"image/color"
	"io"
	"log"
	"math"

	"github.com/MangioneAndrea/airhockey/gamepb"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameMode int

const (
	SinglePlayer GameMode = iota
	MultiPlayer  GameMode = iota
)

var (
	ball         Sprite
	player       Sprite
	opponent     Sprite
	divider      = Rectangle{X: 0, Y: screenHeight/2 - 2, Width: screenWidth, Height: 4, Color: color.White}
	updateStatus gamepb.PositionService_UpdateStatusClient
)

type Game struct {
	token *gamepb.Token
}

func (g *Game) Tick() error {
	cursorX, cursorY := ebiten.CursorPosition()
	delta := ebiten.CurrentTPS() / 60
	if delta == 0 {
		return nil
	}
	player.Rotation += 1 / delta
	player.X = int(math.Min((math.Max(float64(cursorX), 0)), screenWidth))
	player.Y = int(math.Min((math.Max(float64(cursorY), float64(divider.Y))), screenHeight))

	err := updateStatus.Send(&gamepb.UserInput{
		Vector: &gamepb.Vector2D{X: int32(player.X), Y: int32(player.Y)},
		Token:  g.token,
	})

	if err != nil {
		fmt.Printf("Error while sending %v\n", err)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	player.Draw(screen)
	opponent.Draw(screen)
	divider.Draw(screen)
}

func (g *Game) OnConstruction(screenWidth int, screenHeight int, gui *GUI) error {

	stream, streamErr := connection.UpdateStatus(context.Background())
	if streamErr != nil {
		log.Fatal(streamErr)
	}
	updateStatus = stream

	go func() {
		for {
			res, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				log.Fatalf("Error while receiving %v", err)
			}
			if res.Token1.PlayerHash == g.token.PlayerHash {
				opponent.X = int(res.GameStatus.Player2.X)
				opponent.Y = int(res.GameStatus.Player2.Y)
			} else {
				opponent.X = int(res.GameStatus.Player1.X)
				opponent.Y = int(res.GameStatus.Player1.Y)
			}
			//fmt.Printf("%v - %v \n", res.Player1.X, res.Player1.Y)
		}
	}()

	goo, _ := GetImageFromFilePath("client/graphics/gopher.png")

	ball = Sprite{Image: goo}
	ball.Width = ball.Image.Bounds().Size().X
	ball.Height = ball.Image.Bounds().Size().Y

	player = Sprite{
		Image: goo,
	}
	player.Width = ball.Image.Bounds().Size().X
	player.Height = ball.Image.Bounds().Size().Y
	player.X = screenWidth / 2
	player.Y = screenHeight - player.Height - 25
	opponent = Sprite{
		Image: goo,
	}
	opponent.Width = ball.Image.Bounds().Size().X
	opponent.Height = ball.Image.Bounds().Size().Y
	opponent.X = screenWidth / 2
	opponent.Y = opponent.Height + 25

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Animation (Ebiten Demo)")

	return nil
}
