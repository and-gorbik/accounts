package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"accounts/util"
)

type testCaseFilter struct {
	Params   map[string]QueryParam
	Expected string
}

func Test_BuildFilter(t *testing.T) {
	testcases := []testCaseFilter{
		{
			Params: map[string]QueryParam{
				"sex": {
					Field:  "sex",
					Values: []interface{}{"m"},
					Op:     util.PtrString(opEq),
				},
			},
			Expected: "account.sex = $1",
		},
		{
			Params: map[string]QueryParam{
				"email": {
					Field:  "email",
					Values: []interface{}{"test@test.ru"},
					Op:     util.PtrString(opLt),
				},
			},
			Expected: "account.email < $1",
		},
		{
			Params: map[string]QueryParam{
				"email": {
					Field:  "email",
					Values: []interface{}{"test.ru"},
					Op:     util.PtrString(opDomain),
				},
			},
			Expected: fmt.Sprintf("account.email LIKE $1"),
		},
		{
			Params: map[string]QueryParam{
				"birth": {
					Field:  "birth",
					Values: []interface{}{"2005-05-05 00:00:00"},
					Op:     util.PtrString(opGt),
				},
			},
			Expected: "account.birth > $1",
		},
		{
			Params: map[string]QueryParam{
				"status": {
					Field:  "status",
					Values: []interface{}{"заняты"},
					Op:     util.PtrString(opNeq),
				},
			},
			Expected: "account.status <> $1",
		},
		{
			Params: map[string]QueryParam{
				"fname": {
					Field:  "fname",
					Values: []interface{}{true},
					Op:     util.PtrString(opNull),
				},
			},
			Expected: "account.name IS NULL",
		},
		{
			Params: map[string]QueryParam{
				"sname": {
					Field:  "sname",
					Values: []interface{}{false},
					Op:     util.PtrString(opNull),
				},
			},
			Expected: "account.surname IS NOT NULL",
		},
		{
			Params: map[string]QueryParam{
				"sname": {
					Field:  "sname",
					Values: []interface{}{"Ан"},
					Op:     util.PtrString(opStarts),
				},
			},
			Expected: "account.surname LIKE $1",
		},
		{
			Params: map[string]QueryParam{
				"phone": {
					Field:  "phone",
					Values: []interface{}{"985"},
					Op:     util.PtrString(opCode),
				},
			},
			Expected: "account.phone LIKE $1",
		},
		{
			Params: map[string]QueryParam{
				"phone": {
					Field:  "phone",
					Values: []interface{}{"985"},
					Op:     util.PtrString(opCode),
				},
			},
			Expected: "account.phone LIKE $1",
		},
		{
			Params: map[string]QueryParam{
				"city": {
					Field:  "city",
					Values: []interface{}{"Москва", "Питер", "Новосибирск"},
					Op:     util.PtrString(opAny),
				},
			},
			Expected: "city.name IN ($1,$2,$3)",
		},
		{
			Params: map[string]QueryParam{
				"likes": {
					Field:  "likes",
					Values: []interface{}{1, 2, 3},
					Op:     util.PtrString(opContains),
				},
			},
			Expected: "likes.liker_id IN ($1,$2,$3)",
		},
		{
			Params: map[string]QueryParam{
				"premium": {
					Field:  "premium",
					Values: []interface{}{1},
					Op:     util.PtrString(opNow),
				},
			},
			Expected: "(account.prem_start <= $1 AND account.prem_end >= $2)",
		},
		// TODO: test opYear
	}

	for _, tc := range testcases {
		tc.Params[qpLimit] = QueryParam{
			Field:  qpLimit,
			Values: []interface{}{10},
		}

		filter, err := BuildFilter(tc.Params)
		if err != nil {
			t.Fatal(err)
		}

		sql, _, err := filter.Build()
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, tc.Expected, sql)
	}
}
