package repository

import (
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
)

type Filter struct {
	Limit int
	cols  map[string]struct{}
	ops   []squirrel.Sqlizer
}

func NewFilter() *Filter {
	return &Filter{
		ops:  []squirrel.Sqlizer{},
		cols: make(map[string]struct{}),
	}
}

func (f *Filter) Columns() map[string]struct{} {
	return f.cols
}

func (f *Filter) Build() (string, []interface{}, error) {
	predicates := make([]string, 0, len(f.ops))
	totalValues := make([]interface{}, 0)
	for _, op := range f.ops {
		sql, values, err := op.ToSql()
		if err != nil {
			return "", nil, err
		}

		predicates = append(predicates, sql)
		totalValues = append(totalValues, values...)
	}

	sql, err := squirrel.Dollar.ReplacePlaceholders(strings.Join(predicates, " AND "))
	if err != nil {
		return "", nil, err
	}

	return sql, totalValues, nil
}

func (f *Filter) Eq(column string, value interface{}) {
	f.ops = append(f.ops, squirrel.Eq{column: value})
	f.cols[column] = struct{}{}
}

func (f *Filter) Neq(column string, value interface{}) {
	f.ops = append(f.ops, squirrel.NotEq{column: value})
	f.cols[column] = struct{}{}
}

func (f *Filter) Like(column string, value interface{}) {
	f.ops = append(f.ops, squirrel.Like{column: value})
	f.cols[column] = struct{}{}
}

func (f *Filter) Lt(column string, value interface{}) {
	f.ops = append(f.ops, squirrel.Lt{column: value})
	f.cols[column] = struct{}{}
}

func (f *Filter) Gt(column string, value interface{}) {
	f.ops = append(f.ops, squirrel.Gt{column: value})
	f.cols[column] = struct{}{}
}

func (f *Filter) Any(column string, values []interface{}) {
	f.ops = append(f.ops, &opAny{
		Field:  column,
		Values: values,
	})
	f.cols[column] = struct{}{}
}

func (f *Filter) Contains(column string, values []interface{}) {
	f.ops = append(f.ops, &opAll{
		Field:  column,
		Values: values,
	})
	f.cols[column] = struct{}{}
}

func (f *Filter) Null(column string, isNull bool) {
	var op squirrel.Sqlizer
	if isNull {
		op = squirrel.Eq{column: nil}
	} else {
		op = squirrel.NotEq{column: nil}
	}

	f.ops = append(f.ops, op)
	f.cols[column] = struct{}{}
}

func (f *Filter) Starts(column string, value interface{}) {
	f.ops = append(f.ops, squirrel.Like{column: fmt.Sprintf("%v%%", value)})
	f.cols[column] = struct{}{}
}

func (f *Filter) Domain(column string, value interface{}) {
	f.ops = append(f.ops, squirrel.Like{column: fmt.Sprintf("%%@%v", value)})
	f.cols[column] = struct{}{}
}

func (f *Filter) Code(column string, value interface{}) {
	f.ops = append(f.ops, squirrel.Like{column: fmt.Sprintf("%%(%v)%%", value)})
	f.cols[column] = struct{}{}
}

func (f *Filter) Now(value interface{}) {
	f.ops = append(f.ops, squirrel.And{
		squirrel.LtOrEq{AccountPremStart: value},
		squirrel.GtOrEq{AccountPremEnd: value},
	})
}

func (f *Filter) Year(column string, value interface{}) {
	// TODO: add operation
	f.cols[column] = struct{}{}
}

type opAny struct {
	Field  string
	Values []interface{}
}

func (op *opAny) ToSql() (string, []interface{}, error) {
	var b strings.Builder
	b.WriteString(op.Field)
	b.WriteString(" = ANY(")
	b.WriteString(squirrel.Placeholders(len(op.Values)))
	b.WriteRune(')')

	return b.String(), op.Values, nil
}

type opAll struct {
	Field  string
	Values []interface{}
}

func (op *opAll) ToSql() (string, []interface{}, error) {
	var b strings.Builder
	b.WriteString(op.Field)
	b.WriteString(" = ALL(")
	b.WriteString(squirrel.Placeholders(len(op.Values)))
	b.WriteRune(')')

	return b.String(), op.Values, nil
}
