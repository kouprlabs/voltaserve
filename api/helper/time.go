package helper

import (
	"time"
)

func NewTimeString() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func TimeToString(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

func StringToTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

func StringToTimestamp(s string) int64 {
	return StringToTime(s).UnixMilli()
}

func TimeToTimestamp(t time.Time) int64 {
	return t.UnixMilli()
}
