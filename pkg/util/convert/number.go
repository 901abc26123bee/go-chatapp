package convert

import (
	"fmt"
	"strconv"
)

func ToUint64(i interface{}) (uint64, error) {
	switch v := i.(type) {
	case int:
		if v < 0 {
			return 0, fmt.Errorf("negative value %d cannot be converted to uint64", v)
		}
		return uint64(v), nil
	case int8:
		if v < 0 {
			return 0, fmt.Errorf("negative value %d cannot be converted to uint64", v)
		}
		return uint64(v), nil
	case int16:
		if v < 0 {
			return 0, fmt.Errorf("negative value %d cannot be converted to uint64", v)
		}
		return uint64(v), nil
	case int32:
		if v < 0 {
			return 0, fmt.Errorf("negative value %d cannot be converted to uint64", v)
		}
		return uint64(v), nil
	case int64:
		if v < 0 {
			return 0, fmt.Errorf("negative value %d cannot be converted to uint64", v)
		}
		return uint64(v), nil
	case uint:
		return uint64(v), nil
	case uint8:
		return uint64(v), nil
	case uint16:
		return uint64(v), nil
	case uint32:
		return uint64(v), nil
	case uint64:
		return v, nil
	case string:
		n, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, err
		}
		return n, nil
	default:
		return 0, fmt.Errorf("unsupported type: %T", v)
	}
}
