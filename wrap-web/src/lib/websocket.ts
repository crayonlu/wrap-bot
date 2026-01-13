import type { WebSocketEvent } from '@/types/api';

type WebSocketMessageHandler = (event: WebSocketEvent) => void;
type WebSocketStatusHandler = (connected: boolean) => void;

class WebSocketClient {
  private ws: WebSocket | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 3000;
  private messageHandlers: Set<WebSocketMessageHandler> = new Set();
  private statusHandlers: Set<WebSocketStatusHandler> = new Set();
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;

  connect(token: string) {
    const wsUrl = `${import.meta.env.VITE_WS_URL || ''}/ws?token=${token}`;
    
    try {
      this.ws = new WebSocket(wsUrl);
      
      this.ws.onopen = () => {
        console.log('WebSocket connected');
        this.reconnectAttempts = 0;
        this.notifyStatusHandlers(true);
      };

      this.ws.onmessage = (event) => {
        try {
          const data: WebSocketEvent = JSON.parse(event.data);
          this.notifyMessageHandlers(data);
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
      };

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error);
      };

      this.ws.onclose = () => {
        console.log('WebSocket disconnected');
        this.notifyStatusHandlers(false);
        this.attemptReconnect(token);
      };
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error);
      this.attemptReconnect(token);
    }
  }

  disconnect() {
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    
    this.messageHandlers.clear();
    this.statusHandlers.clear();
  }

  private attemptReconnect(token: string) {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error('Max reconnection attempts reached');
      return;
    }

    this.reconnectAttempts++;
    console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`);

    this.reconnectTimer = setTimeout(() => {
      this.connect(token);
    }, this.reconnectDelay);
  }

  onMessage(handler: WebSocketMessageHandler) {
    this.messageHandlers.add(handler);
    return () => {
      this.messageHandlers.delete(handler);
    };
  }

  onStatusChange(handler: WebSocketStatusHandler) {
    this.statusHandlers.add(handler);
    return () => {
      this.statusHandlers.delete(handler);
    };
  }

  private notifyMessageHandlers(event: WebSocketEvent) {
    this.messageHandlers.forEach(handler => {
      try {
        handler(event);
      } catch (error) {
        console.error('Error in message handler:', error);
      }
    });
  }

  private notifyStatusHandlers(connected: boolean) {
    this.statusHandlers.forEach(handler => {
      try {
        handler(connected);
      } catch (error) {
        console.error('Error in status handler:', error);
      }
    });
  }

  isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }
}

export const wsClient = new WebSocketClient();
