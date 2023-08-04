package contime

import (
	"fmt"
	"time"
)

func ParseTime(timeStr string) time.Time {
	// You may want to use a library like time.Parse to convert the string to a time.Time object.
	// For simplicity, here's a basic parsing method:
	layout := "2006-01-02T15:04"
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return time.Time{}
	}

	// Format the time to ISOTime format

	return t
}
