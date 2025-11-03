import { useState, useEffect, useRef } from 'react'
import { useLogs } from '../lib/hooks/useQuery'
import { AlertCircle, Info, AlertTriangle, Pause, Play } from 'lucide-react'
import Selector from '../components/Selector'

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
  const [autoScroll, setAutoScroll] = useState(true)
  const logsEndRef = useRef<HTMLDivElement>(null)
  const { data: logs, isLoading } = useLogs(
    selectedLevel === 'all' ? undefined : selectedLevel,
    limit
  )

  useEffect(() => {
    if (autoScroll && logsEndRef.current) {
      logsEndRef.current.scrollIntoView({ behavior: 'smooth' })
    }
  }, [logs, autoScroll])

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
        <Selector
          value={selectedLevel}
          onChange={(value) => setSelectedLevel(value as string)}
          options={[
            { value: 'all', label: 'All Levels' },
            { value: 'error', label: 'Error' },
            { value: 'warn', label: 'Warning' },
            { value: 'info', label: 'Info' },
            { value: 'debug', label: 'Debug' },
          ]}
        />

        <Selector
          value={limit}
          onChange={(value) => setLimit(value as number)}
          options={[
            { value: 50, label: '50 entries' },
            { value: 100, label: '100 entries' },
            { value: 200, label: '200 entries' },
            { value: 500, label: '500 entries' },
          ]}
        />

        <button
          onClick={() => setAutoScroll(!autoScroll)}
          className="login-page__button"
          style={{ marginLeft: '1rem' }}
        >
          {autoScroll ? <Pause style={{ width: '1rem', height: '1rem' }} /> : <Play style={{ width: '1rem', height: '1rem' }} />}
          {autoScroll ? 'Pause' : 'Resume'}
        </button>
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
          <div ref={logsEndRef} />
        </div>
      </div>
    </div>
  )
}
