package service

type validator interface {
	Validate() error
}
