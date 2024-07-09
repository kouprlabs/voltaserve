package helper

import "time"

func ToUTCString(value *string) string {
	if value == nil {
		return ""
	}
	parsedTime, err := time.Parse(time.RFC3339, *value)
	if err != nil {
		return ""
	}
	return parsedTime.Format(time.RFC1123)
}
