package main

import (
	"context"
	"flag"

	realtime_router "gsm/router/realtime"

	log "github.com/sirupsen/logrus"
)

var (
	sqlConfigPath string
	port          string
)

func init() {
	flag.StringVar(&sqlConfigPath, "sql", "", "sql config path")
	flag.StringVar(&port, "port", ":8081", "service port")
}

func main() {
	flag.Parse()

	// create gin router for realtime service.
	router, err := realtime_router.NewRouter(context.Background(), realtime_router.RouterConfig{
		SqlConfigPath: sqlConfigPath,
	})
	if err != nil {
		log.Fatalf("Init account router error: %v", err)
	}

	if err := router.Run(port); err != nil {
		log.Fatalf("failed to start account server: %v", err)
	}
}
