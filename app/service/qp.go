package service

import (
	"errors"
	"net/url"
	"strings"

	"accounts/domain"
	"accounts/util"
)

type QueryParam struct {
	Type     int
	Field    string
	StrValue string
	Op       *string
}

const (
	typeInt = iota
	typeStr
	typeTimestamp
)

const (
	opEq       = "eq"
	opLt       = "lt"
	opGt       = "gt"
	opNeq      = "neq"
	opAny      = "any"
	opDomain   = "domain"
	opNull     = "null"
	opStarts   = "starts"
	opCode     = "code"
	opYear     = "year"
	opContains = "contains"
	opNow      = "now"

	noOp = ""
)

const (
	qpSex       = "sex"
	qpEmail     = "email"
	qpStatus    = "status"
	qpFirstname = "fname"
	qpSurname   = "sname"
	qpPhone     = "phone"
	qpCountry   = "country"
	qpCity      = "city"
	qpBirth     = "birth"
	qpInterests = "interests"
	qpLikes     = "likes"
	qpPremium   = "premium"

	qpLimit   = "limit"
	qpQueryID = "query_id"
)

var (
	yes = struct{}{}
)

var (
	qpRules = map[string]map[string]struct{}{
		qpSex:       {opEq: yes},
		qpEmail:     {opDomain: yes, opGt: yes, opLt: yes},
		qpStatus:    {opEq: yes, opNeq: yes},
		qpFirstname: {opEq: yes, opAny: yes, opNull: yes},
		qpSurname:   {opEq: yes, opStarts: yes, opNull: yes},
		qpPhone:     {opCode: yes, opNull: yes},
		qpCountry:   {opEq: yes, opNull: yes},
		qpCity:      {opEq: yes, opAny: yes, opNull: yes},
		qpBirth:     {opLt: yes, opGt: yes, opYear: yes},
		qpInterests: {opContains: yes, opAny: yes},
		qpLikes:     {opContains: yes},
		qpPremium:   {opNow: yes, opNull: yes},
	}

	qpTypes = map[string]map[string]int{
		qpSex:       {noOp: typeStr, opEq: typeStr},
		qpEmail:     {noOp: typeStr, opDomain: typeStr, opGt: typeStr, opLt: typeStr},
		qpStatus:    {noOp: typeStr, opEq: typeStr, opNeq: typeStr},
		qpFirstname: {noOp: typeStr, opEq: typeStr, opAny: typeStr, opNull: typeInt},
		qpSurname:   {noOp: typeStr, opEq: typeStr, opStarts: typeStr, opNull: typeInt},
		qpPhone:     {noOp: typeStr, opCode: typeInt, opNull: typeInt},
		qpCountry:   {noOp: typeStr, opEq: typeStr, opNull: typeInt},
		qpCity:      {noOp: typeStr, opEq: typeStr, opAny: typeInt, opNull: typeInt},
		qpBirth:     {noOp: typeTimestamp, opLt: typeTimestamp, opGt: typeTimestamp, opYear: typeInt},
		qpInterests: {noOp: typeStr, opContains: typeStr, opAny: typeStr},
		qpLikes:     {noOp: typeInt, opContains: typeInt},
		qpPremium:   {noOp: typeTimestamp, opNow: typeInt, opNull: typeInt},
	}
)

var (
	errInvalidParam         = errors.New("invalid query param")
	errInvalidValue         = errors.New("invalid value")
	errEmptyValue           = errors.New("empty value")
	errMissingRequiredParam = errors.New("missing required param")
)

func ParseQueryParams(qps url.Values, withOp bool) (map[string]QueryParam, error) {
	limit, err := parseLimit(qps)
	if err != nil {
		return nil, err
	}

	if _, ok := qps[qpQueryID]; !ok {
		return nil, errMissingRequiredParam
	}

	params := make(map[string]QueryParam, len(qps)-1)
	params[limit.Field] = limit

	for param, values := range qps {
		if param == qpLimit || param == qpQueryID {
			continue
		}

		// TODO: pass values []string directly
		qp, err := parseQueryParam(param, strings.Join(values, ","), withOp)
		if err != nil {
			return nil, err
		}

		if err := validateValues(qp.Field, qp.Op, values); err != nil {
			return nil, err
		}

		params[qp.Field] = qp
	}

	return params, nil
}

