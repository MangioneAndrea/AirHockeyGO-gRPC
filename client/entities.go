package main

import (
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"os"

	"github.com/MangioneAndrea/airhockey/client/geometry/figures"
	"github.com/MangioneAndrea/airhockey/client/geometry/vectors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Entity interface {
	Tick()
	AddForce(*vectors.Vector2D, float64)
}

type Clickable interface {
	CheckClicked()
}
type Intersectable interface {
	Intersects(other *Intersectable)
}
type Rectangle struct {
	Position vectors.Vector2D
	Width    int
	Height   int
	Color    color.Color
}

func (rect *Rectangle) Draw(screen *ebiten.Image) {
	if rect.Color == nil {
		rect.Color = color.White
	}
	ebitenutil.DrawRect(screen, float64(rect.Position.X), float64(rect.Position.Y), float64(rect.Width), float64(rect.Height), rect.Color)
}

type Sprite struct {
	Hitbox   *figures.Circle
	Speed    float64
	Rotation float64
	Image    *ebiten.Image
}

type PhisicSprite struct {
	Sprite         *Sprite
	Direction      *vectors.Vector2D
	LineCollisions *[]*figures.Line2D
}

func (phisicSprite *PhisicSprite) Tick() {
	for _, line := range *phisicSprite.LineCollisions {
		if phisicSprite.Sprite.Hitbox.Intersects(line) {

		}
	}

	phisicSprite.Move(
		phisicSprite.Sprite.Hitbox.Center.Plus(
			phisicSprite.Direction.Times(phisicSprite.Sprite.Speed / 100),
		))
}

func (phisicSprite *PhisicSprite) AddForce(force *vectors.Vector2D, speed float64) {
	phisicSprite.Direction = force
	phisicSprite.Sprite.Speed = speed
}
func (phisicSprite *PhisicSprite) Move(where *vectors.Vector2D) {
	phisicSprite.Sprite.Hitbox.Center = where
}

func (phisicSprite *PhisicSprite) Draw(screen *ebiten.Image) {
	if ClientDebug {
		phisicSprite.Sprite.Hitbox.Draw(screen)
	}
	if phisicSprite.Sprite.Image == nil {
		return
	}
	phisicSprite.Sprite.Draw(screen)
}

func (sprite *Sprite) Move(where *vectors.Vector2D) {
	sprite.Speed = math.Abs(sprite.Hitbox.Center.DistanceTo(where))
	sprite.Hitbox.Center = where
}

func (sprite *Sprite) Draw(screen *ebiten.Image) {
	if ClientDebug {
		sprite.Hitbox.Draw(screen)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(sprite.Image.Bounds().Size().X)/2, -float64(sprite.Image.Bounds().Size().X)/2)
	op.GeoM.Rotate(float64(int(sprite.Rotation)%360) * 2 * math.Pi / 360)
	op.GeoM.Translate(float64(sprite.Hitbox.Center.X), float64(sprite.Hitbox.Center.Y))
	screen.DrawImage(sprite.Image, op)
}

type Button struct {
	Position  vectors.Vector2D
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
		if !button.isClicked && button.Intersects(&figures.Rectangle{Position: vectors.Vector2D{X: x, Y: y}}) {
			button.OnClick()
		}
		button.isClicked = true
	} else {
		button.isClicked = false
	}
}

func (button *Button) Intersects(intersectable Intersectable) bool {
	switch t := intersectable.(type) {
	case Point:
		r := button.Image.Bounds().Add(image.Point{int(button.Position.X), int(button.Position.Y)})
		return point.X >= r.Min.X && point.X <= r.Max.X && point.Y >= r.Min.Y && point.Y <= r.Max.Y
	}
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
