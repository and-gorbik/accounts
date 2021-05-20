package service

import (
	"fmt"
	"math/rand"
	"strings"

	"accounts/util"
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

	var itemType int
	if p.Type == util.TypeStrArray {
		itemType = util.TypeStr
	} else {
		itemType = util.TypeInt
	}

	var b strings.Builder
	b.WriteString(p.Field)
	b.WriteString(" = ")
	b.WriteString(sqlOp)
	b.WriteString(" (")
	b.WriteString(buildValue(values[0], itemType))
	for i := 1; i < len(values); i++ {
		b.WriteString(", ")
		b.WriteString(buildValue(values[i], itemType))
	}

	b.WriteByte(')')
	return b.String()
}

func buildValue(value string, typ int) string {
	switch typ {
	case util.TypeStr:
		return escapeString(value)
	case util.TypeTimestamp:
		ts, _ := util.ParseTimestamp(value)
		return fmt.Sprintf("'%s'", util.TimestampToDatetime(ts).Format(TimeLayout))
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
