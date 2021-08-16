package main

import (
	"log"
	"syscall/js"
	"time"

	"github.com/MangioneAndrea/airhockey/client/entities"
	"github.com/MangioneAndrea/airhockey/client/game"
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
	scene  entities.Scene
	canvas js.Value
	ctx    js.Value
}

func (g *GUI) GetConnection() gamepb.PositionServiceClient {
	return connection
}

func (g *GUI) Start() {
	g.FitToWindow()
	g.scene.OnConstruction(g)
	for true {
		time.Sleep(time.Millisecond * 10)
		g.Update()
		g.Draw(g.ctx)
	}
}

func (g *GUI) FitToWindow() {
	screenHeight = js.Global().Get("innerHeight").Float()
	screenWidth = screenHeight / 2
	g.canvas.Set("height", screenHeight)
	g.canvas.Set("width", screenWidth)
}

func (g *GUI) Update() error {
	g.scene.Tick(0)
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

func (g *GUI) Draw(ctx js.Value) {
	g.ctx.Call("clearRect", 0, 0, screenWidth, screenHeight)
	g.scene.Draw(ctx)
}

func (g *GUI) GetHeight() float32 {
	return float32(screenHeight)
}
func (g *GUI) GetWidth() float32 {
	return float32(screenWidth)
}
func (g *GUI) GetCanvas() js.Value {
	return g.canvas
}
func (g *GUI) GetCtx() js.Value {
	return g.ctx
}

func (g *GUI) ChangeScene(scene *entities.Scene) {
	g.scene = *scene
	g.scene.OnConstruction(g)
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
		ctx:    js.Global().Get("document").Call("getElementById", "main").Call("getContext", "2d"),
		scene:  &game.MainMenu{},
	}
	go g.Start()

	<-make(chan int)
}
