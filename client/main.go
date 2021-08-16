package main

import (
	"log"
	"syscall/js"
	"time"

	"github.com/MangioneAndrea/airhockey/client/entity"
	"github.com/MangioneAndrea/airhockey/gamepb"
	"google.golang.org/grpc"
)

const (
	wishedFPS = 60
)

var (
	screenWidth  = 600.
	screenHeight = 1200.
	connection   gamepb.PositionServiceClient
	ClientDebug  = false
)

type GUI struct {
	scene  entity.Scene
	canvas js.Value
}

func (g *GUI) Start() {
	g.FitToWindow()
	g.scene.OnConstruction()
	for true {
		time.Sleep(time.Millisecond * 10)
		g.Update()
		g.Draw(g.canvas)
	}
}

func (g *GUI) FitToWindow() {
	screenHeight = js.Global().Get("innerHeight").Float()
	screenWidth = screenHeight / 2
	g.canvas.Set("height", screenHeight)
	g.canvas.Set("width", screenWidth)
}

func (g *GUI) Update() error {
	g.scene.Tick()
	/*
		if inpututil.IsKeyJustPressed(ebiten.KeyF6) {
			ClientDebug = !ClientDebug
		}
		delta := ebiten.CurrentTPS() / 60
		if delta == 0 {
			return nil
		}
		if err := g.stage.Tick(); err != nil {
			println(err.Error())
		}
	*/
	return nil
}

func (g *GUI) Draw(canvas js.Value) {

}

func (g *GUI) ChangeScene(scene entity.Scene) {
	g.scene = scene
	//scene.OnConstruction(screenWidth, screenHeight, g)
}

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	connection = gamepb.NewPositionServiceClient(cc)

	g := &GUI{
		canvas: js.Global().Get("document").Call("getElementById", "main"),
	}
	/*
		js.Global().Get("document").Call("addEventListener", "unload", js.New(func(args []js.Value) {
			g.alive = false
		}))
	*/
	go g.Start()

	<-make(chan int)
}
