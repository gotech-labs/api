package testing

import (
	"bytes"
	core_http "net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	"github.com/gotech-labs/api"
	"github.com/gotech-labs/api/http"
)

type RequestBuilder struct {
	Method      string
	Path        string
	Body        []byte
	Headers     map[string][]string
	PathParams  map[string]string
	QueryParams map[string]string
}

func (rb RequestBuilder) Build() api.Request {
	req := httptest.NewRequest(
		rb.Method,
		rb.Path,
		bytes.NewBuffer(rb.Body),
	)
	if len(rb.PathParams) > 0 {
		req = mux.SetURLVars(req, rb.PathParams)
	}
	if len(rb.Headers) > 0 {
		req.Header = core_http.Header(rb.Headers)
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
	return http.NewRequest(req)
}
