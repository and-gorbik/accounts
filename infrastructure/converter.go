package infrastructure

import (
	"time"
)

func TimestampToDatetime(ts int64) time.Time {
	return time.Unix(ts, 0)
}
