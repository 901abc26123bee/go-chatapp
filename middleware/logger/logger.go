package logger

import (
	"bytes"
	"gsm/pkg/errors"
	"gsm/pkg/util/convert"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// ResponseWriter is a custom writer that captures the response body
type ResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

// LoggerHandler will write customized log info
func LoggerHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()

		var requestBodyBytes []byte
		var err error
		if ctx.Request.Body != nil {
			requestBodyBytes, err = io.ReadAll(ctx.Request.Body)
			if err != nil {
				ctx.Error(errors.NewError(errors.Unknown, "failed to parse the request body"))
				ctx.Abort()
				return
			}
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(requestBodyBytes))
		}

		// Create a custom ResponseWriter
		rw := &ResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		ctx.Writer = rw

		ctx.Next()
		defer ctx.Request.Body.Close()

		sb := []string{ctx.Request.Method, ctx.Request.RequestURI, ctx.Request.Proto}
		reqLine := strings.Join(sb, " ")
		elapsed := time.Since(startTime).String()
		errMsg := ctx.Errors.ByType(gin.ErrorTypePrivate).String()
		fields := log.Fields{
			"request-header": ctx.Request.Header,
			"request-body":   convert.FormatJsonString(string(requestBodyBytes)),
			"response-body":  convert.FormatJsonString(rw.body.String()),
			"elapsed":        elapsed,
			"status-code":    ctx.Writer.Status(),
		}
		if errMsg == "" {
			entry := log.WithFields(fields)
			entry.Infof("[%s] latency: %v, success", reqLine, elapsed)
			return
		}
		fields["error-message"] = errMsg
		entry := log.WithFields(fields)
		entry.Errorf("[%s] latency: %v, failed", reqLine, elapsed)
	}
}
