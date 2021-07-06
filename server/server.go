package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

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
		LastUpdate: (time.Now().Unix()),
	}
	server.emptyGames = append(server.emptyGames, &game)
	return &game
}

func main() {
	fmt.Println("Server running")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	server := &Server{
		games: make(map[string]*gamepb.Game),
	}

	go func() {
		for {
			time.Sleep(30 * time.Second)
			for key, element := range server.games {
				printDebug("Time elapsed %v \n", time.Since(time.Unix((element.LastUpdate), 0)))
				if time.Since(time.Unix((element.LastUpdate), 0)) > 5*time.Second {
					delete(server.games, key)
				}
			}
			for index, element := range server.emptyGames {
				printDebug("Time elapsed %v \n", time.Since(time.Unix((element.LastUpdate), 0)))
				if time.Since(time.Unix((element.LastUpdate), 0)) > 20*time.Second {
					server.emptyGames = append(server.emptyGames[:index], server.emptyGames[index+1:]...)
				}
			}
			printDebug("Games cleanup -- games alive: %v | emptyGames alive: %v \n", len(server.games), len(server.emptyGames))
		}
	}()

	gamepb.RegisterPositionServiceServer(s, server)
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
		server.games[game.GameHash] = game
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
		// If the token is not valid, ignore the request
		if msg.Token == nil {
			continue
		}
		game := server.games[msg.Token.GameHash]
		if game != nil {
			game.LastUpdate = time.Now().Unix()
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
