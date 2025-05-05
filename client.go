package chartmetric

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/avast/retry-go/v4"
	"golang.org/x/time/rate"
)

const (
	defaultRetryAttempts = 3
	defaultRetryDelay    = 700 * time.Millisecond
)

var (
	errRateLimitExceeded      = errors.New("rate limit exceeded")
	errTemporarilyUnavailable = errors.New("temporarily unavailable")
)

// Client is the Chartmetric API client.
// It handles authentication, rate limiting, and making HTTP requests to the API.
type Client struct {
	refreshToken  string
	accessToken   *accessToken
	httpClient    *http.Client
	baseURL       string
	rateLimiter   *rate.Limiter
	retryAttempts uint
	retryDelay    time.Duration
}

type ClientOption func(*Client)

// NewClient is the constructor for Client. It requires a refresh token
// and can also accept various options to configure the Client.
func NewClient(refreshToken string, options ...ClientOption) *Client {
	client := &Client{
		refreshToken: refreshToken,
		httpClient: &http.Client{
			Timeout: time.Duration(10) * time.Second,
		},
		baseURL:       "https://api.chartmetric.com/api",
		rateLimiter:   rate.NewLimiter(rate.Limit(1), 1),
		retryAttempts: defaultRetryAttempts,
		retryDelay:    defaultRetryDelay,
	}

	for _, option := range options {
		option(client)
	}

	return client
}

// WithHTTPClient allows setting a custom HTTP client for the API requests.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL allows setting a custom base URL for the API requests.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithRateLimitPerSec allows setting a custom client-side rate limit.
// This would correspond to the request per second of the availed Developer API plan.
// See https://chartmetric.com/pricing (Developer Tools).
// The default is set to 1 request per second.
func WithRateLimitPerSec(rateLimitPerSec int) ClientOption {
	return func(c *Client) {
		c.rateLimiter = rate.NewLimiter(rate.Limit(rateLimitPerSec), rateLimitPerSec)
	}
}

// WithRetryAttempts allows setting a custom number of retry attempts for requests.
// Retries are typically triggered by rate limit errors.
func WithRetryAttempts(retryAttempts uint) ClientOption {
	return func(c *Client) {
		c.retryAttempts = retryAttempts
	}
}

// WithRetryDelay allows setting a custom delay between retry attempts.
func WithRetryDelay(retryDelay time.Duration) ClientOption {
	return func(c *Client) {
		c.retryDelay = retryDelay
	}
}

// GetAny is a generic GET request method that can be used to fetch any data from the API.
// This could be useful for testing. For actual API calls, consider using the specific methods provided by the Client.
func (c *Client) GetAny(ctx context.Context, path string, queryParams map[string]any) ([]byte, error) {
	responseData, err := c.requestWithRetry(ctx, http.MethodGet, path, queryParams, nil)
	if err != nil {
		return nil, fmt.Errorf("request with retry: %w", err)
	}

	return responseData, nil
}

func (c *Client) requestWithRetry(ctx context.Context, httpMethod, path string, queryParams map[string]any, body any) ([]byte, error) {
	return retry.DoWithData(
		func() ([]byte, error) {
			return c.request(ctx, httpMethod, path, queryParams, body)
		},
		retry.Context(ctx),
		retry.Attempts(c.retryAttempts),
		retry.Delay(c.retryDelay),
		retry.RetryIf(func(err error) bool {
			return errors.Is(err, errRateLimitExceeded) || errors.Is(err, errTemporarilyUnavailable)
		}),
	)
}

func (c *Client) request(ctx context.Context, httpMethod, path string, queryParams map[string]any, body any) ([]byte, error) {
	accessToken, err := c.resolveAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("resolve access token: %w", err)
	}

	jsonBody, err := buildJSONBody(body)
	if err != nil {
		return nil, fmt.Errorf("build json body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, httpMethod, c.baseURL+path, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Content-Type", "application/json")

	addQueryParams(req, queryParams)

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limiter wait: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http client do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, errRateLimitExceeded
	}
	if resp.StatusCode == http.StatusServiceUnavailable {
		return nil, errTemporarilyUnavailable
	}

	bodyBytes, err := readResponseBody(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if !isStatusSuccess(resp) {
		return nil, fmt.Errorf("received non-success response: [%d] %s", resp.StatusCode, string(bodyBytes))
	}

	return bodyBytes, nil
}

func (c *Client) resolveAccessToken(ctx context.Context) (string, error) {
	if c.accessToken == nil || time.Now().After(c.accessToken.expiresAt) {
		fetchedToken, err := retry.DoWithData(
			func() (*accessToken, error) {
				return c.fetchAccessToken(ctx)
			},
			retry.Context(ctx),
			retry.Attempts(c.retryAttempts),
			retry.Delay(c.retryDelay),
			retry.RetryIf(func(err error) bool {
				return errors.Is(err, errRateLimitExceeded) || errors.Is(err, errTemporarilyUnavailable)
			}),
		)
		if err != nil {
			return "", fmt.Errorf("fetch access token: %w", err)
		}
		c.accessToken = fetchedToken
	}

	return c.accessToken.value, nil
}

func (c *Client) fetchAccessToken(ctx context.Context) (*accessToken, error) {
	jsonBody, err := buildJSONBody(map[string]string{"refreshtoken": c.refreshToken})
	if err != nil {
		return nil, fmt.Errorf("build json body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/token", jsonBody)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http client do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, errRateLimitExceeded
	}
	if resp.StatusCode == http.StatusServiceUnavailable {
		return nil, errTemporarilyUnavailable
	}

	bodyBytes, err := readResponseBody(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if !isStatusSuccess(resp) {
		return nil, fmt.Errorf("received non-success response: [%d] %s", resp.StatusCode, string(bodyBytes))
	}

	var tokenResponse tokenResponse
	if err := json.Unmarshal(bodyBytes, &tokenResponse); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	return &accessToken{
		value:     tokenResponse.Token,
		expiresAt: time.Now().Add(time.Second * time.Duration(tokenResponse.ExpiresIn)).Add(-time.Second * 5), // with 5 sec allowance
	}, nil
}
