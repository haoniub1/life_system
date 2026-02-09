package bark

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client handles Bark push notifications
type Client struct {
	serverURL  string // e.g., "https://api.day.app" or self-hosted URL
	httpClient *http.Client
}

// PushOptions contains optional parameters for push notifications
type PushOptions struct {
	Title     string // Push title (larger font)
	Subtitle  string // Push subtitle
	Sound     string // Notification sound (e.g., "alarm", "bell", "birdsong")
	Icon      string // Custom icon URL
	Group     string // Notification group
	URL       string // Click to open URL
	Level     string // "active", "timeSensitive", "passive", "critical"
	Call      bool   // Repeat sound for 30 seconds (like phone call)
	Badge     int    // App badge number
	Copy      string // Text to copy
	AutoCopy  bool   // Auto copy to clipboard
}

// Response from Bark API
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewClient creates a new Bark client
// serverURL should be like "https://api.day.app" (without trailing slash)
func NewClient(serverURL string) *Client {
	// Remove trailing slash
	serverURL = strings.TrimRight(serverURL, "/")
	
	return &Client{
		serverURL: serverURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Push sends a notification to the device
func (c *Client) Push(deviceKey string, body string, opts *PushOptions) error {
	if opts == nil {
		opts = &PushOptions{}
	}

	// Build form data
	data := url.Values{}
	data.Set("body", body)
	
	if opts.Title != "" {
		data.Set("title", opts.Title)
	}
	if opts.Subtitle != "" {
		data.Set("subtitle", opts.Subtitle)
	}
	if opts.Sound != "" {
		data.Set("sound", opts.Sound)
	}
	if opts.Icon != "" {
		data.Set("icon", opts.Icon)
	}
	if opts.Group != "" {
		data.Set("group", opts.Group)
	}
	if opts.URL != "" {
		data.Set("url", opts.URL)
	}
	if opts.Level != "" {
		data.Set("level", opts.Level)
	}
	if opts.Call {
		data.Set("call", "1")
	}
	if opts.Badge > 0 {
		data.Set("badge", fmt.Sprintf("%d", opts.Badge))
	}
	if opts.Copy != "" {
		data.Set("copy", opts.Copy)
	}
	if opts.AutoCopy {
		data.Set("autoCopy", "1")
	}

	// Send POST request
	pushURL := fmt.Sprintf("%s/%s", c.serverURL, deviceKey)
	
	req, err := http.NewRequest("POST", pushURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var result Response
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w, body: %s", err, string(bodyBytes))
	}

	if result.Code != 200 {
		return fmt.Errorf("bark API error: %s (code: %d)", result.Message, result.Code)
	}

	log.Printf("Bark push sent successfully to %s", deviceKey[:8]+"***")
	return nil
}

// PushAlarm sends an alarm-style notification (repeating sound for 30 seconds)
func (c *Client) PushAlarm(deviceKey string, title, body string) error {
	return c.Push(deviceKey, body, &PushOptions{
		Title: title,
		Sound: "alarm",
		Level: "timeSensitive",
		Call:  true, // Repeat sound for 30 seconds
	})
}

// PushUrgent sends a time-sensitive notification
func (c *Client) PushUrgent(deviceKey string, title, body string) error {
	return c.Push(deviceKey, body, &PushOptions{
		Title: title,
		Sound: "bell",
		Level: "timeSensitive",
	})
}

// PushCritical sends a critical alert (ignores silent/DND mode)
func (c *Client) PushCritical(deviceKey string, title, body string) error {
	return c.Push(deviceKey, body, &PushOptions{
		Title: title,
		Sound: "alarm",
		Level: "critical",
	})
}

// PushSilent sends a passive notification (no sound, no wake screen)
func (c *Client) PushSilent(deviceKey string, title, body string) error {
	return c.Push(deviceKey, body, &PushOptions{
		Title: title,
		Level: "passive",
	})
}
