package esclient

import (
	"context"
	"github.com/olivere/elastic"
)

type IDebug interface {
	// Request() returns full http request
	Request() []byte

	// Response() returns full http response
	Response() []byte
}

type IConn interface {
	Open(bool, ...context.Context) IClient
}

type IClient interface {
	Get() *elastic.Client
	Error() error // connection error
	Debug() IDebug
}
