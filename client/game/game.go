package game

import (
	"math"
	"syscall/js"

	"github.com/MangioneAndrea/airhockey/client/entities"
	"github.com/MangioneAndrea/airhockey/client/entities/actors"
	"github.com/MangioneAndrea/airhockey/client/geometry/figures"
	"github.com/MangioneAndrea/airhockey/gamepb"
)

type GameMode int

const (
	SinglePlayer GameMode = iota
	MultiPlayer  GameMode = iota
)

var (
	constructed  bool = false
	ball         *actors.Sprite
	player       *actors.Sprite
	opponent     *actors.Sprite
	divider      *actors.Sprite // *figures.Rectangle
	contours     *actors.Sprite //*figures.Rectangle
	updateStatus gamepb.PositionService_UpdateStatusClient
)

type Game struct {
	token           *gamepb.Token
	height, width   float32
	actors          []entities.Actor
	sceneController entities.SceneController
}

func (g *Game) GetActors() *[]entities.Actor {
	return &g.actors
}

func (g *Game) Tick(delta int) {
	if !constructed {
		return
	}
	cursorX, cursorY := g.sceneController.GetMousePosition().Values()
	player.Move(&figures.Point{
		X: math.Min((math.Max(float64(cursorX), 0)), float64(g.width)),
		Y: math.Min((math.Max(float64(cursorY), 0)), float64(g.height)),
	})
	/*
		err := updateStatus.Send(&gamepb.UserInput{
			Vector: &gamepb.Vector2D{X: int32(player.Hitbox.GetCenter().X), Y: int32(player.Hitbox.GetCenter().Y)},
			Token:  g.token,
		})
		if err != nil {
			fmt.Printf("Error while sending %v\n", err)
		}
	*/
	if player.Hitbox.Intersects(ball.GetHitbox()) {
		//ball.AddForce(ball.GetHitbox().GetCenter().Vector.Minus(player.Hitbox.GetCenter().Vector), player.Speed)
	}

	ball.Tick(delta)
}

func (g *Game) Draw(ctx js.Value) {
	if /*ClientDebug*/ false {
		s := player.Hitbox.GetCenter().LineTo(ball.GetHitbox().GetCenter()).SnapSegment(contours.Hitbox.(*figures.Rectangle))

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
	for _, actor := range *g.GetActors() {
		actor.Draw(ctx)
	}
}

func (g *Game) OnConstruction(s entities.SceneController) {
	g.sceneController = s
	g.height = s.GetHeight()
	g.width = s.GetWidth()

	divider = actors.NewSprite(
		figures.NewRectangle(figures.NewPoint(0, float64(g.height)/2-2), float64(g.width), 4),
		nil,
		false,
	)

	contours = actors.NewSprite(
		figures.NewRectangle(figures.NewPoint(0, 0), float64(g.width)-2, float64(g.height)-2),
		nil,
		true,
	) /*

		stream, streamErr := s.GetConnection().UpdateStatus(context.Background())
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
					opponent.Hitbox.GetCenter().X = float64(res.GameStatus.Player2.X)
					opponent.Hitbox.GetCenter().Y = float64(res.GameStatus.Player2.Y)
				} else {
					opponent.Hitbox.GetCenter().X = float64(res.GameStatus.Player1.X)
					opponent.Hitbox.GetCenter().Y = float64(res.GameStatus.Player1.Y)
				}
			}
		}()
	*/

	//goo, _ := GetImageFromFilePath("client/graphics/gopher.png")

	ball = actors.NewSprite(
		figures.NewCircle(figures.NewPoint(float64(s.GetWidth())/2, float64(s.GetHeight())/1.3), 15),
		nil,
		true,
	)
	radius := 20. //math.Max(float64(goo.Bounds().Size().X)/2, float64(goo.Bounds().Size().Y)/2)
	player = actors.NewSprite(
		figures.NewCircle(figures.NewPoint(float64(s.GetWidth()/2), float64(s.GetHeight())-radius-25), radius),
		nil,
		true,
	)
	opponent = actors.NewSprite(
		figures.NewCircle(figures.NewPoint(float64(s.GetWidth()/2), float64(radius+25)), radius),
		nil,
		false,
	)

	g.actors = []entities.Actor{player, ball, opponent, divider, contours}
	constructed = true
}
