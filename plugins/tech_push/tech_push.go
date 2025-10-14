package plugins

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/crayon/bot_golang/internal/config"
	"github.com/crayon/bot_golang/pkgs/bot"
	scheduler "github.com/crayon/bot_golang/pkgs/feature"
	"github.com/crayon/bot_golang/pkgs/napcat"
	"github.com/crayon/bot_golang/plugins/tech_push/handlers"
)

type DataSource struct {
	endpoint string
	handler  interface{}
}

type HotAPIClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

var dataSources = map[string]DataSource{
	"juejin": {
		endpoint: "/juejin",
		handler:  handlers.JuejinHandler,
	},
	"bilibili": {
		endpoint: "/bilibili",
		handler:  handlers.BilibiliHandler,
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

func TechPushPlugin(cfg *config.Config, sched *scheduler.Scheduler) bot.HandlerFunc {
	apiHost := cfg.HotApiHost
	apiKey := cfg.HotApiKey
	if apiHost == "" || apiKey == "" {
		log.Println("TechPushPlugin: HOT_API_HOST or HOT_API_KEY is not set, skipping plugin")
		return func(ctx *bot.Context) {}
	}

	client := NewHotAPIClient(apiHost, apiKey)

	fetchedData := make(map[string][]byte)
	for name, source := range dataSources {
		data, err := client.Get(source.endpoint)
		if err != nil {
			log.Printf("TechPushPlugin: failed to fetch %s data: %v", name, err)
		} else {
			fetchedData[name] = data
		}
	}

	if len(fetchedData) == 0 {
		log.Println("TechPushPlugin: all initial data fetch failed, skipping plugin")
		return func(ctx *bot.Context) {}
	}

	sched.At(8, 0, 0).WithID("tech-push-daily").Do(func() {
		sendDailyTechPush(cfg, client, fetchedData)
	})

	return func(ctx *bot.Context) {
		if ctx.Event.RawMessage == "/tech" || ctx.Event.RawMessage == "/techpush" {
			sendDailyTechPush(cfg, client, fetchedData)
			return
		}
		ctx.Next()
	}
}

func sendDailyTechPush(cfg *config.Config, client *HotAPIClient, cachedData map[string][]byte) {
	napcatClient := napcat.NewClient(cfg.NapCatHTTPURL, cfg.NapCatHTTPToken)

	loginInfo, err := napcatClient.GetLoginInfo()
	if err != nil {
		log.Printf("TechPushPlugin: failed to get bot login info: %v", err)
		return
	}
	botQQ := loginInfo.UserID

	freshData := make(map[string][]byte)
	for name, source := range dataSources {
		data, err := client.Get(source.endpoint)
		if err != nil {
			log.Printf("TechPushPlugin: failed to fetch fresh %s data, using cache: %v", name, err)
			if cached, ok := cachedData[name]; ok {
				freshData[name] = cached
			}
		} else {
			freshData[name] = data
			cachedData[name] = data
		}
	}

	forwardNodes := buildForwardNodes(freshData, botQQ)
	if len(forwardNodes) == 0 {
		log.Println("TechPushPlugin: no data to send")
		return
	}

	targetGroups := cfg.TechPushGroups
	for _, groupID := range targetGroups {
		_, err := napcatClient.SendGroupForwardMsg(groupID, forwardNodes)
		if err != nil {
			log.Printf("TechPushPlugin: failed to send to group %d: %v", groupID, err)
		} else {
			log.Printf("TechPushPlugin: sent daily tech push to group %d", groupID)
		}
	}

	targetUsers := cfg.TechPushUsers
	for _, userID := range targetUsers {
		_, err := napcatClient.SendPrivateForwardMsg(userID, forwardNodes)
		if err != nil {
			log.Printf("TechPushPlugin: failed to send to user %d: %v", userID, err)
		} else {
			log.Printf("TechPushPlugin: sent daily tech push to user %d", userID)
		}
	}
}

func buildForwardNodes(data map[string][]byte, botQQ int64) []napcat.ForwardNode {
	var nodes []napcat.ForwardNode

	for name, source := range dataSources {
		rawData, ok := data[name]
		if !ok {
			continue
		}

		switch handler := source.handler.(type) {
		case func([]byte) (*handlers.JuejinRes, error):
			res, err := handler(rawData)
			if err != nil {
				log.Printf("TechPushPlugin: failed to parse %s: %v", name, err)
				continue
			}
			nodes = append(nodes, buildGenericNodes(res.Title, res.Articles, 20, botQQ)...)

		case func([]byte) (*handlers.BilibiliRes, error):
			res, err := handler(rawData)
			if err != nil {
				log.Printf("TechPushPlugin: failed to parse %s: %v", name, err)
				continue
			}
			nodes = append(nodes, buildGenericNodes(res.Title, res.Videos, 20, botQQ)...)
		}
	}

	return nodes
}

func buildGenericNodes(sourceName string, items interface{}, limit int, botQQ int64) []napcat.ForwardNode {
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

		node := napcat.NewMixedForwardNode(
			sourceName,
			botQQ,
			segments...,
		)
		nodes = append(nodes, node)
	}

	return nodes
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
