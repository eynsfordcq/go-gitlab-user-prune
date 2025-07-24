package gitlab

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Client struct {
	client  *http.Client
	baseURL string
	token   string

	// services
	Users *UserService
}

func NewClient(baseURL, token string) *Client {
	c := &Client{
		baseURL: baseURL,
		token:   token,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// init services
	c.Users = &UserService{client: c}

	return c
}

func (c *Client) newRequest(method, path string) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	return req, nil
}

func (c *Client) do(req *http.Request, v any) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// if non-200 response
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return resp, fmt.Errorf(
			"api request failed status: %d %s: %s",
			resp.StatusCode,
			resp.Status,
			body,
		)
	}

	log.Printf("request success with code %s", resp.Status)

	// decode body
	if v != nil {
		decoder := json.NewDecoder(resp.Body)
		if err = decoder.Decode(v); err != nil {
			return resp, fmt.Errorf("fail to decode response body: %w", err)
		}
	}

	return resp, nil
}
