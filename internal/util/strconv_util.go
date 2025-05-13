package util

import (
	"strconv"

	"github.com/rotisserie/eris"
)

func Atoi(a string) (int, error) {
	i, err := strconv.Atoi(a)
	if err != nil {
		return 0, eris.Wrap(err, "error converting string to integer")
	}

	return i, nil
}
