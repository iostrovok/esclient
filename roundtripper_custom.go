package esclient

import (
	"net/http"
)

// httpClientCustom prepares instance of http client
func httpClientCustom(customReq ReqHandler, customRes ResHandler) *http.Client {
	rt := RoundTripperCustom{
		httpRoundTripper: &http.Transport{},
		CustomReqFun:     customReq,
		CustomResFun:     customRes,
	}

	return &http.Client{
		Transport: rt,
	}
}

// see https://golang.org/pkg/net/http/#RoundTripper
type RoundTripperCustom struct {
	httpRoundTripper http.RoundTripper

	CustomReqFun ReqHandler
	CustomResFun ResHandler
}

func (r RoundTripperCustom) RoundTrip(request *http.Request) (*http.Response, error) {
	if !supportedMethods[request.Method] {
		return r.httpRoundTripper.RoundTrip(request)
	}

	if r.CustomReqFun != nil {
		r.CustomReqFun(request)
	}

	if r.CustomResFun == nil {
		// most common case
		return r.httpRoundTripper.RoundTrip(request)
	}

	return r.CustomResFun(r.httpRoundTripper.RoundTrip(request))
}
