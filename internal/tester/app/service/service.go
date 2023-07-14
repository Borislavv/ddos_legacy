package service

type ITester interface {
	Start()
	Stop()
}

type IConsumer interface {
	Consume()
	Stop()
}

type IProvider interface {
	Provide()
	Stop()
}
type IDisplayer interface {
	Start()
	Display(pattern string, args ...interface{})
	DisplayError(err error)
	Stop()
}
