package main

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

// Colors for CLI output
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
)

type Game struct {
	Board       *Board
	CurrentTurn TurnType
	MoveHistory []string
	Players     []Player
}

func NewGame(size int) *Game {
	return &Game{
		Board:       NewBoard(size),
		CurrentTurn: Player1Turn,
		MoveHistory: []string{},
		Players:     []Player{},
	}
}

func (g *Game) AddPlayer(id bson.ObjectId, name string) {
	player := NewPlayer(id, name)
	g.Players = append(g.Players, player)
}

func (g *Game) InitializeBoard() error {
	if len(g.Players) != 2 {
		return errors.New("Game must have 2 players to initialize the board")
	}

	player1Pieces := []Piece{
		{Type: Pawn, Player: Player1Turn},
		{Type: Hero1, Player: Player1Turn},
		{Type: Hero2, Player: Player1Turn},
		{Type: Pawn, Player: Player1Turn},
		{Type: Pawn, Player: Player1Turn},
	}

	player2Pieces := []Piece{
		{Type: Pawn, Player: Player2Turn},
		{Type: Hero1, Player: Player2Turn},
		{Type: Hero2, Player: Player2Turn},
		{Type: Pawn, Player: Player2Turn},
		{Type: Pawn, Player: Player2Turn},
	}

	return g.Board.InitializeBoard(player1Pieces, player2Pieces)
}

func (g *Game) MovePiece(from, to Position) error {
	piece, exists := g.Board.Pieces[from]

	if !exists {
		return errors.New("No piece at the from position")
	}

	if piece.Player != g.CurrentTurn {
		return errors.New("It is not your turn")
	}

	if !g.IsValidMove(from, to) {
		return errors.New("Invalid move")
	}

	if capturedPiece, exists := g.Board.Pieces[to]; exists {
		g.removePiece(to, capturedPiece.Player)
	}

	delete(g.Board.Pieces, from)

	g.Board.Pieces[to] = piece

	g.updateBitBoard(from, to, piece.Player)

	g.appendMoveToHistory(from, to)

	g.switchTurn()

	return nil
}

func (g *Game) IsValidMove(from, to Position) bool {

  if !g.IsWithinBounds(from) || !g.IsWithinBounds(to) {
    return false
  }
  
  if destinationPiece, exists := g.Board.Pieces[to]; exists && destinationPiece.Player == g.CurrentTurn {
    return false
  }

	piece, exists := g.Board.Pieces[from]
	if exists {
		dx := abs(to.X - from.X)
		dy := abs(to.Y - from.Y)
	  switch piece.Type {
      case Pawn:
          return dx <= 1 && dy <= 1 && (dx + dy == 1)
      case Hero1:
          return (dx == 2 && dy == 0) || (dx == 0 && dy == 2)
      case Hero3:
          return (dx == 2 && dy == 2)
      default:
        return false
    }
  }

  return false
}

func (g * Game) IsWithinBounds(position Position) bool {
  return position.X >= 0 && position.X < g.Board.Size && position.Y >= 0 && position.Y < g.Board.Size
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (g *Game) removePiece(position Position, player TurnType) {
	bit := uint(position.X*g.Board.Size + position.Y) // This is a linear index of the 2D board position
	if player == Player1Turn {
		g.Board.Player1Pieces &= ^(1 << bit)
	} else {
		g.Board.Player2Pieces &= ^(1 << bit)
	}

	/*
	  1. 1 << bit creates a bit mask with a 1 at the 'to' position and 0s elsewhere.
	  2. &^= is a bit clear operation (AND NOT). It clears the bit at the give 'bit' position, effectively removing the piece from its original position in the bitboard.
	*/
}

func (g *Game) updateBitBoard(from, to Position, player TurnType) {
	bitFrom := uint(from.X*g.Board.Size + from.Y)
	bitTo := uint(to.X*g.Board.Size + to.Y)

	if player == Player1Turn {
		g.Board.Player1Pieces &= ^(1 << bitFrom)
		g.Board.Player1Pieces |= 1 << bitTo
	} else {
		g.Board.Player2Pieces &= ^(1 << bitFrom)
		g.Board.Player2Pieces |= 1 << bitTo
	}

	/*
	  1. 1 << toBit creates a bit mask with a 1 at the 'to' position and 0s elsewhere.
	  2. |= is a bitwise OR operation. It sets the bit at the 'to' position, effectively adding the piece to its new position in the bitboard.
	*/
}

func (g *Game) switchTurn() {
	g.CurrentTurn = TurnType((int(g.CurrentTurn) + 1) % 2)
}

func (g *Game) appendMoveToHistory(from, to Position) {
	g.MoveHistory = append(g.MoveHistory, fmt.Sprintf("%v -> %v", from, to))
}

func (g *Game) IsGameOver() (bool, Player) {
	if g.Board.Player1Pieces == 0 {
		return true, g.Players[1]
	}

	if g.Board.Player2Pieces == 0 {
		return true, g.Players[0]
	}

	return false, Player{}
}

func (g *Game) DrawBoard() {
  fmt.Println("\nCurrent Board:")
	fmt.Println("   0  1  2  3  4")  // column numbers
	
	for x := 0; x < g.Board.Size; x++ {
		fmt.Printf("%d ", x)  // row numbers
		for y := 0; y < g.Board.Size; y++ {
			pos := Position{X: x, Y: y}
			if piece, exists := g.Board.Pieces[pos]; exists {
				color := Blue
				if piece.Player == Player2Turn {
					color = Red
				}
				
				pieceStr := "P "
				switch piece.Type {
				case Hero1:
					pieceStr = "H1"
				case Hero2:
					pieceStr = "H2"
				}
				
				fmt.Printf("%s[%s]%s", color, pieceStr, Reset)
			} else {
				fmt.Print("[  ]")
			}
		}
		fmt.Println()
	}
}
