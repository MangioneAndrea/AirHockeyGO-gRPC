package main

import (
	"fmt"
	"log"
	"net"

	"github.com/MangioneAndrea/airhockey/positionpb"
	"google.golang.org/grpc"
)

type server struct {
	games []*positionpb.Game
}

func (this *server) CreateGame() positionpb.Game {
	gh := generateId()
	game := positionpb.Game{
		GameHash: gh,
		Token1: &positionpb.Token{
			PlayerHash: generateId(),
			GameHash:   gh,
		},
		Token2: &positionpb.Token{
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

	positionpb.RegisterPositionServiceServer(s, &server{
		games: make([]*positionpb.Game, 0),
	})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	fmt.Println(generateId())
}

func (this *server) RequestGame(v *positionpb.Void, stream positionpb.PositionService_RequestGameServer) error {

	if this.games[len(this.games)-1].P2Ready {
		this.CreateGame()
	}

	stream.Send(this.games[len(this.games)-1])

	return nil
}

func (*server) UpdateStatus(stream positionpb.PositionService_UpdateStatusServer) error {
	msg, err := stream.Recv()
	if err != nil {
		return err
	}
	stream.SendMsg(msg)
	return nil
}
