package chartmetric

type ParamOption func(params map[string]any)

// Limit sets the maximum number of records that should be returned. For pagination.
func Limit(n int) ParamOption {
	return func(params map[string]any) {
		params["limit"] = n
	}
}

// Offset sets the number of records to skip before starting to return records. For pagination.
func Offset(n int) ParamOption {
	return func(params map[string]any) {
		params["offset"] = n
	}
}
