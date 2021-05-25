package service

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"accounts/app/repository"
	"accounts/util"
)

var (
	randScope = fmt.Sprintf("%%%d%%", rand.Int())
	Now       = time.Now().Format((TimeLayout))
)

// TODO: move to repository
func BuildFilter(params map[string]QueryParam) repository.Filter {
	limit := params[qpLimit].StrValue
	delete(params, qpLimit)

	filters := make([]string, 0, len(params))
	fields := make(map[string]struct{}, len(params)-1)
	for _, param := range params {
		filters = append(filters, buildFilter(param))
		if param.Field != qpLimit {
			fields[param.Field] = struct{}{}
		}
	}

	return repository.Filter{
		SQL:    strings.Join(filters, " AND "),
		Limit:  limit,
		Fields: fields,
	}
}

func buildFilter(param QueryParam) string {
	if param.Op == nil {
		eq := opEq
		param.Op = &eq
	}

	switch *param.Op {
	case opEq, opLt, opGt, opNeq:
		return buildBinaryOp(param)
	case opDomain:
		return buildOpLike(param.Field, "%", param.StrValue)
	case opNull:
		return buildOpNull(param)
	case opStarts:
		return buildOpLike(param.Field, param.StrValue, "%")
	case opCode:
		return buildOpLike(param.Field, "%(", param.StrValue, ")%")
	case opAny:
		return buildOpAny(param)
	case opContains:
		return buildOpContains(param)
	case opNow:
		return buildOpNow()
	case opYear:
		// TODO
	}

	return ""
}

func buildBinaryOp(param QueryParam) string {
	var builder strings.Builder
	builder.WriteString(param.Field)

	switch *param.Op {
	case opEq:
		builder.WriteString(" = ")
	case opLt:
		builder.WriteString(" < ")
	case opGt:
		builder.WriteString(" > ")
	case opNeq:
		builder.WriteString(" != ")
	}

	builder.WriteString(buildValue(param.StrValue, param.Type))
	return builder.String()
}

func buildOpNull(param QueryParam) string {
	if param.StrValue == "1" {
		return param.Field + " IS NULL"
	}

	return param.Field + " IS NOT NULL"
}

func buildOpContains(p QueryParam) string {
	return buildOpOverArray(p, "ALL")
}

func buildOpAny(p QueryParam) string {
	return buildOpOverArray(p, "ANY")
}

func buildOpOverArray(p QueryParam, sqlOp string) string {
	values := strings.Split(p.StrValue, ",")

	var b strings.Builder
	b.WriteString(p.Field)
	b.WriteString(" = ")
	b.WriteString(sqlOp)
	b.WriteString(" (")
	b.WriteString(buildValue(values[0], p.Type))
	for i := 1; i < len(values); i++ {
		b.WriteString(", ")
		b.WriteString(buildValue(values[i], p.Type))
	}

	b.WriteByte(')')
	return b.String()
}

func buildValue(value string, typ int) string {
	switch typ {
	case typeStr:
		return escapeString(value)
	case typeTimestamp:
		ts, _ := util.ParseTimestamp(value)
		return wrapDatetime(util.TimestampToDatetime(&ts).Format(TimeLayout))
	default:
	}

	return value
}

func escapeString(s string) string {
	var b strings.Builder
	b.WriteString(randScope)
	b.WriteString(s)
	b.WriteString(randScope)
	return b.String()
}

func wrapDatetime(dt string) string {
	var b strings.Builder
	b.WriteString("'")
	b.WriteString(dt)
	b.WriteString("'")
	return b.String()
}

func buildOpLike(field string, values ...string) string {
	var b strings.Builder
	b.WriteString(field)
	b.WriteString(" LIKE ")
	b.WriteString(randScope)

	for _, val := range values {
		b.WriteString(val)
	}

	b.WriteString(randScope)

	return b.String()
}

func buildOpNow() string {
	var b strings.Builder
	b.WriteString("prem_start <= ")
	b.WriteString(wrapDatetime(Now))
	b.WriteString(" AND prem_end >= ")
	b.WriteString(wrapDatetime(Now))
	return b.String()
}
