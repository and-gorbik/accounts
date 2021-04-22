package controller

import (
	"errors"
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
)

var (
	yes = struct{}{}
)

var (
	qpRules = map[string]map[string]struct{}{
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
)

var (
	errInvalidParam = errors.New("invalid query param")
)

func parseQueryParam(param string, values ...string) (qp service.QueryParam, err error) {
	tokens := strings.Split(param, "_")
	if len(tokens) != 2 {
		err = errInvalidParam
		return
	}

	if _, ok := qpRules[tokens[0]][tokens[1]]; !ok {
		err = errInvalidParam
		return
	}

	qp.Left = tokens[0]
	qp.Op = tokens[1]
	return
}
