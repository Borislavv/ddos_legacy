package service

type ITester interface {
	StartTest()
	StopTest()
}

type IConsumer interface {
	StartConsuming()
	StopConsuming()
}

type IProvider interface {
	StartProviding()
	StopProviding()
}
