package tablecache

import "errors"

var (
	NoFieldNamedValueError = errors.New(`no field named "value"`)
	ValueNotStringError    = errors.New(`value is not string`)
	KeyNotFoundError       = errors.New(`key not found`)
)
