package game

import (
	"image"
	_ "image/png"
	"log"
	"math"
	"os"
	"syscall/js"

	"github.com/MangioneAndrea/airhockey/client/geometry/figures"
	"github.com/MangioneAndrea/airhockey/client/geometry/vectors"
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
	Hitbox   *figures.Circle
	Speed    float64
	Rotation float64
	Image    image.Image
}

type PhisicSprite struct {
	Sprite     *Sprite
	Direction  *vectors.Vector2D
	Collisions *[]figures.Figure
}

func (phisicSprite *PhisicSprite) Tick() {
	for _, item := range *phisicSprite.Collisions {
		if phisicSprite.Sprite.Hitbox.Intersects(item) {
			//fmt.Println("collision with wall")
		}
	}

	phisicSprite.Move(
		phisicSprite.Sprite.Hitbox.Center.Vector.Plus(
			phisicSprite.Direction.Times(phisicSprite.Sprite.Speed / 100),
		))
}

func (phisicSprite *PhisicSprite) AddForce(force *vectors.Vector2D, speed float64) {
	phisicSprite.Direction = force
	phisicSprite.Sprite.Speed = speed
}
func (phisicSprite *PhisicSprite) Move(where *vectors.Vector2D) {
	phisicSprite.Sprite.Hitbox.Center = figures.NewPoint2(where)
}

func (phisicSprite *PhisicSprite) Draw() { /*
		if ClientDebug {
			phisicSprite.Sprite.Hitbox.Draw()
		}
		if phisicSprite.Sprite.Image == nil {
			return
		}*/
	//phisicSprite.Sprite.Draw()
}

func (sprite *Sprite) Move(where *vectors.Vector2D) {
	sprite.Speed = math.Abs(sprite.Hitbox.Center.Vector.DistanceTo(where))
	sprite.Hitbox.Center = figures.NewPoint2(where)
}

func (sprite *Sprite) Draw(ctx js.Value) { /*
		if ClientDebug {
			//sprite.Hitbox.Draw(screen)
		}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(-float64(sprite.Image.Bounds().Size().X)/2, -float64(sprite.Image.Bounds().Size().X)/2)
			op.GeoM.Rotate(float64(int(sprite.Rotation)%360) * 2 * math.Pi / 360)
			op.GeoM.Translate(float64(sprite.Hitbox.Center.X), float64(sprite.Hitbox.Center.Y))*/
	//screen.DrawImage(sprite.Image, op)
}

func GetImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	i, _, err := image.Decode(f)
	if err != nil || i == nil {
		log.Fatalf("Image at %v could not be loaded %v", filePath, err)
	}
	return i, nil
}
