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

func makeHttpClient(in ReqHandler, out ResHandler) (*DebugHandler, *http.Client) {
	debugObject := newDebugHandler()
	return debugObject, httpClient(in, makeReqFunc(debugObject), out, makeResFunc(debugObject))
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

	CustomReqFun ReqHandler
	CustomResFun ResHandler
}

func (r RoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	if !supportedMethods[request.Method] {
		return r.httpRoundTripper.RoundTrip(request)
	}

	if r.CustomReqFun != nil {
		r.CustomReqFun(request)
	}

	if r.ReqFunc != nil {
		r.ReqFunc(request)
	}

	if r.CustomResFun == nil && r.ResFunc == nil {
		// most common case
		return r.httpRoundTripper.RoundTrip(request)
	}

	response, err := r.httpRoundTripper.RoundTrip(request)
	if r.CustomResFun != nil {
		response, err = r.CustomResFun(response, err)
	}

	if r.ResFunc != nil {
		response, err = r.ResFunc(response, err)
	}

	return response, err
}

// httpClient prepares instance of http client
func httpClient(customReq, reqFunc ReqHandler, customRes, resFunc ResHandler) *http.Client {
	rt := RoundTripper{
		httpRoundTripper: &http.Transport{},
		ReqFunc:          reqFunc,
		ResFunc:          resFunc,
		CustomReqFun:     customReq,
		CustomResFun:     customRes,
	}

	return &http.Client{
		Transport: rt,
	}
}
