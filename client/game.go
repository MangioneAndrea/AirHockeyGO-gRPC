package main

import (
	"context"
	"fmt"
	"github.com/MangioneAndrea/gonsole"
	"image/color"
	"io"
	"math"
	"time"

	"github.com/MangioneAndrea/airhockey/client/geometry/figures"
	"github.com/MangioneAndrea/airhockey/client/geometry/vectors"
	"github.com/MangioneAndrea/airhockey/gamepb"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const goalWidth = 50
const goalHeight = 10

var (
	constructed bool = false
	ball        PhisicSprite
	player      Sprite
	opponent    Sprite
	divider     = figures.NewRectangle(figures.NewPoint(0, screenHeight/2-2), screenWidth, 4)
	contours    = figures.NewRectangle(figures.NewPoint(1, 1), screenWidth-2, screenHeight-2)
	goal1       = Sprite{
		Hitbox: figures.NewRectangle(figures.NewPoint(screenWidth/2-goalWidth, 0), goalWidth*2, goalHeight),
	}
	goal2 = Sprite{
		Hitbox: figures.NewRectangle(figures.NewPoint(screenWidth/2-goalWidth, screenHeight-goalHeight), goalWidth*2, goalHeight),
	}
	updateStatus gamepb.PositionService_UpdateStatusClient
	lastUpdate   int64 = 0
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
		X: math.Min(math.Max(float64(cursorX), 0), screenWidth),
		Y: math.Min((math.Max(float64(cursorY), screenHeight/2)), screenHeight),
	})

	update := &gamepb.UserInput{
		Vector: &gamepb.Vector2D{X: int32((player.Hitbox).GetAnchor().X), Y: int32((player.Hitbox).GetAnchor().Y)},
		Token:  g.token,
	}

	_, firstIntersection := player.Intersects(ball.Sprite.Hitbox)
	if firstIntersection {
		ball.AddForce((ball.Sprite.Hitbox).GetAnchor().Vector.Minus((player.Hitbox).GetAnchor().Vector), player.Speed)

		update.DiskStatus = &gamepb.DiskStatus{
			Position:   &gamepb.Vector2D{X: int32((ball.Sprite.Hitbox).GetAnchor().X), Y: int32((ball.Sprite.Hitbox).GetAnchor().Y)},
			LastUpdate: time.Now().Unix(),
			Force:      &gamepb.Vector2D{X: int32(ball.Direction.X), Y: int32(ball.Direction.Y)},
			Speed:      float32(ball.Sprite.Speed),
		}

		if ClientDebug {
			fmt.Println("sending disk position to server")
		}

	}

	err := updateStatus.Send(update)

	if err != nil {
		gonsole.Error(err, "stream.send")
	}
	ball.Tick()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if ClientDebug {
		ebitenutil.DrawLine(screen,
			(player.Hitbox).GetAnchor().X,
			(player.Hitbox).GetAnchor().Y,
			(ball.Sprite.Hitbox).GetAnchor().X,
			(ball.Sprite.Hitbox).GetAnchor().Y,
			color.RGBA{R: 255, G: 0, B: 0, A: 255})
		contours.Draw(screen)
	}
	player.Draw(screen)
	ball.Draw(screen)
	opponent.Draw(screen)
	divider.Draw(screen)
	goal1.Draw(screen)
	goal2.Draw(screen)
}

func (g *Game) OnConstruction(screenWidth int, screenHeight int, gui *GUI) error {

	stream, streamErr := connection.UpdateStatus(context.Background())
	if streamErr != nil {
		gonsole.Error(streamErr, "UpdateStatus")
	}
	updateStatus = stream

	go func() {
		for {
			res, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				gonsole.Error(err, "Error while receiving")
				defer g.OnConstruction(screenWidth, screenHeight, gui)
				return
			}
			if res.Token1.PlayerHash == g.token.PlayerHash {
				(opponent.Hitbox).SetAnchor(figures.NewPoint(float64(res.GameStatus.Player2.X), float64(res.GameStatus.Player2.Y)))
			} else {
				(opponent.Hitbox).SetAnchor(figures.NewPoint(float64(res.GameStatus.Player1.X), float64(res.GameStatus.Player1.Y)))
			}

			if res.GameStatus.Disk != nil && res.GameStatus.Disk.LastUpdate != lastUpdate {
				lastUpdate = res.GameStatus.Disk.LastUpdate

				if ClientDebug {
					fmt.Println("receiving disk update from server")
				}

				force := &vectors.Vector2D{X: float64(res.GameStatus.Disk.Force.X), Y: float64(res.GameStatus.Disk.Force.Y)}

				if res.Token1.PlayerHash != g.token.PlayerHash {
					force = force.Times(-1)
				}

				position := &vectors.Vector2D{X: float64(res.GameStatus.Disk.Position.X), Y: float64(res.GameStatus.Disk.Position.Y)}

				if res.Token1.PlayerHash != g.token.PlayerHash {
					position = &vectors.Vector2D{X: float64(screenWidth) - float64(res.GameStatus.Disk.Position.X), Y: float64(screenHeight) - float64(res.GameStatus.Disk.Position.Y)}
				}

				(ball.Sprite.Hitbox).SetAnchor(figures.NewPoint2(position))
				ball.AddForce(
					force,
					float64(res.GameStatus.Disk.Speed))
			}
		}
	}()

	goo, _ := GetImageFromFilePath("client/graphics/gopher.png")

	bot, right, top, left := contours.Sides()

	radius := math.Max(float64(goo.Bounds().Size().X)/2, float64(goo.Bounds().Size().Y)/2)
	player = Sprite{
		Hitbox: figures.NewCircle(
			figures.NewPoint(float64(screenWidth/2), float64(screenHeight)-radius-25),
			radius,
		),
		Image:                   goo,
		RegisteredIntersections: make(map[figures.Figure]bool),
	}
	opponent = Sprite{
		Image: goo,
		Hitbox: figures.NewCircle(
			figures.NewPoint(float64(screenWidth/2), float64(radius+25)),
			radius,
		),
	}
	ball = PhisicSprite{Sprite: &Sprite{
		Hitbox:                  figures.NewCircle(figures.NewPoint(float64(screenWidth)/2, float64(screenHeight)/2), 15),
		Image:                   goo,
		RegisteredIntersections: make(map[figures.Figure]bool),
	},
		Direction:  &vectors.Vector2D{X: float64(screenWidth) / 2, Y: float64(screenHeight) / 1.3},
		Collisions: &[]figures.Figure{bot, right, top, left, player.Hitbox},
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Airhockey go!")
	constructed = true
	return nil
}
