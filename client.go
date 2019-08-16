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

	// configuration for (re)connection
	connectionType ConnectionType
	cfg            *config.Config
	options        []elastic.ClientOptionFunc

	simpleElasticClient *elastic.Client
	oneClientList       map[Type]*ClientList
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
	return initNewClient(DialType, []elastic.ClientOptionFunc{}, cfg)
}

func Dial(options ...elastic.ClientOptionFunc) (IESClient, error) {
	return initNewClient(DialType, options, nil)
}

func DialContext(ctx context.Context, options ...elastic.ClientOptionFunc) (IESClient, error) {
	return initNewClient(DialType, options, nil)
}

func NewClient(options ...elastic.ClientOptionFunc) (IESClient, error) {
	return initNewClient(ClientType, options, nil)
}

func NewSimpleClient(options ...elastic.ClientOptionFunc) (IESClient, error) {
	return initNewClient(SimpleType, options, nil)
}

func initNewClient(connectionType ConnectionType, options []elastic.ClientOptionFunc, cfg *config.Config) (*Client, error) {

	var es *elastic.Client
	var err error
	switch connectionType {
	case DialType, ClientType:
		es, err = elastic.DialContext(context.Background(), options...)
	case SimpleType:
		es, err = elastic.NewSimpleClient(options...)
	}

	oneClientList := map[Type]*ClientList{
		Error:         NewClientList(),
		Debug:         NewClientList(),
		ErrorAndDebug: NewClientList(),
	}

	return &Client{
		connectionType: connectionType,
		options:        options,
		cfg:            cfg,

		simpleElasticClient: es,
		oneClientList:       oneClientList,
	}, err
}

func (c *Client) Open(options ...Type) IClient {

	t := None
	if len(options) == 1 {
		t = options[0]
	}

	if t == None && c.simpleElasticClient != nil {
		// It just returns wrapper over elasticsearch client
		c.mu.RLock()
		defer c.mu.RUnlock()
		return newOneClient(c.simpleElasticClient, nil, nil, nil)
	} else if c.simpleElasticClient == nil {
		c.mu.Lock()
		defer c.mu.Unlock()
		var err error
		c.simpleElasticClient, err = elastic.NewSimpleClient(c.options...)
		return newOneClient(c.simpleElasticClient, err, nil, nil)
	}

	// Reusing connection
	if cl, find := c.findFree(t); find {
		return cl
	}

	// Creates new connection to elasticsearch
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
		es, err := elastic.DialWithConfig(context.Background(), c.cfg)
		cl := newOneClient(es, err, nil, nil)
		return c.appendOneClient(t, cl)
	}

	errObject, debugObject, httpClient := MakeHttpClient(t)
	options := append(c.options, elastic.SetHttpClient(httpClient))

	var es *elastic.Client
	var err error
	switch c.connectionType {
	case DialType, ClientType:
		es, err = elastic.DialContext(context.Background(), options...)
	case SimpleType:
		es, err = elastic.NewSimpleClient(options...)
	}

	cl := newOneClient(es, err, errObject, debugObject)
	return c.appendOneClient(t, cl)
}
