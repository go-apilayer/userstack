package userstack

import "fmt"

// UnsupportedTypeError describes an unmarshal error when encountering an unknown type.
//
// NOTE: only returned when running client in strict mode. See Options for more info.
type UnsupportedTypeError struct {
	fieldName string
	typ       string
}

func (e *UnsupportedTypeError) Error() string {
	return fmt.Sprintf("go-apilayer/userstack json: unsupported %s type: %s", e.fieldName, e.typ)
}

// ApiErr is a well formatted error returned by the userstack API.
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

// ErrorType represents common userstack API errors.
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

// UnmarshalText satisfies TextUnmarshaler
func (e *ErrorType) UnmarshalText(text []byte) error {
	enum := string(text)
	if !strictUnmarshal {
		*e = ErrorType(enum)
		return nil
	}
	switch enum {
	case "404_not_found":
		*e = ErrNotFound
	case "missing_access_key":
		*e = ErrMissingAccessKey
	case "invalid_access_key":
		*e = ErrInvalidAccessKey
	case "inactive_user":
		*e = ErrInactiveUser
	case "invalid_api_function":
		*e = ErrInvalidAPIFunction
	case "usage_limit_reached":
		*e = ErrUsageLimitReached
	case "function_access_restricted":
		*e = ErrFunctionAccessRestricted
	case "https_access_restricted":
		*e = ErrHTTPSAccessRestricted
	case "missing_user_agent":
		*e = ErrMissingUserAgent
	case "invalid_fields":
		*e = ErrInvalidFields
	case "too_many_user_agents":
		*e = ErrTooManyUserAgents
	case "batch_not_supported_on_plan":
		*e = ErrBatchNotSupportedOnPlan
	default:
		return &UnsupportedTypeError{fieldName: "error", typ: enum}
	}
	return nil
}

// codeFromErrorType maps a userstack error type into a custom code (not to be confused with
// an HTTP status code). Returns 0 if the ErrorType is invalid.
func codeFromErrorType(e ErrorType) int {
	switch e {
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
