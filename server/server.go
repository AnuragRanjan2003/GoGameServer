package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	game "example.com/main/data/game"
	info "example.com/main/data/info"
	"example.com/main/server/room"
	"github.com/gorilla/websocket"
)

type RoomID = string

type Server struct {
	rooms map[RoomID]*room.Room
	ctx   context.Context
}

func NewServer(ctx context.Context) *Server {
	return &Server{
		rooms: make(map[string]*room.Room),
		ctx:   ctx,
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

	s.rooms[room_id].AddPlayer(game.User{
		Uid: uid,
	},
		int8(0),
		conn,
	)
	s.readLooper(conn, uid, s.rooms[room_id])
}

func (s *Server) readLooper(conn *websocket.Conn, uid string, room *room.Room) {
	defer conn.Close()
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			delta := game.GameDelta{}

			err := conn.ReadJSON(&delta)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					s.onDisconnect(uid, room)
				}
				log.Println("json error:", err)
				break
			}
			log.Println("message:", delta)
			go room.BroadcastDelta(delta,s.ctx)
		}

	}
}

func (s *Server) closeRoom(roomid string) {
	delete(s.rooms, roomid)
}

func (s *Server) onDisconnect(uid string, room *room.Room) {
	room.RemovePlayer(uid)
	if room.CurrentSize() == 0 {
		s.closeRoom(room.GetId())
	}

}

func (s Server) GetRoomsCount() uint {
	return uint(len(s.rooms))
}

func (s Server) GetRoomsList() []info.RoomInfo {
	var list = []info.RoomInfo{}
	for _, room := range s.rooms {
		list = append(list, *info.NewRoomInfo(*room))
	}

	return list
}

func (s *Server) HandleInfoRequest(w http.ResponseWriter, r *http.Request) {
	info := &info.ServerInfo{
		Rooms:     s.GetRoomsCount(),
		TimeStamp: time.Now().UTC(),
		RoomInfo:  s.GetRoomsList(),
	}

	_json, err := json.Marshal(*info)

	if err != nil {
		return
	}

	w.Write(_json)
}
