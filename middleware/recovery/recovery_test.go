package recovery_test

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"

	"github.com/gotech-labs/api"
	. "github.com/gotech-labs/api/middleware/recovery"
	apitest "github.com/gotech-labs/api/testing"
	"github.com/gotech-labs/core/system"
)

func TestRecovery(t *testing.T) {
	system.RunTest(t, "handle panic error", func(t *testing.T) {
		var (
			rb = apitest.RequestBuilder{
				Method: http.MethodGet,
				Path:   "/search",
				Body:   nil,
			}
			req     = api.NewRequest(rb.Build())
			handler = func(ctx context.Context, req api.Request) api.Response {
				panic(xerrors.New("connection error"))
			}
			buf = bytes.NewBuffer(nil)
		)
		// call middleware function
		resp := New(buf).Middleware(handler)(context.Background(), req)

		expected := `{
			"level": "error",
			"time": "2022-12-24T00:00:00+09:00",
			"error": "connection error",
			"message": "panic recovered"
		}`
		// assert log message
		assert.JSONEq(t, expected, buf.String())
		// assert error response
		assert.Equal(t, http.StatusInternalServerError, resp.Status())
	})

	system.RunTest(t, "handle panic error (not error type)", func(t *testing.T) {
		var (
			rb = apitest.RequestBuilder{
				Method: http.MethodGet,
				Path:   "/search",
				Body:   nil,
			}
			req     = api.NewRequest(rb.Build())
			handler = func(ctx context.Context, req api.Request) api.Response {
				panic("string message error")
			}
			buf = bytes.NewBuffer(nil)
		)
		// call middleware function
		resp := New(buf).Middleware(handler)(context.Background(), req)

		expected := `{
			"level": "error",
			"time": "2022-12-24T00:00:00+09:00",
			"error": "string message error",
			"message": "panic recovered"
		}`
		// assert log message
		assert.JSONEq(t, expected, buf.String())
		// assert error response
		assert.Equal(t, http.StatusInternalServerError, resp.Status())
	})
}
