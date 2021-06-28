package main

import (
	"context"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"os"

	"andrea.mangione.dev/airhockey/positionpb"
	"google.golang.org/grpc"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type GameMode int

const (
	SinglePlayer GameMode = iota
	MultiPlayer  GameMode = iota
)

const (
	screenWidth  = 600
	screenHeight = 1200
	wishedFPS    = 60
)

type Game struct {
	mode  GameMode
	token positionpb.Token
}

type Rectangle struct {
	x      int
	y      int
	width  int
	height int
	color  color.Color
}

type Actor struct {
	x        int
	y        int
	width    int
	height   int
	rotation float64
	image    *ebiten.Image
}

var (
	connection positionpb.PositionServiceClient
	ball       Actor
	player1    Actor
	player2    Actor
	divider    = Rectangle{x: 0, y: screenHeight/2 - 2, width: screenWidth, height: 4, color: color.White}
)

func (g *Game) Update() error {
	cursorX, cursorY := ebiten.CursorPosition()
	delta := ebiten.CurrentTPS() / 60
	if delta == 0 {
		return nil
	}
	player1.rotation += 1 / delta
	player1.x = int(math.Min((math.Max(float64(cursorX), 0)), screenWidth))
	player1.y = int(math.Min((math.Max(float64(cursorY), float64(divider.y))), screenHeight))

	connection.UpdateStatus(context.Background(), positionpb.UserInput{
		Vector: &positionpb.Vector2D{X: int32(player1.x), Y: int32(player1.y)},
		Token:  &g.token,
	})
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawActor(screen, &player1)
	drawActor(screen, &player2)
	drawRect(screen, &divider)
}

func drawActor(screen *ebiten.Image, actor *Actor) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(actor.width)/2, -float64(actor.height)/2)
	op.GeoM.Rotate(float64(int(actor.rotation)%360) * 2 * math.Pi / 360)
	op.GeoM.Translate(float64(actor.x), float64(actor.y))
	screen.DrawImage(player1.image, op)
}

func drawRect(screen *ebiten.Image, rect *Rectangle) {
	ebitenutil.DrawRect(screen, float64(rect.x), float64(rect.y), float64(rect.width), float64(rect.height), rect.color)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	connection = positionpb.NewPositionServiceClient(cc)

	ball = Actor{image: getImageFromFilePath("gopher.png")}
	ball.width = ball.image.Bounds().Size().X
	ball.height = ball.image.Bounds().Size().Y

	player1 = Actor{
		image: getImageFromFilePath("gopher.png"),
	}
	player1.width = ball.image.Bounds().Size().X
	player1.height = ball.image.Bounds().Size().Y
	player1.x = screenWidth / 2
	player1.y = screenHeight - player1.height - 25

	player2 = Actor{
		image: getImageFromFilePath("gopher.png"),
	}
	player2.width = ball.image.Bounds().Size().X
	player2.height = ball.image.Bounds().Size().Y
	player2.x = screenWidth / 2
	player2.y = player2.height + 25

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Animation (Ebiten Demo)")
	if err := ebiten.RunGame(&Game{mode: SinglePlayer}); err != nil {
		log.Fatal(err)
	}

}

func getImageFromFilePath(filePath string) *ebiten.Image {
	f, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(image)
}
