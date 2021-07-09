package main

import (
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Entity interface {
	Draw()
}

type Clickable interface {
	CheckClicked()
}
type Intersectable interface {
	ContainsPoint(point image.Point)
}

type Rectangle struct {
	Position Vector2D
	Width    int
	Height   int
	Color    color.Color
}

func (rect *Rectangle) Draw(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, float64(rect.Position.X), float64(rect.Position.Y), float64(rect.Width), float64(rect.Height), rect.Color)
}

type Sprite struct {
	Position *Vector2D
	Speed    float64
	Width    int
	Height   int
	Rotation float64
	Image    *ebiten.Image
}

func (sprite *Sprite) Move(where *Vector2D) {
	sprite.Speed = sprite.Position.DistanceTo(where)
	sprite.Position = where
}

func (sprite *Sprite) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(sprite.Width)/2, -float64(sprite.Height)/2)
	op.GeoM.Rotate(float64(int(sprite.Rotation)%360) * 2 * math.Pi / 360)
	op.GeoM.Translate(float64(sprite.Position.X), float64(sprite.Position.Y))
	screen.DrawImage(sprite.Image, op)
}

type Button struct {
	Position  Vector2D
	Image     *ebiten.Image
	OnClick   func()
	isClicked bool
}

func (button *Button) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(button.Position.X), float64(button.Position.Y))
	screen.DrawImage(button.Image, op)
}

func (button *Button) CheckClicked() {
	if button.OnClick == nil {
		return
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if !button.isClicked && button.ContainsPoint(image.Point{X: x, Y: y}) {
			button.OnClick()
		}
		button.isClicked = true
	} else {
		button.isClicked = false
	}
}

func (button *Button) ContainsPoint(point image.Point) bool {
	r := button.Image.Bounds().Add(image.Point{int(button.Position.X), int(button.Position.Y)})
	return point.X >= r.Min.X && point.X <= r.Max.X && point.Y >= r.Min.Y && point.Y <= r.Max.Y
}

func GetImageFromFilePath(filePath string) (*ebiten.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	if err != nil || image == nil {
		log.Fatalf("Image at %v could not be loaded %v", filePath, err)
	}
	return ebiten.NewImageFromImage(image), nil
}
