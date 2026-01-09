import { useEffect } from 'react'
import { useStatus } from '../lib/hooks/useQuery'
import { useWebSocket } from '../lib/hooks/useWebSocket'
import { useWebSocketStore } from '../stores/websocket'
import { Activity, MessageSquare, Send, Wifi, WifiOff } from 'lucide-react'

export default function DashboardPage() {
  const { data: initialStatus, isLoading } = useStatus()
  const { connected } = useWebSocketStore()
  const status = useWebSocketStore((state) => state.status)
  const setStatus = useWebSocketStore((state) => state.setStatus)

  // Initialize status from API response
  useEffect(() => {
    if (initialStatus) {
      setStatus(initialStatus)
    }
  }, [initialStatus, setStatus])

  // WebSocket integration
  useWebSocket({
    enabled: true,
    onMessage: (message: any) => {
      if (message.type === 'status') {
        const updatedStatus = message.data as any
        setStatus(updatedStatus)
      }
    },
  })

  const displayStatus = status || initialStatus

  if (isLoading) {
    return (
      <div className="loading">
        <div className="loading__spinner"></div>
      </div>
    )
  }

  return (
    <div className="dashboard">
      <div className="dashboard__header">
        <div>
          <h1>Dashboard</h1>
          <p>Bot status and statistics</p>
        </div>
        <div className="dashboard__connection-status">
          {connected ? (
            <span className="dashboard__status dashboard__status--connected">
              <Wifi size={16} />
              Connected
            </span>
          ) : (
            <span className="dashboard__status dashboard__status--disconnected">
              <WifiOff size={16} />
              Disconnected
            </span>
          )}
        </div>
      </div>

      <div className="dashboard__grid">
        <div className="dashboard__card">
          <div className="dashboard__card-content">
            <div className="dashboard__icon">
              <Activity />
            </div>
            <div>
              <p className="dashboard__label">Status</p>
              <p className="dashboard__value">
                {displayStatus?.running ? 'Running' : 'Stopped'}
              </p>
            </div>
          </div>
        </div>

        <div className="dashboard__card">
          <div className="dashboard__card-content">
            <div className="dashboard__icon">
              <MessageSquare />
            </div>
            <div>
              <p className="dashboard__label">Uptime (seconds)</p>
              <p className="dashboard__value">
                {displayStatus?.uptime || 0}
              </p>
            </div>
          </div>
        </div>

        <div className="dashboard__card">
          <div className="dashboard__card-content">
            <div className="dashboard__icon">
              <Send />
            </div>
            <div>
              <p className="dashboard__label">Version</p>
              <p className="dashboard__value">
                {displayStatus?.version || 'N/A'}
              </p>
            </div>
          </div>
        </div>
      </div>

      {displayStatus && (
        <div className="dashboard__info-card">
          <h2>System Info</h2>
          <div className="dashboard__info-row">
            <span>Go Version:</span>
            <span>{displayStatus.go_version}</span>
          </div>
          <div className="dashboard__info-row">
            <span>Bot Version:</span>
            <span>{displayStatus.version}</span>
          </div>
        </div>
      )}
    </div>
  )
}
