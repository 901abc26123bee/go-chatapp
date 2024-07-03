package errors_middleware

import (
	"gsm/pkg/errors"

	"github.com/gin-gonic/gin"
)

// ErrorHandler is a gin middleware which catch first error and error string to response in text/plain
func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// execute pending handlers
		ctx.Next()

		// do nothing if not error
		if len(ctx.Errors) == 0 {
			return
		}

		// get the first error in ctx.Errors
		ctxFirstError := ctx.Errors[0]
		firstErr := ctxFirstError.Err
		errContent, ok := firstErr.(*errors.ErrContent)
		if !ok {
			errContent = errors.NewError(errors.Unknown, firstErr.Error())
		}
		errCode := errContent.StatusCode

		// error without detail for response
		errorWithoutDetail := &errors.ErrContent{
			StatusCode: errContent.StatusCode,
			Code:       errContent.Code,
			Message:    errContent.Message,
		}

		ctx.JSON(errCode, errorWithoutDetail)
	}
}
