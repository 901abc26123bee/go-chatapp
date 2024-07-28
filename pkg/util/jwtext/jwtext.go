package jwtext

import (
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

// JwtMapClaims defines as an alias of jwt MapClaims
type JwtMapClaims = jwt.MapClaims

const (
	// SubjectAccess implies the access token category.
	SubjectAccess = "access_token"
	// SubjectRefresh implies the refresh token category.
	SubjectRefresh = "refresh_token"
	// AccessTokenDuration defines access token duration as 1 hours
	AccessTokenDuration = time.Hour
	// RefreshTokenDuration defines refresh token duration as 30 days
	RefreshTokenDuration = 720 * time.Hour
)

var (
	// TokenIssuer defines the required issuer in jwt-token.
	TokenIssuer = fmt.Sprintf("gsm-%s", "dev") // TODO: get flavor "dev" from evn variable
	// TokenAudience defines the required audience in jwt-token.
	TokenAudience = fmt.Sprintf("gsm-%s", "dev") // TODO: get flavor "dev" from evn variable
)

// Claims defines fields used for JWT
const (
	ClaimAudience  = "aud"
	ClaimExpiresAt = "exp"
	ClaimID        = "jti"
	ClaimIssuedAt  = "iat"
	ClaimIssuer    = "iss"
	ClaimNotBefore = "nbf"
	ClaimSubject   = "sub"
)

// CreateAccessToken creates access token with access subject and with duration limit (10 minutes) and return jwt token and expire time
func CreateAccessToken(uid, secret string, claims JwtMapClaims, duration time.Duration) (string, time.Duration, error) {
	token, err := createToken(uid, secret, SubjectAccess, minDuration(AccessTokenDuration, duration), claims)
	return token, minDuration(AccessTokenDuration, duration), err
}

// CreateRefreshToken creates refresh token with refresh subject and return jwt token and expire time
func CreateRefreshToken(uid, secret string, claims JwtMapClaims) (string, time.Duration, error) {
	token, err := createToken(uid, secret, SubjectRefresh, RefreshTokenDuration, claims)
	return token, RefreshTokenDuration, err
}

func createToken(uid, secret, subject string, duration time.Duration, customClaims JwtMapClaims) (string, error) {
	// create jwt claims which is composed of standard infos and custom infos.
	now := time.Now()
	exp := now.Add(duration)
	claims := JwtMapClaims{
		ClaimAudience:  TokenAudience,
		ClaimID:        uid,
		ClaimIssuedAt:  now.Unix(),
		ClaimIssuer:    TokenIssuer,
		ClaimSubject:   subject,
		ClaimNotBefore: now.Unix(),
		ClaimExpiresAt: exp.Unix(),
	}
	for key, val := range customClaims {
		claims[key] = val
	}

	// sign the claim with secret to generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}

func minDuration(a, b time.Duration) time.Duration {
	if a <= b {
		return a
	} else {
		return b
	}
}

// ParseToken parses token string and returns a claim map.
func ParseToken(tokenString, secret string) (jwt.MapClaims, error) {
	// parse the token string.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token %s: %v", tokenString, err)
	} else if !token.Valid {
		return nil, fmt.Errorf("token %s not valid", tokenString)
	}

	// convert the claims to map claims.
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unable to cast jwt claims")
	}

	// do verification.
	if err := token.Claims.Valid(); err != nil {
		return nil, fmt.Errorf("invalid jwt claim: %v", err)
	}
	if !claims.VerifyAudience(TokenAudience, true) {
		return nil, fmt.Errorf("invalid token audience")
	}
	if !claims.VerifyIssuer(TokenIssuer, true) {
		return nil, fmt.Errorf("invalid token issuer")
	}

	return claims, nil
}

// GetFieldFromClaims get field from claims
func GetFieldFromClaims(claims JwtMapClaims, field string) (interface{}, error) {
	if value, ok := claims[field]; ok {
		return value, nil
	} else {
		return "", fmt.Errorf("claims does not contain %s", field)
	}
}
