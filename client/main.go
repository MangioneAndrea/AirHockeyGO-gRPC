package main

import (
	"image"
	_ "image/png"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

var (
	gopher *ebiten.Image
)

type Game struct {
	count        int
	goopherSpeed float32
}

func (g *Game) Update() error {
	g.count++
	g.goopherSpeed += 0.1
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	w, h := gopher.Size()
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(float64(g.count*2%360) * 2 * math.Pi / 360)
	op.GeoM.Translate(screenWidth/2, screenHeight/2)
	screen.DrawImage(gopher, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	gopher = getImageFromFilePath("gopher.png")

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Animation (Ebiten Demo)")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}

}

func getImageFromFilePath(filePath string) *ebiten.Image {
	f, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(image)
}
