package main

import (
	"context"
	"fmt"
	"github.com/MangioneAndrea/gonsole"
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
			time.Sleep(5 * time.Second)
			for key, element := range server.games {
				printDebug("Time elapsed %v \n", time.Since(time.Unix((element.LastUpdate), 0)))
				if time.Since(time.Unix((element.LastUpdate), 0)) > 5*time.Second {
					gonsole.Error(key, "Removing game due to inactivity")
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
		gonsole.Success(game.Token2.GameHash, "Joining game", gonsole.ShowIfNotNil)
		return game.Token2, nil
	} else {
		token := server.CreateGame().Token1
		gonsole.Success(token.GameHash, "Creating game", gonsole.ShowIfNotNil)
		return token, nil
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
			if msg.Vector != nil {
				if msg.Token.PlayerHash == game.Token1.PlayerHash {
					game.GameStatus.Player1.X = 600 - msg.Vector.X
					game.GameStatus.Player1.Y = 1200 - msg.Vector.Y
				} else if msg.Token.PlayerHash == game.Token2.PlayerHash {
					game.GameStatus.Player2.X = 600 - msg.Vector.X
					game.GameStatus.Player2.Y = 1200 - msg.Vector.Y
				}
			}

			if msg.DiskStatus != nil {
				game.GameStatus.Disk = msg.DiskStatus
				if msg.Token.PlayerHash == game.Token2.PlayerHash {
					game.GameStatus.Disk.Force.Y = -game.GameStatus.Disk.Force.Y
					game.GameStatus.Disk.Force.X = -game.GameStatus.Disk.Force.X

					game.GameStatus.Disk.Position.X = 600 - game.GameStatus.Disk.Position.X
					game.GameStatus.Disk.Position.Y = 1200 - game.GameStatus.Disk.Position.Y
				}
			}
			err := stream.Send(game)
			gonsole.Error(err, "Stream.send", gonsole.ShowIfNotNil)
		}
	}
	return nil
}
