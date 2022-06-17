package health

import (
	"context"

	"github.com/gotech-labs/api"
)

func New(path string, body interface{}) *health {
	return &health{
		path: path,
		body: body,
	}
}

type health struct {
	path string
	body interface{}
}

func (mw *health) Middleware(next api.HandlerFunc) api.HandlerFunc {
	return func(ctx context.Context, req api.Request) api.Response {
		if mw.path == req.Path() {
			return api.OK(mw.body)
		}
		// call next handler function
		return next(ctx, req)
	}
}
