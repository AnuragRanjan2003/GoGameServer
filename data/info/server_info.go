package data

import (
	"time"
)

type ServerInfo struct {
	Rooms            uint       `json:"room_nos"`
	RoomInfo         []RoomInfo `json:"room_info"`
	TimeStamp        time.Time  `json:"time_stamp"`
	TotalConnections uint        `json:"total"`
}
