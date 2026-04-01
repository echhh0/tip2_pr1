package authclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"tip2_pr1/shared/httpx"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type verifyResponse struct {
	Valid   bool   `json:"valid"`
	Subject string `json:"subject,omitempty"`
	Error   string `json:"error,omitempty"`
}

func New(baseURL string) *Client {
	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: httpx.NewClient(3 * time.Second),
	}
}

func (c *Client) Verify(ctx context.Context, authorization string, requestID string) (bool, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v1/auth/verify", nil)
	if err != nil {
		return false, http.StatusInternalServerError, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", authorization)
	if requestID != "" {
		req.Header.Set("X-Request-ID", requestID)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, http.StatusBadGateway, fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var vr verifyResponse
		if err := json.NewDecoder(resp.Body).Decode(&vr); err != nil {
			return false, http.StatusBadGateway, fmt.Errorf("decode auth response: %w", err)
		}
		return vr.Valid, http.StatusOK, nil
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return false, http.StatusUnauthorized, nil
	}

	if resp.StatusCode >= 500 {
		return false, http.StatusBadGateway, fmt.Errorf("auth service returned %d", resp.StatusCode)
	}

	return false, http.StatusUnauthorized, nil
}
