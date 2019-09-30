package userstack

import "net/http"

var (
	// Disabled with OptionDisableStrictMode. By default client runs in "strict" mode,
	// where custom UnmarshalText ensures types are well-defined.
	// In strict mode, if the API returns a string value that cannot be expressed as
	// a typed string constant, we fail and return a UnsupportedTypeError.
	strictUnmarshal = true
)

// HTTPClient is the interface used to send HTTP requests. Users can provide their own implementation.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Option is a functional option to modify the underlying Client.
type Option func(*Client)

// OptionHTTPClient - provide a custom http client to the client.
func OptionHTTPClient(client HTTPClient) func(*Client) {
	return func(c *Client) {
		c.client = client
	}
}

// OptionDebug enable debugging for the client.
func OptionDebug(b bool) func(*Client) {
	return func(c *Client) {
		c.debug = b
	}
}

// OptionDisableStrictMode stops custom UnmarshalText on well-known types.
// Will not cause client to error if userstack API returns a type
// that cannot be expressed as a typed string constant.
func OptionDisableStrictMode() func(*Client) {
	return func(c *Client) { strictUnmarshal = false }
}
