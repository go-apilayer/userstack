package userstack

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	apiHost = "api.userstack.com"
)

// defaultClient is an http client with sane defaults.
var defaultClient = &http.Client{
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,

		ExpectContinueTimeout: 10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	},
	Timeout: 60 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

// NewClient returns a userstack client. Users can modify clients with functional options.
//
// Note: if you have a non-paying account, you must specify secure: false. Only paid accounts
// get access to `https`.
func NewClient(apiKey string, secure bool, options ...Option) (*Client, error) {
	if apiKey == "" {
		b := false
		err := ApiErr{Success: &b}
		err.Err.Type = ErrMissingAccessKey
		err.Err.Code = codeFromErrorType(ErrMissingAccessKey)
		err.Err.Info = "User did not supply an access key."
		return nil, &err
	}

	c := &Client{client: defaultClient}

	u := &url.URL{
		Scheme: "http", // Unpaid accounts do not have access to https, sadly.
		Host:   apiHost,
	}
	if secure {
		u.Scheme = "https"
	}
	q := u.Query()
	q.Add("access_key", apiKey)
	u.RawQuery = q.Encode()

	c.url = u

	for _, opt := range options {
		opt(c)
	}

	return c, nil
}

type Client struct {
	client HTTPClient
	url    *url.URL

	debug bool
}

func (c *Client) debugf(format string, v ...interface{}) {
	if c.debug {
		log.Printf("go-apilayer/userstack: %v\n", v)
	}
}

func (c *Client) Detect(userAgent string) (*Stack, error) {
	// deep copy URL
	u := new(url.URL)
	*u = *c.url

	u.Path = "detect"
	q := u.Query()
	q.Add("ua", userAgent)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	c.debugf("HTTP GET:%d header:%+v", resp.StatusCode, resp.Header)

	// Invalid requests may return 200 OK with custom error body.
	// Not safe to rely on HTTP response codes for decoding.

	var st *Stack
	if err := c.decode(resp.Body, &st); err != nil {
		return nil, err
	}

	if st.ApiErr != nil && !st.ApiErr.Success {
		return nil, st.ApiErr

	}

	return st, nil
}

func (c *Client) decode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

type Stack struct {
	Ua    string     `json:"ua"`
	Type  EntityType `json:"type"`
	Brand string     `json:"brand"` // Is this device.brand ?
	Name  string     `json:"name"`  // Is this device.name ?
	URL   string     `json:"url"`
	Os    struct {
		Name         string `json:"name"`
		Code         string `json:"code"`
		URL          string `json:"url"`
		Family       string `json:"family"`
		FamilyCode   string `json:"family_code"`
		FamilyVendor string `json:"family_vendor"`
		Icon         string `json:"icon"`
		IconLarge    string `json:"icon_large"`
	} `json:"os"`
	Device struct {
		IsMobileDevice bool       `json:"is_mobile_device"`
		Type           DeviceType `json:"type"`
		Brand          string     `json:"brand"`
		BrandCode      string     `json:"brand_code"`
		BrandURL       string     `json:"brand_url"`
		Name           string     `json:"name"`
	} `json:"device"`
	Browser struct {
		Name         string `json:"name"`
		Version      string `json:"version"`
		VersionMajor string `json:"version_major"`
		Engine       string `json:"engine"`
	} `json:"browser"`
	Crawler struct {
		IsCrawler bool         `json:"is_crawler"`
		Category  CategoryType `json:"category"`
		LastSeen  string       `json:"last_seen"` // "2019-09-15 20:35:33"
	} `json:"crawler"`

	*ApiErr `json:",omitempty"`
}

// This is a workaround to prevent embedded ApiError.Error from
// shadowing Stack struct prints due to Go's native priority: error > Stringer.
// TODO(mf): github.com/golang/go/issues/34464
func (s *Stack) Error() string {
	by, _ := json.Marshal(s)
	return string(by)
}
