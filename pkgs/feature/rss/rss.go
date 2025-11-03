package rss

import (
	"fmt"

	"github.com/crayon/wrap-bot/pkgs/logger"
)

type RssService struct {
	client *RssAPIClient
}

func NewRssService(baseURL string) *RssService {
	return &RssService{
		client: NewRssAPIClient(baseURL),
	}
}

func (rs *RssService) FetchAllFeeds() (map[string]*RSS, error) {
	feeds := make(map[string]*RSS)

	for _, config := range RssConfigs {
		data, err := rs.client.Get(config.Route)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to fetch RSS feed %s: %v", config.ID, err))
			continue
		}

		rss, err := Parse(data)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to parse RSS feed %s: %v", config.ID, err))
			continue
		}

		feeds[config.ID] = rss
		logger.Info(fmt.Sprintf("Successfully fetched RSS feed: %s with %d items", config.ID, len(rss.Channel.Items)))
	}

	if len(feeds) == 0 {
		return nil, fmt.Errorf("no RSS feeds fetched successfully")
	}

	return feeds, nil
}

func (rs *RssService) FetchFeed(configID string) (*RSS, error) {
	var route string
	for _, config := range RssConfigs {
		if config.ID == configID {
			route = config.Route
			break
		}
	}

	if route == "" {
		return nil, fmt.Errorf("RSS config not found: %s", configID)
	}

	data, err := rs.client.Get(route)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS feed %s: %w", configID, err)
	}

	rss, err := Parse(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed %s: %w", configID, err)
	}

	return rss, nil
}
