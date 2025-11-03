package napcat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/crayon/wrap-bot/pkgs/bot"
	"github.com/crayon/wrap-bot/pkgs/logger"
	"github.com/gorilla/websocket"
)

type WSClient struct {
	url       string
	token     string
	conn      *websocket.Conn
	mu        sync.Mutex
	connected bool
	done      chan struct{}
}

func NewWSClient(url, token string) *WSClient {
	return &WSClient{
		url:   url,
		token: token,
		done:  make(chan struct{}),
	}
}

func (ws *WSClient) Connect(url string, token string) error {
	ws.url = url
	ws.token = token
	return ws.connect()
}

func (ws *WSClient) connect() error {
	header := http.Header{}
	if ws.token != "" {
		header.Set("Authorization", "Bearer "+ws.token)
	}

	conn, _, err := websocket.DefaultDialer.Dial(ws.url, header)
	if err != nil {
		return fmt.Errorf("failed to connect to websocket: %w", err)
	}

	ws.mu.Lock()
	ws.conn = conn
	ws.connected = true
	ws.mu.Unlock()

	logger.Info(fmt.Sprintf("WebSocket connected to %s", ws.url))
	return nil
}

func (ws *WSClient) Start(eventChan chan<- *bot.Event) error {
	if err := ws.connect(); err != nil {
		return err
	}

	go ws.heartbeat()

	for {
		select {
		case <-ws.done:
			return nil
		default:
			_, message, err := ws.conn.ReadMessage()
			if err != nil {
				logger.Error(fmt.Sprintf("WebSocket read error: %v", err))
				if ws.reconnect() != nil {
					return err
				}
				continue
			}

			var event bot.Event
			if err := json.Unmarshal(message, &event); err != nil {
				logger.Error(fmt.Sprintf("Failed to unmarshal event: %v", err))
				continue
			}

			select {
			case eventChan <- &event:
			case <-ws.done:
				return nil
			}
		}
	}
}

func (ws *WSClient) reconnect() error {
	ws.mu.Lock()
	ws.connected = false
	if ws.conn != nil {
		ws.conn.Close()
	}
	ws.mu.Unlock()

	for i := 0; i < 10; i++ {
		logger.Warn(fmt.Sprintf("Attempting to reconnect... (attempt %d/10)", i+1))
		if err := ws.connect(); err == nil {
			logger.Info("Reconnected successfully")
			return nil
		}
		time.Sleep(time.Second * time.Duration(i+1))
	}

	return fmt.Errorf("failed to reconnect after 10 attempts")
}

func (ws *WSClient) heartbeat() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ws.mu.Lock()
			if ws.connected && ws.conn != nil {
				if err := ws.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					logger.Warn(fmt.Sprintf("Heartbeat failed: %v", err))
					ws.connected = false
				}
			}
			ws.mu.Unlock()
		case <-ws.done:
			return
		}
	}
}

func (ws *WSClient) Close() error {
	close(ws.done)

	ws.mu.Lock()
	defer ws.mu.Unlock()

	if ws.conn != nil {
		ws.connected = false
		return ws.conn.Close()
	}

	return nil
}
