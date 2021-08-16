package actors

import (
	"math"
	"syscall/js"

	"github.com/MangioneAndrea/airhockey/client/entities"
	"github.com/MangioneAndrea/airhockey/client/geometry/figures"
)

type Sprite struct {
	Hitbox     figures.Figure
	Speed      float64
	Rotation   float64
	Skin       figures.Figure
	collisions bool
}

func NewSprite(hitbox figures.Figure, skin figures.Figure, collisions bool) *Sprite {
	return &Sprite{
		Hitbox:     hitbox,
		Skin:       skin,
		collisions: collisions,
	}
}

func (sprite *Sprite) Move(where *figures.Point) {
	sprite.Speed = math.Abs(sprite.GetHitbox().GetCenter().DistanceTo(where))
	sprite.Hitbox.MoveTo(where)
}

func (sprite *Sprite) OnConstruction(s entities.SceneController) {

}
func (sprite *Sprite) Tick(delta int) {

}

func (sprite *Sprite) Draw(ctx js.Value) {
	if sprite.Skin != nil {
		sprite.Skin.Draw(ctx)
	}
	if sprite.Hitbox != nil {
		sprite.Hitbox.Draw(ctx)
	}
}

func (sprite *Sprite) HasCollisions() bool {
	return sprite.collisions
}
func (sprite *Sprite) GetHitbox() figures.Figure {
	return sprite.Hitbox
}
