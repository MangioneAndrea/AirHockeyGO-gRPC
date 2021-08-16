package figures

import "syscall/js"

type Figure interface {
	Intersects(other Figure) bool
	Draw(ctx js.Value)
}
