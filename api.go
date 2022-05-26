package api

import "context"

type (
	// Handler - api handler function
	HandlerFunc func(context.Context, Request) Response

	// Middleware - middleware for api handler
	MiddlewareFunc func(HandlerFunc) HandlerFunc
)