func parseQueryParam(param string, value string, withOp bool) (qp QueryParam, err error) {
	if withOp {
		tokens := strings.Split(param, "_")
		if len(tokens) != 2 {
			err = errInvalidParam
			return
		}

		if _, ok := qpRules[tokens[0]][tokens[1]]; !ok {
			err = errInvalidParam
			return
		}

		qp.Field = tokens[0]
		qp.Op = &tokens[1]
	} else {
		qp.Field = param
	}

	if value == "" {
		err = errEmptyValue
		return
	}

	var op string
	if qp.Op == nil {
		op = noOp
	}
	qp.Type = qpTypes[param][op]
	qp.StrValue = value

	return
}

func parseLimit(qps url.Values) (QueryParam, error) {
	if _, ok := qps[qpLimit]; !ok {
		return QueryParam{}, errMissingRequiredParam
	}

	if len(qps[qpLimit]) != 1 {
		return QueryParam{}, errInvalidValue
	}

	return QueryParam{
		Field:    qpLimit,
		StrValue: qps[qpLimit][0],
	}, nil
}

// TODO: refact to make it more readable
func validateValues(param string, op *string, values []string) error {
	if len(values) == 0 {
		return nil
	}

	switch param {
	case qpSex:
		return (*domain.FieldSex)(&values[0]).Validate()
	case qpEmail:
		if op == nil {
			return (*domain.FieldEmail)(&values[0]).Validate()
		}
	case qpStatus:
		return (*domain.FieldStatus)(&values[0]).Validate()
	case qpFirstname:
		if op == nil {
			return (*domain.FieldFirstname)(&values[0]).Validate()
		}
		if *op == opNull {
			return validateBoolValue(values[0])
		}
		for _, val := range values {
			if err := (*domain.FieldFirstname)(&val).Validate(); err != nil {
				return err
			}
		}
	case qpSurname:
		if op == nil || (op != nil && *op == opEq) {
			return (*domain.FieldSurname)(&values[0]).Validate()
		}
		if op != nil && *op == opNull {
			return validateBoolValue(values[0])
		}
	case qpPhone:
		if op == nil {
			return (*domain.FieldPhone)(&values[0]).Validate()
		}
		if *op == opNull {
			return validateBoolValue(values[0])
		}
	case qpCountry:
		if op != nil && *op == opNull {
			return validateBoolValue(values[0])
		}
		return (*domain.FieldCountry)(&values[0]).Validate()
	case qpCity:
		if op != nil && *op == opNull {
			return validateBoolValue(values[0])
		}
		for _, val := range values {
			if err := (*domain.FieldCity)(&val).Validate(); err != nil {
				return err
			}
		}
	case qpBirth:
		if op != nil && *op == opYear {
			if _, err := util.ParseInt(values[0]); err != nil {
				return err
			}
		}
		ts, err := util.ParseTimestamp(values[0])
		if err != nil {
			return err
		}

		return (*domain.FieldBirth)(&ts).Validate()
	case qpPremium:
		if op != nil {
			return validateBoolValue(values[0])
		}

		ts, err := util.ParseTimestamp(values[0])
		if err != nil {
			return err
		}

		return (*domain.FieldPremium)(&ts).Validate()
	case qpInterests:
		for _, value := range values {
			if err := (*domain.FieldInterest)(&value).Validate(); err != nil {
				return err
			}
		}

		return nil
	case qpLikes:
		for _, value := range values {
			intVal, err := util.ParseInt(value)
			if err != nil {
				return err
			}

			int32Val := int32(intVal)
			if err := (*domain.FieldID)(&int32Val).Validate(); err != nil {
				return err
			}
		}

		return nil
	case qpLimit:
		_, err := util.ParseInt(values[0])
		if err != nil {
			return err
		}
	default:
	}

	return nil
}

func validateBoolValue(val string) error {
	intVal, err := util.ParseInt(val)
	if err != nil {
		return err
	}

	if intVal != 0 && intVal != 1 {
		return errInvalidValue
	}

	return nil
}
