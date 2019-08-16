package esclient

import (
	"context"
	"sync"

	"github.com/olivere/elastic"
	"github.com/olivere/elastic/config"
)

type ConnectionType int

const (
	SimpleType ConnectionType = iota
	DialType
	ClientType
	ConfigType
)

type Client struct {
	mu sync.RWMutex

	cfg            *config.Config
	connectionType ConnectionType

	options []elastic.ClientOptionFunc
	ctx     context.Context

	simplesClient *oneClient
	oneClientList map[Type]*ClientList
}

/*
Wrappers over github.com/olivere/elastic functions:
	- DialWithConfig(...)
	- Dial(...)
	- DialContext(...)
	- NewClient(...)
	- NewSimpleClient(...)
*/

func DialWithConfig(ctx context.Context, cfg *config.Config) (IESClient, error) {
	c, err := elastic.DialWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	out := initNewClient(DialType, c, []elastic.ClientOptionFunc{}, cfg)
	return out, nil
}

func Dial(options ...elastic.ClientOptionFunc) (IESClient, error) {
	c, err := elastic.Dial(options...)
	if err != nil {
		return nil, err
	}

	out := initNewClient(DialType, c, options, nil)
	return out, nil
}

func DialContext(ctx context.Context, options ...elastic.ClientOptionFunc) (IESClient, error) {
	c, err := elastic.NewClient(options...)
	if err != nil {
		return nil, err
	}

	out := initNewClient(DialType, c, options, nil)
	return out, nil
}

func NewClient(options ...elastic.ClientOptionFunc) (IESClient, error) {
	c, err := elastic.NewClient(options...)
	if err != nil {
		return nil, err
	}

	out := initNewClient(ClientType, c, options, nil)
	return out, nil
}

func NewSimpleClient(options ...elastic.ClientOptionFunc) (IESClient, error) {
	c, err := elastic.NewSimpleClient(options...)
	if err != nil {
		return nil, err
	}

	out := initNewClient(SimpleType, c, options, nil)
	return out, nil
}

func initNewClient(connectionType ConnectionType, elasticClient *elastic.Client, options []elastic.ClientOptionFunc, cfg *config.Config) *Client {

	oneClientList := map[Type]*ClientList{
		Error:         NewClientList(),
		Debug:         NewClientList(),
		ErrorAndDebug: NewClientList(),
	}

	return &Client{
		connectionType: connectionType,
		options:        options,
		cfg:            cfg,

		simplesClient: newOneClient(elasticClient, nil, nil),
		oneClientList: oneClientList,
	}

}

func (c *Client) Open(options ...Type) IClient {

	t := None
	if len(options) == 1 {
		t = options[0]
	}

	if t == None {
		return newOneClient(c.simplesClient.Get(), nil, nil)
	}

	if cl, find := c.findFree(t); find {
		return cl
	}

	return c.addNewClient(t).lock()
}

// next just return client
func (c *Client) findFree(t Type) (IClient, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.oneClientList[t].findFree()
}

// nextSimple just return oliver client
func (c *Client) appendOneClient(t Type, o *oneClient) *oneClient {
	c.mu.RLock()
	defer c.mu.RUnlock()

	c.oneClientList[t].appendTo(o)

	return o
}

// addNewClient just return oliver client
func (c *Client) addNewClient(t Type) *oneClient {

	// DialWithConfig does not support option SetHttpClient
	if c.connectionType == ConfigType {
		es, _ := elastic.DialWithConfig(context.Background(), c.cfg)
		cl := newOneClient(es, nil, nil)
		return c.appendOneClient(t, cl)
	}

	errObject, debugObject, httpClient := MakeHttpClient(t)
	options := append(c.options, elastic.SetHttpClient(httpClient))

	var es *elastic.Client
	switch c.connectionType {
	case DialType, ClientType:
		es, _ = elastic.DialContext(context.Background(), options...)
	case SimpleType:
		es, _ = elastic.NewSimpleClient(options...)
	}

	cl := newOneClient(es, errObject, debugObject)
	return c.appendOneClient(t, cl)
}
