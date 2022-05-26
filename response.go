package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gotech-labs/core/errors"
)

// Response is ...
type Response interface {
	Status() int
	Body() interface{}
	BodyJSON() []byte
	Headers() map[string]string
	WithHeader(string, string) Response
}

// response is ...
type response struct {
	status  int
	body    interface{}
	headers map[string]string
}

// Status is ...
func (r *response) Status() int {
	return r.status
}

// Body is ...
func (r *response) Body() interface{} {
	switch body := r.body.(type) {
	case errors.Error:
		return body
	case error:
		return errors.UnexpectedError.Wrap(body)
	case []byte:
		return json.RawMessage(body)
	case string:
		if strings.HasPrefix(body, "{") && strings.HasSuffix(body, "}") {
			return json.RawMessage(body)
		}
		return map[string]string{
			"message": body,
		}
	default:
		return body
	}
}

// BodyJSON is ...
func (r *response) BodyJSON() []byte {
	body := r.Body()
	if body == nil {
		return nil
	}
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(r.Body())
	if err != nil {
		panic(JSONEncodeError.Wrapf(err,
			"Failed to encode json object: error=%v", err.Error()))
	}
	println(buf.String())
	return buf.Bytes()
}

// Headers is ...
func (r *response) Headers() map[string]string {
	return r.headers
}

// WithHeader is ...
func (r *response) WithHeader(key, value string) Response {
	r.headers[key] = value
	return r
}

// OK is ...
func OK(body interface{}) Response {
	return newResponse(http.StatusOK, body)
}

// Created is ...
func Created(body interface{}) Response {
	return newResponse(http.StatusCreated, body)
}

// NoContent is ...
func NoContent() Response {
	return newResponse(http.StatusNoContent, nil)
}

// BadRequest is ...
func BadRequest(err error) Response {
	return newResponse(http.StatusBadRequest, err)
}

// Unauthorized is ...
func Unauthorized(err error) Response {
	return newResponse(http.StatusUnauthorized, err)
}

// NotFound is ...
func NotFound(err error) Response {
	return newResponse(http.StatusNotFound, err)
}

// ProxyAuthRequired is ...
func ProxyAuthRequired(err error) Response {
	return newResponse(http.StatusProxyAuthRequired, err)
}

// Conflict is ...
func Conflict(err error) Response {
	return newResponse(http.StatusConflict, err)
}

// InternalServerError is ...
func InternalServerError(err error) Response {
	return newResponse(http.StatusInternalServerError, err)
}

func newResponse(status int, body interface{}) Response {
	return &response{
		status: status,
		body:   body,
		headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}
