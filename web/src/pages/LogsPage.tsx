import { useState } from 'react'
import { useLogs } from '../lib/hooks/useQuery'
import { AlertCircle, Info, AlertTriangle } from 'lucide-react'

type LogLevel = 'error' | 'warn' | 'info' | 'debug'

const levelIcons: Record<LogLevel, typeof AlertCircle> = {
  error: AlertCircle,
  warn: AlertTriangle,
  info: Info,
  debug: Info,
}

export default function LogsPage() {
  const [selectedLevel, setSelectedLevel] = useState<string>('all')
  const [limit, setLimit] = useState(100)
  const { data: logs, isLoading } = useLogs(
    selectedLevel === 'all' ? undefined : selectedLevel,
    limit
  )

  if (isLoading) {
    return (
      <div className="loading">
        <div className="loading__spinner"></div>
      </div>
    )
  }

  return (
    <div className="logs">
      <div className="logs__header">
        <h1>Logs</h1>
        <p>System logs and events</p>
      </div>

      <div className="logs__filters">
        <select
          value={selectedLevel}
          onChange={(e) => setSelectedLevel(e.target.value)}
        >
          <option value="all">All Levels</option>
          <option value="error">Error</option>
          <option value="warn">Warning</option>
          <option value="info">Info</option>
          <option value="debug">Debug</option>
        </select>

        <select
          value={limit}
          onChange={(e) => setLimit(Number(e.target.value))}
        >
          <option value={50}>50 entries</option>
          <option value={100}>100 entries</option>
          <option value={200}>200 entries</option>
          <option value={500}>500 entries</option>
        </select>
      </div>

      <div className="logs__container">
        <div className="logs__list">
          {logs?.map((log, index) => {
            const level = (log.level || 'info') as LogLevel
            const Icon = levelIcons[level]
            return (
              <div key={index} className="logs__item">
                <div className="logs__item-content">
                  <div className={`logs__icon logs__icon--${level}`}>
                    <Icon />
                  </div>
                  <div className="logs__details">
                    <div className="logs__meta">
                      <span className="logs__timestamp">
                        {new Date(log.timestamp).toLocaleString()}
                      </span>
                      <span className={`logs__badge logs__badge--${level}`}>
                        {level.toUpperCase()}
                      </span>
                    </div>
                    <p className="logs__message">{log.message}</p>
                    {log.context && (
                      <div className="logs__context">
                        <pre>{JSON.stringify(log.context, null, 2)}</pre>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            )
          })}
        </div>
      </div>
    </div>
  )
}
