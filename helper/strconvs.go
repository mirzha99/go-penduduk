package helper

import (
	"strconv"
)

func Atoi(s string) any {
	r, err := strconv.Atoi(s)

	if err != nil {
		return err.Error()
	}
	return r
}
func Itoa(s int) string {
	r := strconv.Itoa(s)
	return r
}
