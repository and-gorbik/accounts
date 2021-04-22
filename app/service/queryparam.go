package service

type action = func()

type QueryParam struct {
	Left   string
	Right  string
	Op     string
	Action action
}
