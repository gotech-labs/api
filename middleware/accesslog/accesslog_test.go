package accesslog_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"

	apitest "github.com/gotech-labs/api/http/testing"
	"github.com/stretchr/testify/assert"

	"github.com/gotech-labs/api"
	. "github.com/gotech-labs/api/middleware/accesslog"
	"github.com/gotech-labs/core/system"
)

func TestAccessLog(t *testing.T) {
	var (
		host, _ = os.Hostname()
	)
	system.RunTest(t, "skip logging", func(t *testing.T) {
		var (
			rb = apitest.RequestBuilder{
				Method: http.MethodGet,
				Path:   "/search",
				Body: []byte(`{
					"keyword": "hello"
				}`),
				Headers: map[string][]string{
					"Content-Type": {"application/json"},
					"X-Request-Id": {"req-12345"},
					"User-Agent":   {"TestUserAgent"},
				},
				QueryParams: map[string]string{
					"req_id": "req-XXXXXXXXXXXXXXXXX",
					"pretty": "true",
				},
			}
			req      = rb.Build()
			buf      = bytes.NewBuffer(nil)
			skipPath = []string{rb.Path} // skip path
			filter   = func(req api.Request) bool {
				return true
			}
			handler = func(ctx context.Context, req api.Request) api.Response {
				return api.OK(map[string]string{"message": "ok"})
			}
			middleware = New(buf).
					WithSkipPath(skipPath...).
					WithLoggingReqBodyFilter(filter).
					Middleware()
		)
		// call middleware function
		resp := middleware(handler)(context.Background(), req)
		assert.Equal(t, http.StatusOK, resp.Status())

		// assert log message
		assert.Empty(t, buf.String())
	})

	system.RunTest(t, "ok response (without request body)", func(t *testing.T) {
		var (
			rb = apitest.RequestBuilder{
				Method: http.MethodPost,
				Path:   "/search",
				Body: []byte(`{
					"keyword": "hello"
				}`),
				Headers: map[string][]string{
					"Content-Type": {"application/json"},
					"X-Request-Id": {"req-12345"},
					"User-Agent":   {"TestUserAgent"},
					"Referer":      {"TestReferer"},
				},
				QueryParams: map[string]string{
					"req_id": "req-XXXXXXXXXXXXXXXXX",
					"pretty": "true",
				},
			}
			req      = rb.Build()
			buf      = bytes.NewBuffer(nil)
			skipPath = []string{}
			filter   = func(req api.Request) bool {
				return false
			}
			handler = func(ctx context.Context, req api.Request) api.Response {
				return api.OK(map[string]string{"message": "ok"})
			}
			middleware = New(buf).
					WithSkipPath(skipPath...).
					WithLoggingReqBodyFilter(filter).
					Middleware()
		)
		// call middleware function
		resp := middleware(handler)(context.Background(), req)
		assert.Equal(t, http.StatusOK, resp.Status())

		expected := fmt.Sprintf(`{
			"level": "info",
			"time": "2022-12-24T00:00:00+09:00",
			"server": "%s",
			"status": 200,
			"method": "POST",
			"path": "/search",
			"query": {
				"pretty": "true",
				"req_id": "req-XXXXXXXXXXXXXXXXX"
			},
			"protocol": "HTTP/1.1",
			"client_ip": "127.0.0.1",
			"header": {
				"Content-Type": ["application/json"],
				"X-Request-Id": ["req-12345"],
				"User-Agent":   ["TestUserAgent"],
				"Referer":      ["TestReferer"]
			},
			"useragent": "TestUserAgent",
			"referer": "TestReferer",
			"latency": 0,
			"target": "localhost"
		}`, host)
		// assert log message
		assert.JSONEq(t, expected, buf.String())
	})

	system.RunTest(t, "ok response (with body)", func(t *testing.T) {
		var (
			rb = apitest.RequestBuilder{
				Method: http.MethodPost,
				Path:   "/search",
				Body:   []byte(`{"keyword": "hello"}`),
				Headers: map[string][]string{
					"Content-Type": {"application/json"},
					"X-Request-Id": {"req-12345"},
					"User-Agent":   {"TestUserAgent"},
					"Referer":      {"TestReferer"},
				},
				QueryParams: map[string]string{
					"req_id": "req-XXXXXXXXXXXXXXXXX",
					"pretty": "true",
				},
			}
			req      = rb.Build()
			buf      = bytes.NewBuffer(nil)
			skipPath = []string{}
			filter   = func(req api.Request) bool {
				return true
			}
			handler = func(ctx context.Context, req api.Request) api.Response {
				return api.OK(map[string]string{"message": "ok"})
			}
			middleware = New(buf).
					WithSkipPath(skipPath...).
					WithLoggingReqBodyFilter(filter).
					Middleware()
		)
		// call middleware function
		resp := middleware(handler)(context.Background(), req)
		assert.Equal(t, http.StatusOK, resp.Status())

		expected := fmt.Sprintf(`{
			"level": "info",
			"time": "2022-12-24T00:00:00+09:00",
			"server": "%s",
			"status": 200,
			"method": "POST",
			"path": "/search",
			"query": {
				"pretty": "true",
				"req_id": "req-XXXXXXXXXXXXXXXXX"
			},
			"body": {
				"keyword": "hello"
			},
			"protocol": "HTTP/1.1",
			"client_ip": "127.0.0.1",
			"header": {
				"Content-Type": ["application/json"],
				"X-Request-Id": ["req-12345"],
				"User-Agent":   ["TestUserAgent"],
				"Referer":      ["TestReferer"]
			},
			"useragent": "TestUserAgent",
			"referer": "TestReferer",
			"latency": 0,
			"target": "localhost"
		}`, host)
		// assert log message
		assert.JSONEq(t, expected, buf.String())
	})

	system.RunTest(t, "bad request response", func(t *testing.T) {
		var (
			rb = apitest.RequestBuilder{
				Method: http.MethodGet,
				Path:   "/search",
				Headers: map[string][]string{
					"Content-Type": {"application/json"},
				},
			}
			req      = rb.Build()
			buf      = bytes.NewBuffer(nil)
			skipPath = []string{}
			filter   = func(req api.Request) bool {
				return true
			}
			handler = func(ctx context.Context, req api.Request) api.Response {
				return api.BadRequest(errors.New("bad request"))
			}
			middleware = New(buf).
					WithSkipPath(skipPath...).
					WithLoggingReqBodyFilter(filter).
					Middleware()
		)
		// call middleware function
		resp := middleware(handler)(context.Background(), req)
		assert.Equal(t, http.StatusBadRequest, resp.Status())

		expected := fmt.Sprintf(`{
			"level": "warn",
			"time": "2022-12-24T00:00:00+09:00",
			"server": "%s",
			"status": 400,
			"method": "GET",
			"path": "/search",
			"query": {},
			"protocol": "HTTP/1.1",
			"client_ip": "127.0.0.1",
			"header": {
				"Content-Type": ["application/json"]
			},
			"useragent": "",
			"referer": "",
			"latency": 0,
			"target": "localhost"
		}`, host)
		// assert log message
		assert.JSONEq(t, expected, buf.String())
	})

	system.RunTest(t, "internal server error response", func(t *testing.T) {
		var (
			rb = apitest.RequestBuilder{
				Method: http.MethodGet,
				Path:   "/search",
				Headers: map[string][]string{
					"Content-Type": {"application/json"},
				},
			}
			req      = rb.Build()
			buf      = bytes.NewBuffer(nil)
			skipPath = []string{}
			filter   = func(req api.Request) bool {
				return false
			}
			handler = func(ctx context.Context, req api.Request) api.Response {
				return api.InternalServerError(errors.New("unexpected error"))
			}
			middleware = New(buf).
					WithSkipPath(skipPath...).
					WithLoggingReqBodyFilter(filter).
					Middleware()
		)
		// call middleware function
		resp := middleware(handler)(context.Background(), req)
		assert.Equal(t, http.StatusInternalServerError, resp.Status())

		expected := fmt.Sprintf(`{
			"level": "error",
			"time": "2022-12-24T00:00:00+09:00",
			"server": "%s",
			"status": 500,
			"method": "GET",
			"path": "/search",
			"query": {},
			"protocol": "HTTP/1.1",
			"client_ip": "127.0.0.1",
			"header": {
				"Content-Type": ["application/json"]
			},
			"useragent": "",
			"referer": "",
			"latency": 0,
			"target": "localhost"
		}`, host)
		// assert log message
		assert.JSONEq(t, expected, buf.String())
	})
}
