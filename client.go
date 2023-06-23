package sourcify

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// RetryOptions represents options for configuring retry settings.
type RetryOptions struct {
	MaxRetries int           // The maximum number of retries.
	Delay      time.Duration // The delay between retries.
}

// RetryOption sets a configuration option for retry settings.
type RetryOption func(*RetryOptions)

// WithMaxRetries sets the maximum number of retries.
func WithMaxRetries(maxRetries int) RetryOption {
	return func(options *RetryOptions) {
		options.MaxRetries = maxRetries
	}
}

// WithDelay sets the delay between retries.
func WithDelay(delay time.Duration) RetryOption {
	return func(options *RetryOptions) {
		options.Delay = delay
	}
}

type ClientOption func(*Client)

type Client struct {
	BaseURL      string       // The base URL of the Sourcify API.
	HTTPClient   *http.Client // The HTTP client to use for making requests.
	RetryOptions RetryOptions // The retry options for the client.
	RateLimiter  *RateLimiter // The rate limiter for the client.
}

// WithHTTPClient allows you to provide your own http.Client for the Sourcify client.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.HTTPClient = client
	}
}

// WithBaseURL allows you to provide your own base URL for the Sourcify client.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.BaseURL = baseURL
	}
}

// WithRetryOptions allows you to configure retry settings for the Sourcify client.
func WithRetryOptions(options ...RetryOption) ClientOption {
	return func(c *Client) {
		for _, opt := range options {
			opt(&c.RetryOptions)
		}
	}
}

// WithRateLimit allows you to configure rate limits for the Sourcify client.
func WithRateLimit(max int, duration time.Duration) ClientOption {
	return func(c *Client) {
		c.RateLimiter = NewRateLimiter(max, duration)
	}
}

// NewClient initializes a new Sourcify client with optional configurations.
// By default, it uses the Sourcify API's base URL (https://sourcify.dev/server),
// the default http.Client, and no retry options.
func NewClient(options ...ClientOption) *Client {
	c := &Client{
		BaseURL:      "https://sourcify.dev/server",
		HTTPClient:   http.DefaultClient,
		RetryOptions: RetryOptions{},
	}

	for _, option := range options {
		option(c)
	}

	return c
}

// CallMethod calls the specified method function with the provided parameters.
// It returns the response body as a byte slice and an error if any.
func (c *Client) CallMethod(method Method) (io.ReadCloser, int, error) {
	if method.ParamType == MethodParamTypeUri {
		return c.callURIMethod(method)
	} else if method.ParamType == MethodParamTypeQueryString {
		return c.callQueryMethod(method)
	} else {
		return nil, 0, fmt.Errorf("invalid MethodParamType: %v", method.ParamType)
	}
}

// callURIMethod calls the URI-based method function with the provided parameters.
func (c *Client) callURIMethod(method Method) (io.ReadCloser, int, error) {
	// Build the URL for the API endpoint
	requestUrl, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse API base URL: %w", err)
	}

	uri, err := method.ParseUri()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse method parameters: %w", err)
	}

	requestPath, err := url.JoinPath(requestUrl.Path, uri)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse full API URL: %w", err)
	}
	requestUrl.Path = requestPath

	// Prepare the request
	req, err := http.NewRequest(method.Method, requestUrl.String(), nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	return c.doRequestWithRetry(req)
}

// callQueryMethod calls the query-based method function with the provided parameters.
func (c *Client) callQueryMethod(method Method) (io.ReadCloser, int, error) {
	// Build the URL for the API endpoint
	requestUrl, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse API base URL: %w", err)
	}

	requestPath, err := url.JoinPath(requestUrl.Path, method.URI)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse full API URL: %w", err)
	}
	requestUrl.Path = requestPath

	queryParams := method.GetQueryParams()
	requestUrl.RawQuery = queryParams.Encode()

	// Prepare the request
	req, err := http.NewRequest(method.Method, requestUrl.String(), nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	return c.doRequestWithRetry(req)
}

// doRequestWithRetry sends the HTTP request with retry according to the configured retry options.
func (c *Client) doRequestWithRetry(req *http.Request) (io.ReadCloser, int, error) {
	attempt := 0

	for {
		if c.RateLimiter != nil {
			c.RateLimiter.Wait()
		}

		attempt++
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			if attempt <= c.RetryOptions.MaxRetries {
				time.Sleep(c.RetryOptions.Delay)
				continue
			}
			return nil, 0, fmt.Errorf("failed to send HTTP request: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			if attempt <= c.RetryOptions.MaxRetries {
				time.Sleep(c.RetryOptions.Delay)
				continue
			}
			return nil, resp.StatusCode, fmt.Errorf("unexpected response status: %s", resp.Status)
		}

		return resp.Body, resp.StatusCode, nil
	}
}
