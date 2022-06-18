package health_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"

	"github.com/gotech-labs/api"
	apitest "github.com/gotech-labs/api/http/testing"
	. "github.com/gotech-labs/api/middleware/health"
	"github.com/gotech-labs/core/system"
)

func TestHealth(t *testing.T) {
	system.RunTest(t, "success response", func(t *testing.T) {
		var (
			rb = apitest.RequestBuilder{
				Method: http.MethodGet,
				Path:   "/health",
				Body:   nil,
			}
			req     = rb.Build()
			handler = func(ctx context.Context, req api.Request) api.Response {
				return api.InternalServerError(xerrors.New("unexpected error")) // not called
			}
			middleware = New("/health", nil).Middleware()
		)
		// call middleware function
		resp := middleware(handler)(context.Background(), req)

		// assert response
		assert.Equal(t, http.StatusOK, resp.Status())
	})

	system.RunTest(t, "skip health check", func(t *testing.T) {
		var (
			rb = apitest.RequestBuilder{
				Method: http.MethodGet,
				Path:   "/health",
				Body:   nil,
			}
			req     = rb.Build()
			handler = func(ctx context.Context, req api.Request) api.Response {
				return api.BadRequest(xerrors.New("validation error")) // not called
			}
			middleware = New("/search", map[string]string{"status": "ok"}).Middleware()
		)
		// call middleware function
		resp := middleware(handler)(context.Background(), req)

		// assert response
		assert.Equal(t, http.StatusBadRequest, resp.Status())
	})
}
