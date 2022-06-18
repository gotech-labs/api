package http_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	. "github.com/gotech-labs/api/http"
	apitest "github.com/gotech-labs/api/testing"
)

func TestNewRequest(t *testing.T) {
	var (
		rb = apitest.RequestBuilder{
			Method: http.MethodGet,
			Path:   "/events/123",
			PathParams: map[string]string{
				"id": "123",
			},
			QueryParams: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			Body: []byte(`{
					"level": "%s",
					"time": "2022-12-24T00:00:00+09:00",
					"message": "%s message"
				}`),
			Headers: map[string][]string{
				"Content-Type": {"application/json"},
				"X-Request-Id": {"req-12345"},
				"User-Agent":   {"TestUserAgent"},
				"Referer":      {"TestReferer"},
			},
		}
	)
	actual := NewRequest(rb.Build())

	// assert log message
	assert.Equal(t, rb.Method, actual.Method())
	assert.Equal(t, rb.Path, actual.Path())
	assert.Equal(t, rb.Body, actual.Body())
	assert.Equal(t, rb.Headers, actual.Headers())
	assert.Equal(t, rb.Headers["X-Request-Id"], actual.Header("X-Request-Id"))
	assert.Equal(t, rb.QueryParams, actual.QueryParameters())
	assert.Equal(t, rb.QueryParams["key1"], actual.QueryParameter("key1"))
	assert.Equal(t, "", actual.QueryParameter("unknown"))
	assert.Equal(t, "123", actual.PathParameter("id"))
	assert.Equal(t, "", actual.PathParameter("unknown"))
	assert.Equal(t, "127.0.0.1", actual.ClientIP())
	assert.Equal(t, "TestUserAgent", actual.UserAgent())
	assert.Equal(t, "TestReferer", actual.Referer())
	assert.Equal(t, "localhost", actual.Domain())
	assert.Equal(t, "HTTP/1.1", actual.Protocol())
	assert.Equal(t, "localhost", actual.Host())
	assert.Equal(t, int64(len(rb.Body)), actual.ContentLength())
}

func TestBindRequest(t *testing.T) {
	input := struct {
		ID   int
		Name string
		Time time.Time
	}{}

	t.Run("success binding", func(t *testing.T) {
		req := NewRequest(apitest.RequestBuilder{
			Method: http.MethodPost,
			Path:   "/events",
			Body: []byte(`{
				"id": 12345,
				"name": "Michael Jordan",
				"time": "2022-12-24T00:00:00+09:00"
			}`),
		}.Build())

		err := req.Bind(&input)
		if assert.NoError(t, err) {
			assert.Equal(t, 12345, input.ID)
			assert.Equal(t, "Michael Jordan", input.Name)
			assert.Equal(t, "2022-12-24T00:00:00+09:00", input.Time.Format(time.RFC3339))
		}
	})

	t.Run("binding error (Syntax error)", func(t *testing.T) {
		req := NewRequest(apitest.RequestBuilder{
			Method: http.MethodPost,
			Path:   "/events",
			// invalid field type
			Body: []byte(`{
				"id": "12345",
				"name": "Michael Jordan",
				"time": "2022-12-24T00:00:00+09:00"
			}`),
		}.Build())

		err := req.Bind(&input)
		if assert.Error(t, err) {
			assert.Equal(t, "Unmarshal type error: expected=int, got=string, field=ID, offset=19", err.Error())
		}
	})

	t.Run("binding error (Syntax error)", func(t *testing.T) {
		req := NewRequest(apitest.RequestBuilder{
			Method: http.MethodPost,
			Path:   "/events",
			// invalid json format
			Body: []byte(`{
				"id": 12345,
				"name": "Michael Jordan",
				"time": "2022-12-24T00:00:00+09:00",,,,
			}`),
		}.Build())

		err := req.Bind(&input)
		if assert.Error(t, err) {
			assert.Equal(t, "Syntax error: offset=90, error=invalid character ',' looking for beginning of object key string", err.Error())
		}
	})

	t.Run("binding error (empty json)", func(t *testing.T) {
		req := NewRequest(apitest.RequestBuilder{
			Method: http.MethodPost,
			Path:   "/events",
			// invalid json format
			Body: []byte(``),
		}.Build())

		err := req.Bind(&input)
		if assert.Error(t, err) {
			assert.Equal(t, "Failed to binding object: error=EOF", err.Error())
		}
	})
}
