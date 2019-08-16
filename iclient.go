package esclient

import (
	"github.com/olivere/elastic"
	"net/http"
)

type Type int

const (
	None Type = iota
	Error
	Debug
	ErrorAndDebug
)

type IDebug interface {
	// Request() returns full http request
	Request() []byte

	// Response() returns full http response
	Response() []byte

	// WasUpdated() indicates that request was processed
	WasUpdated() bool

	// simple functions for collect debug data
	Add(string, interface{})
	Get(string) interface{}
	All() map[string]interface{}

	// internal functions
	SetHttpRequest(req *http.Request)
	SetHttpResponse(*http.Response,  error)
}

type IError interface {
	// Error() returns description from elasticsearch error
	Error() string

	// Reason() returns reason of elasticsearch error
	Reason() string

	// Type() returns type of elasticsearch error
	Type() string

	// Status() returns elasticsearch response status
	Status() int

	// PasredError() returns parsed elasticsearch error reason
	PasredError() *FullError

	/*
		Code() returns "correct" status code.
		For example, if request contents the wrong sorting field,
		we get Status() == 200 and Code() == 500
	*/
	Code() Code

	// WasUpdated() indicates that request was processed
	WasUpdated() bool

	// internal functions
	SetHttpResponse(*http.Response,  error)
	SetHttpBody([]byte)
}

type IESClient interface {
	Open(...Type) IClient
}

type IClient interface {
	Get() *elastic.Client
	Close()
	ConnError() error // connection error
	GetDebug() IDebug
	GetError() IError
}
