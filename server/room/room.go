package room

import (
	"log"
	"example.com/main/data"
	"example.com/main/types"
	"github.com/gorilla/websocket"
)

type RoomID = string
type UserID = string
type Role = int8

type Room struct {
	roomid RoomID
	roles  map[UserID]Role
	conns  map[UserID]*websocket.Conn
}

func NewRoom(rid string) *Room {
	log.Println("creating new room :", rid)
	return &Room{
		roles:  make(map[UserID]Role),
		conns:  make(map[UserID]*websocket.Conn),
		roomid: rid,
	}
}

func (r *Room) AddPlayer(user data.User, role Role, conn *websocket.Conn) {
	r.roles[user.Uid] = role
	r.conns[user.Uid] = conn
}

func (r Room) BroadcastDelta(delta types.Delta){
	// uid := delta.GetProducer()

	for _,conn := range r.conns {
		conn.WriteJSON(delta)
	}
}


