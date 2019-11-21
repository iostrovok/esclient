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
	pingService        []*elastic.PingService

	logger        ILogger
	loggerSet     bool
	sniffDuration time.Duration
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

func (c *conn) setPingService() {
	c.mc.Lock()
	defer c.mc.Unlock()

	c.pingService = make([]*elastic.PingService, 0)
	urls := extractURLs(c.firstElasticClient.String())
	for _, url := range urls {
		c.pingService = append(c.pingService, c.firstElasticClient.Ping(url))
	}
}

func (c *conn) Sniff(ctx context.Context) {
	c.setPingService()
	go c.run(ctx)
}

func (c *conn) run(ctx context.Context) {

	// sleep before first attempt
	time.Sleep(c.sniffDuration)

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(c.sniffDuration):
			c.checkConnections(ctx)
		}
	}
}

func (c *conn) checkConnections(ctx context.Context) {
	c.mc.Lock()
	defer c.mc.Unlock()

	for i := range c.pingService {
		_, statusCode, err := c.pingService[i].Do(ctx)
		if err != nil || statusCode < 200 || statusCode >= 300 {
			c.Printf("Reconnect statusCode: %d. err: %v. res: %v", statusCode, err)
			c.reConnect()
			return
		}
	}
}

func (c *conn) reConnect() error {

	c.Printf("Reconnect.....")

	var es *elastic.Client
	var err error

	switch c.connectionType {
	case ClientType:
		es, err = elastic.NewClient(c.options...)
	case DialContextType:
		es, err = elastic.DialContext(c.ctx, c.options...)
	case DialType:
		es, err = elastic.Dial(c.options...)
	case SimpleType:
		es, err = elastic.NewSimpleClient(c.options...)
	}

	if err == nil {
		c.firstElasticClient = es
	}

	return err
}
