package userstack

import "net/http"

var (
	// Enabled with OptionStrictMode. Enabling running this client in strict mode,
	// where custom UnmarshalText ensures types continue to be well-defined.
	// In strict mode, if the API returns a string value we cannot express as
	// a typed string constant, we fail and return
	strictUnmarshal = false
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

// OptionStrictMode forces custom UnmarshalText on well-known types.
// Will cause client to error if userstack API returns a type the
// client cannot express as a typed string constant.
func OptionStrictMode() func(*Client) {
	return func(c *Client) {
		strictUnmarshal = true
	}
}
