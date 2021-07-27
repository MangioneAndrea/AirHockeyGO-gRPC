package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"

	"github.com/MangioneAndrea/airhockey/client/geometry/figures"
	"github.com/MangioneAndrea/airhockey/client/geometry/vectors"
	"github.com/MangioneAndrea/airhockey/gamepb"

	"github.com/hajimehoshi/ebiten/v2"
)

type GameMode int

const (
	SinglePlayer GameMode = iota
	MultiPlayer  GameMode = iota
)

var (
	constructed bool = false
	ball        PhisicSprite
	player      Sprite
	opponent    Sprite
	divider     = figures.NewRectangle(figures.NewPoint(0, screenHeight/2-2), screenWidth, 4)
	contours    = []*figures.Line{figures.NewLine(figures.NewPoint(1, 1), figures.NewPoint(1, screenHeight)),
		figures.NewLine(figures.NewPoint(1, 1), figures.NewPoint(screenWidth, 1)),
		figures.NewLine(figures.NewPoint(screenWidth, 1), figures.NewPoint(screenWidth, screenHeight)),
		figures.NewLine(figures.NewPoint(1, screenHeight-1), figures.NewPoint(screenWidth, screenHeight-1)),
	}
	updateStatus gamepb.PositionService_UpdateStatusClient
)

type Game struct {
	token *gamepb.Token
}

func (g *Game) Tick() error {
	if !constructed {
		return nil
	}
	cursorX, cursorY := ebiten.CursorPosition()
	delta := ebiten.CurrentTPS() / 60
	if delta == 0 {
		return nil
	}
	player.Rotation += 1 / delta
	player.Move(&vectors.Vector2D{
		X: math.Min((math.Max(float64(cursorX), 0)), screenWidth),
		Y: math.Min((math.Max(float64(cursorY), divider.Start.Y)), screenHeight),
	})

	err := updateStatus.Send(&gamepb.UserInput{
		Vector: &gamepb.Vector2D{X: int32(player.Hitbox.Center.X), Y: int32(player.Hitbox.Center.Y)},
		Token:  g.token,
	})

	if player.Hitbox.Intersects(ball.Sprite.Hitbox) {
		ball.AddForce(ball.Sprite.Hitbox.Center.Vector.Minus(player.Hitbox.Center.Vector), player.Speed)
	}

	if err != nil {
		fmt.Printf("Error while sending %v\n", err)
	}

	ball.Tick()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if ClientDebug {
		player.Hitbox.Center.LineTo(ball.Sprite.Hitbox.Center).DrawAxis(screen, figures.NewRectangle2(figures.NewPoint(0, 0), figures.NewPoint(screenWidth, screenHeight)))
		for _, rect := range contours {
			rect.Draw(screen)
		}
	}
	player.Draw(screen)
	ball.Draw(screen)
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
				opponent.Hitbox.Center.X = float64(res.GameStatus.Player2.X)
				opponent.Hitbox.Center.Y = float64(res.GameStatus.Player2.Y)
			} else {
				opponent.Hitbox.Center.X = float64(res.GameStatus.Player1.X)
				opponent.Hitbox.Center.Y = float64(res.GameStatus.Player1.Y)
			}
		}
	}()

	goo, _ := GetImageFromFilePath("client/graphics/gopher.png")

	ball = PhisicSprite{Sprite: &Sprite{
		Hitbox: figures.NewCircle(figures.NewPoint(float64(screenWidth)/2, float64(screenHeight)/1.3), 15),
		Image:  goo,
	},
		Direction:      &vectors.Vector2D{X: float64(screenWidth) / 2, Y: float64(screenHeight) / 1.3},
		LineCollisions: &contours,
	}
	player = Sprite{
		Hitbox: figures.NewCircle(
			figures.NewPoint(float64(screenWidth/2), float64(screenHeight)-player.Hitbox.Radius-25),
			math.Max(float64(goo.Bounds().Size().X)/2, float64(goo.Bounds().Size().Y)/2),
		),
		Image: goo,
	}
	opponent = Sprite{
		Image: goo,
		Hitbox: figures.NewCircle(
			figures.NewPoint(float64(screenWidth/2), float64(opponent.Hitbox.Radius+25)),
			math.Max(float64(player.Image.Bounds().Size().X)/2, float64(player.Image.Bounds().Size().Y)/2),
		),
	}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Airhockey go!")
	constructed = true
	return nil
}
