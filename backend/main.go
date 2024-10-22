package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	// Create a new game
	game := NewGame(5)
	// Add players to the game
	game.AddPlayer(bson.NewObjectId(), "Player Reaulo")
	game.AddPlayer(bson.NewObjectId(), "Player Roefstedt")
	// Initialize the board
	err := game.InitializeBoard()
	if err != nil {
		fmt.Println(err)
	}

  // List of Moves
  // Example moves
	moves := []struct{ from, to Position }{
		{Position{0, 0}, Position{1, 0}}, // Player1 Pawn moves
		{Position{4, 0}, Position{3, 0}}, // Player2 Pawn moves
		{Position{0, 1}, Position{2, 1}}, // Player1 Hero1 moves
	}
  
  game.DrawBoard()
	for _, move := range moves {
		err := game.MovePiece(move.from, move.to)
		if err != nil {
			fmt.Printf("Error making move %v -> %v: %v\n", move.from, move.to, err)
		} else {
			fmt.Printf("Move made: %v -> %v\n", move.from, move.to)
		}

		if isGameOver, player := game.IsGameOver(); isGameOver {
			fmt.Printf("Player %s wins!\n", player.Name)
			break
		}

    game.DrawBoard()
	}
}
