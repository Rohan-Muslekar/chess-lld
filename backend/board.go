package main

import (
  "errors"
)

type Position struct { 
  X int
  Y int
}

type Board struct {
  Player1Pieces uint64
  Player2Pieces uint64
  Size int
  Pieces map[Position]Piece
}


func NewBoard (size int) *Board {
  return &Board{
    Player1Pieces: 0,
    Player2Pieces: 0,
    Size: size,
    Pieces: make(map[Position]Piece),
  }
}

func (b *Board) InitializeBoard(player1Pieces, player2Pieces []Piece) error {
  if len(player1Pieces) != len(player2Pieces) {
    return errors.New("player1Pieces and player2Pieces must have the same length")
  }

  if len(player1Pieces) != b.Size || len(player2Pieces) != b.Size {
    return errors.New("player1Pieces and player2Pieces must have the same length as the board size")
  }
  
  for i := 0; i < b.Size; i++ {
    b.Pieces[Position{X: 0, Y: i}] = player1Pieces[i]
    b.Player1Pieces |= 1 << uint(i)
  }

  for i := 0; i < b.Size; i++ {
    b.Pieces[Position{X: b.Size - 1, Y: i}] = player2Pieces[i]
    b.Player2Pieces |= 1 << uint((b.Size - 1) * (b.Size) + i)
  }

  return nil
}
