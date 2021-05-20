package domain

import (
	"errors"

	"accounts/util"
)

var (
	errInvalidValue = errors.New("invalid value")
)

type FieldID int32

func (field *FieldID) Validate() error {
	if field == nil {
		return nil
	}

	if *field < 0 {
		return errInvalidValue
	}

	return nil
}

type FieldEmail string

func (field *FieldEmail) Validate() error {
	if field == nil {
		return nil
	}

	return util.ValidateString(string(*field), maxLenEmail, nil, regexpEmail)
}

type FieldInterest string

func (field *FieldInterest) Validate() error {
	if field == nil {
		return nil
	}

	return util.ValidateString(string(*field), maxLenInterest, nil, nil)
}

type FieldFirstname string

func (field *FieldFirstname) Validate() error {
	if field == nil {
		return nil
	}

	return util.ValidateString(string(*field), maxLenFirstname, nil, nil)
}

type FieldSurname string

func (field *FieldSurname) Validate() error {
	if field == nil {
		return nil
	}

	return util.ValidateString(string(*field), maxLenSurname, nil, nil)
}

type FieldPhone string

func (field *FieldPhone) Validate() error {
	if field == nil {
		return nil
	}

	return util.ValidateString(string(*field), maxLenPhone, nil, regexpPhone)
}

type FieldSex string

func (field *FieldSex) Validate() error {
	if field == nil {
		return nil
	}

	if !sexIsValid(string(*field)) {
		return errInvalidValue
	}

	return nil
}

type FieldStatus string

func (field *FieldStatus) Validate() error {
	if field == nil {
		return nil
	}

	if !statusIsValid(string(*field)) {
		return errInvalidValue
	}

	return nil
}

type FieldCountry string

func (field *FieldCountry) Validate() error {
	if field == nil {
		return nil
	}

	return util.ValidateString(string(*field), maxLenCountry, nil, nil)
}

type FieldCity string

func (field *FieldCity) Validate() error {
	if field == nil {
		return nil
	}

	return util.ValidateString(string(*field), maxLenCity, nil, nil)
}

type FieldBirth int64

func (field *FieldBirth) Validate() error {
	if field == nil {
		return nil
	}

	return util.ValidateTimestamp(int64(*field), &minBirth, &maxBirth)
}

type FieldJoined int64

func (field *FieldJoined) Validate() error {
	if field == nil {
		return nil
	}

	return util.ValidateTimestamp(int64(*field), &minJoined, &maxJoined)
}

type FieldPremium int64

func (field *FieldPremium) Validate() error {
	if field == nil {
		return nil
	}

	return util.ValidateTimestamp(int64(*field), &minPremiumStart, &minPremiumEnd)
}

type FieldTimestamp int64

func (field *FieldTimestamp) Validate() error {
	if field == nil {
		return nil
	}

	return util.ValidateTimestamp(int64(*field), nil, nil)
}
