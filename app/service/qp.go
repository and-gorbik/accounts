package service

import (
	"fmt"
	"net/url"
	"strings"

	"accounts/domain"
	"accounts/util"
)

type QueryParam struct {
	Field  string
	Values []interface{}
	Op     *string
}

type parserFunc = func(op *string, strValues []string) ([]interface{}, error)

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
	paramParsers = map[string]parserFunc{
		qpSex:       parseSex,
		qpEmail:     parseEmail,
		qpStatus:    parseStatus,
		qpFirstname: parseFirstname,
		qpSurname:   parseSurname,
		qpPhone:     parsePhone,
		qpCountry:   parseCountry,
		qpCity:      parseCity,
		qpBirth:     parseBirth,
		qpInterests: parseInterests,
		qpLikes:     parseLikes,
		qpPremium:   parsePremium,
	}
)

var (
	errInvalidParam         = "invalid query param: %s"
	errInvalidParamWithOp   = "invalid query param with operation: %s"
	errInvalidValue         = "invalid value: %v"
	errInvalidOp            = "invalid operation: %s"
	errEmptyValue           = "empty value"
	errValuesLen            = "invalid number of values: %d"
	errMissingRequiredParam = "missing required param: %s"
	errEmptyOp              = "empty operation"
)

func ParseQueryParams(qps url.Values, withOp bool) (map[string]QueryParam, error) {
	limit, err := parseLimit(qps)
	if err != nil {
		return nil, err
	}

	if _, ok := qps[qpQueryID]; !ok {
		return nil, fmt.Errorf(errMissingRequiredParam, qpQueryID)
	}

	params := make(map[string]QueryParam, len(qps)-1)
	params[limit.Field] = limit

	for param, values := range qps {
		if param == qpLimit || param == qpQueryID {
			continue
		}

		qp, err := parseQueryParam(param, values, withOp)
		if err != nil {
			return nil, err
		}

		params[qp.Field] = qp
	}

	return params, nil
}

func parseQueryParam(param string, strValues []string, withOp bool) (qp QueryParam, err error) {
	if !withOp {
		parser, ok := paramParsers[param]
		if !ok {
			return qp, fmt.Errorf(errInvalidParam, param)
		}

		qp.Field = param
		qp.Values, err = parser(nil, strValues)
		return
	}

	tokens := strings.Split(param, "_")
	if len(tokens) != 2 {
		return qp, fmt.Errorf(errInvalidParamWithOp, param)
	}

	qp.Field = tokens[0]
	qp.Op = &tokens[1]
	parser, ok := paramParsers[qp.Field]
	if !ok {
		return qp, fmt.Errorf(errInvalidParamWithOp, param)
	}
	qp.Values, err = parser(qp.Op, strValues)
	return
}

func parseLimit(qps url.Values) (QueryParam, error) {
	if _, ok := qps[qpLimit]; !ok {
		return QueryParam{}, fmt.Errorf(errMissingRequiredParam, qpLimit)
	}

	if len(qps[qpLimit]) != 1 {
		return QueryParam{}, fmt.Errorf(errInvalidValue, len(qps[qpLimit]))
	}

	limit, err := util.ParseInt(qps[qpLimit][0])
	if err != nil {
		return QueryParam{}, err
	}

	return QueryParam{
		Field:  qpLimit,
		Values: []interface{}{limit},
	}, nil
}

// opEq
func parseSex(op *string, strValues []string) ([]interface{}, error) {
	if op != nil && *op != opEq {
		return nil, fmt.Errorf(errInvalidOp, *op)
	}

	values, err := NewParser(strValues).SingleValue().String().Parse()
	if err != nil {
		return nil, err
	}

	sex := domain.FieldSex(values[0].(string))
	return values, sex.Validate()
}

// opEq, opAny, opNull
func parseFirstname(op *string, strValues []string) ([]interface{}, error) {
	if op == nil || (op != nil && *op == opEq) {
		values, err := NewParser(strValues).SingleValue().String().Parse()
		if err != nil {
			return nil, err
		}

		fname := domain.FieldFirstname(values[0].(string))
		if err = fname.Validate(); err != nil {
			return nil, err
		}

		return values, nil
	}

	switch *op {
	case opNull:
		return NewParser(strValues).SingleValue().Bool().Parse()

	case opAny:
		values, err := NewParser(strValues).String().Parse()
		if err != nil {
			return nil, err
		}

		for _, v := range values {
			fname := domain.FieldFirstname(v.(string))
			if err = fname.Validate(); err != nil {
				return nil, err
			}
		}

		return values, nil

	default:
		return nil, fmt.Errorf(errInvalidOp, *op)
	}
}

// opDomain, opGt, opLt
func parseEmail(op *string, strValues []string) ([]interface{}, error) {
	values, err := NewParser(strValues).SingleValue().String().Parse()
	if err != nil {
		return nil, err
	}

	if op == nil {
		email := domain.FieldEmail(values[0].(string))
		return values, email.Validate()
	}

	switch *op {
	case opDomain, opGt, opLt:
		return values, nil
	default:
		return nil, fmt.Errorf(errInvalidOp, *op)
	}
}

