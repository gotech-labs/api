package http

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gotech-labs/api"
	"github.com/gotech-labs/core/log"
)

func NewRequest(req *http.Request) api.Request {
	return &request{
		Request:    req,
		body:       nil,
		query:      req.URL.Query(),
		pathParams: mux.Vars(req),
	}
}

type request struct {
	*http.Request
	body       []byte
	query      url.Values
	pathParams map[string]string
}

func (r *request) Method() string {
	return r.Request.Method
}

func (r *request) Path() string {
	return r.URL.Path
}

func (r *request) Body() []byte {
	b, err := io.ReadAll(r.Request.Body)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to read request body")
	}
	return b
}

func (r *request) Headers() map[string][]string {
	params := make(map[string][]string, 0)
	for key, values := range r.Request.Header {
		params[key] = values
	}
	return params
}

func (r *request) Header(key string) []string {
	return r.Request.Header.Values(key)
}

func (r *request) QueryParameter(key string) string {
	return r.query.Get(key)
}

func (r *request) QueryParameters() map[string]string {
	params := make(map[string]string, 0)
	for key, values := range r.query {
		params[key] = strings.Join(values, ",")
	}
	return params
}

func (r *request) PathParameter(key string) string {
	if value, ok := r.pathParams[key]; ok {
		return value
	}
	return ""
}

func (r *request) ClientIP() string {
	if ip := r.Request.Header.Get(headerXForwardedFor); ip != "" {
		i := strings.IndexAny(ip, ",")
		if i > 0 {
			return strings.TrimSpace(ip[:i])
		}
		return ip
	}
	if ip := r.Request.Header.Get(headerXRealIP); ip != "" {
		return ip
	}
	ra, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ra
}

func (r *request) UserAgent() string {
	return r.Request.UserAgent()
}

func (r *request) Referer() string {
	return r.Request.Referer()
}

func (r *request) Domain() string {
	return r.Request.URL.Host
}

func (r *request) Protocol() string {
	return r.Request.Proto
}

func (r *request) Host() string {
	return r.Request.URL.Host
}

func (r *request) ContentLength() int64 {
	return r.Request.ContentLength
}

func (r *request) Bind(obj interface{}) error {
	err := json.NewDecoder(r.Request.Body).Decode(obj)
	if err != nil {
		if ute, ok := err.(*json.UnmarshalTypeError); ok {
			return api.BindingError.Wrapf(err,
				"Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset)
		} else if se, ok := err.(*json.SyntaxError); ok {
			return api.BindingError.Wrapf(err,
				"Syntax error: offset=%v, error=%v", se.Offset, se.Error())
		}
		return api.BindingError.Wrapf(err,
			"Failed to binding object: error=%v", err.Error())
	}
	return nil
}

const (
	headerXRealIP       = "X-Real-Ip"
	headerXForwardedFor = "X-Forwarded-For"
)
