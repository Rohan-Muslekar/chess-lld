package main

type PieceType int

const (
  Pawn PieceType = iota
  Hero1
  Hero2
  Hero3
)

type Piece struct {
  Type PieceType
  Player TurnType
}
