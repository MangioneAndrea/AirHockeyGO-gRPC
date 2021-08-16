package figures

import (
	"syscall/js"
)

type Figure interface {
	Intersects(other Figure) bool
	Draw(ctx js.Value)
	GetCenter() *Point
	MoveTo(where *Point)
}

// Empty figure which can be compared to nil
type empty struct{}

func (e *empty) Intersects(other Figure) bool { return false }
func (e *empty) Draw(ctx js.Value)            {}
func (e *empty) MoveTo(where *Point)          {}
func (e *empty) GetCenter() *Point            { return nil }

func Empty() Figure { return &empty{} }
