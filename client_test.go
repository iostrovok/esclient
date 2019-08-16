package esclient

import (
	"github.com/olivere/elastic"
	"net/http"
	"regexp"
	"testing"

	. "gopkg.in/check.v1"
)

var pathReg *regexp.Regexp = regexp.MustCompile(`^/clutch_provider_.*?/`)

type testSuite struct{}

var _ = Suite(&testSuite{})

func TestService(t *testing.T) { TestingT(t) }

// TestErrorHandler_NoSuchIndex
func (s *testSuite) TestErrorHandler_NoSuchIndex(c *C) {

	data := []byte(`{"error":{"root_cause":[{"type":"index_not_found_exception","reason":"no such index","resource.type":"index_or_alias","resource.id":"providers-","index_uuid":"_na_","index":"providers-"}],"type":"index_not_found_exception","reason":"no such index","resource.type":"index_or_alias","resource.id":"providers-","index_uuid":"_na_","index":"providers-"},"status":404}`)

	errObject := NewErrorHandler()
	errObject.SetHttpResponse(&http.Response{StatusCode: 404}, nil)
	errObject.parseBody(data)

	c.Assert(errObject.Error(), Equals, "no such index")
	c.Assert(errObject.Reason(), Equals, "no such index")
	c.Assert(errObject.Type(), Equals, "index_not_found_exception")
	c.Assert(errObject.Status(), Equals, 404)
	c.Assert(errObject.Code(), Equals, Internal)
}

// TestErrorHandler_NoSuchIndex
func (s *testSuite) TestErrorHandler_NoMappingFound(c *C) {

	data := []byte(`{"error":{"root_cause":[{"type":"query_shard_exception","reason":"No mapping found for [title-] in order to sort on","index_uuid":"gN4QNi-UQU-S0fGabv8sCw","index":"clutch_provider_20190806_v1"}],"type":"search_phase_execution_exception","reason":"all shards failed","phase":"query","grouped":true,"failed_shards":[{"shard":0,"index":"clutch_provider_20190806_v1","node":"t0_af3KASy2DxRRYx0RDhw","reason":{"type":"query_shard_exception","reason":"No mapping found for [title-] in order to sort on","index_uuid":"gN4QNi-UQU-S0fGabv8sCw","index":"clutch_provider_20190806_v1"}}]},"status":400}`)

	errObject := NewErrorHandler()
	errObject.SetHttpResponse(&http.Response{StatusCode: 400}, nil)
	errObject.parseBody([]byte(data))

	//c.Assert(errObject.Error(), Equals, "No mapping found for [title-] in order to sort on")
	c.Assert(errObject.Reason(), Equals, "No mapping found for [title-] in order to sort on")
	c.Assert(errObject.Status(), Equals, 400)
	c.Assert(errObject.Code(), Equals, Internal)
	c.Assert(errObject.Type(), Equals, "query_shard_exception")
}

// TestErrorHandler_NoSuchIndex
func (s *testSuite) TestErrorHandler_ConnectionError(c *C) {

	testClient, err := NewClient(elastic.SetURL("bla-bla-bla"))
	c.Assert(err, NotNil)
	c.Assert(testClient, NotNil)

	cl := testClient.Open(ErrorAndDebug)
	c.Assert(cl.ConnError(), NotNil)

	cl = testClient.Open(Debug)
	c.Assert(cl.ConnError(), NotNil)

	cl = testClient.Open(Error)
	c.Assert(cl.ConnError(), NotNil)

	cl = testClient.Open(None)
	c.Assert(cl.ConnError(), NotNil)
}
