package util

import (
	"time"
)

func TimestampToDatetime(ts int64) time.Time {
	return time.Unix(ts, 0)
}

func PtrString(s string) *string {
	return &s
}

func PtrInt64(val int64) *int64 {
	return &val
}

func PtrInt32(val int32) *int32 {
	return &val
}