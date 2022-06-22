package main

import (
	"fmt"
	"image"
	_ "image/png"
	"log"
	"math"
	"os"

	"github.com/MangioneAndrea/airhockey/client/geometry/figures"
	"github.com/MangioneAndrea/airhockey/client/geometry/vectors"
	"github.com/hajimehoshi/ebiten/v2"
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

type Sprite struct {
	Hitbox                  *figures.Circle
	Speed                   float64
	Rotation                float64
	Image                   *ebiten.Image
	RegisteredIntersections map[figures.Figure]bool
}

type PhisicSprite struct {
	Sprite     *Sprite
	Direction  *vectors.Vector2D
	Collisions *[]figures.Figure
}

func (phisicSprite *PhisicSprite) Tick() {
	phisicSprite.Sprite.Speed *= 0.995
	if phisicSprite.Sprite.Speed > 0 {
		phisicSprite.Sprite.Speed -= 0.01
	} else {
		phisicSprite.Sprite.Speed += 0.01
	}
	for index, item := range *phisicSprite.Collisions {
		if _, firstIntersection := phisicSprite.Sprite.Intersects(item); firstIntersection {

			switch coll := (item).(type) {
			case *figures.Circle:
			case *figures.Segment:
				slope := math.Abs(float64(int(coll.Slope()) % 1000))
				x := slope / 1000
				y := 1 - math.Min(slope, 1)
				phisicSprite.Bounce(x, y)
				if ClientDebug {
					fmt.Printf("first collision with segment at index %d with slope %f \n", index, slope)
				}
			case *figures.Line:
				slope := math.Abs(float64(int(coll.Slope()) % 1000))
				x := slope / 1000
				y := 1 - math.Min(slope, 1)
				phisicSprite.Bounce(x, y)
				if ClientDebug {
					fmt.Printf("first collision with line at index %d with slope %f \n", index, slope)
				}
			}

			//phisicSprite.AddForce(phisicSprite.Sprite.Hitbox.Center.Vector.Minus(player.Hitbox.Center.Vector), player.Speed)

		}
	}

	phisicSprite.Move(
		phisicSprite.Sprite.Hitbox.Center.Vector.Plus(
			phisicSprite.Direction.Times(phisicSprite.Sprite.Speed / 100),
		))
}

func (phisicSprite *PhisicSprite) Bounce(x, y float64) {
	if x > 0.05 {
		phisicSprite.Direction.X *= -x
	}
	if y > 0.05 {
		phisicSprite.Direction.Y *= -y
	}
}

func (phisicSprite *PhisicSprite) AddForce(force *vectors.Vector2D, speed float64) {
	phisicSprite.Direction = force
	phisicSprite.Sprite.Speed = speed
}
func (phisicSprite *PhisicSprite) Move(where *vectors.Vector2D) {
	phisicSprite.Sprite.Hitbox.Center = figures.NewPoint2(where)
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

func (sprite *Sprite) Intersects(with figures.Figure) (intersects, isFirstIntersection bool) {
	intersects = sprite.Hitbox.Intersects(with)
	if intersects && !sprite.isIntersectionRegistered(with) {
		if ClientDebug {
			fmt.Println("registering intersection")
		}
		sprite.RegisteredIntersections[with] = true
		isFirstIntersection = true
	}
	if !intersects && sprite.isIntersectionRegistered(with) {
		if ClientDebug {
			fmt.Println("unregistering intersection")
		}
		sprite.RegisteredIntersections[with] = false
	}

	return intersects, isFirstIntersection
}

func (sprite *Sprite) isIntersectionRegistered(with figures.Figure) bool {
	if intersects, ok := sprite.RegisteredIntersections[with]; ok {
		return ok && intersects
	}
	return false
}

func (sprite *Sprite) Move(where *vectors.Vector2D) {
	sprite.Speed = math.Abs(sprite.Hitbox.Center.Vector.DistanceTo(where))
	sprite.Hitbox.Center = figures.NewPoint2(where)
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
		if !button.isClicked &&
			figures.NewRectangle(
				figures.NewPoint2(&button.Position),
				float64(button.Image.Bounds().Dx()),
				float64(button.Image.Bounds().Dy())).Intersects(figures.NewPoint(float64(x), float64(y))) {
			button.OnClick()
		}
		button.isClicked = true
	} else {
		button.isClicked = false
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
