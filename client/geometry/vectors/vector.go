package vectors

type vector interface {
	Plus(*vector) *vector
	Minus(*vector) *vector
	Dot(*vector) *vector
	Times(*float64) *vector
	Size() float64
	Normalize() *vector
}
