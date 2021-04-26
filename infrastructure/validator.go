package infrastructure

import (
	"errors"
	"regexp"
	"strconv"
)

var (
	errLessThanMin    = errors.New("less than min")
	errGreaterThanMax = errors.New("greater than max")
	errInvalidString  = errors.New("invalid string")
)

func ValidateInt(s string, min, max int) error {
	val, err := strconv.Atoi(s)
	if err != nil {
		return err
	}

	if val < min {
		return errLessThanMin
	}

	if val > max {
		return errGreaterThanMax
	}

	return nil
}

func ValidateString(s string, maxlen int, like *string, pattern *regexp.Regexp) error {
	if len(s) > maxlen {
		return errGreaterThanMax
	}

	if pattern != nil && !pattern.MatchString(s) {
		return errInvalidString
	}

	if like != nil && s != *like {
		return errInvalidString
	}

	return nil
}

func ValidateTimestamp(s string, min, max int64) error {
	timestamp, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	if timestamp < min {
		return errLessThanMin
	}

	if timestamp > max {
		return errGreaterThanMax
	}

	return nil
}
