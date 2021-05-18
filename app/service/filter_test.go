package service

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"accounts/infrastructure"
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
					Type:     infrastructure.TypeStr,
					Field:    "sex",
					StrValue: "m",
					Op:       infrastructure.PtrString(opEq),
				},
			},
			Expected: fmt.Sprintf("sex = %s", escapeString("m")),
		},
		{
			Params: map[string]QueryParam{
				"email": {
					Type:     infrastructure.TypeStr,
					Field:    "email",
					StrValue: "test@test.ru",
					Op:       infrastructure.PtrString(opLt),
				},
			},
			Expected: fmt.Sprintf("email < %s", escapeString("test@test.ru")),
		},
		{
			Params: map[string]QueryParam{
				"email": {
					Type:     infrastructure.TypeStr,
					Field:    "email",
					StrValue: "test.ru",
					Op:       infrastructure.PtrString(opDomain),
				},
			},
			Expected: fmt.Sprintf("email LIKE %s", escapeString("%test.ru")),
		},
		{
			Params: map[string]QueryParam{
				"birth": {
					Type:     infrastructure.TypeTimestamp,
					Field:    "birth",
					StrValue: "1485724260",
					Op:       infrastructure.PtrString(opGt),
				},
			},
			Expected: "birth > '2017-01-30 00:11:00'",
		},
		{
			Params: map[string]QueryParam{
				"status": {
					Type:     infrastructure.TypeStr,
					Field:    "status",
					StrValue: "заняты",
					Op:       infrastructure.PtrString(opNeq),
				},
			},
			Expected: fmt.Sprintf("status != %s", escapeString("заняты")),
		},
		{
			Params: map[string]QueryParam{
				"fname": {
					Type:     infrastructure.TypeInt,
					Field:    "fname",
					StrValue: "1",
					Op:       infrastructure.PtrString(opNull),
				},
			},
			Expected: fmt.Sprintf("fname IS NULL"),
		},
		{
			Params: map[string]QueryParam{
				"sname": {
					Type:     infrastructure.TypeInt,
					Field:    "sname",
					StrValue: "0",
					Op:       infrastructure.PtrString(opNull),
				},
			},
			Expected: fmt.Sprintf("sname IS NOT NULL"),
		},
		{
			Params: map[string]QueryParam{
				"sname": {
					Type:     infrastructure.TypeStr,
					Field:    "sname",
					StrValue: "Ан",
					Op:       infrastructure.PtrString(opStarts),
				},
			},
			Expected: fmt.Sprintf("sname LIKE %s", escapeString("Ан%")),
		},
		{
			Params: map[string]QueryParam{
				"phone": {
					Type:     infrastructure.TypeInt,
					Field:    "phone",
					StrValue: "985",
					Op:       infrastructure.PtrString(opCode),
				},
			},
			Expected: fmt.Sprintf("phone LIKE %s", escapeString("%(985)%")),
		},
		{
			Params: map[string]QueryParam{
				"phone": {
					Type:     infrastructure.TypeInt,
					Field:    "phone",
					StrValue: "985",
					Op:       infrastructure.PtrString(opCode),
				},
			},
			Expected: fmt.Sprintf("phone LIKE %s", escapeString("%(985)%")),
		},
		{
			Params: map[string]QueryParam{
				"city": {
					Type:     infrastructure.TypeStrArray,
					Field:    "city",
					StrValue: "Москва,Питер,Новосибирск",
					Op:       infrastructure.PtrString(opAny),
				},
			},
			Expected: fmt.Sprintf("city = ANY (%s, %s, %s)", escapeString("Москва"), escapeString("Питер"), escapeString("Новосибирск")),
		},
		{
			Params: map[string]QueryParam{
				"likes": {
					Type:     infrastructure.TypeIntArray,
					Field:    "likes",
					StrValue: "1,2,3",
					Op:       infrastructure.PtrString(opContains),
				},
			},
			Expected: "likes = ALL (1, 2, 3)",
		},
		// TODO: test opNow и opYear
	}

	for _, tc := range testcases {
		assert.Equal(t, tc.Expected, BuildFilter(tc.Params))
	}
}
