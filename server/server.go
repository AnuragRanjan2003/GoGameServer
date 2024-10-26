package server

import (
	"context"
	// "encoding/json"
	// "fmt"
	"html/template"
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

		delta := game.GameDelta{}

		err := conn.ReadJSON(&delta)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				s.onDisconnect(uid, room)
			} else {
				log.Println("json error:", err)
			}
			break
		}
		log.Println("message:", delta)
		go room.BroadcastDelta(delta, s.ctx)

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
	status := &info.ServerInfo{
		Rooms:     s.GetRoomsCount(),
		TimeStamp: time.Now().UTC(),
		RoomInfo:  s.GetRoomsList(),
	}
	tmpl, err := template.ParseFiles("./public/index.html")
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, status)
	if err != nil {
		http.Error(w, "Unable to execute template : "+err.Error(), http.StatusInternalServerError)
	}

}

func (s *Server) ServeTemplate(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./public/index.html")
	if err != nil {
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
		return
	}

	info := &info.ServerInfo{
		Rooms:            s.GetRoomsCount(),
		TimeStamp:        time.Now().UTC(),
		RoomInfo:         s.GetRoomsList(),
		TotalConnections: s.GetActiveConnections(),
	}

	err = tmpl.Execute(w, info)
	if err != nil {
		http.Error(w, "Unable to execute template", http.StatusInternalServerError)
	}
}

func (s Server) GetActiveConnections() uint {
	var sum uint = 0

	for _, r := range s.GetRoomsList() {
		sum += r.Players
	}

	return sum
}
