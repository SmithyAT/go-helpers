package utils

import "strconv"

// ConvertStrToInt converts a string to an integer
func ConvertStrToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
