package service

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"accounts/util"
)

type testCaseQueryParam struct {
	Field    string
	StrValue string
	Expected QueryParam
}

var testGoodQueryParams = []testCaseQueryParam{
	{
		Field:    qpSex + "_" + opEq,
		StrValue: "m",
		Expected: QueryParam{
			Field:  qpSex,
			Values: []interface{}{"m"},
			Op:     util.PtrString(opEq),
		},
	},
	{
		Field:    qpEmail + "_" + opDomain,
		StrValue: "test.com",
		Expected: QueryParam{
			Field:  qpEmail,
			Values: []interface{}{"test.com"},
			Op:     util.PtrString(opDomain),
		},
	},
	{
		Field:    qpEmail + "_" + opGt,
		StrValue: "aaa",
		Expected: QueryParam{
			Field:  qpEmail,
			Values: []interface{}{"aaa"},
			Op:     util.PtrString(opGt),
		},
	},
	{
		Field:    qpEmail + "_" + opLt,
		StrValue: "zzz",
		Expected: QueryParam{
			Field:  qpEmail,
			Values: []interface{}{"zzz"},
			Op:     util.PtrString(opLt),
		},
	},
	{
		Field:    qpStatus + "_" + opEq,
		StrValue: "свободны",
		Expected: QueryParam{
			Field:  qpStatus,
			Values: []interface{}{"свободны"},
			Op:     util.PtrString(opEq),
		},
	},
	{
		Field:    qpStatus + "_" + opNeq,
		StrValue: "заняты",
		Expected: QueryParam{
			Field:  qpStatus,
			Values: []interface{}{"заняты"},
			Op:     util.PtrString(opNeq),
		},
	},
	{
		Field:    qpFirstname + "_" + opEq,
		StrValue: "Андрей",
		Expected: QueryParam{
			Field:  qpFirstname,
			Values: []interface{}{"Андрей"},
			Op:     util.PtrString(opEq),
		},
	},
	{
		Field:    qpFirstname + "_" + opAny,
		StrValue: "Андрей,Олег,Борис",
		Expected: QueryParam{
			Field:  qpFirstname,
			Values: []interface{}{"Андрей", "Олег", "Борис"},
			Op:     util.PtrString(opAny),
		},
	},
	{
		Field:    qpFirstname + "_" + opNull,
		StrValue: "1",
		Expected: QueryParam{
			Field:  qpFirstname,
			Values: []interface{}{true},
			Op:     util.PtrString(opNull),
		},
	},
	{
		Field:    qpSurname + "_" + opEq,
		StrValue: "Иванов",
		Expected: QueryParam{
			Field:  qpSurname,
			Values: []interface{}{"Иванов"},
			Op:     util.PtrString(opEq),
		},
	},
	{
		Field:    qpSurname + "_" + opStarts,
		StrValue: "Ива",
		Expected: QueryParam{
			Field:  qpSurname,
			Values: []interface{}{"Ива"},
			Op:     util.PtrString(opStarts),
		},
	},
	{
		Field:    qpSurname + "_" + opNull,
		StrValue: "0",
		Expected: QueryParam{
			Field:  qpSurname,
			Values: []interface{}{false},
			Op:     util.PtrString(opNull),
		},
	},
	{
		Field:    qpPhone + "_" + opNull,
		StrValue: "1",
		Expected: QueryParam{
			Field:  qpPhone,
			Values: []interface{}{true},
			Op:     util.PtrString(opNull),
		},
	},
	{
		Field:    qpPhone + "_" + opCode,
		StrValue: "999",
		Expected: QueryParam{
			Field:  qpPhone,
			Values: []interface{}{999},
			Op:     util.PtrString(opCode),
		},
	},
	{
		Field:    qpCountry + "_" + opEq,
		StrValue: "Россия",
		Expected: QueryParam{
			Field:  qpCountry,
			Values: []interface{}{"Россия"},
			Op:     util.PtrString(opEq),
		},
	},
	{
		Field:    qpCountry + "_" + opNull,
		StrValue: "1",
		Expected: QueryParam{
			Field:  qpCountry,
			Values: []interface{}{true},
			Op:     util.PtrString(opNull),
		},
	},
	{
		Field:    qpCity + "_" + opEq,
		StrValue: "Москва",
		Expected: QueryParam{
			Field:  qpCity,
			Values: []interface{}{"Москва"},
			Op:     util.PtrString(opEq),
		},
	},
	{
		Field:    qpCity + "_" + opAny,
		StrValue: "Москва,Питер",
		Expected: QueryParam{
			Field:  qpCity,
			Values: []interface{}{"Москва", "Питер"},
			Op:     util.PtrString(opAny),
		},
	},
	{
		Field:    qpCity + "_" + opNull,
		StrValue: "1",
		Expected: QueryParam{
			Field:  qpCity,
			Values: []interface{}{true},
			Op:     util.PtrString(opNull),
		},
	},
	{
		Field:    qpBirth + "_" + opLt,
		StrValue: strconv.FormatInt(time.Date(1994, 3, 24, 0, 0, 0, 0, time.Local).Unix(), 10),
		Expected: QueryParam{
			Field:  qpBirth,
			Values: []interface{}{time.Date(1994, 3, 24, 0, 0, 0, 0, time.Local).Unix()},
			Op:     util.PtrString(opLt),
		},
	},
	{
		Field:    qpBirth + "_" + opGt,
		StrValue: strconv.FormatInt(time.Date(1994, 3, 24, 0, 0, 0, 0, time.Local).Unix(), 10),
		Expected: QueryParam{
			Field:  qpBirth,
			Values: []interface{}{time.Date(1994, 3, 24, 0, 0, 0, 0, time.Local).Unix()},
			Op:     util.PtrString(opGt),
		},
	},
	{
		Field:    qpBirth + "_" + opYear,
		StrValue: "1994",
		Expected: QueryParam{
			Field:  qpBirth,
			Values: []interface{}{1994},
			Op:     util.PtrString(opYear),
		},
	},
	{
		Field:    qpPremium + "_" + opNull,
		StrValue: "1",
		Expected: QueryParam{
			Field:  qpPremium,
			Values: []interface{}{true},
			Op:     util.PtrString(opNull),
		},
	},
	{
		Field:    qpPremium + "_" + opNow,
		StrValue: "1",
		Expected: QueryParam{
			Field:  qpPremium,
			Values: []interface{}{true},
			Op:     util.PtrString(opNow),
		},
	},
	{
		Field:    qpInterests + "_" + opContains,
		StrValue: "фортепиано,фотография,дайвинг",
		Expected: QueryParam{
			Field:  qpInterests,
			Values: []interface{}{"фортепиано", "фотография", "дайвинг"},
			Op:     util.PtrString(opContains),
		},
	},
	{
		Field:    qpInterests + "_" + opAny,
		StrValue: "фортепиано,фотография,дайвинг",
		Expected: QueryParam{
			Field:  qpInterests,
			Values: []interface{}{"фортепиано", "фотография", "дайвинг"},
			Op:     util.PtrString(opAny),
		},
	},
	{
		Field:    qpLikes + "_" + opContains,
		StrValue: "1,2,3",
		Expected: QueryParam{
			Field:  qpLikes,
			Values: []interface{}{1, 2, 3},
			Op:     util.PtrString(opContains),
		},
	},
}

func Test_parseQueryParam_Success(t *testing.T) {
	for _, tc := range testGoodQueryParams {
		qp, err := parseQueryParam(tc.Field, strings.Split(tc.StrValue, ","), true)
		require.NoError(t, err)
		assert.Equal(t, tc.Expected, qp)
	}
}
