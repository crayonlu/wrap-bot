package tech_push

import (
	"fmt"
	"log"

	"github.com/crayon/bot_golang/internal/config"
	"github.com/crayon/bot_golang/pkgs/feature/tech_push/handlers"
	"github.com/crayon/bot_golang/pkgs/napcat"
)

type TechPush struct {
	cfg       *config.Config
	aiService AIAnalyzer
}

type AIAnalyzer interface {
	Analyze(content string) (string, error)
}

func NewTechPush(cfg *config.Config, aiService AIAnalyzer) *TechPush {
	return &TechPush{
		cfg:       cfg,
		aiService: aiService,
	}
}

func (tp *TechPush) SendTechPush(cachedData map[string][]byte) error {
	client := NewHotAPIClient(tp.cfg.HotApiHost, tp.cfg.HotApiKey)
	napcatClient := napcat.NewClient(tp.cfg.NapCatHTTPURL, tp.cfg.NapCatHTTPToken)

	loginInfo, err := napcatClient.GetLoginInfo()
	if err != nil {
		return fmt.Errorf("failed to get bot login info: %w", err)
	}
	botQQ := loginInfo.UserID

	freshData := make(map[string][]byte)
	for name, source := range dataSources {
		data, err := client.Get(source.Endpoint)
		if err != nil {
			log.Printf("Failed to fetch %s data, using cache: %v", name, err)
			if cached, ok := cachedData[name]; ok {
				freshData[name] = cached
			}
		} else {
			freshData[name] = data
			cachedData[name] = data
		}
	}

	forwardNodes := tp.buildForwardNodes(freshData, botQQ)
	if len(forwardNodes) == 0 {
		return fmt.Errorf("no data to send")
	}

	var sendErr error
	for _, groupID := range tp.cfg.TechPushGroups {
		_, err := napcatClient.SendGroupForwardMsg(groupID, forwardNodes)
		if err != nil {
			log.Printf("Failed to send to group %d: %v", groupID, err)
			sendErr = err
		}
	}

	for _, userID := range tp.cfg.TechPushUsers {
		_, err := napcatClient.SendPrivateForwardMsg(userID, forwardNodes)
		if err != nil {
			log.Printf("Failed to send to user %d: %v", userID, err)
			sendErr = err
		}
	}

	return sendErr
}

func (tp *TechPush) buildForwardNodes(data map[string][]byte, botQQ int64) []napcat.ForwardNode {
	var nodes []napcat.ForwardNode
	var allContent string

	for name, source := range dataSources {
		rawData, ok := data[name]
		if !ok {
			continue
		}

		switch handler := source.Handler.(type) {
		case func([]byte) (*handlers.JuejinRes, error):
			res, err := handler(rawData)
			if err != nil {
				log.Printf("Failed to parse %s: %v", name, err)
				continue
			}
			nodes = append(nodes, buildGenericNodes(res.Title, res.Articles, 10, botQQ)...)
			if tp.cfg.AIEnabled {
				allContent += formatContentForAI(res.Title, res.Articles, 10)
			}

		case func([]byte) (*handlers.BilibiliRes, error):
			res, err := handler(rawData)
			if err != nil {
				log.Printf("Failed to parse %s: %v", name, err)
				continue
			}
			nodes = append(nodes, buildGenericNodes(res.Title, res.Videos, 10, botQQ)...)
			if tp.cfg.AIEnabled {
				allContent += formatContentForAI(res.Title, res.Videos, 10)
			}
		}
	}

	if tp.cfg.AIEnabled && tp.aiService != nil && allContent != "" {
		if analysis, err := tp.aiService.Analyze(allContent); err == nil {
			aiNode := napcat.NewMixedForwardNode(
				"AI 今日热点分析",
				botQQ,
				napcat.NewTextSegment("📝 "+analysis),
			)
			nodes = append([]napcat.ForwardNode{aiNode}, nodes...)
		} else {
			log.Printf("AI analysis failed: %v", err)
		}
	}

	return nodes
}
