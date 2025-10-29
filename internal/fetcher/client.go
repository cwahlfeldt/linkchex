package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client wraps an HTTP client with retry logic and timeouts
type Client struct {
	httpClient  *http.Client
	maxRetries  int
	retryDelay  time.Duration
	userAgent   string
	rateLimiter *RateLimiter
}

// NewClient creates a new HTTP client with the specified configuration
func NewClient(timeout int, maxRetries int) *Client {
	// Configure transport for better connection handling
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   false,
	}

	return &Client{
		httpClient: &http.Client{
			Timeout:   time.Duration(timeout) * time.Second,
			Transport: transport,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// Limit redirects to 10
				if len(via) >= 10 {
					return fmt.Errorf("stopped after 10 redirects")
				}
				return nil
			},
		},
		maxRetries:  maxRetries,
		retryDelay:  1 * time.Second,
		userAgent:   "Linkchex/0.1.0 (Link Validator)",
		rateLimiter: NewRateLimiter(0), // No rate limiting by default
	}
}

// SetRateLimit sets the rate limit for requests (requests per second)
func (c *Client) SetRateLimit(requestsPerSecond float64) {
	if c.rateLimiter != nil {
		c.rateLimiter.Stop()
	}
	c.rateLimiter = NewRateLimiter(requestsPerSecond)
}

// Response contains the result of an HTTP request
type Response struct {
	StatusCode int
	Status     string
	URL        string
	FinalURL   string // After redirects
	Body       []byte
	Error      error
	Duration   time.Duration
}

// Get performs an HTTP GET request with retry logic
func (c *Client) Get(url string) *Response {
	var lastErr error
	startTime := time.Now()

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			// Wait before retrying
			time.Sleep(c.retryDelay * time.Duration(attempt))
		}

		// Apply rate limiting
		if c.rateLimiter != nil {
			c.rateLimiter.Wait()
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			lastErr = err
			continue
		}

		req.Header.Set("User-Agent", c.userAgent)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		// Success - read response
		body := make([]byte, 0)
		if resp.Body != nil {
			defer resp.Body.Close()
			// Read response body efficiently
			body, _ = io.ReadAll(resp.Body)
		}

		duration := time.Since(startTime)

		return &Response{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			URL:        url,
			FinalURL:   resp.Request.URL.String(),
			Body:       body,
			Error:      nil,
			Duration:   duration,
		}
	}

	// All retries failed
	duration := time.Since(startTime)
	return &Response{
		StatusCode: 0,
		Status:     "Failed",
		URL:        url,
		FinalURL:   url,
		Body:       nil,
		Error:      fmt.Errorf("failed after %d attempts: %w", c.maxRetries+1, lastErr),
		Duration:   duration,
	}
}

// Head performs an HTTP HEAD request (lightweight check)
func (c *Client) Head(url string) *Response {
	var lastErr error
	startTime := time.Now()

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(c.retryDelay * time.Duration(attempt))
		}

		// Apply rate limiting
		if c.rateLimiter != nil {
			c.rateLimiter.Wait()
		}

		req, err := http.NewRequest("HEAD", url, nil)
		if err != nil {
			lastErr = err
			continue
		}

		req.Header.Set("User-Agent", c.userAgent)

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		if resp.Body != nil {
			resp.Body.Close()
		}

		duration := time.Since(startTime)

		return &Response{
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			URL:        url,
			FinalURL:   resp.Request.URL.String(),
			Body:       nil,
			Error:      nil,
			Duration:   duration,
		}
	}

	duration := time.Since(startTime)
	return &Response{
		StatusCode: 0,
		Status:     "Failed",
		URL:        url,
		FinalURL:   url,
		Body:       nil,
		Error:      fmt.Errorf("failed after %d attempts: %w", c.maxRetries+1, lastErr),
		Duration:   duration,
	}
}
