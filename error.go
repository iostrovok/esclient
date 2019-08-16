package esclient

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type ErrorHandler struct {
	httpStatusCode int
	httpError      error
	httpBody       []byte
	fullError      *FullError

	wasUpdated bool

	esStatus         int
	code             Code
	esType, esReason string
}

func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{
		httpStatusCode: 0,
		httpError:      nil,
		httpBody:       []byte{},
		fullError:      &FullError{},
	}
}

// >>>>>>>>>> Interface function
func (e *ErrorHandler) Error() string {
	return e.esReason
}

func (e *ErrorHandler) Code() Code {
	return e.code
}

func (e *ErrorHandler) Status() int {
	return e.esStatus
}

func (e *ErrorHandler) WasUpdated() bool {
	return e.wasUpdated
}

func (e *ErrorHandler) Reason() string {
	return e.esReason
}

func (e *ErrorHandler) Type() string {
	return e.esType
}

func (e *ErrorHandler) PasredError() *FullError {
	return e.fullError
}

// <<<<<<<<<< Interface function

func (e *ErrorHandler) SetHttpResponse(resp *http.Response, err error) {
	if resp != nil {
		e.httpStatusCode = resp.StatusCode
	}
	e.httpError = err
	e.wasUpdated = true
}

func (e *ErrorHandler) SetHttpBody(body []byte) {

	bodyParts := bytes.SplitN(body, []byte("\r\n\r\n"), 2)

	if len(bodyParts) != 2 || len(bodyParts[1]) == 0 {
		return
	}

	e.parseBody(bodyParts[1])
	e.wasUpdated = true
}

func (e *ErrorHandler) parseBody(body []byte) {

	err := json.Unmarshal(body, &e.fullError)
	if err != nil {
		if e.httpError == nil {
			e.httpError = err
		}
	} else {
		e.esStatus, e.esType, e.esReason = extractErrorMessage(e.fullError)
		e.checkExceptionType()
	}
}

func extractErrorMessage(fullError *FullError) (int, string, string) {

	if fullError.ErrorData == nil || fullError.Status == 0 {
		return 0, "", ""
	}

	if fullError.ErrorData.RootCause != nil && len(fullError.ErrorData.RootCause) > 0 {
		return fullError.Status, fullError.ErrorData.RootCause[0].Type, fullError.ErrorData.RootCause[0].Reason
	}

	return fullError.Status, fullError.ErrorData.Type, fullError.ErrorData.Reason
}

/// Elasticsearch error
type FullError struct {
	ErrorData *ErrorData `json:"error"`
	Status    int        `json:"status"`

	// for by ID searching
	Index string `json:"_index"`
	ID    string `json:"_id"`
	Type  string `json:"_type"`
	Found bool   `json:"found"`
}

type Reason struct {
	Type         string `json:"type"`
	Reason       string `json:"reason"`
	IndexUUID    string `json:"index_uuid"`
	Index        string `json:"index"`
	ResourceType string `json:"resource.type"`
	ResourceID   string `json:"resource.id"`
}

type FailedShards struct {
	Shard  int     `json:"shard"`
	Index  string  `json:"index"`
	Node   string  `json:"node"`
	Reason *Reason `json:"reason"`
}

type ErrorData struct {
	RootCause []*Reason `json:"root_cause"`
	Type      string    `json:"type"`
	Reason    string    `json:"reason"`

	ResourceType string          `json:"resource.type"`
	ResourceID   string          `json:"resource.id"`
	IndexUUID    string          `json:"index_uuid"`
	Index        string          `json:"index"`
	Phase        string          `json:"phase"`
	Grouped      bool            `json:"grouped"`
	FailedShards []*FailedShards `json:"failed_shards"`
}
