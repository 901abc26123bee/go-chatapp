package timeout

import (
	"context"
	"gsm/pkg/errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const timeoutPeriod = 60 * time.Second

// RequestTimeoutHandler is a gin middleware which set default timeout to gin request
func RequestTimeoutHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Check if ctx.Request.Context is nil
		if ctx.Request.Context() == nil {
			ctx.Request = ctx.Request.WithContext(context.Background())
		}
		// Check if the deadline is already set
		deadline, ok := ctx.Deadline()
		if ok {
			log.Infof("timeout is already set in the context, context will be canceled at: %v", deadline)
		} else {
			// replace the request context with the new context.
			cancelCtx, cancel := context.WithTimeout(ctx.Request.Context(), timeoutPeriod)
			defer cancel()
			ctx.Request = ctx.Request.WithContext(cancelCtx)
		}

		ctx.Next()

		// update the error message to indicate a timeout occurred.
		if ctx.Err() == context.DeadlineExceeded {
			ctx.AbortWithError(http.StatusRequestTimeout, errors.NewError(errors.RequestTimeout, "request timed out"))
			return
		}
	}
}