// opEq, opNeq
func parseStatus(op *string, strValues []string) ([]interface{}, error) {
	values, err := NewParser(strValues).SingleValue().String().Parse()
	if err != nil {
		return nil, err
	}

	email := domain.FieldStatus(values[0].(string))
	if err = email.Validate(); err != nil {
		return nil, err
	}

	if op == nil {
		return values, nil
	}

	switch *op {
	case opEq, opNeq:
		return values, nil
	default:
		return nil, fmt.Errorf(errInvalidOp, *op)
	}

}

// opEq, opStarts, opNull
func parseSurname(op *string, strValues []string) ([]interface{}, error) {
	if op == nil || (op != nil && *op == opEq) {
		values, err := NewParser(strValues).SingleValue().String().Parse()
		if err != nil {
			return nil, err
		}

		surname := domain.FieldSurname(values[0].(string))
		return values, surname.Validate()
	}

	switch *op {
	case opNull:
		return NewParser(strValues).SingleValue().Bool().Parse()
	case opStarts:
		return NewParser(strValues).SingleValue().String().Parse()
	default:
		return nil, fmt.Errorf(errInvalidOp, *op)
	}
}

// opCode, opNull
func parsePhone(op *string, strValues []string) ([]interface{}, error) {
	if op == nil {
		values, err := NewParser(strValues).SingleValue().String().Parse()
		if err != nil {
			return nil, err
		}

		phone := domain.FieldPhone(values[0].(string))
		return values, phone.Validate()
	}

	switch *op {
	case opCode:
		return NewParser(strValues).SingleValue().Int().Parse()
	case opNull:
		return NewParser(strValues).SingleValue().Bool().Parse()
	default:
		return nil, fmt.Errorf(errInvalidOp, *op)
	}
}

// opEq, opNull
func parseCountry(op *string, strValues []string) ([]interface{}, error) {
	if op == nil || (op != nil && *op == opEq) {
		values, err := NewParser(strValues).SingleValue().String().Parse()
		if err != nil {
			return nil, err
		}

		country := domain.FieldCountry(values[0].(string))
		return values, country.Validate()
	}

	switch *op {
	case opNull:
		return NewParser(strValues).SingleValue().Bool().Parse()
	default:
		return nil, fmt.Errorf(errInvalidOp, *op)
	}
}

// opEq, opAny, opNull
func parseCity(op *string, strValues []string) ([]interface{}, error) {
	if op == nil || (op != nil && *op == opEq) {
		values, err := NewParser(strValues).SingleValue().String().Parse()
		if err != nil {
			return nil, err
		}

		city := domain.FieldCity(values[0].(string))
		return values, city.Validate()
	}

	switch *op {
	case opNull:
		return NewParser(strValues).SingleValue().Bool().Parse()
	case opAny:
		values, err := NewParser(strValues).String().Parse()
		if err != nil {
			return nil, err
		}

		for _, v := range values {
			city := domain.FieldCity(v.(string))
			if err = city.Validate(); err != nil {
				return nil, err
			}
		}

		return values, nil
	default:
		return nil, fmt.Errorf(errInvalidOp, *op)
	}
}

// opLt, opGt, opYear
func parseBirth(op *string, strValues []string) ([]interface{}, error) {
	if op == nil {
		values, err := NewParser(strValues).SingleValue().Timestamp().Parse()
		if err != nil {
			return nil, err
		}

		birth := domain.FieldBirth(values[0].(int64))
		return values, birth.Validate()
	}

	switch *op {
	case opLt, opGt:
		values, err := NewParser(strValues).SingleValue().Timestamp().Parse()
		if err != nil {
			return nil, err
		}

		birth := domain.FieldBirth(values[0].(int64))
		return values, birth.Validate()
	case opYear:
		return NewParser(strValues).SingleValue().Int().Parse()
	default:
		return nil, fmt.Errorf(errInvalidOp, *op)
	}
}

// opNow, opNull
func parsePremium(op *string, strValues []string) ([]interface{}, error) {
	if op == nil {
		values, err := NewParser(strValues).SingleValue().Timestamp().Parse()
		if err != nil {
			return nil, err
		}

		premium := domain.FieldPremium(values[0].(int64))
		return values, premium.Validate()
	}

	switch *op {
	case opNow, opNull:
		return NewParser(strValues).SingleValue().Bool().Parse()
	default:
		return nil, fmt.Errorf(errInvalidOp, *op)
	}
}

// opContains, opAny
func parseInterests(op *string, strValues []string) ([]interface{}, error) {
	if op == nil {
		return nil, fmt.Errorf(errEmptyOp)
	}

	if op != nil && *op != opContains && *op != opAny {
		return nil, fmt.Errorf(errInvalidOp, *op)
	}

	values, err := NewParser(strValues).String().Parse()
	if err != nil {
		return nil, err
	}

	for _, v := range values {
		interest := domain.FieldInterest(v.(string))
		if err = interest.Validate(); err != nil {
			return nil, err
		}
	}

	return values, nil
}

// opContains
func parseLikes(op *string, strValues []string) ([]interface{}, error) {
	if op == nil {
		return nil, fmt.Errorf(errEmptyOp)
	}

	if op != nil && *op != opContains {
		return nil, fmt.Errorf(errInvalidOp, *op)
	}

	values, err := NewParser(strValues).Int().Parse()
	if err != nil {
		return nil, err
	}

	for _, v := range values {
		likeeID := domain.FieldID(v.(int))
		if err = likeeID.Validate(); err != nil {
			return nil, err
		}
	}

	return values, nil
}
