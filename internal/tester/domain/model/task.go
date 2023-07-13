package model

type Task struct {
	request *Request
}

func NewTask(request *Request) *Task {
	return &Task{request: request}
}
