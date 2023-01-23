package esclient

import (
	"net/http"
)

type ReqHandler func(*http.Request)
type ResHandler func(*http.Response, error) (*http.Response, error)

var supportedMethods = map[string]bool{
	http.MethodGet:  true,
	http.MethodPost: true,
	http.MethodPut:  true,
}

func makeHttpClient() (*DebugHandler, *http.Client) {
	debugObject := newDebugHandler()
	return debugObject, httpClient(makeReqFunc(debugObject), makeResFunc(debugObject))
}

func makeReqFunc(debugObject *DebugHandler) ReqHandler {
	return func(req *http.Request) {
		debugObject.setHttpRequest(req)
	}
}

func makeResFunc(debugObject *DebugHandler) ResHandler {
	return func(resp *http.Response, err error) (*http.Response, error) {
		debugObject.setHttpResponse(resp, err)
		return resp, err
	}
}

// see https://golang.org/pkg/net/http/#RoundTripper
type RoundTripper struct {
	httpRoundTripper http.RoundTripper
	ReqFunc          ReqHandler
	ResFunc          ResHandler
}

func (r RoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	if !supportedMethods[request.Method] {
		r.httpRoundTripper.RoundTrip(request)
	}

	if r.ReqFunc != nil {
		r.ReqFunc(request)
	}

	if r.ResFunc != nil {
		return r.ResFunc(r.httpRoundTripper.RoundTrip(request))
	}

	return r.httpRoundTripper.RoundTrip(request)
}

// httpClient prepares instance of http client
func httpClient(reqFunc ReqHandler, resFunc ResHandler) *http.Client {
	rt := RoundTripper{
		httpRoundTripper: &http.Transport{},
		ReqFunc:          reqFunc,
		ResFunc:          resFunc,
	}

	return &http.Client{
		Transport: rt,
	}
}
