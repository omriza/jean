package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type UsageLimit struct {
	Utilization float64   `json:"utilization"`
	ResetsAt    time.Time `json:"resets_at"`
}

type UsageResponse struct {
	FiveHour         *UsageLimit `json:"five_hour"`
	SevenDay         *UsageLimit `json:"seven_day"`
	SevenDayOAuthApps *UsageLimit `json:"seven_day_oauth_apps"`
	SevenDayOpus     *UsageLimit `json:"seven_day_opus"`
	SevenDaySonnet   *UsageLimit `json:"seven_day_sonnet"`
	IguanaNecktie    *UsageLimit `json:"iguana_necktie"`
	ExtraUsage       *UsageLimit `json:"extra_usage"`
}

type Client struct {
	sessionKey string
	orgID      string
	httpClient *http.Client
}

func NewClient(sessionKey, orgID string) *Client {
	return &Client{
		sessionKey: sessionKey,
		orgID:      orgID,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetUsage() (*UsageResponse, error) {
	url := fmt.Sprintf("https://claude.ai/api/organizations/%s/usage", c.orgID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Jean/1.0")
	req.AddCookie(&http.Cookie{
		Name:  "sessionKey",
		Value: c.sessionKey,
	})

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var usage UsageResponse
	if err := json.NewDecoder(resp.Body).Decode(&usage); err != nil {
		return nil, err
	}

	return &usage, nil
}
