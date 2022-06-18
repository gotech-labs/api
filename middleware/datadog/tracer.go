package datadog

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/gotech-labs/api"
	"github.com/gotech-labs/core/log"
)

func New(serviceName, env string) *datadog {
	return &datadog{
		serviceName: serviceName,
		env:         env,
		tracerOpts: []tracer.StartOption{
			tracer.WithService(serviceName),
			tracer.WithServiceVersion("1.0.0"),
			tracer.WithEnv(env),
		},
		resourceNameFunc: func(req api.Request) string {
			return fmt.Sprintf("%s %s", req.Method(), req.Path())
		},
	}
}

type datadog struct {
	serviceName      string
	env              string
	tracerOpts       []tracer.StartOption
	resourceNameFunc func(api.Request) string
}

func (mw *datadog) Middleware(next api.HandlerFunc) api.HandlerFunc {
	return func(ctx context.Context, req api.Request) (resp api.Response) {
		var (
			opts = []ddtrace.StartSpanOption{
				tracer.SpanType(ext.SpanTypeHTTP),
				tracer.ServiceName(mw.serviceName),
				tracer.ResourceName(mw.resourceNameFunc(req)),
				tracer.Tag(tags.Method, req.Method()),
				tracer.Tag(tags.URL, req.Path()),
				tracer.Measured(),
			}
			carrier = tracer.HTTPHeadersCarrier(http.Header(req.Headers()))
			span    ddtrace.Span
		)
		if spanCtx, err := tracer.Extract(carrier); err == nil {
			opts = append(opts, tracer.ChildOf(spanCtx))
		}
		// pass the span through the request context
		span, ctx = tracer.StartSpanFromContext(ctx, tags.Operation, opts...)

		defer func() {
			status := resp.Status()
			span.SetTag(tags.Status, status)
			if status >= 400 {
				span.SetTag(tags.Error, fmt.Sprintf("%d %s", status, http.StatusText(status)))
			}
			span.Finish()
		}()
		// call next handler function
		return next(ctx, req)
	}
}

func (mw *datadog) WithEnabledRuntimeMetrics() *datadog {
	mw.tracerOpts = append(mw.tracerOpts, tracer.WithRuntimeMetrics())
	return mw
}

func (mw *datadog) WithEnabledTraceLogger(writer io.Writer) *datadog {
	mw.tracerOpts = append(mw.tracerOpts, tracer.WithLogger(&datadogTraceLogger{
		Logger: log.New(writer),
	}))
	return mw
}

func (mw *datadog) StartTracer() {
	tracer.Start(mw.tracerOpts...)
}

func (mw *datadog) StopTracer() {
	tracer.Stop()
}

var tags = struct {
	Operation string
	Method    string
	URL       string
	Status    string
	Error     string
}{
	Operation: "http.request",
	Method:    ext.HTTPMethod,
	URL:       ext.HTTPURL,
	Status:    ext.HTTPCode,
	Error:     ext.Error,
}

type datadogTraceLogger struct {
	*log.Logger
}

func (l *datadogTraceLogger) Log(msg string) {
	l.Logger.Info().RawJSON("datadog", []byte(msg)).Send()
}
