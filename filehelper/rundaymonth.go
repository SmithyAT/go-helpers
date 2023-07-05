package filehelper

import (
	"os"
	"time"
)

// ShouldRunDayMonth checks if the current day or month is different from the last execution.
// It checks the timestamp of the last execution stored in the file at the provided path.
// If the file does not exist, it will be created.
// The function returns two boolean values indicating whether there was a change in day or month
// since the last execution, and an error value.
// If the day or month has changed since the last execution, respective bool will be true.
// If not, it will be false. If there was an error reading from or writing to the file,
// the error return value will be non-nil.
func ShouldRunDayMonth(dateFilePath string) (runDay bool, runMonth bool, retErr error) {

	// Get the current time
	now := time.Now()

	// Read the last execution timestamp from the file
	content, err := os.ReadFile(dateFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// If the file does not exist, create it and write the current timestamp
			err := os.WriteFile(dateFilePath, []byte(now.Format(time.RFC3339)), 0644)
			if err != nil {
				return false, false, err
			}
			// Read the file again
			content, err = os.ReadFile(dateFilePath)
			if err != nil {
				return false, false, err
			}
		} else {
			return false, false, err
		}
	}

	// Parse the last execution timestamp
	lastExecution, err := time.Parse(time.RFC3339, string(content))
	if err != nil {
		return false, false, err
	}

	// Check if there was a day change
	dayChanged := lastExecution.YearDay() != now.YearDay()

	// Check if there was a month change
	monthChanged := lastExecution.Month() != now.Month()

	// Update the file with the current timestamp
	if err := os.WriteFile(dateFilePath, []byte(now.Format(time.RFC3339)), 0644); err != nil {
		return dayChanged, monthChanged, err
	}

	return dayChanged, monthChanged, nil
}
