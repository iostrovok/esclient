package esclient

import (
	"context"
	"net/http"
	"strings"
	"testing"

	. "github.com/iostrovok/check"
	"github.com/olivere/elastic/v7"
)

const (
	testURL = "http://127.0.0.1:9200"
	testIndex          = "my_test_index_123"
	testMappings       = `{"mappings": {"dynamic": true}}`
	testDeleteMappings = `{"query": {"match_all":}}`
)

var options []elastic.ClientOptionFunc

type testSuite struct{}

var _ = Suite(&testSuite{})

func TestService(t *testing.T) { TestingT(t) }

//Run once when the suite starts running.
func insertRecord(c *C, id, record string) {

	client := &http.Client{}

	req, err := http.NewRequest("POST", testURL+"/"+testIndex+"/_doc"+"/"+id, strings.NewReader(record))
	c.Assert(err, IsNil)

	req.Header.Add("Content-Type", "application/json")
	_, err = client.Do(req)
	c.Assert(err, IsNil)

}

//Run once when the suite starts running.
func (s *testSuite) SetUpSuite(c *C) {
	client := &http.Client{}

	req, err := http.NewRequest("PUT", testURL+"/"+testIndex, strings.NewReader(testMappings))
	c.Assert(err, IsNil)

	req.Header.Add("Content-Type", "application/json")
	_, err = client.Do(req)
	c.Assert(err, IsNil)

	record := `{
		"user" : "Lelik",
		"post_date" : "2009-11-15T14:12:12",
		"message" : "lelik - trying out Elasticsearch"
	}`
	insertRecord(c, "one", record)

	record = `{
		"user" : "Bolik",
		"post_date" : "2009-11-16T12:10:01",
		"message" : "bolik - trying out Elasticsearch"
	}`
	insertRecord(c, "two", record)

	options = []elastic.ClientOptionFunc{
		elastic.SetURL(testURL),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	}
}

//Run before each test or benchmark starts running.
func (s *testSuite) SetUpTest(c *C) {}

//Run after each test or benchmark runs.
func (s *testSuite) TearDownTest(c *C) {}

//Run once after all tests or benchmarks have finished running.
func (s *testSuite) TearDownSuite(c *C) {
	//client := &http.Client{}
	//
	//req, err := http.NewRequest("DELETE", testURL+"/"+testIndex, strings.NewReader(testDeleteMappings))
	//c.Assert(err, IsNil)
	//
	//req.Header.Add("Content-Type", "application/json")
	//_, err = client.Do(req)
	//c.Assert(err, IsNil)
}

//// TestErrorHandler_NoSuchIndex
func (s *testSuite) TestConnectionError(c *C) {

	connection, err := NewClient(elastic.SetURL("bla-bla-bla"))
	c.Assert(err, NotNil)
	c.Assert(connection, NotNil)

	cl := connection.Open(true)
	c.Assert(cl.Error(), NotNil)

	cl = connection.Open(false)
	c.Assert(cl.Error(), NotNil)
}

func (s *testSuite) TestClientDebug_1(c *C) {

	// Create an Elasticsearch connection
	connection, err := NewSimpleClient(options...)
	c.Assert(err, IsNil)
	cl := connection.Open(true)

	result, err := cl.Get().Get().
		Index(testIndex).
		Id("one").
		Do(context.Background())
	c.Assert(err, IsNil)

	c.Assert(result.Id, Equals, "one")
}

func (s *testSuite) TestClientDebug_2(c *C) {

	// Create an Elasticsearch connection
	connection, err := Dial(options...)
	c.Assert(err, IsNil)
	cl := connection.Open(true)

	result, err := cl.Get().Get().
		Index(testIndex).
		Id("one").
		Do(context.Background())
	c.Assert(err, IsNil)

	c.Assert(result.Id, Equals, "one")
}

func (s *testSuite) TestClientDebug_3(c *C) {

	// Create an Elasticsearch connection
	ctx := context.Background()
	connection, err := DialContext(ctx, options...)
	c.Assert(err, IsNil)
	cl := connection.Open(true)

	result, err := cl.Get().Get().
		Index(testIndex).
		Id("one").
		Do(ctx)
	c.Assert(err, IsNil)

	c.Assert(result.Id, Equals, "one")

	debug := cl.Debug()
	c.Assert(cl.Debug(), NotNil)

	req := string(debug.Request())
	c.Assert(strings.HasPrefix(req, "GET /"), Equals, true)

	res := string(debug.Response())
	c.Assert(strings.HasPrefix(res, "HTTP/"), Equals, true)
}

func (s *testSuite) TestClientDebug_Empty(c *C) {

	// Create an Elasticsearch connection
	ctx := context.Background()
	connection, err := DialContext(ctx, options...)
	c.Assert(err, IsNil)
	cl := connection.Open(true)

	debug := cl.Debug()
	c.Assert(cl.Debug(), NotNil)
	c.Assert(string(debug.Request()), Equals, "")
	c.Assert(string(debug.Response()), Equals, "")

}

func (s *testSuite) TestClient_1(c *C) {

	// Create an Elasticsearch connection
	connection, err := NewSimpleClient(options...)
	c.Assert(err, IsNil)
	cl := connection.Open(false)

	result, err := cl.Get().Get().
		Index(testIndex).
		Id("one").
		Do(context.Background())
	c.Assert(err, IsNil)

	c.Assert(result.Id, Equals, "one")
}
