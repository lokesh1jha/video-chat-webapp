package handlers

import (
	"github.com/gofiber/fiber"
	"github.com/gofiber/websocket"
)

func Stream(c *fiber.Ctx) error {
	return c.Render("stream", fiber.Map{}, "layout/main")
}

func Streamwebsocket(c *websocket.Conn) error {
	return nil
}
