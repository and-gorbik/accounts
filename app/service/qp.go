package service

import (
	"errors"
	"net/url"
	"strings"
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
)

var (
	errInvalidParam         = errors.New("invalid query param")
	errInvalidValue         = errors.New("invalid value")
	errEmptyValue           = errors.New("empty value")
	errMissingRequiredParam = errors.New("missing required param")
)

func ParseQueryParamsWithOp(qps url.Values) (map[string]QueryParamWithOp, error) {
	limit, err := parseLimit(qps)
	if err != nil {
		return nil, err
	}

	if _, ok := qps[qpQueryID]; !ok {
		return nil, errMissingRequiredParam
	}

	params := make(map[string]QueryParamWithOp, len(qps)-1)
	params[limit.Field] = QueryParamWithOp{
		Field:    limit.Field,
		Op:       opEq,
		StrValue: limit.StrValue,
	}

	for param, values := range qps {
		if param == qpLimit || param == qpQueryID {
			continue
		}

		qp, err := parseQueryParamWithOp(param, strings.Join(values, ","))
		if err != nil {
			return nil, err
		}

		params[qp.Field] = qp
	}

	return params, nil
}

func ParseQueryParams(qps url.Values) (map[string]QueryParam, error) {
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

		qp, err := parseQueryParam(param, strings.Join(values, ","))
		if err != nil {
			return nil, err
		}

		params[qp.Field] = qp
	}

	return params, nil
}

func parseQueryParam(param string, value string) (qp QueryParam, err error) {
	if _, ok := qpRules[param]; !ok {
		err = errInvalidParam
		return
	}

	if value == "" {
		err = errEmptyValue
		return
	}

	qp.Field = param
	qp.StrValue = value
	return
}

func parseQueryParamWithOp(param string, value string) (qp QueryParamWithOp, err error) {
	tokens := strings.Split(param, "_")
	if len(tokens) != 2 {
		err = errInvalidParam
		return
	}

	if _, ok := qpWithOpRules[tokens[0]][tokens[1]]; !ok {
		err = errInvalidParam
		return
	}

	if value == "" {
		err = errEmptyValue
		return
	}

	qp.Field = tokens[0]
	qp.Op = tokens[1]
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
