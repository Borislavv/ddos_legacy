package service

import "github.com/Borislavv/ddos/internal/shared/infrastructure/network/safehttp"

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

type IMeter interface {
	Start()
	CommitReq(req *safehttp.Req)
	Stop()
	Summary()
}
