package esclient

import (
	"context"
	"time"

	. "github.com/iostrovok/check"
	"github.com/olivere/elastic/v7"
)

// TestExtractURLs
func (s *testSuite) TestExtractURLs(c *C) {
	line := "https://vpc.us-east-1.es.amazonaws.com [dead=false,failures=0,deadSince=<nil>], http://127.0.0.1:9200 [dead=false,failures=0,deadSince=<nil>]"
	urls := extractURLs(line)

	checkList := []string{
		"https://vpc.us-east-1.es.amazonaws.com",
		"http://127.0.0.1:9200",
	}
	c.Assert(urls, DeepEquals, checkList)
}

func (s *testSuite) TestReConnect_1(c *C) {
	// Create an Elasticsearch connection
	connection, err := NewSimpleClient(options...)
	c.Assert(err, IsNil)
	cl := connection.Open(false, context.Background())

	result, err := cl.Get().Get().
		Index(testIndex).
		Id("one").
		Do(context.Background())
	c.Assert(err, IsNil)

	c.Assert(result.Id, Equals, "one")

	connection.SniffTimeout(1 * time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	connection.SniffTimeout(1 * time.Second)
	connection.Sniff(ctx)
	time.Sleep(3 * time.Second)

	err = connection.reConnect()
	c.Assert(err, IsNil)

	cl = connection.Open(false, context.Background())
	result, err = cl.Get().Get().
		Index(testIndex).
		Id("one").
		Do(context.Background())
	c.Assert(err, IsNil)

	c.Assert(result.Id, Equals, "one")
}

func (s *testSuite) TestReConnect_2(c *C) {
	connection, err := NewSimpleClient(elastic.SetURL("bla-bla-bla"))
	c.Assert(err, NotNil)
	connection.SniffTimeout(1 * time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	connection.SniffTimeout(1 * time.Second)
	connection.Sniff(ctx)
	time.Sleep(3 * time.Second)

	err = connection.reConnect()
	c.Assert(err, NotNil)
}
