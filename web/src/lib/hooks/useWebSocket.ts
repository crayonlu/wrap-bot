import { useEffect, useRef, useState } from 'react'

export type WebSocketMessage = {
  type: 'status' | 'log' | 'plugins' | 'tasks'
  data: unknown
}

export type WebSocketStatus = 'connecting' | 'connected' | 'disconnected' | 'error'

interface UseWebSocketOptions {
  enabled?: boolean
  onMessage?: (message: WebSocketMessage) => void
  onStatusChange?: (status: WebSocketStatus) => void
}

export function useWebSocket(options: UseWebSocketOptions = {}) {
  const { enabled = true, onMessage, onStatusChange } = options
  const [status, setStatus] = useState<WebSocketStatus>('disconnected')
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimeoutRef = useRef<number | null>(null)
  const reconnectAttemptsRef = useRef(0)
  const MAX_RECONNECT_ATTEMPTS = 5
  const RECONNECT_DELAY = 3000

  const connect = () => {
    if (!enabled) {
      return
    }

    try {
      const wsUrl = import.meta.env.VITE_WS_URL || `ws://${window.location.host}/api/ws`
      const token = localStorage.getItem('token')

      const ws = new WebSocket(`${wsUrl}?token=${token}`)
      wsRef.current = ws

      setStatus('connecting')
      onStatusChange?.('connecting')

      ws.onopen = () => {
        console.log('WebSocket connected')
        setStatus('connected')
        onStatusChange?.('connected')
        reconnectAttemptsRef.current = 0
      }

      ws.onmessage = (event) => {
        try {
          const message: WebSocketMessage = JSON.parse(event.data)
          onMessage?.(message)
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error)
        }
      }

      ws.onerror = (error) => {
        console.error('WebSocket error:', error)
        setStatus('error')
        onStatusChange?.('error')
      }

      ws.onclose = () => {
        console.log('WebSocket disconnected')
        setStatus('disconnected')
        onStatusChange?.('disconnected')

        // Attempt to reconnect
        if (enabled && reconnectAttemptsRef.current < MAX_RECONNECT_ATTEMPTS) {
          reconnectAttemptsRef.current++
          console.log(`Reconnecting... Attempt ${reconnectAttemptsRef.current}/${MAX_RECONNECT_ATTEMPTS}`)
          reconnectTimeoutRef.current = setTimeout(() => {
            connect()
          }, RECONNECT_DELAY)
        }
      }
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error)
      setStatus('error')
      onStatusChange?.('error')
    }
  }

  const disconnect = () => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
    }
    if (wsRef.current) {
      wsRef.current.close()
      wsRef.current = null
    }
    setStatus('disconnected')
    onStatusChange?.('disconnected')
  }

  useEffect(() => {
    if (enabled) {
      connect()
    } else {
      disconnect()
    }

    return () => {
      disconnect()
    }
  }, [enabled])

  return {
    status,
    connect,
    disconnect,
  }
}
