package main

import (
	"context"
	"log"
	"math"

	"github.com/MangioneAndrea/airhockey/gamepb"
	"github.com/hajimehoshi/ebiten/v2"
	"google.golang.org/grpc"
)

const (
	screenWidth  = 600
	screenHeight = 1200
	wishedFPS    = 60
)

var (
	connection   gamepb.PositionServiceClient
	updateStatus gamepb.PositionService_UpdateStatusClient
)

type Actor interface {
	Tick() error
	OnConstruction(int, int) error
}

type Stage interface {
	Tick() error
	OnConstruction(int, int, *GUI) error
	Draw(screen *ebiten.Image)
}

type GUI struct {
	stage Stage
}

func (g *GUI) Update() error {
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
		Token:  nil,
	})

	if err := g.stage.Tick(); err != nil {
		println(err.Error())
	}
	return nil
}

func (g *GUI) Draw(screen *ebiten.Image) {
	g.stage.Draw(screen)
}

func (g *GUI) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *GUI) ChangeStage(stage Stage) {
	g.stage = stage
	stage.OnConstruction(screenWidth, screenHeight, g)
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Airhockey")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	connection = gamepb.NewPositionServiceClient(cc)

	stream, streamErr := connection.UpdateStatus(context.Background())
	updateStatus = stream
	if streamErr != nil {
		log.Fatal(streamErr)
	}

	g := &GUI{}
	g.ChangeStage(&MainMenu{})
	guiError := ebiten.RunGame(g)
	if guiError != nil {
		log.Fatal(guiError)
	}

}
