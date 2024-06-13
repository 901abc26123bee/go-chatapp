package main

import (
	"context"
	"flag"

	log "github.com/sirupsen/logrus"

	account_router "gsm/router/account"
)

var (
	sqlConfigPath string
	dbkey string
	port          string
)

func init() {
	flag.StringVar(&sqlConfigPath, "sql", "", "sql config path")
	flag.StringVar(&dbkey, "db-key", "", "da encrypted key")
	flag.StringVar(&port, "port", ":8080", "service port")
}

func main() {
	flag.Parse()

	// create gin router for realtime service.
	router, err := account_router.NewRouter(context.Background(), account_router.RouterConfig{
		SqlConfigPath: sqlConfigPath,
		DBKey: dbkey,
	})
	if err != nil {
		log.Fatalf("Init account router error: %v", err)
	}

	if err := router.Run(port); err != nil {
		log.Fatalf("failed to start account server: %v", err)
	}
}
