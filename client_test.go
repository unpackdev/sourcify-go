package sourcify

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client := NewClient()

	assert.Equal(t, "https://sourcify.dev/server", client.BaseURL)
	assert.Equal(t, http.DefaultClient, client.HTTPClient)
	assert.Equal(t, RetryOptions{}, client.RetryOptions)
}

func TestWithBaseURL(t *testing.T) {
	client := NewClient(WithBaseURL("https://api.example.com"))

	assert.Equal(t, "https://api.example.com", client.BaseURL)
}

func TestWithHTTPClient(t *testing.T) {
	httpClient := &http.Client{Timeout: 5 * time.Second}
	client := NewClient(WithHTTPClient(httpClient))

	assert.Equal(t, httpClient, client.HTTPClient)
}

func TestWithRetryOptions(t *testing.T) {
	retryOpts := []RetryOption{
		WithMaxRetries(3),
		WithDelay(1 * time.Second),
	}
	client := NewClient(WithRetryOptions(retryOpts...))

	expectedRetryOpts := RetryOptions{
		MaxRetries: 3,
		Delay:      1 * time.Second,
	}
	assert.Equal(t, expectedRetryOpts, client.RetryOptions)
}

func TestCallMethod_URIMethod(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, world!")
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))

	method := Method{
		Method:    "GET",
		ParamType: MethodParamTypeUri,
		URI:       "/test",
	}

	resp, statusCode, err := client.CallMethod(method)
	defer func() {
		_ = resp.Close()
	}()

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)

	body, err := io.ReadAll(resp)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, world!", string(body))
}

func TestCallMethod_QueryMethod(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, world!")
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))

	method := Method{
		Method:    "GET",
		ParamType: MethodParamTypeQueryString,
		URI:       "/test",
		Params: []MethodParam{
			{Key: "param1", Value: "value1"},
			{Key: "param2", Value: "value2"},
		},
	}

	resp, statusCode, err := client.CallMethod(method)
	defer func() {
		_ = resp.Close()
	}()

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)

	body, err := io.ReadAll(resp)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, world!", string(body))
}

func TestDoRequestWithRetry_SuccessfulRequest(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, world!")
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, _, err := client.doRequestWithRetry(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	body, err := io.ReadAll(resp)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, world!", string(body))
}

func TestDoRequestWithRetry_RetriesExceeded(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	retryOpts := []RetryOption{
		WithMaxRetries(2),
		WithDelay(1 * time.Second),
	}
	client := NewClient(
		WithBaseURL(server.URL),
		WithRetryOptions(retryOpts...),
	)

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, statusCode, err := client.doRequestWithRetry(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, http.StatusInternalServerError, statusCode)
}

func TestDoRequestWithRetry_SuccessfulRetry(t *testing.T) {
	count := 0
	handler := func(w http.ResponseWriter, r *http.Request) {
		if count < 2 {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		} else {
			fmt.Fprint(w, "Hello, world!")
		}
		count++
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	retryOpts := []RetryOption{
		WithMaxRetries(2),
		WithDelay(1 * time.Second),
	}
	client := NewClient(
		WithBaseURL(server.URL),
		WithRetryOptions(retryOpts...),
	)

	req, _ := http.NewRequest("GET", server.URL, nil)
	resp, _, err := client.doRequestWithRetry(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	body, err := io.ReadAll(resp)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, world!", string(body))
}

func TestWithRateLimiting(t *testing.T) {
	client := NewClient(WithRateLimit(10, 1*time.Second))

	assert.NotNil(t, client.RateLimiter)
	assert.Equal(t, 10, client.RateLimiter.Max)
	assert.Equal(t, 1*time.Second, client.RateLimiter.Duration)
}

func TestRateLimiting(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, world!")
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithRateLimit(1, 1*time.Second),
	)

	req, _ := http.NewRequest("GET", server.URL, nil)

	// Perform first request - should pass
	resp, _, err := client.doRequestWithRetry(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
