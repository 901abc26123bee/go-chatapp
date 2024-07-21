package rediscache

import (
	"github.com/go-redis/redis"

	"gsm/pkg/errors"
)

// IsKeyNotExistError check if is redis nil error.
func IsKeyNotExistError(errMessage error) error {
	if errMessage == redis.Nil {
		return errors.NewErrorf(errors.NotFound,
			"redis: nil")
	}

	return nil
}
