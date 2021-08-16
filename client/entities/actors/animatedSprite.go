package actors

import (
	"syscall/js"

	"github.com/MangioneAndrea/airhockey/client/entities"
	"github.com/MangioneAndrea/airhockey/client/geometry/figures"
)

type AnimatedSprite struct {
	Sprite    *Sprite
	Direction *figures.Line
	Force     float64
}

func NewAnimatedSprite(hitbox figures.Figure, skin figures.Figure, collisions bool) *AnimatedSprite {
	return &AnimatedSprite{
		Sprite: NewSprite(hitbox, skin, collisions),
	}
}

func (as *AnimatedSprite) Move(where *figures.Point) {
	as.Sprite.Move(where)
}

func (as *AnimatedSprite) OnConstruction(s entities.SceneController) {
	as.Sprite.OnConstruction(s)
}
func (as *AnimatedSprite) Tick(delta int) {
	if as.Force != 0 {
		// TODO: Walk along line towards the direction or the origin
	}
	as.Sprite.Tick(delta)
}

func (as *AnimatedSprite) AddForce(d *figures.Point, f float64) {
	as.Direction = figures.NewLine(
		as.Direction.Start.Avg(as.Sprite.Hitbox.GetCenter()),
		as.Direction.Direction.Avg(d),
	)
	as.Force = as.Force + f
}

func (as *AnimatedSprite) Draw(ctx js.Value) {
	as.Sprite.Draw(ctx)
}

func (as *AnimatedSprite) HasCollisions() bool {
	return as.Sprite.HasCollisions()
}
func (as *AnimatedSprite) GetHitbox() figures.Figure {
	return as.Sprite.GetHitbox()
}
