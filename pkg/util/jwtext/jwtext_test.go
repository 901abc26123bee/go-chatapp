package jwtext

import (
	"testing"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

const (
	testUID   = "testUID"
	secret    = "secret"
	testField = "testField"
)

func TestCreateAccessToken(t *testing.T) {
	// test without claim
	token, dur, err := CreateAccessToken(testUID, secret, nil, AccessTokenDuration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, AccessTokenDuration, dur)

	// test with claim
	tokenClaim, dur, errClaim := CreateAccessToken(testUID, secret, JwtMapClaims{testField: testField}, AccessTokenDuration)
	assert.NoError(t, errClaim)
	assert.NotEmpty(t, tokenClaim)
	assert.Equal(t, AccessTokenDuration, dur)

	assert.NotEqual(t, token, tokenClaim)
}

func TestCreateRefreshToken(t *testing.T) {
	// Test without claim
	token, dur, err := CreateRefreshToken(testUID, secret, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, RefreshTokenDuration, dur)

	// Test with claim
	tokenClaim, dur, errClaim := CreateRefreshToken(testUID, secret, JwtMapClaims{testField: testField})
	assert.NoError(t, errClaim)
	assert.NotEmpty(t, tokenClaim)
	assert.Equal(t, RefreshTokenDuration, dur)

	assert.NotEqual(t, token, tokenClaim)
}

func createTestTokenWithCustomTimes(s string, issueAt, start, end time.Time) (string, error) {
	claims := JwtMapClaims{
		ClaimAudience:  TokenAudience,
		ClaimID:        testUID,
		ClaimIssuedAt:  issueAt.Unix(),
		ClaimIssuer:    TokenIssuer,
		ClaimSubject:   SubjectAccess,
		ClaimNotBefore: start.Unix(),
		ClaimExpiresAt: end.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s))

	return tokenString, err
}

func createTestToken(s string) (string, error) {
	now := time.Now()
	exp := now.Add(time.Hour)
	return createTestTokenWithCustomTimes(s, now, now, exp)
}

func createExpiredTestToken(s string) (string, error) {
	// create jwt claims.
	issueAt := time.Now()
	start := issueAt.Add(-2 * time.Hour)
	end := issueAt.Add(-time.Hour) // 1 hour before

	return createTestTokenWithCustomTimes(s, issueAt, start, end)
}

func TestVerifyToken(t *testing.T) {
	// valid token
	token, err := createTestToken(secret)
	if err != nil {
		assert.FailNowf(t, "failed to generate test token: %s", err.Error())
	}
	claims, err := ParseToken(token, secret)
	assert.NoError(t, err)
	assert.NotEmpty(t, claims)
	assert.Equal(t, testUID, claims[ClaimID])

	// invalid token (signed by fakeSecret)
	invalidToken, err := createTestToken("fakeSecret")
	if err != nil {
		assert.FailNowf(t, "failed to generate test token: %s", err.Error())
	}
	claims, err = ParseToken(invalidToken, secret)
	assert.Empty(t, claims)
	assert.Error(t, err)

	// invalid token (empty token)
	claims, err = ParseToken("", secret)
	assert.Empty(t, claims)
	assert.Error(t, err)

	// expired token
	expiredToken, err := createExpiredTestToken(secret)
	if err != nil {
		assert.FailNowf(t, "failed to generate test token: %s", err.Error())
	}
	claims, err = ParseToken(expiredToken, secret)
	assert.Empty(t, claims)
	assert.Error(t, err)
}
