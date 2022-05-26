package api

import "github.com/gotech-labs/core/errors"

var (
	BindingError    = errors.TypedError("binding_error")
	JSONEncodeError = errors.TypedError("json_encode_error")
)
