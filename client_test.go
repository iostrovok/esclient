package esclient

import (
	"testing"

	. "github.com/iostrovok/check"
	"github.com/olivere/elastic/v7"
)

type testSuite struct{}

var _ = Suite(&testSuite{})

func TestService(t *testing.T) { TestingT(t) }

//// TestErrorHandler_NoSuchIndex
func (s *testSuite) TestErrorHandler_ConnectionError(c *C) {

	testClient, err := NewClient(elastic.SetURL("bla-bla-bla"))
	c.Assert(err, NotNil)
	c.Assert(testClient, NotNil)

	cl := testClient.Open(true)
	c.Assert(cl.Error(), NotNil)

	cl = testClient.Open(false)
	c.Assert(cl.Error(), NotNil)

}
