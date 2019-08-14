package esclient

import (
	"context"
	"github.com/olivere/elastic"
)

type ConnectionType int

const (
	SimpleType ConnectionType = iota
	DialType
	DialContextType
	ClientType
)

type conn struct {
	// configuration for (re)connection
	connectionType     ConnectionType
	options            []elastic.ClientOptionFunc
	firstElasticClient *elastic.Client
	connectionError    error
}

func newConn(c ConnectionType, options []elastic.ClientOptionFunc, es *elastic.Client, err error) *conn {
	return &conn{
		connectionType:     c,
		options:            options,
		firstElasticClient: es,
		connectionError:    err,
	}
}

/*
Wrappers over github.com/olivere/elastic functions:
	- Dial(...)
	- DialContext(...)
	- NewClient(...)
	- NewSimpleClient(...)
*/

func Dial(options ...elastic.ClientOptionFunc) (IConn, error) {
	es, err := elastic.Dial(options...)
	return newConn(DialType, options, es, err), err
}

func DialContext(ctx context.Context, options ...elastic.ClientOptionFunc) (IConn, error) {
	es, err := elastic.DialContext(ctx, options...)
	return newConn(DialContextType, options, es, err), err
}

func NewClient(options ...elastic.ClientOptionFunc) (IConn, error) {
	es, err := elastic.NewClient(options...)
	return newConn(ClientType, options, es, err), err
}

func NewSimpleClient(options ...elastic.ClientOptionFunc) (IConn, error) {
	es, err := elastic.NewSimpleClient(options...)
	return newConn(SimpleType, options, es, err), err
}

func (c *conn) Open(useDebug bool, ctxs ...context.Context) IClient {
	if useDebug {
		return c.newDebugClient(ctxs...)
	}

	return newClient(c.firstElasticClient, c.connectionError, newDebugHandler())
}

func (c *conn) newDebugClient(ctxs ...context.Context) IClient {

	debugObject, httpClient := makeHttpClient()

	options := append(c.options,
		elastic.SetHttpClient(httpClient),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetMaxRetries(1),
		elastic.SetHealthcheckTimeoutStartup(0),
	)

	var es *elastic.Client
	var err error

	switch c.connectionType {
	case ClientType:
		es, err = elastic.NewClient(options...)
	case DialContextType:
		if len(ctxs) > 0 {
			es, err = elastic.DialContext(ctxs[0], options...)
		} else {
			es, err = elastic.DialContext(context.Background(), options...)
		}
	case DialType:
		es, err = elastic.Dial(options...)
	case SimpleType:
		es, err = elastic.NewSimpleClient(options...)
	}

	return newClient(es, err, debugObject)
}
