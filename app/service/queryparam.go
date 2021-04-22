package service

type action = func()

type QueryParam struct {
	Field    string
	StrValue string
}

type QueryParamWithOp struct {
	Field    string
	StrValue string
	Op       string
	Action   action
}
