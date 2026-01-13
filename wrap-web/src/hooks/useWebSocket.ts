import { useEffect, useCallback } from 'react';
import { useAuthStore } from '@/stores/auth';
import { useBotStore } from '@/stores/bot';
import { wsClient } from '@/lib/websocket';
import type { WebSocketEvent } from '@/types/api';

export function useWebSocket() {
  const { token } = useAuthStore();
  const { setStatus, setPlugins, setTasks, addLog } = useBotStore();

  const handleMessage = useCallback((event: WebSocketEvent) => {
    switch (event.type) {
      case 'status':
        setStatus(event.data);
        break;
      case 'plugins':
        setPlugins(event.data);
        break;
      case 'tasks':
        setTasks(event.data);
        break;
      case 'log':
        addLog(event.data);
        break;
    }
  }, [setStatus, setPlugins, setTasks, addLog]);

  useEffect(() => {
    if (!token) return;

    wsClient.connect(token);

    const unsubscribe = wsClient.onMessage(handleMessage);

    return () => {
      unsubscribe();
      wsClient.disconnect();
    };
  }, [token, handleMessage]);

  return {
    isConnected: wsClient.isConnected(),
  };
}
