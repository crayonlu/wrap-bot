package rss

import (
	"fmt"
	"net/http"
	"time"
	"io"
)

type RssAPIClient struct {
	baseURL string
	client  *http.Client
}

func NewRssAPIClient(baseURL string) *RssAPIClient {
	return &RssAPIClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *RssAPIClient) Get(route string) ([]byte, error) {
	url := c.baseURL + route
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}
	
	return body, nil
}