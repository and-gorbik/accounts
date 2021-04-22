package controller

import (
	"testing"
)

func Test_parseQueryParam(t *testing.T) {
	qp, err := parseQueryParam("sex_eq")
	if err != nil {
		t.Fatal(err)
	}

	if qp.Left != qpSex || qp.Op != eq {
		t.Fatal(errInvalidParam)
	}
}
