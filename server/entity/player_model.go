package entity

import "github.com/gorilla/websocket"

type Player struct {
	PlayerId    string          `json:"playerId"`
	DisplayName string          `json:"displayName"`
	Connection  *websocket.Conn `json:"-"`
}
