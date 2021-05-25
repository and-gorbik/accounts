package service

import (
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"accounts/util"
)

var testGoodQueryParams = []QueryParam{
	{
		Field:    qpSex,
		Op:       util.PtrString(opEq),
		Type:     typeStr,
		StrValue: "m",
	},
	{
		Field:    qpEmail,
		Op:       util.PtrString(opDomain),
		Type:     typeStr,
		StrValue: "test.com",
	},
	{
		Field:    qpEmail,
		Op:       util.PtrString(opGt),
		Type:     typeStr,
		StrValue: "aaa",
	},
	{
		Field:    qpEmail,
		Op:       util.PtrString(opLt),
		Type:     typeStr,
		StrValue: "zzz",
	},
	{
		Field:    qpEmail,
		Op:       nil,
		Type:     typeStr,
		StrValue: "test@test.com",
	},
	{
		Field:    qpStatus,
		Op:       util.PtrString(opEq),
		Type:     typeStr,
		StrValue: "свободны",
	},
	{
		Field:    qpStatus,
		Op:       util.PtrString(opNeq),
		Type:     typeStr,
		StrValue: "заняты",
	},
	{
		Field:    qpFirstname,
		Op:       util.PtrString(opEq),
		Type:     typeStr,
		StrValue: "Андрей",
	},
	{
		Field:    qpFirstname,
		Op:       nil,
		Type:     typeStr,
		StrValue: "Андрей",
	},
	{
		Field:    qpFirstname,
		Op:       util.PtrString(opAny),
		Type:     typeStr,
		StrValue: "Андрей,Олег,Борис",
	},
	{
		Field:    qpFirstname,
		Op:       util.PtrString(opNull),
		Type:     typeInt,
		StrValue: "1",
	},
	{
		Field:    qpSurname,
		Op:       util.PtrString(opEq),
		Type:     typeStr,
		StrValue: "Иванов",
	},
	{
		Field:    qpSurname,
		Op:       nil,
		Type:     typeStr,
		StrValue: "Иванов",
	},
	{
		Field:    qpSurname,
		Op:       util.PtrString(opStarts),
		Type:     typeStr,
		StrValue: "Ива",
	},
	{
		Field:    qpSurname,
		Op:       util.PtrString(opNull),
		Type:     typeInt,
		StrValue: "0",
	},
	{
		Field:    qpPhone,
		Op:       util.PtrString(opNull),
		Type:     typeInt,
		StrValue: "1",
	},
	{
		Field:    qpPhone,
		Op:       util.PtrString(opCode),
		Type:     typeInt,
		StrValue: "999",
	},
	{
		Field:    qpPhone,
		Op:       nil,
		Type:     typeStr,
		StrValue: "8(999)7654321",
	},
	{
		Field:    qpCountry,
		Op:       nil,
		Type:     typeStr,
		StrValue: "Россия",
	},
	{
		Field:    qpCountry,
		Op:       util.PtrString(opEq),
		Type:     typeStr,
		StrValue: "Россия",
	},
	{
		Field:    qpCountry,
		Op:       util.PtrString(opNull),
		Type:     typeInt,
		StrValue: "1",
	},
	{
		Field:    qpCity,
		Op:       nil,
		Type:     typeStr,
		StrValue: "Москва",
	},
	{
		Field:    qpCity,
		Op:       util.PtrString(opEq),
		Type:     typeStr,
		StrValue: "Москва",
	},
	{
		Field:    qpCity,
		Op:       util.PtrString(opAny),
		Type:     typeStr,
		StrValue: "Москва,Питер",
	},
	{
		Field:    qpCity,
		Op:       util.PtrString(opNull),
		Type:     typeInt,
		StrValue: "1",
	},
	{
		Field:    qpBirth,
		Op:       util.PtrString(opLt),
		Type:     typeTimestamp,
		StrValue: strconv.FormatInt(time.Date(1994, 3, 24, 0, 0, 0, 0, time.Local).Unix(), 10),
	},
	{
		Field:    qpBirth,
		Op:       util.PtrString(opGt),
		Type:     typeTimestamp,
		StrValue: strconv.FormatInt(time.Date(1994, 3, 24, 0, 0, 0, 0, time.Local).Unix(), 10),
	},
	{
		Field:    qpBirth,
		Op:       util.PtrString(opYear),
		Type:     typeInt,
		StrValue: "1994",
	},
	{
		Field:    qpBirth,
		Op:       nil,
		Type:     typeTimestamp,
		StrValue: strconv.FormatInt(time.Date(1994, 3, 24, 0, 0, 0, 0, time.Local).Unix(), 10),
	},
	{
		Field:    qpPremium,
		Op:       nil,
		Type:     typeTimestamp,
		StrValue: strconv.FormatInt(time.Date(2018, 1, 1, 0, 0, 0, 0, time.Local).Unix(), 10),
	},
	{
		Field:    qpPremium,
		Op:       util.PtrString(opNow),
		Type:     typeInt,
		StrValue: "1",
	},
	{
		Field:    qpPremium,
		Op:       util.PtrString(opNull),
		Type:     typeInt,
		StrValue: "1",
	},
	{
		Field:    qpInterests,
		Op:       util.PtrString(opContains),
		Type:     typeStr,
		StrValue: "фортепиано,фотография,дайвинг",
	},
	{
		Field:    qpInterests,
		Op:       util.PtrString(opAny),
		Type:     typeStr,
		StrValue: "фортепиано,фотография,дайвинг",
	},
	{
		Field:    qpLikes,
		Op:       util.PtrString(opContains),
		Type:     typeInt,
		StrValue: "1,2,3",
	},
}

func Test_parseQueryParam_Success(t *testing.T) {
	qp, err := parseQueryParam("sex_eq", "m", true)
	if err != nil {
		t.Fatal(err)
	}

	if qp.Field != qpSex || *qp.Op != opEq || qp.StrValue != "m" {
		t.Fatal(errInvalidParam)
	}
}

func Test_validateValues_Success(t *testing.T) {
	for _, qp := range testGoodQueryParams {
		assert.Nil(t, validateValues(qp.Field, qp.Op, strings.Split(qp.StrValue, ",")))
	}
}
