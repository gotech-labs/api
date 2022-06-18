package accesslog

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gotech-labs/api"
	"github.com/gotech-labs/core/log"
	"github.com/gotech-labs/core/system"
	"github.com/rs/zerolog"
)

func New(writer io.Writer) *accessLog {
	return &accessLog{
		loggingFilter:        func(_ string) bool { return false },
		loggingReqBodyFilter: func(api.Request) bool { return false },
		logger:               log.New(writer),
	}
}

type accessLog struct {
	loggingFilter        func(string) bool
	loggingReqBodyFilter func(api.Request) bool
	logger               *log.Logger
}

func (mw *accessLog) Middleware() api.MiddlewareFunc {
	return func(next api.HandlerFunc) api.HandlerFunc {
		return func(ctx context.Context, req api.Request) (resp api.Response) {
			var (
				body []byte
			)
			if mw.loggingFilter(req.Path()) {
				if mw.loggingReqBodyFilter(req) && req.ContentLength() > 0 {
					body = req.Body()
				}
				defer func(begin time.Time) {
					evt := mw.logEvent(resp.Status()).
						Int("status", resp.Status()).
						Str("method", req.Method()).
						Str("path", req.Path()).
						Interface("query", req.QueryParameters()).
						Interface("header", req.Headers()).
						Str("protocol", req.Protocol()).
						Str("client_ip", req.ClientIP()).
						Str("useragent", req.UserAgent()).
						Str("referer", req.Referer()).
						Dur("latency", system.CurrentTime().Sub(begin)).
						Str("target", req.Host()).
						Str("server", host)
					if len(body) > 0 {
						evt = evt.RawJSON("body", body)
					}
					// write access log
					evt.Send()
				}(system.CurrentTime())
			}
			// call next handler function
			return next(ctx, req)
		}
	}
}

func (mw *accessLog) WithSkipPath(paths ...string) *accessLog {
	mw.loggingFilter = func(target string) bool {
		for _, path := range paths {
			if path == target {
				return false
			}
		}
		return true
	}
	return mw
}

func (mw *accessLog) WithLoggingReqBodyFilter(filter func(api.Request) bool) *accessLog {
	mw.loggingReqBodyFilter = filter
	return mw
}

func (mw *accessLog) logEvent(status int) *zerolog.Event {
	switch {
	case status >= http.StatusBadRequest && status < http.StatusInternalServerError:
		return mw.logger.Warn()
	case status >= http.StatusInternalServerError:
		return mw.logger.Error()
	default:
		return mw.logger.Info()
	}
}

var (
	host, _ = os.Hostname()
)
