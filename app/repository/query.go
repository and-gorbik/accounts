package repository

import (
	"strconv"
	"strings"

	"accounts/domain"
	"accounts/util"
)

type Filter struct {
	Fields map[string]struct{}
	Limit  string
	SQL    string
}

func buildAccountSearchQuery(filter Filter) string {
	joins := make([]string, 0)
	resultFields := []string{"account.id", "email"}

	for field := range filter.Fields {
		switch field {
		case "sex", "status", "birth", "fname", "sname", "phone":
			resultFields = append(resultFields, field)
		case "city":
			resultFields = append(resultFields, "city.name")
			joins = append(joins, "JOIN city ON city.id = account.city_id")
		case "country":
			resultFields = append(resultFields, "country.name")
			joins = append(joins, "JOIN country ON country.id = account.country_id")
		case "likes":
			joins = append(joins, "JOIN like ON like.liker_id = account.id")
		case "interests":
			joins = append(joins, "JOIN interest ON interest.account_id = account.id")
		}
	}

	var b strings.Builder
	b.WriteString("SELECT ")
	b.WriteString(strings.Join(resultFields, ", "))
	b.WriteString(" FROM account ")
	b.WriteString(strings.Join(joins, " "))
	if filter.SQL != "" {
		b.WriteString(" WHERE ")
		b.WriteString(filter.SQL)
	}
	b.WriteString(" ORDER BY account.id DESC LIMIT ")
	b.WriteString(filter.Limit)

	return b.String()
}

func buildAccountUpdateQuery(a domain.AccountUpdate, cityID, countryID int32) (string, []interface{}) {
	fields := []string{}
	values := []interface{}{}

	if a.Email != nil {
		fields = append(fields, "email = $"+strconv.Itoa(len(fields)+1))
		values = append(values, string(*a.Email))
	}
	if a.Birth != nil {
		fields = append(fields, "birth = $"+strconv.Itoa(len(fields)+1))
		values = append(values, *util.TimestampToDatetime((*int64)(a.Birth)))
	}
	if a.Status != nil {
		fields = append(fields, "status = $"+strconv.Itoa(len(fields)+1))
		values = append(values, string(*a.Status))
	}
	fields = append(fields, "city_id = $"+strconv.Itoa(len(fields)+1))
	values = append(values, cityID)
	fields = append(fields, "country_id = $"+strconv.Itoa(len(fields)+1))
	values = append(values, countryID)
	values = append(values, a.ID)

	var builder strings.Builder
	builder.WriteString("UPDATE person SET ")
	builder.WriteString(strings.Join(fields, ", "))
	builder.WriteString(" WHERE id = $" + strconv.Itoa(len(fields)+1))

	return builder.String(), values
}
