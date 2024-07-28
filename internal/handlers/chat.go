package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"videochat/pkg/chat"
	w "videochat/pkg/webrtc"
)
fun RoomChat(c *fiber.Ctx) error {
	return c.render("chat", fiber.Map{}, "layout/main")
}


func RoomChatWebsocket(c *websocket.Conn) error {
	uuid := c.Params("uuid")
	if uuid == "" {
		return
	}

	w.RoomsLock.Lock()
	room := w.Rooms[uuid]
	w.RoomsLock.Unlock()

	if room == nil {
		return
	}

	if room.Hub == nil {
		return
	}

	chat.PeerChatConn(c.Conn, room.Hub)
}

func StreamChatWebsocket(c *websocket.Conn) error {
	suuid := c.Params("ssuid")

	if suuid == "" {
		return
	}

	w.RommsLock.Lock()

	if stream, ok := w.Streams[suuid]; ok {
		w.RommsLock.Unlock()
		
		if stream.Hub == nil {
			hub := chat.NewHub()
			stream.Hub = hub
			go hub.Run()
		}
		chat.PeerChatConn(c.Conn, stream.Hub)
		return 
	}
	w.RoomsLock.Unlock()
}


