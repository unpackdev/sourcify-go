package sourcify

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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
func (c *Client) CallMethod(method Method) ([]byte, int, error) {
	if method.ParamType == MethodParamTypeUri {
		return c.callURIMethod(method)
	} else if method.ParamType == MethodParamTypeQueryString {
		return c.callQueryMethod(method)
	} else {
		return nil, 0, fmt.Errorf("invalid MethodParamType: %v", method.ParamType)
	}
}

// callURIMethod calls the URI-based method function with the provided parameters.
func (c *Client) callURIMethod(method Method) ([]byte, int, error) {
	// Build the URL for the API endpoint
	apiURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse API base URL: %w", err)
	}

	requestPath, err := url.JoinPath(apiURL.Path, method.URI)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse full API URL: %w", err)
	}
	apiURL.Path = requestPath

	// Replace URI placeholders with actual values
	for _, param := range method.Params {
		if v, ok := param.Value.(string); ok {
			apiURL.Path = strings.ReplaceAll(apiURL.Path, ":"+param.Key, v)
		} else {
			return nil, 0, fmt.Errorf("invalid parameter value for URI: %s", param.Key)
		}
	}

	// Prepare the request
	req, err := http.NewRequest(method.Method, apiURL.String(), nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	return c.doRequestWithRetry(req)
}

// callQueryMethod calls the query-based method function with the provided parameters.
func (c *Client) callQueryMethod(method Method) ([]byte, int, error) {
	// Build the URL for the API endpoint
	apiURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse API base URL: %w", err)
	}
	apiURL.Path = method.URI
	queryParams := method.GetQueryParams()
	apiURL.RawQuery = queryParams.Encode()

	// Prepare the request
	req, err := http.NewRequest(method.Method, apiURL.String(), nil)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	return c.doRequestWithRetry(req)
}

// doRequestWithRetry sends the HTTP request with retry according to the configured retry options.
func (c *Client) doRequestWithRetry(req *http.Request) ([]byte, int, error) {
	attempt := 0
	for {
		attempt++
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			if attempt <= c.RetryOptions.MaxRetries {
				time.Sleep(c.RetryOptions.Delay)
				continue
			}
			return nil, 0, fmt.Errorf("failed to send HTTP request: %w", err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to read response body: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			if attempt <= c.RetryOptions.MaxRetries {
				time.Sleep(c.RetryOptions.Delay)
				continue
			}
			return nil, resp.StatusCode, fmt.Errorf("unexpected response status: %s", resp.Status)
		}

		return body, resp.StatusCode, nil
	}
}
