package tech_push

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/crayon/bot_golang/internal/config"
	"github.com/crayon/bot_golang/pkgs/feature/ai"
	"github.com/crayon/bot_golang/pkgs/feature/tech_push/handlers"
	"github.com/crayon/bot_golang/pkgs/napcat"
)

type DataSource struct {
	Endpoint string
	Handler  interface{}
}

type HotAPIClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

var dataSources = map[string]DataSource{
	"juejin": {
		Endpoint: "/juejin",
		Handler:  handlers.JuejinHandler,
	},
	"bilibili": {
		Endpoint: "/bilibili",
		Handler:  handlers.BilibiliHandler,
	},
}

func NewHotAPIClient(baseURL, apiKey string) *HotAPIClient {
	return &HotAPIClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *HotAPIClient) Get(endpoint string) ([]byte, error) {
	url := c.baseURL + endpoint
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
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

func SendTechPush(cfg *config.Config, cachedData map[string][]byte) error {
	client := NewHotAPIClient(cfg.HotApiHost, cfg.HotApiKey)
	napcatClient := napcat.NewClient(cfg.NapCatHTTPURL, cfg.NapCatHTTPToken)

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

	forwardNodes := buildForwardNodes(freshData, botQQ, cfg.AIEnabled)
	if len(forwardNodes) == 0 {
		return fmt.Errorf("no data to send")
	}

	var sendErr error
	for _, groupID := range cfg.TechPushGroups {
		_, err := napcatClient.SendGroupForwardMsg(groupID, forwardNodes)
		if err != nil {
			log.Printf("Failed to send to group %d: %v", groupID, err)
			sendErr = err
		}
	}

	for _, userID := range cfg.TechPushUsers {
		_, err := napcatClient.SendPrivateForwardMsg(userID, forwardNodes)
		if err != nil {
			log.Printf("Failed to send to user %d: %v", userID, err)
			sendErr = err
		}
	}

	return sendErr
}

func AnalyzeTechWithAI(title, content string) (string, error) {
	prompt := fmt.Sprintf("分析一下今天的热点\n标题：%s\n内容：%s", title, content)
	return ai.Chat(prompt, false)
}

func buildForwardNodes(data map[string][]byte, botQQ int64, aiEnabled bool) []napcat.ForwardNode {
	var nodes []napcat.ForwardNode

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
			nodes = append(nodes, buildGenericNodes(res.Title, res.Articles, 20, botQQ, aiEnabled)...)

		case func([]byte) (*handlers.BilibiliRes, error):
			res, err := handler(rawData)
			if err != nil {
				log.Printf("Failed to parse %s: %v", name, err)
				continue
			}
			nodes = append(nodes, buildGenericNodes(res.Title, res.Videos, 20, botQQ, aiEnabled)...)
		}
	}

	return nodes
}

func buildGenericNodes(sourceName string, items interface{}, limit int, botQQ int64, aiEnabled bool) []napcat.ForwardNode {
	var nodes []napcat.ForwardNode

	val := reflect.ValueOf(items)
	if val.Kind() != reflect.Slice {
		return nodes
	}

	maxItems := val.Len()
	if limit > 0 && limit < maxItems {
		maxItems = limit
	}

	for i := 0; i < maxItems; i++ {
		item := val.Index(i)
		segments := structToSegments(item)

		if aiEnabled {
			title := extractTitle(item)
			content := extractContent(segments)
			if analysis, err := AnalyzeTechWithAI(title, content); err == nil {
				segments = append(segments, napcat.NewTextSegment("\n📝 AI: "+analysis+"\n"))
			} else {
				log.Printf("AI analysis failed: %v", err)
			}
		}

		node := napcat.NewMixedForwardNode(
			sourceName,
			botQQ,
			segments...,
		)
		nodes = append(nodes, node)
	}

	return nodes
}

func extractTitle(val reflect.Value) string {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return ""
	}

	titleField := val.FieldByName("Title")
	if titleField.IsValid() && titleField.Kind() == reflect.String {
		return titleField.String()
	}
	return ""
}

func extractContent(segments []napcat.MessageSegment) string {
	var content string
	for _, seg := range segments {
		if seg.Type == "text" {
			if text, ok := seg.Data["text"].(string); ok {
				content += text
			}
		}
	}
	if len(content) > 300 {
		content = content[:300]
	}
	return content
}

func structToSegments(val reflect.Value) []napcat.MessageSegment {
	var segments []napcat.MessageSegment

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return segments
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		if !fieldValue.CanInterface() {
			continue
		}

		fieldName := field.Name

		if fieldName == "Cover" && fieldValue.Kind() == reflect.String && fieldValue.String() != "" {
			segments = append(segments, napcat.NewImageSegment(fieldValue.String()))
			segments = append(segments, napcat.NewTextSegment("\n"))
			continue
		}

		if fieldName == "MobileUrl" {
			continue
		}

		var valueStr string
		switch fieldValue.Kind() {
		case reflect.String:
			valueStr = fieldValue.String()
		case reflect.Int, reflect.Int64:
			valueStr = fmt.Sprintf("%d", fieldValue.Int())
		case reflect.Float64:
			valueStr = fmt.Sprintf("%f", fieldValue.Float())
		default:
			valueStr = fmt.Sprintf("%v", fieldValue.Interface())
		}

		segments = append(segments, napcat.NewTextSegment(fmt.Sprintf("%s: %s\n", fieldName, valueStr)))
	}

	return segments
}
