package chat

import "github.com/gofiber/websocket/v2"

const (
	writeWait
)

type Client struct {
	Hub        *Hub
	Connection *websocket.Conn
	Send       chan []byte
}
