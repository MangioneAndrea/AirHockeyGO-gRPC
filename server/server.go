package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/MangioneAndrea/airhockey/gamepb"
	"google.golang.org/grpc"
)

type Server struct {
	games      map[string]*gamepb.Game
	emptyGames []*gamepb.Game
}

func (server *Server) CreateGame() *gamepb.Game {
	gh := generateId()
	game := gamepb.Game{
		GameHash: gh,
		Token1: &gamepb.Token{
			PlayerHash: generateId(),
			GameHash:   gh,
		},
		GameStatus: &gamepb.GameStatus{
			Player1: &gamepb.Vector2D{X: 0, Y: 0},
			Player2: &gamepb.Vector2D{X: 0, Y: 0},
		},
	}
	server.emptyGames = append(server.emptyGames, &game)
	server.games[gh] = &game
	return &game
}

func main() {
	fmt.Println("Server running")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	gamepb.RegisterPositionServiceServer(s, &Server{
		games: make(map[string]*gamepb.Game),
	})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
}

func (server *Server) RequestGame(ctx context.Context, v *gamepb.GameRequest) (*gamepb.Token, error) {
	if len(server.emptyGames) > 0 {
		game := server.emptyGames[0]
		server.emptyGames = server.emptyGames[1:]
		game.Token2 = &gamepb.Token{
			PlayerHash: generateId(),
			GameHash:   game.GameHash,
		}
		return game.Token2, nil
	} else {
		return server.CreateGame().Token1, nil
	}
}

func (server *Server) UpdateStatus(stream gamepb.PositionService_UpdateStatusServer) error {
	for {
		msg, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if msg.Token == nil {
			continue
		}
		game := server.games[msg.Token.GameHash]
		if game != nil {
			if msg.Token.PlayerHash == game.Token1.PlayerHash {
				game.GameStatus.Player1.X = 600 - msg.Vector.X
				game.GameStatus.Player1.Y = 1200 - msg.Vector.Y
			} else if msg.Token.PlayerHash == game.Token2.PlayerHash {
				game.GameStatus.Player2.X = 600 - msg.Vector.X
				game.GameStatus.Player2.Y = 1200 - msg.Vector.Y
			}
			stream.Send(game)
		}
	}
	return nil
}
