package recovery

import (
	"context"
	"fmt"
	"io"

	"github.com/gotech-labs/api"
	"github.com/gotech-labs/core/log"
)

func New(writer io.Writer) *recovery {
	return &recovery{
		logger: log.New(writer),
	}
}

type recovery struct {
	logger *log.Logger
}

func (mw *recovery) Middleware() api.MiddlewareFunc {
	return func(next api.HandlerFunc) api.HandlerFunc {
		return func(ctx context.Context, req api.Request) (resp api.Response) {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					mw.logger.Error().Stack().Err(err).Msg("panic recovered")
					resp = api.InternalServerError(err)
				}
			}()
			// call next handler function
			return next(ctx, req)
		}
	}
}
