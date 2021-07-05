package main

import (
	"context"
	"image/color"
	"log"
	"math"

	"github.com/MangioneAndrea/airhockey/gamepb"
	"google.golang.org/grpc"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameMode int

const (
	SinglePlayer GameMode = iota
	MultiPlayer  GameMode = iota
)

var (
	ball    Sprite
	player1 Sprite
	player2 Sprite
	divider = Rectangle{X: 0, Y: screenHeight/2 - 2, Width: screenWidth, Height: 4, Color: color.White}
)

type Game struct {
}

func (g *Game) Tick() error {
	cursorX, cursorY := ebiten.CursorPosition()
	delta := ebiten.CurrentTPS() / 60
	if delta == 0 {
		return nil
	}
	player1.Rotation += 1 / delta
	player1.X = int(math.Min((math.Max(float64(cursorX), 0)), screenWidth))
	player1.Y = int(math.Min((math.Max(float64(cursorY), float64(divider.Y))), screenHeight))

	updateStatus.Send(&gamepb.UserInput{
		Vector: &gamepb.Vector2D{X: int32(player1.X), Y: int32(player1.Y)},
		//Token:  &g.token,
	})

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	player1.Draw(screen)
	player2.Draw(screen)
	divider.Draw(screen)
}

func (g *Game) OnConstruction(screenWidth int, screenHeight int, gui *GUI) error {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	connection = gamepb.NewPositionServiceClient(cc)
	goo, _ := GetImageFromFilePath("client/graphics/gopher.png")

	ball = Sprite{Image: goo}
	ball.Width = ball.Image.Bounds().Size().X
	ball.Height = ball.Image.Bounds().Size().Y

	player1 = Sprite{
		Image: goo,
	}
	player1.Width = ball.Image.Bounds().Size().X
	player1.Height = ball.Image.Bounds().Size().Y
	player1.X = screenWidth / 2
	player1.Y = screenHeight - player1.Height - 25
	player2 = Sprite{
		Image: goo,
	}
	player2.Width = ball.Image.Bounds().Size().X
	player2.Height = ball.Image.Bounds().Size().Y
	player2.X = screenWidth / 2
	player2.Y = player2.Height + 25

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Animation (Ebiten Demo)")

	stream, err := connection.UpdateStatus(context.Background())
	updateStatus = stream
	return nil
}
