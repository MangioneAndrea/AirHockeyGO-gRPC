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
	divider      = figures.NewRectangle(figures.NewPoint(0, screenHeight/2-2), screenWidth, 4)
	contours     = figures.NewRectangle(figures.NewPoint(1, 1), screenWidth-2, screenHeight-2)
	updateStatus gamepb.PositionService_UpdateStatusClient
)

type Game struct {
	token *gamepb.Token
}

func (g *Game) Tick() error {
	if !constructed {
		return nil
	}
	cursorX, cursorY := .0, .0 //ebiten.CursorPosition()
	delta := 30.               //ebiten.CurrentTPS() / 60
	if delta == 0 {
		return nil
	}
	player.Rotation += 1 / delta
	player.Move(&vectors.Vector2D{
		X: math.Min((math.Max(float64(cursorX), 0)), screenWidth),
		Y: math.Min((math.Max(float64(cursorY), 0)), screenHeight),
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

func (g *Game) Draw() {
	if ClientDebug {
		s := player.Hitbox.Center.LineTo(ball.Sprite.Hitbox.Center).SnapSegment(contours)

		if s != nil {
			/*
				ebitenutil.DrawLine(screen,
					s.Start.X,
					s.Start.Y,
					s.End.X,
					s.End.Y,
					color.White)
			*/
		}
		//contours.Draw(screen)
	}
	//player.Draw(screen)
	//ball.Draw(screen)
	//opponent.Draw(screen)
	//divider.Draw(screen)
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

	bot, right, top, left := contours.Sides()

	ball = PhisicSprite{Sprite: &Sprite{
		Hitbox: figures.NewCircle(figures.NewPoint(float64(screenWidth)/2, float64(screenHeight)/1.3), 15),
		Image:  goo,
	},
		Direction:  &vectors.Vector2D{X: float64(screenWidth) / 2, Y: float64(screenHeight) / 1.3},
		Collisions: &[]figures.Figure{bot, right, top, left},
	}
	radius := math.Max(float64(goo.Bounds().Size().X)/2, float64(goo.Bounds().Size().Y)/2)
	player = Sprite{
		Hitbox: figures.NewCircle(
			figures.NewPoint(float64(screenWidth/2), float64(screenHeight)-radius-25),
			radius,
		),
		Image: goo,
	}
	opponent = Sprite{
		Image: goo,
		Hitbox: figures.NewCircle(
			figures.NewPoint(float64(screenWidth/2), float64(radius+25)),
			radius,
		),
	}
	//ebiten.SetWindowSize(screenWidth, screenHeight)
	//ebiten.SetWindowTitle("Airhockey go!")
	constructed = true
	return nil
}
