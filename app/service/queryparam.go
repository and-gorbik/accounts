package service

type QueryParam struct {
	Type     int
	Field    string
	StrValue string
}

type QueryParamWithOp struct {
	Type     int
	Field    string
	StrValue string
	Op       string
}
