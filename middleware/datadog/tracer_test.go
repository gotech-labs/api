package datadog_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"

	"github.com/gotech-labs/api"
	apitest "github.com/gotech-labs/api/http/testing"
	. "github.com/gotech-labs/api/middleware/datadog"
)

func TestHealth(t *testing.T) {
	for _, test := range []struct {
		name     string
		method   string
		path     string
		response api.Response
		errorMsg string
	}{
		{
			name:     "ok response",
			method:   http.MethodGet,
			path:     "/health",
			response: api.OK("ok"),
			errorMsg: "",
		},
		{
			name:     "bad request response",
			method:   http.MethodPost,
			path:     "/search",
			response: api.BadRequest(fmt.Errorf("validation error")),
			errorMsg: "400 Bad Request",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var (
				rb = apitest.RequestBuilder{
					Method: test.method,
					Path:   test.path,
				}
				req     = rb.Build()
				handler = func(ctx context.Context, req api.Request) api.Response {
					return test.response
				}
				middleware = New("api test", "test").
						WithEnabledTraceLogger(bytes.NewBuffer(nil)).
						WithEnabledRuntimeMetrics().
						Middleware()
			)
			defer StopTracer()

			mt := mocktracer.Start()
			defer mt.Stop()

			// call middleware function
			resp := middleware(handler)(context.Background(), req)

			// assert response
			assert.Equal(t, test.response.Status(), resp.Status())

			// assert datadog tracer log
			spans := mt.FinishedSpans()
			if assert.Len(t, spans, 1) {
				assert.Equal(t, "http.request", spans[0].OperationName())
				assert.Equal(t, test.method, spans[0].Tag("http.method"))
				assert.Equal(t, test.path, spans[0].Tag("http.url"))
				assert.Equal(t, test.response.Status(), spans[0].Tag("http.status_code"))
				if test.errorMsg == "" {
					assert.Empty(t, spans[0].Tag("error"))
				} else {
					assert.Equal(t, test.errorMsg, spans[0].Tag("error"))
				}
			}
		})
	}
}
