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
	AddForce(*Vector3D)
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

type Circle struct {
	Center *Vector2D
	Radius int
}

func (circle *Circle) Intersects(other *Circle) bool {
	return false
}

func (circle *Circle) Draw(screen *ebiten.Image) {
	previousX := .0
	previousY := .0
	for theta := float64(0); theta < 2*math.Pi; theta += math.Pi * 0.1 {
		x := float64(circle.Center.X) + float64(circle.Radius)*math.Cos(theta)
		y := float64(circle.Center.Y) - float64(circle.Radius)*math.Sin(theta)
		if previousX != 0 {
			ebitenutil.DrawLine(screen, previousX, previousY, x, y, color.White)
		}
		previousX = x
		previousY = y
	}
}

type Sprite struct {
	Hitbox   *Circle
	Speed    float64
	Rotation float64
	Image    *ebiten.Image
}

type PhisicSprite struct {
	Hitbox    *Circle
	Rotation  float64
	Direction *Vector3D
	Image     *ebiten.Image
}

func (sprite *PhisicSprite) AddForce(force *Vector3D) {
	sprite.Direction = force
}
func (sprite *PhisicSprite) Move(where *Vector2D) {
	sprite.Direction.Z = sprite.Hitbox.Center.DistanceTo(where)
	sprite.Hitbox.Center = where
}

func (sprite *PhisicSprite) Draw(screen *ebiten.Image) {
	if ClientDebug {
		sprite.Hitbox.Draw(screen)
	}
	if sprite.Image == nil {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(sprite.Hitbox.Radius)/2, -float64(sprite.Hitbox.Radius)/2)
	op.GeoM.Rotate(float64(int(sprite.Rotation)%360) * 2 * math.Pi / 360)
	op.GeoM.Translate(float64(sprite.Hitbox.Center.X), float64(sprite.Hitbox.Center.Y))
	screen.DrawImage(sprite.Image, op)
}

func (sprite *Sprite) Move(where *Vector2D) {
	sprite.Speed = sprite.Hitbox.Center.DistanceTo(where)
	sprite.Hitbox.Center = where
}

func (sprite *Sprite) Draw(screen *ebiten.Image) {
	if ClientDebug {
		sprite.Hitbox.Draw(screen)
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(sprite.Hitbox.Radius)/2, -float64(sprite.Hitbox.Radius)/2)
	op.GeoM.Rotate(float64(int(sprite.Rotation)%360) * 2 * math.Pi / 360)
	op.GeoM.Translate(float64(sprite.Hitbox.Center.X), float64(sprite.Hitbox.Center.Y))
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
