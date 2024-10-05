package server

import (
	"log"
	"net/http"

	"example.com/main/data"
	"example.com/main/server/room"
	"github.com/gorilla/websocket"
)

type RoomID = string

type Server struct {
	rooms map[RoomID]*room.Room
}

func NewServer() *Server {
	return &Server{
		rooms: make(map[string]*room.Room),
	}
}

func (s *Server) HandleWSConnection(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	room_id := r.URL.Query().Get("rid")
	uid := r.URL.Query().Get("uid")

	if s.rooms[room_id] == nil {
		s.rooms[room_id] = room.NewRoom(room_id)
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade errors:", err)
	}

	s.rooms[room_id].AddPlayer(data.User{
		Uid: uid,
	},
		int8(0),
		conn,
	)
	readLooper(conn,s.rooms[room_id])
}

func readLooper(conn *websocket.Conn , room *room.Room) {
	defer conn.Close()
	for {
		delta := data.GameDelta{}

		err := conn.ReadJSON(&delta)
		if err != nil {
			log.Println("json error:", err)
		}
		log.Println("message:", delta)
		go room.BroadcastDelta(delta)
		
	}
}
