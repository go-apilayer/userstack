[![](https://godoc.org/github.com/go-apilayer/userstack?status.svg)](http://godoc.org/github.com/go-apilayer/userstack)

## userstack

`userstack` is a Go client library for the [apilayer userstack](https://userstack.com/) service, which provides instant user-agent string lookups.

To use this client you will need an **API access key**. The free tier supports 10,000 monthly lookups, but over `http-only`. To get HTTPS Encryption you need a Basic plan or higher 🤷‍♂️.

The official documentation can be found [here](https://userstack.com/documentation).

If you are on a free account, initialize the client with secure `false`. Otherwise a `105` error:

> Access Restricted - Your current Subscription Plan does not support HTTPS Encryption.

---

### Technical bits: 

Well-defined types are available: `entity`, `device`, `category` and `api error`.

Users can supply their own HTTP Client implementation through options, otherwise a default client is used.

This library will return a custom error, `*ApiErr`, which callers can assert to get at the raw code, type and info. If using go1.13 use `errors.As` for assertions, otherwise a regular type switch will do. This is especially useful for `104` or `ErrUsageLimitReached` errors: monthly usage limit has been exceeded.

Crawler information is only available for Basic Plan or higher. So if you don't see this in the response, check plan.

### Example usage - single lookup (HTTP GET)

```go
client, err := userstack.NewClient(key, false)
if err != nil {
    // handler err
}

ua := "Mozilla/5.0 (Linux; Android 9; Pixel 2 Build/PQ3A.190801.002; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/76.0.3809.132 Mobile Safari/537.36 Instagram 109.0.0.18.124 Android (28/9; 420dpi; 1080x1794; Google/google; Pixel 2; walleye; walleye; en_US; 170693979)"

stack, err := client.Detect(ua)
if err != nil {
    var e *userstack.ApiErr
    if errors.As(err, &e) { 
        // handler error of type ApiErr
    }
    // handler err
}

fmt.Printf("device: %s\nentity: %s\ncrawler: %s", stack.Device.Type, stack.Type, stack.Crawler.Category)
// device: smartphone
// entity: mobile-browser
// crawler: unknown
```
