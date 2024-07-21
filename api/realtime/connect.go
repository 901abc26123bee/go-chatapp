package realtime

import (
	"gsm/pkg/cache"
	"net/http"
	// "github.com/gorilla/websocket"
	// log "github.com/sirupsen/logrus"
)

// connectService defines the implementation of ConnectService interface
type connectService struct {
	redisClient cache.Client
}

// ConnectService defines the connect service interface
type ConnectService interface {
	HandleWebsocketIO(w http.ResponseWriter, r *http.Request)
}

// NewConnectService init the connect service
func NewConnectService(redisClient cache.Client) ConnectService {
	return &connectService{
		redisClient: redisClient,
	}
}

// var (
// 	upgrader = websocket.Upgrader{
// 		ReadBufferSize:  1024,
// 		WriteBufferSize: 1024,
// 		CheckOrigin: func(r *http.Request) bool {
// 			return true // Allow all origins for simplicity, adjust as needed
// 		},
// 	}
// )

func (impl *connectService) HandleWebsocketIO(w http.ResponseWriter, r *http.Request) {

}
