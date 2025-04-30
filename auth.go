package chartmetric

import "time"

type tokenResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

type accessToken struct {
	value     string
	expiresAt time.Time
}
