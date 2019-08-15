package esclient

import (
	"github.com/olivere/elastic"
)

/*

type IClient interface {
	Get() *elastic.Client
	Close()
	GetDebug() IDebug
	GetError() IError
}

*/

type oneClient struct {
	isFree bool

	client *elastic.Client

	errorObject IError
	debugObject IDebug
}

func newOneClient(client *elastic.Client, errObj IError, debObj IDebug) *oneClient {

	if errObj == nil {
		errObj = &ErrorHandler{}
	}
	if errObj == nil {
		debObj = &DebugHandler{}
	}

	return &oneClient{
		client: client,

		errorObject: errObj,
		debugObject: debObj,
	}
}

func (o oneClient) Get() *elastic.Client {
	return o.client
}

func (o *oneClient) GetDebug() IDebug {
	if o.debugObject == nil {
		return NeWDebugHandler()
	}
	return o.debugObject
}

func (o *oneClient) GetError() IError {
	if o.errorObject == nil {
		return NewErrorHandler()
	}
	return o.errorObject
}

func (o *oneClient) lock() *oneClient {
	o.isFree = false
	return o
}

func (o *oneClient) isLock() bool {
	return !o.isFree
}

func (o *oneClient) Close() {
	o.isFree = false
	o.errorObject = &ErrorHandler{}
	o.debugObject = &DebugHandler{}
}
