package plugins

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type SerpAPIClient struct {
	apiKey string
	client *http.Client
}

func NewSerpAPIClient(apiKey string) *SerpAPIClient {
	return &SerpAPIClient{
		apiKey: apiKey,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *SerpAPIClient) Search(ctx context.Context, query string, numResults int) (*SearchResult, error) {
	if numResults == 0 {
		numResults = 5
	}

	params := url.Values{}
	params.Set("q", query)
	params.Set("num", fmt.Sprintf("%d", numResults))
	params.Set("engine", "google")
	params.Set("api_key", c.apiKey)

	fullURL := "https://serpapi.com/search?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := readAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("SerpAPI error %d: %s", resp.StatusCode, string(body))
	}

	var searchResp struct {
		OrganicResults []struct {
			Title   string `json:"title"`
			Link    string `json:"link"`
			Snippet string `json:"snippet"`
		} `json:"organic_results"`
	}

	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, err
	}

	result := &SearchResult{
		Query:   query,
		Results: make([]SearchItem, 0, len(searchResp.OrganicResults)),
	}

	for _, item := range searchResp.OrganicResults {
		result.Results = append(result.Results, SearchItem{
			Title:   item.Title,
			URL:     item.Link,
			Snippet: item.Snippet,
		})
	}

	return result, nil
}

type SearchResult struct {
	Query   string
	Results []SearchItem
}

type SearchItem struct {
	Title   string
	URL     string
	Snippet string
}

func readAll(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}
