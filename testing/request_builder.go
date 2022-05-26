package testing

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
)

type RequestBuilder struct {
	Method      string
	Path        string
	Body        []byte
	Headers     map[string]string
	PathParams  map[string]string
	QueryParams map[string]string
}

func (rb RequestBuilder) Build() *http.Request {
	req := httptest.NewRequest(
		rb.Method,
		rb.Path,
		bytes.NewBuffer(rb.Body),
	)
	if len(rb.PathParams) > 0 {
		req = mux.SetURLVars(req, rb.PathParams)
	}
	if len(rb.Headers) > 0 {
		for k, v := range rb.Headers {
			req.Header.Add(k, v)
		}
	}
	if len(rb.QueryParams) > 0 {
		values := req.URL.Query()
		for k, v := range rb.QueryParams {
			values.Add(k, v)
		}
		req.URL.RawQuery = values.Encode()
	}
	req.RemoteAddr = "127.0.0.1:80"
	req.URL.Host = "localhost"
	return req
}
