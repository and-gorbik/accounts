package controller

import (
	"errors"
	"net/url"
	"strings"

	"accounts/app/service"
)

const (
	eq       = "eq"
	lt       = "lt"
	gt       = "gt"
	neq      = "neq"
	any      = "any"
	domain   = "domain"
	null     = "null"
	starts   = "starts"
	code     = "code"
	year     = "year"
	contains = "contains"
	now      = "now"
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
		qpSex:       {eq: yes},
		qpEmail:     {domain: yes, gt: yes, lt: yes},
		qpStatus:    {eq: yes, neq: yes},
		qpFirstname: {eq: yes, any: yes, null: yes},
		qpSurname:   {eq: yes, starts: yes, null: yes},
		qpPhone:     {code: yes, null: yes},
		qpCountry:   {eq: yes, null: yes},
		qpCity:      {eq: yes, any: yes, null: yes},
		qpBirth:     {lt: yes, gt: yes, year: yes},
		qpInterests: {contains: yes, any: yes},
		qpLikes:     {contains: yes},
		qpPremium:   {now: yes, null: yes},
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

func ParseQueryParamsWithOp(qps url.Values) ([]service.QueryParamWithOp, error) {
	limit, err := parseLimit(qps)
	if err != nil {
		return nil, err
	}

	if _, ok := qps[qpQueryID]; !ok {
		return nil, errMissingRequiredParam
	}

	params := []service.QueryParamWithOp{
		{Field: limit.Field, Op: eq, StrValue: limit.StrValue},
	}

	for param, values := range qps {
		if param == qpLimit || param == qpQueryID {
			continue
		}

		qp, err := parseQueryParamWithOp(param, strings.Join(values, ","))
		if err != nil {
			return nil, err
		}

		params = append(params, qp)
	}

	return params, nil
}

func ParseQueryParams(qps url.Values) ([]service.QueryParam, error) {
	limit, err := parseLimit(qps)
	if err != nil {
		return nil, err
	}

	if _, ok := qps[qpQueryID]; !ok {
		return nil, errMissingRequiredParam
	}

	params := []service.QueryParam{limit}

	for param, values := range qps {
		if param == qpLimit || param == qpQueryID {
			continue
		}

		qp, err := parseQueryParam(param, strings.Join(values, ","))
		if err != nil {
			return nil, err
		}

		params = append(params, qp)
	}

	return params, nil
}

func parseQueryParam(param string, value string) (qp service.QueryParam, err error) {
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

func parseQueryParamWithOp(param string, value string) (qp service.QueryParamWithOp, err error) {
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

func parseLimit(qps url.Values) (service.QueryParam, error) {
	if _, ok := qps[qpLimit]; !ok {
		return service.QueryParam{}, errMissingRequiredParam
	}

	if len(qps[qpLimit]) != 1 {
		return service.QueryParam{}, errInvalidValue
	}

	return service.QueryParam{
		Field:    qpLimit,
		StrValue: qps[qpLimit][0],
	}, nil
}
