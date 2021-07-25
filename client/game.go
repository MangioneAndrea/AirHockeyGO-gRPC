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
	constructed  bool = false
	ball         PhisicSprite
	player       Sprite
	opponent     Sprite
	divider      = Rectangle{Position: Vector2D{X: 0, Y: screenHeight/2 - 2}, Width: screenWidth, Height: 4, Color: color.White}
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
	player.Move(&Vector2D{
		X: math.Min((math.Max(float64(cursorX), 0)), screenWidth),
		Y: math.Min((math.Max(float64(cursorY), divider.Position.Y)), screenHeight),
	})

	err := updateStatus.Send(&gamepb.UserInput{
		Vector: &gamepb.Vector2D{X: int32(player.Hitbox.Center.X), Y: int32(player.Hitbox.Center.Y)},
		Token:  g.token,
	})

	if player.Hitbox.Intersects(ball.Sprite.Hitbox) {
		ball.AddForce(ball.Sprite.Hitbox.Center.Minus(player.Hitbox.Center), 5)
	}

	if err != nil {
		fmt.Printf("Error while sending %v\n", err)
	}

	ball.Tick()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if ClientDebug {
		player.Hitbox.Center.To(ball.Sprite.Hitbox.Center).DrawAxis(screen)
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
		Hitbox: &Circle{Center: &Vector2D{X: float64(screenWidth) / 2, Y: float64(screenHeight) / 1.3}, Radius: 15},

		Image: goo,
	},
		Direction: &Vector2D{X: float64(screenWidth) / 2, Y: float64(screenHeight) / 1.3},
	}
	player = Sprite{
		Hitbox: &Circle{
			Center: &Vector2D{X: 0, Y: 0},
		},
		Image: goo,
	}
	player.Hitbox.Radius = int(math.Max(float64(player.Image.Bounds().Size().X)/2, float64(player.Image.Bounds().Size().Y)/2))
	player.Hitbox.Center.X = float64(screenWidth / 2)
	player.Hitbox.Center.Y = float64(screenHeight - player.Hitbox.Radius - 25)
	opponent = Sprite{
		Image: goo,
		Hitbox: &Circle{
			Radius: int(math.Max(float64(player.Image.Bounds().Size().X)/2, float64(player.Image.Bounds().Size().Y)/2)),
		},
	}
	opponent.Hitbox.Center = &Vector2D{X: float64(screenWidth / 2), Y: float64(opponent.Hitbox.Radius + 25)}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Airhockey go!")
	constructed = true
	return nil
}
