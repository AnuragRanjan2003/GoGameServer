package data

import (
	"time"

	"example.com/main/server/room"
)

type RoomInfo struct {
	RoomId    string `json:"room_id"`
	Players   uint   `json:"players"`
	Timestamp time.Time `json:"timestamp"`
}

func NewRoomInfo(room room.Room) *RoomInfo {
	return &RoomInfo{
		RoomId: room.GetId(),
		Players: uint(room.CurrentSize()),
		Timestamp: time.Now().UTC(),
	}
}