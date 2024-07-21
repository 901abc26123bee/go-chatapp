package sonyflake

import (
	"gsm/pkg/errors"

	"github.com/sony/sonyflake"
)

// IDGenerator defines interface of id generator
type IDGenerator interface {
	NextID() (uint64, error)
}

// NewSonyFlake return IDGenerator with sony flake
func NewSonyFlake() (IDGenerator, error) {
	var st sonyflake.Settings
	sf := sonyflake.NewSonyflake(st)
	if sf == nil {
		return nil, errors.Errorf("sonyflake not created")
	}
	return sf, nil
}
