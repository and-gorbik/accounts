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
					Type:     typeStr,
					Field:    "sex",
					StrValue: "m",
					Op:       util.PtrString(opEq),
				},
			},
			Expected: fmt.Sprintf("sex = %s", escapeString("m")),
		},
		{
			Params: map[string]QueryParam{
				"email": {
					Type:     typeStr,
					Field:    "email",
					StrValue: "test@test.ru",
					Op:       util.PtrString(opLt),
				},
			},
			Expected: fmt.Sprintf("email < %s", escapeString("test@test.ru")),
		},
		{
			Params: map[string]QueryParam{
				"email": {
					Type:     typeStr,
					Field:    "email",
					StrValue: "test.ru",
					Op:       util.PtrString(opDomain),
				},
			},
			Expected: fmt.Sprintf("email LIKE %s", escapeString("%test.ru")),
		},
		{
			Params: map[string]QueryParam{
				"birth": {
					Type:     typeTimestamp,
					Field:    "birth",
					StrValue: "1485724260",
					Op:       util.PtrString(opGt),
				},
			},
			Expected: "birth > '2017-01-30 00:11:00'",
		},
		{
			Params: map[string]QueryParam{
				"status": {
					Type:     typeStr,
					Field:    "status",
					StrValue: "заняты",
					Op:       util.PtrString(opNeq),
				},
			},
			Expected: fmt.Sprintf("status != %s", escapeString("заняты")),
		},
		{
			Params: map[string]QueryParam{
				"fname": {
					Type:     typeInt,
					Field:    "fname",
					StrValue: "1",
					Op:       util.PtrString(opNull),
				},
			},
			Expected: fmt.Sprintf("fname IS NULL"),
		},
		{
			Params: map[string]QueryParam{
				"sname": {
					Type:     typeInt,
					Field:    "sname",
					StrValue: "0",
					Op:       util.PtrString(opNull),
				},
			},
			Expected: fmt.Sprintf("sname IS NOT NULL"),
		},
		{
			Params: map[string]QueryParam{
				"sname": {
					Type:     typeStr,
					Field:    "sname",
					StrValue: "Ан",
					Op:       util.PtrString(opStarts),
				},
			},
			Expected: fmt.Sprintf("sname LIKE %s", escapeString("Ан%")),
		},
		{
			Params: map[string]QueryParam{
				"phone": {
					Type:     typeInt,
					Field:    "phone",
					StrValue: "985",
					Op:       util.PtrString(opCode),
				},
			},
			Expected: fmt.Sprintf("phone LIKE %s", escapeString("%(985)%")),
		},
		{
			Params: map[string]QueryParam{
				"phone": {
					Type:     typeInt,
					Field:    "phone",
					StrValue: "985",
					Op:       util.PtrString(opCode),
				},
			},
			Expected: fmt.Sprintf("phone LIKE %s", escapeString("%(985)%")),
		},
		{
			Params: map[string]QueryParam{
				"city": {
					Type:     typeStr,
					Field:    "city",
					StrValue: "Москва,Питер,Новосибирск",
					Op:       util.PtrString(opAny),
				},
			},
			Expected: fmt.Sprintf("city = ANY (%s, %s, %s)", escapeString("Москва"), escapeString("Питер"), escapeString("Новосибирск")),
		},
		{
			Params: map[string]QueryParam{
				"likes": {
					Type:     typeInt,
					Field:    "likes",
					StrValue: "1,2,3",
					Op:       util.PtrString(opContains),
				},
			},
			Expected: "likes = ALL (1, 2, 3)",
		},
		{
			Params: map[string]QueryParam{
				"premium": {
					Type:     typeInt,
					Field:    "premium",
					StrValue: "1",
					Op:       util.PtrString(opNow),
				},
			},
			Expected: fmt.Sprintf("prem_start <= '%s' AND prem_end >= '%s'", Now, Now),
		},
		// TODO: test opYear
	}

	for _, tc := range testcases {
		assert.Equal(t, tc.Expected, BuildFilter(tc.Params).SQL)
	}
}
