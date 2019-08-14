package esclient

import (
	"net/http"
	"net/http/httputil"
)

type DebugHandler struct {
	requestData  []byte
	responseData []byte
}

func newDebugHandler() *DebugHandler {
	return &DebugHandler{
		requestData:  []byte{},
		responseData: []byte{},
	}
}

// >>>>>>>>>> Interface function

func (d *DebugHandler) Request() []byte {
	return d.requestData
}

func (d *DebugHandler) Response() []byte {
	return d.responseData
}

func (d *DebugHandler) setHttpRequest(req *http.Request) {
	if req != nil {
		if body, err := httputil.DumpRequestOut(req, true); err == nil {
			d.requestData = body
		}
	}
}

func (d *DebugHandler) setHttpResponse(resp *http.Response, err error) {
	if resp != nil {
		if body, errParsing := httputil.DumpResponse(resp, true); errParsing == nil {
			d.responseData = body
		}
	}
}
