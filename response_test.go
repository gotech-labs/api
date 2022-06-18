package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/gotech-labs/api"
	"github.com/gotech-labs/core/errors"
)

func TestResponseStatus(t *testing.T) {
	for _, test := range []struct {
		name     string
		response Response
		status   int
	}{
		{
			name:     "status ok",
			response: OK("OK"),
			status:   http.StatusOK,
		},
		{
			name:     "status created",
			response: Created(`{"message": "OK"}`),
			status:   http.StatusCreated,
		},
		{
			name:     "status no content",
			response: NoContent(),
			status:   http.StatusNoContent,
		},
		{
			name:     "status bad request",
			response: BadRequest(errors.ValidationError.New("error")),
			status:   http.StatusBadRequest,
		},
		{
			name:     "status unauthorized",
			response: Unauthorized(fmt.Errorf("error")),
			status:   http.StatusUnauthorized,
		},
		{
			name:     "status not found",
			response: NotFound(fmt.Errorf("error")),
			status:   http.StatusNotFound,
		},
		{
			name:     "status proxy auth required",
			response: ProxyAuthRequired(fmt.Errorf("error")),
			status:   http.StatusProxyAuthRequired,
		},
		{
			name:     "status conflict",
			response: Conflict(fmt.Errorf("error")),
			status:   http.StatusConflict,
		},
		{
			name:     "status internal server error",
			response: InternalServerError(fmt.Errorf("error")),
			status:   http.StatusInternalServerError,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			actual := test.response
			assert.Equal(t, test.status, actual.Status())
			assert.Equal(t, map[string]string{"Content-Type": "application/json"}, actual.Headers())
		})
	}
}

func TestResponseBody(t *testing.T) {
	t.Run("string result", func(t *testing.T) {
		actual := OK("OK")
		expected := map[string]string{
			"message": "OK",
		}
		assert.Equal(t, expected, actual.Body())
	})
	t.Run("string json result", func(t *testing.T) {
		actual := OK(`{"message": "OK"}`)
		expected := json.RawMessage(`{"message": "OK"}`)
		assert.Equal(t, expected, actual.Body())
	})
	t.Run("bytes result", func(t *testing.T) {
		actual := OK([]byte(`{"message": "OK"}`))
		expected := json.RawMessage(`{"message": "OK"}`)
		assert.Equal(t, expected, actual.Body())
	})
	t.Run("struct result", func(t *testing.T) {
		type response struct {
			ID     int
			Status string
		}
		result := &response{ID: 1000, Status: "succeeded"}
		actual := OK(result)
		assert.Equal(t, result, actual.Body())
	})
	t.Run("core error result", func(t *testing.T) {
		err := errors.ValidationError.New("validation error")
		actual := BadRequest(err)
		assert.Equal(t, err, actual.Body())
	})
	t.Run("unknown error result", func(t *testing.T) {
		err := fmt.Errorf("unknown error aaaa")
		actual := InternalServerError(err)
		expected := errors.UnexpectedError.Wrap(err)
		actualErr, ok := actual.Body().(errors.Error)
		assert.True(t, ok)
		assert.Error(t, actualErr)
		assert.Equal(t, expected.Error(), actualErr.Error())
	})
}

func TestResponseCustomHeader(t *testing.T) {
	actual := OK("OK").WithHeader("X-Custom-Id", "123")
	assert.Equal(t, map[string]string{
		"Content-Type": "application/json",
		"X-Custom-Id":  "123",
	}, actual.Headers())
}
