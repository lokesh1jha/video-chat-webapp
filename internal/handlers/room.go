package handlers

import (
	"fmt"
	"os"
	w "videochat/pkg/webrtc"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	guuid "github.com/google/uuid"
)

type websocketMessage struct {
	Event string `json:"event"`
	Data  string `josn: "data"`
}

func RoomCreate(c *fiber.Ctx) error {
	return c.Redirect(fmt.Sprintf("/room/%s", guuid.New().String()))
}

func Room(c *fiber.Ctx) error {
	uuid := c.Params("uuid")

	if uuid == "" {
		c.Status(400)
		return nil
	}
	ws := "ws"
	if os.Getenv("ENVIRONMENT") == "production" {
		ws = "wss"
	}

	uuid, suuid, _ := createOrGetRoom(uuid)
	return c.Render("peer", fiber.Map{
		"RoomWebSocketAddr": fmt.Sprintf("%s://%s/room/%s/websocket", ws, c.Hostname(), uuid),
		"RoomLink":          fmt.Sprintf("%s://%s/room/%s", c.Protocol(), , c.Hostname(), uuid),
		"ChatWebsocketAddr": fmt.Sprintf("%s://%s/room/%s/chat/websocket", ws, c.Hostname(), uuid),
		"ViewerWebsocket":   fmt.Sprintf("%s://%s/room/%s/viewer/websocket", ws, c.Hostname(), uuid),
		"StreamLink":        fmt.Sprintf("%s://%s/stream/%s", ws, c.Hostname(), suuid),
		"Type":              "room",
	}, "layout/main")
}

func RoomWebsocket(c *websocket.Conn) error {

	uuid := c.Params("uuid")
	if uuid == "" {
		return nil
	}

	_, _, room := createOrGetRoom(uuid)
	w.RoomConn(c, room.Peer)
	return room
}

func createOrGetRoom(uuid string) (string, string, Room) {
	w.RoomsLock.Lock()
	
	defer w.RoomsLock.Unlock()

	h := sha256.new()
	h.write([]byte(uuid))
	suuid := fmt.Sprintf("%x", h.sum(nil))

	if room, ok := w.Rooms[uuid]; ok {
		if _, ok := w.Stream[suuid]; !ok {
			w.Stream[suuid] = room
		}
		return uuid, suuid, room
	}

	hub := chat.NewHub()
	peer := &w.Peer{}
	p.TrackLocals = make(map[string]*webrtc.TrackLocalStaticRTP)
	room := &w.Room{
		Peers: p,
		Hub: hub,
	}
	w.Rooms[uuid] = room
	w.Stream[uuid] = room

	go hub.Run()
	return uuid, suuid, room
	
}

func RoomViewerWebsocket(c *websocket.Conn) {
	uuid := c.Params("uuid")
	if uuid == "" {
		return
	}

	w.RoomsLock.Lock()
	if peer, ok := w.Rooms[uuid]; ok {
		w.RoomsLock.Unlock()
		roomViewerConn(c, peer)
		return
	}

	w.RoomsLock.Unlock()
}

func roomViewerConn(c *websocket.Conn, p *w.Peer) {
	uuid := c.Params("uuid")
	if uuid == "" {
		return
	}

	ticker := time.NewTicker(1 * time.Second)
	defer func () {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
			case <-ticker.C:
				w, err := c.Conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}
				w.Write([]byte(fmt.Sprintf("%d", len(p.Connection))))
		}
	}
}


func StreamwebSocket(c *websocket.Conn) error {
	suuid := c.Params("ssuid")
	if suuid == "" {
		return
	}
	w.StreamLock.Lock()
	if stream, ok := w.Stream[suuid]; ok {
		w.StreamLock.Unlock()
		StreamConn(c, stream)
		return
	}
	w.StreamLock.Unlock()
}

func StreamViewerWebSocket(c *fiber.Ctx) error {
	suuid := c.Params("ssuid")
	if suuid == "" {
		return
	}
	w.StreamLock.Lock()
	if stream, ok := w.Stream[suuid]; ok {
		w.StreamLock.Unlock()
		viewerConn(c, stream)
		return
	}
	w.StreamLock.Unlock()
	return nil
}

func viewerConn(c websocket.Conn, p *w.Peers){
	ticker := time.NewTicker(1 * time.Second)
	defer func () {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
			case <-ticker.C:
				w, err := c.Conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}
				w.Write([]byte(fmt.Sprintf("%d", len(p.Connection))))
		}
	}
}

