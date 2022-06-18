package api

// Request is ...
type Request interface {
	Method() string
	Path() string
	Body() []byte
	Headers() map[string][]string
	Header(key string) []string
	QueryParameter(key string) string
	QueryParameters() map[string]string
	PathParameter(key string) string
	ClientIP() string
	UserAgent() string
	Referer() string
	Domain() string
	Protocol() string
	Host() string
	ContentLength() int64
	Bind(obj interface{}) error
}
