package realtime_service

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// RealtimeController is the interface for realtime api
type RealtimeController interface {
	HandleWebsocketIO(*gin.Context)
	Echo(*gin.Context)
	ServeWebSocket() error
}

// realtimeController defines the implementation of RealtimeController interface
type realtimeController struct {
}

// NewRealtimeController creates a new realtime service
func NewRealtimeController(ctx context.Context) (RealtimeController, error) {
	return &realtimeController{}, nil
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

func(impl *realtimeController) Echo(ctx *gin.Context) {
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

func(impl *realtimeController) HandleWebsocketIO(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
			return
	}
	defer conn.Close()

	for {

	}
}

func(impl *realtimeController) ServeWebSocket() error {
	return nil
}