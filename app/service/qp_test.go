package service

import (
	"testing"
)

func Test_parseQueryParamWithOp(t *testing.T) {
	qp, err := parseQueryParamWithOp("sex_eq", "m")
	if err != nil {
		t.Fatal(err)
	}

	if qp.Field != qpSex || qp.Op != opEq || qp.StrValue != "m" {
		t.Fatal(errInvalidParam)
	}
}
