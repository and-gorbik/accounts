package domain

import (
	"time"
)

const (
	statusFree        = "свободны"
	statusBusy        = "заняты"
	statusComplicated = "все сложно"
)

const (
	maxLenInterest  = 100
	maxLenEmail     = 100
	maxLenFirstname = 50
	maxLenSurname   = 50
	maxLenPhone     = 16
	maxLenSex       = 1
	maxLenCountry   = 50
	maxLenCity      = 50
)

var (
	minBirth        = time.Date(1950, 1, 1, 0, 0, 0, 0, time.Local).Unix()
	maxBirth        = time.Date(2005, 1, 1, 0, 0, 0, 0, time.Local).Unix()
	minJoined       = time.Date(2011, 1, 1, 0, 0, 0, 0, time.Local).Unix()
	maxJoined       = time.Date(2018, 1, 1, 0, 0, 0, 0, time.Local).Unix()
	minPremiumStart = time.Date(2018, 1, 1, 0, 0, 0, 0, time.Local).Unix()
	minPremiumEnd   = time.Date(2018, 1, 1, 0, 0, 0, 0, time.Local).Unix()
)

func statusIsValid(s string) bool {
	return s == statusBusy || s == statusComplicated || s == statusFree
}
