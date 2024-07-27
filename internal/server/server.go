package server

import (
	"flag"
	"os"
	"time"
	"videochat/internal/handlers"
	w "videochat/pkg/webrtc"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/gofiber/websocket/v2"
	// "golang.org/x/net/websocket"
)

var (
	addr = flag.String("addr", ":", os.Getenv("PORT"), "http service address")
	cert = flag.String("cert", "", "certificate file")
	key  = flag.String("key", "", "key file")
)

func Run() error {
	flag.Parse()

	if *addr == ":" {
		*addr = ":8080"
	}

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/", handlers.Welcome)
	app.Get("/room/create", handlers.RoomCreate)
	app.Get("/room/:uuid", handlers.Room)
	app.Get("/room/:uuid/websocket", websocket.new(handlers.RoomWebsocket, websocket.Config{
		HandshakeTimeout: 10 * time.Second,
	}))
	app.Get("/room/:uuid/chat", handlers.RoomChat)
	app.Get("room/:uuid/chat/websocket", websocket.new(handlers.RoomChatWebsocket))
	app.Get("/room/:uuid/viewer/websocket", websocket.new(handlers.RoomViewerWebsocket))

	app.Get("/stream/:ssuid", handlers.Stream)
	app.Get("/stream/:ssuid/websocket", websocket.new(handlers.Streamwebsocket, websocket.Config{
		HandshakeTimeout: 10*time.Second,
	}))
	app.Get("/stream/:ssuid/chat/websocket", websocket.new(handler.StreamChatWebsocket))
	app.Get("/stream/:ssuid/viewer/websocket", websocket.new(handler.StreamViewerWebsocket))
	app.Static("/", "./assets")

	w.Rooms = make(map[string]*w.Room)
	w.Stream = make(map[string]*w.Room)
	go dispatchKeyFrames()

	if *cert != nil {
		return app.ListenTLS(*addr, *cert, *key)
	}
	return app.Listen(*addr)
	

}


fun dispatchKeyFrames() {
	for range time.NewTicker(time.Second*3).C {
		for _, room := range w.Rooms {
			room.Peers.DispatchKeyFrame()
		}
	}
	room.Peers.DispatchKeyFrame()
}
