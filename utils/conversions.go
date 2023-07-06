package utils

import "strconv"

// ConvertStrToInt takes a string as input and converts it to an integer.
// If the string cannot be converted to an integer, the function returns 0.
// The function uses the Atoi function from the strconv package for conversion.
//
// Parameters:
//
//	s: a string value that should be converted to integer
//
// Returns:
//
//	an integer representing the string input, or 0 if conversion is not possible
func ConvertStrToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
