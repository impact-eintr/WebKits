package esvc

import (
	"context"
	"os/signal"
)

var signalNotify = signal.Notify

type Service interface {
	Init() error
	Start() error
	Stop() error
}

type Context interface {
	Context() context.Context
}
