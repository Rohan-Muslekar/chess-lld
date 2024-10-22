package main

// import bson object id generator
import (
  "gopkg.in/mgo.v2/bson"
)


type TurnType int

const (
  Player1Turn TurnType = iota
  Player2Turn
)

type Player struct {
	Id bson.ObjectId `bson:"_id"`
  Name string
  TurnId TurnType
}

func NewPlayer(id bson.ObjectId, name string) Player {
  return Player{
    Id: id,
    Name: name,
    TurnId: Player1Turn,
  }
}
