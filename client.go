package esclient

import (
	"context"

	"github.com/olivere/elastic/v7"
)

type ConnectionType int

const (
	SimpleType ConnectionType = iota
	DialType
	DialContextType
	ClientType
)

/*
Wrappers over github.com/olivere/elastic functions:
	- Dial(...)
	- DialContext(...)
	- NewClient(...)
	- NewSimpleClient(...)
*/

func Dial(options ...elastic.ClientOptionFunc) (IConn, error) {
	es, err := elastic.Dial(options...)
	return newConn(DialType, options, es, err, context.Background()), err
}

func DialContext(ctx context.Context, options ...elastic.ClientOptionFunc) (IConn, error) {
	es, err := elastic.DialContext(ctx, options...)
	return newConn(DialContextType, options, es, err, ctx), err
}

func NewClient(options ...elastic.ClientOptionFunc) (IConn, error) {
	es, err := elastic.NewClient(options...)
	return newConn(ClientType, options, es, err, context.Background()), err
}

func NewSimpleClient(options ...elastic.ClientOptionFunc) (IConn, error) {
	es, err := elastic.NewSimpleClient(options...)
	return newConn(SimpleType, options, es, err, context.Background()), err
}

func (c *conn) Open(useDebug bool, ctxs ...context.Context) IClient {
	if useDebug {
		return c.newDebugClient(ctxs...)
	}

	c.mc.RLock()
	defer c.mc.RUnlock()

	return newClient(c.firstElasticClient, c.connectionError, newDebugHandler())
}

func (c *conn) newDebugClient(ctxs ...context.Context) IClient {
	debugObject, httpClient := makeHttpClient(c.requestHandler, c.responseHandler)

	options := append(c.options,
		elastic.SetHttpClient(httpClient),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetHealthcheckTimeoutStartup(0),
	)

	var es *elastic.Client
	var err error

	switch c.connectionType {
	case ClientType:
		es, err = elastic.NewClient(options...)
	case DialContextType:
		if len(ctxs) == 0 {
			ctxs = append(ctxs, context.Background())
		}
		es, err = elastic.DialContext(ctxs[0], options...)
	case DialType:
		es, err = elastic.Dial(options...)
	case SimpleType:
		es, err = elastic.NewSimpleClient(options...)
	}

	return newClient(es, err, debugObject)
}
