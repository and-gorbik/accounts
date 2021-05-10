package service

import (
	"fmt"
	"math/rand"
	"strings"

	"accounts/infrastructure"
)

var (
	randScope = fmt.Sprintf("%%%d%%", rand.Int())
)

func BuildFilter(params map[string]QueryParam) string {
	filters := make([]string, 0, len(params))
	for _, param := range params {
		filters = append(filters, buildFilter(param))
	}

	return strings.Join(filters, " AND ")
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
		return fmt.Sprintf("%s LIKE '%%@%s'", param.Field, param.StrValue)
	case opNull:
		return buildOpNull(param)
	case opStarts:
		return fmt.Sprintf("%s LIKE '%s%%'", param.Field, param.StrValue)
	case opCode:
		return fmt.Sprintf("%s LIKE '%%(%s)%%'", param.Field, param.StrValue)
	case opAny:
		return buildOpAny(param)
	case opContains:
		return buildOpContains(param)
	case opYear:
		// TODO
	case opNow:
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
		return fmt.Sprintf("%s IS NULL", param.Field)
	}

	return fmt.Sprintf("%s IS NOT NULL", param.Field)
}

func buildOpContains(p QueryParam) string {
	return buildOpOverArray(p, "ALL")
}

func buildOpAny(p QueryParam) string {
	return buildOpOverArray(p, "ANY")
}

func buildOpOverArray(p QueryParam, sqlOp string) string {
	values := strings.Split(p.StrValue, ",")
	if len(values) == 1 {
		return p.StrValue
	}

	var itemType int
	if p.Type == infrastructure.TypeStrArray {
		itemType = infrastructure.TypeStr
	} else {
		itemType = infrastructure.TypeInt
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("%s = %s (", p.Field, sqlOp))
	b.WriteString(buildValue(values[0], itemType))
	for i := 1; i < len(values); i++ {
		b.WriteString(", ")
		b.WriteString(buildValue(values[0], itemType))
	}

	b.WriteByte(')')
	return b.String()
}

func buildValue(value string, typ int) string {
	switch typ {
	case infrastructure.TypeStr:
		return escapeString(value)
	case infrastructure.TypeTimestamp:
		ts, _ := infrastructure.ParseTimestamp(value)
		return fmt.Sprintf("'%s'", infrastructure.TimestampToDatetime(ts).Format(TimeLayout))
	default:
	}

	return value
}

func escapeString(s string) string {
	return fmt.Sprintf("%s%s%s", randScope, s, randScope)
}
