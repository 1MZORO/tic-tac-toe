package models

import (
	"github.com/1MZORO/tiktactoe/game"
	"github.com/gorilla/websocket"
)
type Player struct {
    Conn *websocket.Conn
    ID   string 
}

type Room struct {
    ID      string
    PlayerX *Player
    PlayerO *Player
    Game    *game.Game
}

var Rooms = make(map[string]*Room)
