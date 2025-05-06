package chartmetric

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/avast/retry-go/v4"
)

type tokenResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

type accessToken struct {
	value     string
	expiresAt time.Time
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
