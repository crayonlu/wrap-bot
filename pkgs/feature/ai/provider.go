package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Provider interface {
	Complete(req ChatRequest) (*ChatResponse, error)
}

type HTTPProvider struct {
	apiURL string
	apiKey string
	client *http.Client
}

func NewHTTPProvider(apiURL, apiKey string) *HTTPProvider {
	return &HTTPProvider{
		apiURL: apiURL,
		apiKey: apiKey,
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

func (p *HTTPProvider) Complete(reqBody ChatRequest) (*ChatResponse, error) {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", p.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
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

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, err
	}

	return &chatResp, nil
}
