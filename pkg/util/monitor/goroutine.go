package monitor

import (
	"context"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

// LogGoroutineCount log goroutine count fro connection name
func LogGoroutineCount(ctx context.Context, duration time.Duration, serviceName string) {
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("context canceled for logGoroutineCount")
			return
		case <-ticker.C:
			log.Infof("Current goroutine count %d for service %s ", runtime.NumGoroutine(), serviceName)
		}
	}
}
