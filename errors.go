package userstack

import "fmt"

type ApiErr struct {
	// Most unfortunate this "success" is also not returned in the successful state.
	Success *bool `json:"success"`

	Err struct {
		Code int       `json:"code,omitempty"`
		Type ErrorType `json:"type,omitempty"`
		Info string    `json:"info,omitempty"`
	} `json:"error,omitempty"`
}

func (e *ApiErr) Error() string {
	return fmt.Sprintf("%d: %s", e.Err.Code, e.Err.Info)
}

// ErrorType represents a userstack error type.
type ErrorType string

const (
	ErrNotFound                 ErrorType = "404_not_found"
	ErrMissingAccessKey         ErrorType = "missing_access_key"
	ErrInvalidAccessKey         ErrorType = "invalid_access_key"
	ErrInactiveUser             ErrorType = "inactive_user"
	ErrInvalidAPIFunction       ErrorType = "invalid_api_function"
	ErrUsageLimitReached        ErrorType = "usage_limit_reached"
	ErrFunctionAccessRestricted ErrorType = "function_access_restricted"
	ErrHTTPSAccessRestricted    ErrorType = "https_access_restricted"
	ErrMissingUserAgent         ErrorType = "missing_user_agent"
	ErrInvalidFields            ErrorType = "invalid_fields"
	ErrTooManyUserAgents        ErrorType = "too_many_user_agents"
	ErrBatchNotSupportedOnPlan  ErrorType = "batch_not_supported_on_plan"
)

func (e ErrorType) Error() string {
	return string(e)
}

// MarshalText satisfies TextMarshaler
func (e ErrorType) MarshalText() ([]byte, error) {
	return []byte(e.Error()), nil
}

// UnmarshalText satisfies TextUnmarshaler
func (e *ErrorType) UnmarshalText(text []byte) error {
	typ := ErrorType(text)
	switch typ {
	case ErrNotFound:
		*e = typ
	case ErrMissingAccessKey:
		*e = typ
	case ErrInvalidAccessKey:
		*e = typ
	case ErrInactiveUser:
		*e = typ
	case ErrInvalidAPIFunction:
		*e = typ
	case ErrUsageLimitReached:
		*e = typ
	case ErrFunctionAccessRestricted:
		*e = typ
	case ErrHTTPSAccessRestricted:
		*e = typ
	case ErrMissingUserAgent:
		*e = typ
	case ErrInvalidFields:
		*e = typ
	case ErrTooManyUserAgents:
		*e = typ
	case ErrBatchNotSupportedOnPlan:
		*e = typ
	default:
		return fmt.Errorf("unknown userstack api error type: %s", typ)
	}

	return nil
}

// codeFromErrorType maps a userstack error type into a custom code (not to be confused with
// an HTTP status code). Returns 0 if the ErrorType is invalid.
func codeFromErrorType(typ ErrorType) int {
	switch typ {
	case ErrNotFound:
		return 404 // User requested a resource which does not exist.
	case ErrMissingAccessKey, ErrInvalidAccessKey:
		// User supplied an invalid access key.
		// or
		// User did not supply an access key.
		return 101
	case ErrInactiveUser:
		return 102 // User account is inactive or blocked.
	case ErrInvalidAPIFunction:
		return 103 // User requested a non-existent API function.
	case ErrUsageLimitReached:
		return 104 // User has reached his subscription's monthly request allowance.
	case ErrFunctionAccessRestricted, ErrHTTPSAccessRestricted:
		// The user's current subscription plan does not support HTTPS.
		// or
		// The user's current subscription does not support this API function.
		return 105
	case ErrMissingUserAgent:
		return 301 // No User-Agent string has been specified.
	case ErrInvalidFields:
		return 302 // One or more invalid output fields have been specified.
	case ErrTooManyUserAgents:
		return 303 // Too many User-Agent strings have been specified in a single Bulk request.
	case ErrBatchNotSupportedOnPlan:
		return 304 // Requests to the Bulk endpoint are not supported at this subscription level.
	default:
		return 0 // Not valid!
	}
}
