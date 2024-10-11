package room

import (
	"context"
	"log"

	game "example.com/main/data/game"
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

func (r *Room) AddPlayer(user game.User, role Role, conn *websocket.Conn) {
	r.roles[user.Uid] = role
	r.conns[user.Uid] = conn
}

func (r Room) BroadcastDelta(delta types.Delta, ctx context.Context) {
	// uid := delta.GetProducer()

	for _, conn := range r.conns {
		select {
		case <-ctx.Done():
			return
		default:
			conn.WriteJSON(delta)
		}
	}
}
func (r *Room) RemovePlayer(uid string) {
	delete(r.roles, uid)
	delete(r.conns, uid)
}

func (r Room) CurrentSize() int {
	return len(r.conns)
}

func (r Room) GetId() string {
	return r.roomid
}
