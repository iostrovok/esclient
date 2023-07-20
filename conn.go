package esclient

import (
	"context"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/olivere/elastic/v7"
)

type conn struct {
	mc sync.RWMutex

	// configuration for (re)connection
	connectionType     ConnectionType
	options            []elastic.ClientOptionFunc
	firstElasticClient *elastic.Client
	connectionError    error
	ctx                context.Context
	pingService        *pingService

	logger        ILogger
	loggerSet     bool
	sniffDuration time.Duration

	// custom savers
	responseHandler ResHandler
	requestHandler  ReqHandler
}

func newConn(c ConnectionType, options []elastic.ClientOptionFunc, es *elastic.Client, err error, ctx context.Context) *conn {
	return &conn{
		connectionType:     c,
		options:            options,
		firstElasticClient: es,
		connectionError:    err,
		ctx:                ctx,
		sniffDuration:      PingDuration,
	}
}

const (
	PingDuration = 5 * time.Second
)

type FakeClient struct {
	elastic.Client
}

var splitURLs = regexp.MustCompile(`\s+\[[^]]+\],?`)

func extractURLs(line string) []string {
	// "%s [dead=%v,failures=%d,deadSince=%v]"
	out := make([]string, 0)
	for _, u := range splitURLs.Split(line, -1) {
		u = strings.TrimSpace(u)
		if u != "" {
			out = append(out, u)
		}
	}
	return out
}

func (c *conn) SetCustomHandler(req ReqHandler, res ResHandler) error {
	c.mc.Lock()
	c.requestHandler = req
	c.responseHandler = res
	c.mc.Unlock()

	return c.reConnect()
}

func (c *conn) SetLogger(logger ILogger) {
	c.logger = logger
	c.loggerSet = true
}

func (c *conn) Fatal(v ...interface{}) {
	if c.loggerSet {
		c.logger.Fatal(v...)
	}
}

func (c *conn) Print(v ...interface{}) {
	if c.loggerSet {
		c.logger.Print(v...)
	}
}

func (c *conn) Printf(format string, v ...interface{}) {
	if c.loggerSet {
		c.logger.Printf(format, v...)
	}
}

func (c *conn) SniffTimeout(duration time.Duration) {
	c.mc.Lock()
	defer c.mc.Unlock()

	c.sniffDuration = duration
}

func (c *conn) Sniff(ctx context.Context) {
	if c.firstElasticClient == nil {
		return
	}

	c.pingService = newPingService(c.sniffDuration, c.reConnect, c.Printf)
	c.mc.RLock()
	urls := extractURLs(c.firstElasticClient.String())
	c.mc.RUnlock()

	for _, url := range urls {
		c.pingService.Add(c.firstElasticClient.Ping(url))
	}

	go c.pingService.runSniff(ctx)
}

func (c *conn) reConnect() error {
	c.mc.Lock()
	defer c.mc.Unlock()

	var es *elastic.Client
	var err error

	options := c.options
	if c.responseHandler != nil || c.requestHandler != nil {
		httpClient := httpClientCustom(c.requestHandler, c.responseHandler)
		options = append(options, elastic.SetHttpClient(httpClient))
	}

	switch c.connectionType {
	case ClientType:
		es, err = elastic.NewClient(options...)
	case DialContextType:
		es, err = elastic.DialContext(c.ctx, options...)
	case DialType:
		es, err = elastic.Dial(options...)
	case SimpleType:
		es, err = elastic.NewSimpleClient(options...)
	}

	if err == nil {
		c.firstElasticClient = es
	}

	return err
}
