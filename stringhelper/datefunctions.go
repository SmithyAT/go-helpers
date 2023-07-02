package stringhelper

import (
	"errors"
	"regexp"
	"strings"
)

// ExtractDate is a function that takes a string and returns a formatted date string and an error.
// This can be useful if your workflow includes parsing dates from filenames.
//
// The function first compiles a regex that matches either a date formatted as "YYYYMMDD" or a date with delimiters,
// such as "YYYY-MM-DD" or any combination of "_", "-", or ".".
//
// If a matching string is found within `inputString`, all delimiters ("_", "-", and ".") in the string are then removed
// to maintain a consistent "YYYYMMDD" format.
//
// If no matching string is found (date is empty), an error message is returned declaring that no date was found.
func ExtractDate(inputString string) (date string, err error) {
	re := regexp.MustCompile(`(20\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01]))|(20\d{2}[._-](0[1-9]|1[0-2])[._-](0[1-9]|[12]\d|3[01]))`)
	date = re.FindString(inputString)
	charsToRemove := []string{"-", "_", "."}
	for _, char := range charsToRemove {
		date = strings.ReplaceAll(date, char, "")
	}
	if date == "" {
		err = errors.New("no date found in the string")
	}

	return
}

// SplitDate takes a date string in the "YYYYMMDD" format and returns individual string components
// for the year, month, and day.
//
// The function expects a date string of exactly 8 characters in length. The first four characters
// are treated as the "year", the next two as the "month", and the final two as the "day".
// If the input string does not meet these criteria, the function's output will be inconsistent.
//
// The function's return values are in the following order: year, month, day. Each returned string
// contains the corresponding part of the date.
//
// Example usage:
//
//	year, month, day := SplitDate("20231225")
//	fmt.Println(year)  // Outputs: "2023"
//	fmt.Println(month) // Outputs: "12"
//	fmt.Println(day)   // Outputs: "25"
func SplitDate(date string) (year string, month string, day string, err error) {
	if len(date) != 8 {
		return "", "", "", errors.New("date should be in format YYYYMMDD --")
	}

	year = date[0:4]
	month = date[4:6]
	day = date[6:8]
	return year, month, day, nil
}
