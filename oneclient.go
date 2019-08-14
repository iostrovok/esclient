package esclient

import (
	"github.com/olivere/elastic/v7"
)

/*

type IClient interface {
	Get() *elastic.Client
	Error() error // connection error
	Debug() IDebug
}

*/

type Client struct {
	client   *elastic.Client
	conError error
	debug    *DebugHandler
}

func newClient(client *elastic.Client, conError error, debug *DebugHandler) *Client {
	return &Client{
		client:   client,
		conError: conError,
		debug:    debug,
	}
}

func (o Client) Get() *elastic.Client {
	return o.client
}

func (o Client) Error() error {
	return o.conError
}

func (o *Client) Debug() IDebug {
	if o.debug == nil {
		return newDebugHandler()
	}
	return o.debug
}
