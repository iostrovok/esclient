package esclient

import (
	"net/http"
	"net/http/httputil"
)

type ReqHandler func(*http.Request)
type ResHandler func(*http.Response, error) (*http.Response, error)

func MakeHttpClient(t Type) (*ErrorHandler, *DebugHandler, *http.Client) {

	debugObject := NeWDebugHandler()
	errObject := NewErrorHandler()

	reqFunc := MakeReqFunc(t, debugObject)
	resFunc := MakeResFunc(t, errObject, debugObject)

	return errObject, debugObject, httpClient(reqFunc, resFunc)
}

func MakeReqFunc(t Type, debugObject *DebugHandler) ReqHandler {

	if t == None || t == Error {
		return nil
	}

	return func(req *http.Request) {
		if body, err := httputil.DumpRequestOut(req, true); err == nil {
			debugObject.SetRequest(body)
		}
	}
}

func MakeResFunc(t Type, errObject *ErrorHandler, debugObject *DebugHandler) ResHandler {

	if t == None {
		return nil
	}

	if t == Debug {
		return func(resp *http.Response, err error) (*http.Response, error) {
			if body, errParsing := httputil.DumpResponse(resp, true); errParsing == nil {
				debugObject.SetResponse(body)
			}

			debugObject.Update()

			return resp, err
		}
	}

	if t == Error {
		return func(resp *http.Response, err error) (*http.Response, error) {

			errObject.SetHttpError(err)
			errObject.SetHttpStatusCode(resp.StatusCode)
			if body, errParsing := httputil.DumpResponse(resp, true); errParsing == nil {
				errObject.SetHttpBody(body)
			}

			errObject.Update()

			return resp, err
		}
	}

	if t == ErrorAndDebug {
		return func(resp *http.Response, err error) (*http.Response, error) {

			errObject.SetHttpError(err)
			errObject.SetHttpStatusCode(resp.StatusCode)
			if body, errParsing := httputil.DumpResponse(resp, true); errParsing == nil {
				errObject.SetHttpBody(body)
				debugObject.SetResponse(body)
			}

			errObject.Update()
			debugObject.Update()

			return resp, err
		}
	}

	// we never will be here
	return nil
}

// see https://golang.org/pkg/net/http/#RoundTripper
type RoundTripper struct {
	httpRoundTripper http.RoundTripper
	ReqFunc          ReqHandler
	ResFunc          ResHandler
}

func (r RoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
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
