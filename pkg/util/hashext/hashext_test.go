package hashext

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBcryptPasswordAndcomparePasswords test BcryptPassword and ComparePasswords
func TestBcryptPasswordAndComparePasswords(t *testing.T) {
	password := "qwe123QWE"
	hashedPassword, err := BcryptPassword(password)
	res := ComparePasswords(hashedPassword, password)
	assert.NoError(t, err)
	assert.NoError(t, res)
}

func TestSha256(t *testing.T) {
	s := Sha256("qwe123QWE")
	assert.Equal(t, "92314C299E28E811FCBC64E6AD1209CCF4439A962B4384CD0D47D23A3067D574", s)
}
