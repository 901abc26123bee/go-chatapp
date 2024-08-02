package main

import (
	"flag"

	log "github.com/sirupsen/logrus"

	realtime_router "gsm/router/realtime"
)

var (
	sqlConfigPath   string
	redisConfigPath string
	port            string
	jwtSecret       string
	mongodbPath     string
)

func init() {
	flag.StringVar(&sqlConfigPath, "sql", "", "sql config path")
	flag.StringVar(&redisConfigPath, "redis", "", "redis config path")
	flag.StringVar(&port, "port", ":8081", "service port")
	flag.StringVar(&jwtSecret, "jwt", "", "jwt secret")
	flag.StringVar(&mongodbPath, "mongodb", "", "mongodb config path")
}

func main() {
	flag.Parse()

	// create gin router for realtime service.
	router, err := realtime_router.NewRouter(realtime_router.RouterConfig{
		SqlConfigPath:     sqlConfigPath,
		RedisConfigPath:   redisConfigPath,
		MongodbConfigPath: mongodbPath,
		JwtSecret:         jwtSecret,
	})
	if err != nil {
		log.Fatalf("Init realtime router error: %v", err)
	}

	if err := router.Run(port); err != nil {
		log.Fatalf("failed to start realtime server: %v", err)
	}
}
