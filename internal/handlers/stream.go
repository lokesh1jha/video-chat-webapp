package handlers

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func Stream(c *fiber.Ctx) error {
	suuid := c.Params("ssuid")
	if suuid == "" {
		c.Status(400)
		return nil
	}

	ws := "ws"
	if os.Getenv("ENVIRONMENT") == "production" {
		ws = "wss"
	}
	w.RoomsLock.Lock()
	if _, ok := w.Stream[suuid]; !ok {
		w.RoomsLock.Unlock()
		return c.Render("stream", fiber.Map{
			"StreamWebSocketAddr": fmt.Sprintf("%s://%s/stream/%s", ws, c.Hostname(), suuid),
			"ChatWebSocketAddr":   fmt.Sprintf("%s://%s/stream/%s/chat/websocket", ws, c.Hostname(), suuid),
			"ViewerWebSocketAddr": fmt.Sprintf("%s://%s/stream/%s/viewer/websocket", ws, c.Hostname(), suuid),
			"Type":                "stream",
		}, "layout/main")
	}

	w.RoomsLock.Unlock()

	return c.Render("viewer", fiber.Map{
		"NoStream": true,
		"Leave":    true,
	}, "layout/main")

}

func StreamwebSocket(c *websocket.Conn) error {
	return nil
}

func StreamViewerWebSocket(c *fiber.Ctx) error {
	return c.Render("viewer", fiber.Map{}, "layout/main")
}

func viewerConn(c websocket.Conn, p *w.Peers)
