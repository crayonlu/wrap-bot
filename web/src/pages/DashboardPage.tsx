import { useStatus } from '../lib/hooks/useQuery'
import { Activity, MessageSquare, Send } from 'lucide-react'

export default function DashboardPage() {
  const { data: status, isLoading } = useStatus()

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
        <h1>Dashboard</h1>
        <p>Bot status and statistics</p>
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
                {status?.running ? 'Running' : 'Stopped'}
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
                {status?.uptime || 0}
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
                {status?.version || 'N/A'}
              </p>
            </div>
          </div>
        </div>
      </div>

      {status && (
        <div className="dashboard__info-card">
          <h2>System Info</h2>
          <div className="dashboard__info-row">
            <span>Go Version:</span>
            <span>{status.go_version}</span>
          </div>
          <div className="dashboard__info-row">
            <span>Bot Version:</span>
            <span>{status.version}</span>
          </div>
        </div>
      )}
    </div>
  )
}
