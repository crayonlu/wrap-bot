package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/crayon/wrap-bot/pkgs/logger"
)

type HandlerFunc func(ctx *Context)

type Engine struct {
	handlers   []HandlerFunc
	mu         sync.RWMutex
	eventChan  chan *Event
	ctx        context.Context
	cancel     context.CancelFunc
	wsClient   WebSocketClient
	apiClient  APIClient
	maxWorkers int
	workerPool chan struct{}
	plugins    map[string]*PluginInfo
	pluginsMu  sync.RWMutex
	startTime  time.Time
}

type BotStatus struct {
	Running   bool   `json:"running"`
	Uptime    int64  `json:"uptime"`
	Version   string `json:"version"`
	GoVersion string `json:"go_version"`
}

type PluginInfo struct {
	Name        string
	Description string
	Enabled     bool
	Handler     HandlerFunc
}

type WebSocketClient interface {
	Connect(url string, token string) error
	Start(eventChan chan<- *Event) error
	Close() error
}

type APIClient interface {
	SendGroupMessage(groupID int64, message interface{}) (int32, error)
	SendPrivateMessage(userID int64, message interface{}) (int32, error)
	DeleteMessage(messageID int32) error
	GetGroupList() ([]Group, error)
	GetGroupInfo(groupID int64) (*GroupInfo, error)
	GetGroupMemberList(groupID int64) ([]GroupMember, error)
	GetFriendList() ([]Friend, error)
}

type Group struct {
	GroupID        int64  `json:"group_id"`
	GroupName      string `json:"group_name"`
	MemberCount    int32  `json:"member_count"`
	MaxMemberCount int32  `json:"max_member_count"`
}

type GroupInfo struct {
	GroupID        int64  `json:"group_id"`
	GroupName      string `json:"group_name"`
	MemberCount    int32  `json:"member_count"`
	MaxMemberCount int32  `json:"max_member_count"`
}

type GroupMember struct {
	GroupID  int64  `json:"group_id"`
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Card     string `json:"card"`
	Role     string `json:"role"`
}

type Friend struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Remark   string `json:"remark"`
}

func New() *Engine {
	ctx, cancel := context.WithCancel(context.Background())
	e := &Engine{
		handlers:   make([]HandlerFunc, 0),
		eventChan:  make(chan *Event, 100),
		ctx:        ctx,
		cancel:     cancel,
		maxWorkers: 10,
		workerPool: make(chan struct{}, 10),
		plugins:    make(map[string]*PluginInfo),
		startTime:  time.Now(),
	}

	logger.Info("Bot engine initialized")
	return e
}

func (e *Engine) GetStatus() *BotStatus {
	uptime := time.Since(e.startTime).Seconds()
	return &BotStatus{
		Running:   e.ctx.Err() == nil,
		Uptime:    int64(uptime),
		Version:   "1.0.0",
		GoVersion: runtime.Version(),
	}
}

func (e *Engine) Use(handler HandlerFunc) *Engine {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.handlers = append(e.handlers, handler)
	return e
}

func (e *Engine) SetWebSocketClient(client WebSocketClient) {
	e.wsClient = client
}

func (e *Engine) SetAPIClient(client APIClient) {
	e.apiClient = client
}

func (e *Engine) GetAPIClient() APIClient {
	return e.apiClient
}

func (e *Engine) SetMaxWorkers(max int) {
	e.maxWorkers = max
	e.workerPool = make(chan struct{}, max)
}

func (e *Engine) handleEvent(event *Event) {
	e.workerPool <- struct{}{}
	go func() {
		defer func() {
			<-e.workerPool
			if err := recover(); err != nil {
				logger.Error(fmt.Sprintf("Panic recovered in event handler: %v", err))
				log.Printf("Panic recovered in event handler: %v", err)
			}
		}()

		e.mu.RLock()
		handlers := make([]HandlerFunc, len(e.handlers))
		copy(handlers, e.handlers)
		e.mu.RUnlock()

		ctx := newContext(event, handlers)
		ctx.Next()
	}()
}

func (e *Engine) Run() error {
	if e.wsClient == nil {
		return ErrWebSocketClientNotSet
	}

	go func() {
		if err := e.wsClient.Start(e.eventChan); err != nil {
			logger.Error(fmt.Sprintf("WebSocket client error: %v", err))
			log.Printf("WebSocket client error: %v", err)
			e.cancel()
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case event := <-e.eventChan:
			e.handleEvent(event)
		case <-sigChan:
			log.Println("Shutting down gracefully...")
			e.Shutdown()
			return nil
		case <-e.ctx.Done():
			log.Println("Context cancelled, shutting down...")
			e.Shutdown()
			return e.ctx.Err()
		}
	}
}

func (e *Engine) Shutdown() {
	e.cancel()
	if e.wsClient != nil {
		e.wsClient.Close()
	}
	close(e.eventChan)
}

func (e *Engine) InjectEvent(eventData []byte) error {
	var event Event
	if err := json.Unmarshal(eventData, &event); err != nil {
		return err
	}
	e.eventChan <- &event
	return nil
}

var (
	ErrWebSocketClientNotSet = &BotError{Code: "WS_CLIENT_NOT_SET", Message: "WebSocket client not set"}
)

type BotError struct {
	Code    string
	Message string
}

func (e *BotError) Error() string {
	return e.Code + ": " + e.Message
}

func (e *Engine) RegisterPlugin(name, description string, handler HandlerFunc) {
	e.pluginsMu.Lock()
	defer e.pluginsMu.Unlock()
	e.plugins[name] = &PluginInfo{
		Name:        name,
		Description: description,
		Enabled:     true,
		Handler:     handler,
	}
	e.Use(handler)
	e.broadcastPlugins()
}

func (e *Engine) GetPlugins() map[string]*PluginInfo {
	e.pluginsMu.RLock()
	defer e.pluginsMu.RUnlock()
	result := make(map[string]*PluginInfo)
	for k, v := range e.plugins {
		result[k] = v
	}
	return result
}

func (e *Engine) TogglePlugin(name string) bool {
	e.pluginsMu.Lock()
	defer e.pluginsMu.Unlock()
	if plugin, exists := e.plugins[name]; exists {
		plugin.Enabled = !plugin.Enabled
		e.broadcastPlugins()
		return true
	}
	return false
}

func (e *Engine) IsPluginEnabled(name string) bool {
	e.pluginsMu.RLock()
	defer e.pluginsMu.RUnlock()
	if plugin, exists := e.plugins[name]; exists {
		return plugin.Enabled
	}
	return false
}

func (e *Engine) broadcastPlugins() {
}
