package utils

import (
	"time"
)

func TimeToString(t *time.Time) *string {
	if t != nil && !t.IsZero() {
		str := t.Format(time.RFC3339Nano)
		return &str
	}
	return nil
}

func ParseTime(datePtr *string) *time.Time {
	if datePtr != nil {
		t, err := time.Parse(time.RFC3339Nano, *datePtr)
		if err == nil {
			return &t
		}
	}
	return nil
}
