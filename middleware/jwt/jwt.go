package jwt

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"

	pkgerrors "gsm/pkg/errors"
	"gsm/pkg/util/jwtext"
)

const (
	// JWTClaimID defines the key of jwt claims ID save in context variable
	JWTClaimID = "jwt_claim_id"
)

// HeaderAuthorizationHandler verify bearer token in header is a valid jwt token
func HeaderAuthorizationHandler(secret string) gin.HandlerFunc {
	return handleAuthorizationWithGin(secret, func(r *http.Request) string {
		return strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	})
}

func HttpRequestQueryParamsAuthorizationHandler(secret string) gin.HandlerFunc {
	return handleAuthorizationWithGin(secret, func(r *http.Request) string {
		return r.URL.Query().Get(jwtext.AccessTokenQueryParamKey)
	})
}

func handleAuthorizationWithGin(secret string, getToken func(req *http.Request) string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// check jwt secret is empty to avoid using an empty string to create token by accident
		if secret == "" {
			ctx.Error(pkgerrors.NewError(pkgerrors.InternalServerError, "jwt secret is empty"))
			ctx.Abort()
			return
		}

		// get jwt token from bearer token
		token := getToken(ctx.Request)
		if token == "" {
			ctx.Error(pkgerrors.NewError(pkgerrors.TokenEmpty, "token is empty"))
			ctx.Abort()
			return
		}

		// verify token is valid
		claims, err := jwtext.ParseToken(token, secret)
		if err != nil {
			errCode := pkgerrors.TokenInvalid
			if errors.Is(err, jwt.ErrTokenExpired) {
				errCode = pkgerrors.TokenExpired
			}
			ctx.Error(pkgerrors.NewErrorf(errCode, "failed to verify token: %v", err))
			ctx.Abort()
			return
		}

		// parse jwt claims id and set as context variable for later usage
		claimsID, err := jwtext.GetFieldFromClaims(claims, jwtext.ClaimID)
		if err != nil {
			ctx.Error(pkgerrors.NewErrorf(pkgerrors.TokenInvalid, "failed to get id from calim: %v", err))
			ctx.Abort()
			return
		}
		ctx.Set(JWTClaimID, claimsID)
	}
}
