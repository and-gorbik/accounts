package service

import (
	"fmt"

	"accounts/app/repository"
	repo "accounts/app/repository"
)

var qpOnColumns = map[string]string{
	qpSex:       repo.AccountSex,
	qpEmail:     repo.AccountEmail,
	qpStatus:    repo.AccountStatus,
	qpFirstname: repo.AccountFirstname,
	qpSurname:   repo.AccountSurname,
	qpPhone:     repo.AccountPhone,
	qpCountry:   repo.CountryName,
	qpCity:      repo.CityName,
	qpBirth:     repo.AccountBirth,
	qpInterests: repo.InterestName,
	qpLikes:     repo.LikesLikerID,     // ?
	qpPremium:   repo.AccountPremStart, // ?
}

func BuildFilter(params map[string]QueryParam) (*repo.Filter, error) {
	filter := repository.NewFilter()
	limit := params[qpLimit].Values[0].(int)
	delete(params, qpLimit)
	filter.Limit = limit

	for _, param := range params {
		column := qpOnColumns[param.Field]

		if param.Op == nil {
			filter.Eq(column, param.Values[0])
			continue
		}

		if len(param.Values) == 0 {
			return nil, fmt.Errorf("param.Values is empty. Field: %s, Op: %#v", param.Field, param.Op)
		}

		switch *param.Op {
		case opEq:
			filter.Eq(column, param.Values[0])
		case opLt:
			filter.Lt(column, param.Values[0])
		case opGt:
			filter.Gt(column, param.Values[0])
		case opNeq:
			filter.Neq(column, param.Values[0])
		case opAny:
			filter.Any(column, param.Values)
		case opDomain:
			filter.Domain(column, param.Values[0])
		case opNull:
			filter.Null(column, param.Values[0].(bool))
		case opStarts:
			filter.Starts(column, param.Values[0])
		case opCode:
			filter.Code(column, param.Values[0])
		case opYear:
			filter.Year(column, param.Values[0])
		case opNow:
			filter.Now(column)
		case opContains:
			filter.Contains(column, param.Values)
		}
	}

	return filter, nil
}
