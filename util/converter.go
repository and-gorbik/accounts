package util

import (
	"time"

	"github.com/google/uuid"
)

func TimestampToDatetime(ts *int64) *time.Time {
	if ts == nil {
		return nil
	}

	res := time.Unix(*ts, 0)
	return &res
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

func PtrUUID(val uuid.UUID) *uuid.UUID {
	if val == uuid.Nil {
		return nil
	}

	return &val
}
