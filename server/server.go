package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/MangioneAndrea/airhockey/gamepb"
	"google.golang.org/grpc"
)

type server struct {
	games []*gamepb.Game
}

func (this *server) CreateGame() gamepb.Game {
	gh := generateId()
	game := gamepb.Game{
		GameHash: gh,
		Token1: &gamepb.Token{
			PlayerHash: generateId(),
			GameHash:   gh,
		},
		Token2: &gamepb.Token{
			PlayerHash: generateId(),
			GameHash:   gh,
		},
		P1Ready: false,
		P2Ready: false,
	}
	this.games = append(this.games, &game)
	return game
}

func main() {
	fmt.Println("Server running")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	gamepb.RegisterPositionServiceServer(s, &server{
		games: make([]*gamepb.Game, 0),
	})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	fmt.Println(generateId())
}

func (this *server) RequestGame(v *gamepb.Void, stream gamepb.PositionService_RequestGameServer) error {

	if this.games[len(this.games)-1].P2Ready {
		this.CreateGame()
	}

	stream.Send(this.games[len(this.games)-1])

	return nil
}

func (*server) UpdateStatus(stream gamepb.PositionService_UpdateStatusServer) error {
	for {
		msg, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			println(err)
			return err
		}
		fmt.Printf("player position: %v | %v\n", msg.Vector.X, msg.Vector.Y)
		stream.Send(&gamepb.GameStatus{Player1: msg.Vector})
	}
	return nil
}
