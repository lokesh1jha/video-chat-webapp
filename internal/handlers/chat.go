package handlers

import (
	"videochat/pkg/chat"
	w "videochat/pkg/webrtc"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func RoomChat(c *fiber.Ctx) error {
	return c.Render("chat", fiber.Map{}, "layout/main")
}

func RoomChatWebsocket(c *websocket.Conn) error {
	uuid := c.Params("uuid")
	if uuid == "" {
		return nil
	}

	w.RoomsLock.Lock()
	room := w.Rooms[uuid]
	w.RoomsLock.Unlock()

	if room == nil {
		return nil
	}

	if room.Hub == nil {
		return nil
	}

	chat.PeerChatConn(c.Conn, room.Hub)
}

func StreamChatWebsocket(c *websocket.Conn) error {
	suuid := c.Params("ssuid")

	if suuid == "" {
		return nil
	}

	w.RoomsLock.Lock()

	if stream, ok := w.Streams[suuid]; ok {
		w.RoomsLock.Unlock()

		if stream.Hub == nil {
			hub := chat.NewHub()
			stream.Hub = hub
			go hub.Run()
		}
		chat.PeerChatConn(c.Conn, stream.Hub)
		return nil
	}
	w.RoomsLock.Unlock()

	return nil
}
