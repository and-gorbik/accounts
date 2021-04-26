package domain

import (
	"time"
)

const (
	statusFree        = "свободны"
	statusBusy        = "заняты"
	statusComplicated = "все сложно"
)

var (
	maxLenInterest = 100
	minBirth       = time.Date(1950, 1, 1, 0, 0, 0, 0, time.Local).Unix()
	maxBirth       = time.Date(2005, 1, 1, 0, 0, 0, 0, time.Local).Unix()
)
