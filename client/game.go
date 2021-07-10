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
	DEBUG                 = false
	SinglePlayer GameMode = iota
	MultiPlayer  GameMode = iota
)

var (
	ball         Circle
	player       Sprite
	opponent     Sprite
	divider      = Rectangle{Position: Vector2D{X: 0, Y: screenHeight/2 - 2}, Width: screenWidth, Height: 4, Color: color.White}
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
	player.Move(&Vector2D{
		X: math.Min((math.Max(float64(cursorX), 0)), screenWidth),
		Y: math.Min((math.Max(float64(cursorY), divider.Position.Y)), screenHeight),
	})

	err := updateStatus.Send(&gamepb.UserInput{
		Vector: &gamepb.Vector2D{X: int32(player.Position.X), Y: int32(player.Position.Y)},
		Token:  g.token,
	})

	if err != nil {
		fmt.Printf("Error while sending %v\n", err)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	player.Draw(screen)
	ball.Draw(screen)
	opponent.Draw(screen)
	divider.Draw(screen)

	player.Position.To(&ball.Center).DrawAxis(screen)
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
				opponent.Position.X = float64(res.GameStatus.Player2.X)
				opponent.Position.Y = float64(res.GameStatus.Player2.Y)
			} else {
				opponent.Position.X = float64(res.GameStatus.Player1.X)
				opponent.Position.Y = float64(res.GameStatus.Player1.Y)
			}
		}
	}()

	goo, _ := GetImageFromFilePath("client/graphics/gopher.png")

	ball = Circle{Center: Vector2D{X: float64(screenWidth / 2), Y: float64(screenHeight / 2)}, Radius: 15}
	/*
		ball = Sprite{Image: goo, Position: &Vector2D{X: float64(screenWidth / 2), Y: float64(screenHeight / 2)}}
		ball.Width = ball.Image.Bounds().Size().X
		ball.Height = ball.Image.Bounds().Size().Y
	*/
	player = Sprite{
		Image: goo, Position: &Vector2D{X: 0, Y: 0},
	}
	player.Width = player.Image.Bounds().Size().X
	player.Height = player.Image.Bounds().Size().Y
	player.Position.X = float64(screenWidth / 2)
	player.Position.Y = float64(screenHeight - player.Height - 25)
	opponent = Sprite{
		Image: goo, Position: &Vector2D{X: float64(screenWidth / 2), Y: float64(opponent.Height + 25)},
	}
	opponent.Width = opponent.Image.Bounds().Size().X
	opponent.Height = opponent.Image.Bounds().Size().Y

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Airhockey go!")

	return nil
}
