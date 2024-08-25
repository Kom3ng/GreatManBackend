package util

import "strconv"

func ParseUint(s string) (uint, error) {
	i, err := strconv.ParseUint(s, 10, 32)

	if err != nil {
		return 0, err
	}

	return uint(i), nil
}
