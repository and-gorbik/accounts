package domain

import (
	"errors"

	"accounts/infrastructure"
)

var (
	errInvalidValue = errors.New("invalid value")
)

type FieldID int32

func (id FieldID) Validate() error {
	if id < 0 {
		return errInvalidValue
	}

	return nil
}

type FieldEmail string

func (s FieldEmail) Validate() error {
	return infrastructure.ValidateString(string(s), maxLenEmail, nil, regexpEmail)
}

type FieldInterest string

func (s FieldInterest) Validate() error {
	return infrastructure.ValidateString(string(s), maxLenInterest, nil, nil)
}

type FieldFirstname string

func (s FieldFirstname) Validate() error {
	return infrastructure.ValidateString(string(s), maxLenFirstname, nil, nil)
}

type FieldSurname string

func (s FieldSurname) Validate() error {
	return infrastructure.ValidateString(string(s), maxLenSurname, nil, nil)
}

type FieldPhone string

func (s FieldPhone) Validate() error {
	return infrastructure.ValidateString(string(s), maxLenPhone, nil, regexpPhone)
}

type FieldSex string

func (s FieldSex) Validate() error {
	if !sexIsValid(string(s)) {
		return errInvalidValue
	}

	return nil
}

type FieldStatus string

func (s FieldStatus) Validate() error {
	if !statusIsValid(string(s)) {
		return errInvalidValue
	}

	return nil
}

type FieldCountry string

func (s FieldCountry) Validate() error {
	return infrastructure.ValidateString(string(s), maxLenCountry, nil, nil)
}

type FieldCity string

func (s FieldCity) Validate() error {
	return infrastructure.ValidateString(string(s), maxLenCity, nil, nil)
}

type FieldBirth int64

func (t FieldBirth) Validate() error {
	return infrastructure.ValidateTimestamp(int64(t), &minBirth, &maxBirth)
}

type FieldJoined int64

func (t FieldJoined) Validate() error {
	return infrastructure.ValidateTimestamp(int64(t), &minJoined, &maxJoined)
}

type FieldPremium int64

func (t FieldPremium) Validate() error {
	return infrastructure.ValidateTimestamp(int64(t), &minPremiumStart, &minPremiumEnd)
}

type FieldTimestamp int64

func (t FieldTimestamp) Validate() error {
	return infrastructure.ValidateTimestamp(int64(t), nil, nil)
}
