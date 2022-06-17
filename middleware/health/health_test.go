package health_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"

	"github.com/gotech-labs/api"
	. "github.com/gotech-labs/api/middleware/health"
	apitest "github.com/gotech-labs/api/testing"
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
			req     = api.NewRequest(rb.Build())
			handler = func(ctx context.Context, req api.Request) api.Response {
				return api.InternalServerError(xerrors.New("unexpected error")) // not called
			}
		)
		// call middleware function
		resp := New("/health", nil).Middleware(handler)(context.Background(), req)

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
			req     = api.NewRequest(rb.Build())
			handler = func(ctx context.Context, req api.Request) api.Response {
				return api.BadRequest(xerrors.New("validation error")) // not called
			}
			body interface{}
		)
		// call middleware function
		resp := New("/search", body).Middleware(handler)(context.Background(), req)

		// assert response
		assert.Equal(t, http.StatusBadRequest, resp.Status())
	})
}
