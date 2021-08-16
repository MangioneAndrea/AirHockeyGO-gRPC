package game

import (
	"image"
	_ "image/png"
	"log"
	"os"

	"github.com/MangioneAndrea/airhockey/client/geometry/figures"
)

type Entity interface {
	Tick(delta int)
	AddForce(*figures.Point, float64)
}

type Clickable interface {
	CheckClicked()
}
type Intersectable interface {
	Intersects(other *Intersectable)
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
