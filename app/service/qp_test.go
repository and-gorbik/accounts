package service

import (
	"testing"
)

func Test_parseQueryParamWithOp(t *testing.T) {
	qp, err := parseQueryParam("sex_eq", "m", true)
	if err != nil {
		t.Fatal(err)
	}

	if qp.Field != qpSex || *qp.Op != opEq || qp.StrValue != "m" {
		t.Fatal(errInvalidParam)
	}
}
