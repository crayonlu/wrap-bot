package rss

import (
	"fmt"
	"html"
	"regexp"
	"strings"

	"github.com/crayon/wrap-bot/internal/config"
	"github.com/crayon/wrap-bot/pkgs/logger"
	"github.com/crayon/wrap-bot/pkgs/napcat"
)

type AIAnalyzer interface {
	Analyze(content string) (string, error)
}

type RssPush struct {
	cfg        *config.Config
	rssService *RssService
	aiService  AIAnalyzer
}

func NewRssPush(cfg *config.Config, aiService AIAnalyzer) *RssPush {
	return &RssPush{
		cfg:        cfg,
		rssService: NewRssService(cfg.RSSApiHost),
		aiService:  aiService,
	}
}

func (rp *RssPush) SendRssPush() error {
	napcatClient := napcat.NewClient(rp.cfg.NapCatHTTPURL, rp.cfg.NapCatHTTPToken)

	loginInfo, err := napcatClient.GetLoginInfo()
	if err != nil {
		return fmt.Errorf("failed to get bot login info: %w", err)
	}
	botQQ := loginInfo.UserID

	feeds, err := rp.rssService.FetchAllFeeds()
	if err != nil {
		return fmt.Errorf("failed to fetch RSS feeds: %w", err)
	}

	if len(feeds) == 0 {
		return fmt.Errorf("no RSS data to send")
	}

	var sendErr error
	for feedID, rss := range feeds {
		if rss.Channel == nil || len(rss.Channel.Items) == 0 {
			logger.Warn(fmt.Sprintf("Skipping empty feed: %s", feedID))
			continue
		}

		forwardNodes := rp.buildFeedNodes(rss, botQQ)
		if len(forwardNodes) == 0 {
			continue
		}

		if rp.cfg.AIEnabled && rp.aiService != nil {
			content := rp.formatFeedForAI(rss)
			if analysis, err := rp.aiService.Analyze(content); err == nil {
				aiNode := napcat.NewMixedForwardNode(
					rss.Channel.Title+" - AI分析",
					botQQ,
					napcat.NewTextSegment("📊 "+analysis),
				)
				forwardNodes = append([]napcat.ForwardNode{aiNode}, forwardNodes...)
			} else {
				logger.Error(fmt.Sprintf("AI analysis failed for %s: %v", feedID, err))
			}
		}

		for _, groupID := range rp.cfg.RssPushGroups {
			_, err := napcatClient.SendGroupForwardMsg(groupID, forwardNodes)
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to send RSS %s to group %d: %v", feedID, groupID, err))
				sendErr = err
			}
		}

		for _, userID := range rp.cfg.RssPushUsers {
			_, err := napcatClient.SendPrivateForwardMsg(userID, forwardNodes)
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to send RSS %s to user %d: %v", feedID, userID, err))
				sendErr = err
			}
		}
	}

	return sendErr
}

func (rp *RssPush) buildFeedNodes(rss *RSS, botQQ int64) []napcat.ForwardNode {
	var nodes []napcat.ForwardNode

	maxItems := len(rss.Channel.Items)
	if maxItems > 10 {
		maxItems = 10
	}

	for i := 0; i < maxItems; i++ {
		item := rss.Channel.Items[i]

		segments := []napcat.MessageSegment{
			napcat.NewTextSegment(fmt.Sprintf("📌 %s\n", item.Title)),
			napcat.NewTextSegment(fmt.Sprintf("🔗 %s\n", item.Link)),
		}

		if item.Description != "" {
			if img := extractFirstImageURL(item.Description); img != "" {
				segments = append(segments, napcat.NewImageSegment(img))
			}

			descText := stripHTML(item.Description)
			descText = strings.TrimSpace(descText)
			if descText != "" {
				segments = append(segments, napcat.NewTextSegment(fmt.Sprintf("📝 %s\n", descText)))
			}
		}

		if item.PubDate != "" {
			segments = append(segments, napcat.NewTextSegment(fmt.Sprintf("🕒 %s", item.PubDate)))
		}

		node := napcat.NewMixedForwardNode(
			rss.Channel.Title,
			botQQ,
			segments...,
		)
		nodes = append(nodes, node)
	}

	return nodes
}

func extractFirstImageURL(htmlStr string) string {
	re := regexp.MustCompile(`(?i)<img[^>]+src=["']?([^"' >]+)["' >]`)
	m := re.FindStringSubmatch(htmlStr)
	if len(m) >= 2 {
		return html.UnescapeString(m[1])
	}
	return ""
}

func stripHTML(htmlStr string) string {
	brRe := regexp.MustCompile(`(?i)<br\s*/?>`)
	s := brRe.ReplaceAllString(htmlStr, "\n")

	tagRe := regexp.MustCompile(`(?s)<[^>]*>`)
	s = tagRe.ReplaceAllString(s, "")

	s = html.UnescapeString(s)

	s = strings.ReplaceAll(s, "\r", "")
	s = regexp.MustCompile(`\n{2,}`).ReplaceAllString(s, "\n\n")
	return strings.TrimSpace(s)
}

func (rp *RssPush) formatFeedForAI(rss *RSS) string {
	content := fmt.Sprintf("【%s】\n", rss.Channel.Title)

	maxItems := len(rss.Channel.Items)
	if maxItems > 5 {
		maxItems = 5
	}

	for i := 0; i < maxItems; i++ {
		item := rss.Channel.Items[i]
		content += fmt.Sprintf("%d. %s\n", i+1, item.Title)
	}

	return content
}
