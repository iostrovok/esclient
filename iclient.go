package esclient

import (
	"context"
	"time"

	"github.com/olivere/elastic/v7"
)

type ILogger interface {
	Fatal(v ...interface{})
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}

type IDebug interface {
	// Request  returns full http request
	Request() []byte

	// Response returns full http response
	Response() []byte
}

type IConn interface {
	Open(bool, ...context.Context) IClient
	Sniff(context.Context)
	SetLogger(ILogger)
	SniffTimeout(time.Duration)

	// SetCustomHandler is a setter
	SetCustomHandler(ReqHandler, ResHandler) error

	// internal
	reConnect() error
}

type IClient interface {
	Get() *elastic.Client
	Error() error // connection error
	Debug() IDebug
}
