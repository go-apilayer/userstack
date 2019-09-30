package userstack

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		// obfuscate the access key when printing debug lines.
		for i := range v {
			if u, ok := v[i].(*url.URL); ok && u != nil {
				copy := new(url.URL)
				*copy = *u
				q := copy.Query()
				q.Set("access_key", "hidden")
				copy.RawQuery = q.Encode()
				v[i] = copy
			}
		}
		msg := fmt.Sprintf(format, v...)
		log.Println("go-apilayer/userstack:", msg)
	}
}

// deep copy URL
func (c *Client) deepCopyURL() *url.URL {
	u := new(url.URL)
	*u = *c.url
	return u
}

type RequestParam struct {
	// Set to your preferred output field(s). Follow docs to compose this.
	Fields string
}

// Detect looks a single User-Agent string. Context may be nil, in which case context.Background() is used.
//
// Only first param struct is used, do not pass more than one.
func (c *Client) Detect(ctx context.Context, userAgent string, params ...RequestParam) (*Stack, error) {
	u := c.deepCopyURL()

	u.Path = "detect"
	q := u.Query()
	q.Add("ua", userAgent)

	if len(params) > 0 {
		param := params[0]
		if param.Fields != "" {
			q.Add("fields", param.Fields)
		}
	}

	u.RawQuery = q.Encode()

	c.debugf("HTTP request: %v", u)

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	if ctx == nil {
		ctx = context.Background()
	}
	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		switch e := err.(type) {
		case *url.Error:
			c.debugf("HTTP %s url error: %s", e.Op, e.Err.Error())
		}
		return nil, err
	}
	defer resp.Body.Close()

	c.debugf("HTTP GET:%d header:%+v", resp.StatusCode, resp.Header)

	// Not safe to rely on HTTP status codes for unmarshalling.
	// API returns 200 for both successful and error states.

	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiErr *ApiErr
	if err := json.Unmarshal(by, &apiErr); err != nil {
		return nil, err
	}
	if apiErr.Success != nil && !*apiErr.Success {
		return nil, apiErr
	}

	var st *Stack
	if err := json.Unmarshal(by, &st); err != nil {
		return nil, err
	}

	return st, nil
}

type Stack struct {
	Ua    string     `json:"ua,omitempty"`
	Type  EntityType `json:"type,omitempty"`
	Brand string     `json:"brand,omitempty"` // Is this device.brand ?
	Name  string     `json:"name,omitempty"`  // Is this device.name ?
	URL   string     `json:"url,omitempty"`
	Os    struct {
		Name         string `json:"name,omitempty"`
		Code         string `json:"code,omitempty"`
		URL          string `json:"url,omitempty"`
		Family       string `json:"family,omitempty"`
		FamilyCode   string `json:"family_code,omitempty"`
		FamilyVendor string `json:"family_vendor,omitempty"`
		Icon         string `json:"icon,omitempty"`
		IconLarge    string `json:"icon_large,omitempty"`
	} `json:"os"`
	Device struct {
		IsMobileDevice bool       `json:"is_mobile_device,omitempty"`
		Type           DeviceType `json:"type,omitempty"`
		Brand          string     `json:"brand,omitempty"`
		BrandCode      string     `json:"brand_code,omitempty"`
		BrandURL       string     `json:"brand_url,omitempty"`
		Name           string     `json:"name,omitempty"`
	} `json:"device"`
	Browser struct {
		Name         string `json:"name,omitempty"`
		Version      string `json:"version,omitempty"`
		VersionMajor string `json:"version_major,omitempty"`
		Engine       string `json:"engine,omitempty"`
	} `json:"browser"`
	Crawler struct {
		IsCrawler bool         `json:"is_crawler,omitempty"`
		Category  CategoryType `json:"category,omitempty"`
		LastSeen  string       `json:"last_seen,omitempty"` // "2019-09-15 20:35:33"
	} `json:"crawler"`
}
