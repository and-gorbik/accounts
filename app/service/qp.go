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
	qpJoined    = "joined"

	qpLimit   = "limit"
	qpQueryID = "query_id"
)

var (
	yes = struct{}{}
)

var (
	qpWithOpRules = map[string]map[string]struct{}{
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

	qpRules = map[string]struct{}{
		qpEmail:     yes,
		qpSex:       yes,
		qpBirth:     yes,
		qpFirstname: yes,
		qpSurname:   yes,
		qpPhone:     yes,
		qpCountry:   yes,
		qpCity:      yes,
		qpJoined:    yes,
		qpStatus:    yes,
		qpInterests: yes,
		qpPremium:   yes,
		qpLikes:     yes,
	}

	qpTypes = map[string]int{
		qpSex:       util.TypeStr,
		qpEmail:     util.TypeStr,
		qpStatus:    util.TypeStr,
		qpFirstname: util.TypeStr,
		qpSurname:   util.TypeStr,
		qpPhone:     util.TypeStr,
		qpCountry:   util.TypeStr,
		qpCity:      util.TypeStr,
		qpBirth:     util.TypeTimestamp,
		qpInterests: util.TypeStrArray,
		qpLikes:     util.TypeIntArray,
		qpPremium:   util.TypeTimestamp,
		qpJoined:    util.TypeTimestamp,
		qpLimit:     util.TypeInt,
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

		qp, err := parseQueryParam(param, strings.Join(values, ","), withOp)
		if err != nil {
			return nil, err
		}

		if err := validateValues(qp.Field, values); err != nil {
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

		if _, ok := qpWithOpRules[tokens[0]][tokens[1]]; !ok {
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

	qp.StrValue = value
	qp.Type = qpTypes[param]

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

func validateValues(param string, values []string) error {
	if len(values) < 1 {
		return nil
	}

	switch param {
	case qpSex:
		return (*domain.FieldSex)(&values[0]).Validate()
	case qpEmail:
		return (*domain.FieldEmail)(&values[0]).Validate()
	case qpStatus:
		return (*domain.FieldStatus)(&values[0]).Validate()
	case qpFirstname:
		return (*domain.FieldFirstname)(&values[0]).Validate()
	case qpSurname:
		return (*domain.FieldSurname)(&values[0]).Validate()
	case qpPhone:
		return (*domain.FieldPhone)(&values[0]).Validate()
	case qpCountry:
		return (*domain.FieldCountry)(&values[0]).Validate()
	case qpCity:
		return (*domain.FieldCity)(&values[0]).Validate()
	case qpBirth:
		ts, err := util.ParseTimestamp(values[0])
		if err != nil {
			return err
		}

		return (*domain.FieldBirth)(&ts).Validate()
	case qpPremium:
		ts, err := util.ParseTimestamp(values[0])
		if err != nil {
			return err
		}

		return (*domain.FieldPremium)(&ts).Validate()
	case qpJoined:
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
