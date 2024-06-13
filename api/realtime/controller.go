package realtime_service

import (
	"context"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"gsm/pkg/orm"
)

// RealtimeController is the interface for realtime api
type RealtimeController interface {
	HandleWebsocketIO(*gin.Context)
	ServeWebSocket() error
	Echo(*gin.Context)
}

// realtimeController defines the implementation of RealtimeController interface
type realtimeController struct {
	db orm.DB
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for simplicity, adjust as needed
		},
	}
)

// NewRealtimeController creates a new realtime service
func NewRealtimeController(ctx context.Context, db orm.DB) (RealtimeController, error) {
	return &realtimeController{db: db}, nil
}

func (impl *realtimeController) Echo(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		// Read message from client
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// Echo message back to client
		if err := conn.WriteMessage(messageType, p); err != nil {
			return
		}
	}
}

func (impl *realtimeController) HandleWebsocketIO(ctx *gin.Context) {
	handleConnections(ctx.Writer, ctx.Request)
}

func (impl *realtimeController) ServeWebSocket() error {
	return nil
}

// Store connected clients
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var mutex = sync.Mutex{}

// Define a message object
type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}
