package figures

type Figure interface {
	Intersects(other Figure) bool
}
