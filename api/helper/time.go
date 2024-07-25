package helper

import "time"

func NewTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}
