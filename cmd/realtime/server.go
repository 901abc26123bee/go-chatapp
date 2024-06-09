package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

var (
    upgrader = websocket.Upgrader{
        ReadBufferSize:  1024,
        WriteBufferSize: 1024,
        CheckOrigin: func(r *http.Request) bool {
            return true // Allow all origins for simplicity, adjust as needed
        },
    }
)

func main() {
    // Initialize Gin
    r := gin.Default()

    // Routes
    r.GET("/realtime", corsMiddleware(), func(c *gin.Context) {
        serveWs(c.Writer, c.Request)
    })

    // Start server
    r.Run(":8081")
}

func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "DNT,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Range")
        c.Writer.Header().Set("Access-Control-Max-Age", "1728000")
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusOK)
            return
        }
        c.Next()
    }
}

func serveWs(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
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
