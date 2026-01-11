import { useEffect } from 'react'
import { useAIStats } from '../lib/hooks/useQuery'
import { useAIStore } from '../stores/ai'
import { Activity, CheckCircle, XCircle, Clock, TrendingUp } from 'lucide-react'

export default function AIStatsPage() {
  const { data: stats, isLoading } = useAIStats()
  const { setStats } = useAIStore()

  useEffect(() => {
    if (stats) {
      setStats(stats)
    }
  }, [stats, setStats])

  if (isLoading) {
    return (
      <div className="loading">
        <div className="loading__spinner"></div>
      </div>
    )
  }

  const displayStats = stats || {
    total_calls: 0,
    tool_usage: {},
    success_rate: 100,
    recent_calls: [],
  }

  const sortedTools = Object.entries(displayStats.tool_usage)
    .sort(([, a], [, b]) => b - a)
    .slice(0, 10)

  return (
    <div className="config">
      <div className="config__header">
        <div>
          <h1>AI Statistics</h1>
          <p>AI usage statistics and tool performance</p>
        </div>
      </div>

      <div className="dashboard__grid">
        <div className="dashboard__card">
          <div className="dashboard__card-content">
            <div className="dashboard__icon">
              <Activity />
            </div>
            <div>
              <p className="dashboard__label">Total Calls</p>
              <p className="dashboard__value">{displayStats.total_calls}</p>
            </div>
          </div>
        </div>

        <div className="dashboard__card">
          <div className="dashboard__card-content">
            <div className="dashboard__icon" style={{ color: '#22c55e' }}>
              <CheckCircle />
            </div>
            <div>
              <p className="dashboard__label">Success Rate</p>
              <p className="dashboard__value">{displayStats.success_rate.toFixed(1)}%</p>
            </div>
          </div>
        </div>

        <div className="dashboard__card">
          <div className="dashboard__card-content">
            <div className="dashboard__icon" style={{ color: '#3b82f6' }}>
              <TrendingUp />
            </div>
            <div>
              <p className="dashboard__label">Active Tools</p>
              <p className="dashboard__value">{Object.keys(displayStats.tool_usage).length}</p>
            </div>
          </div>
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(400px, 1fr))', gap: '1.5rem', marginTop: '1.5rem' }}>
        <div className="logs__container">
          <h2 style={{ marginBottom: '1rem' }}>Tool Usage</h2>
          <div className="logs__list">
            {sortedTools.length > 0 ? (
              sortedTools.map(([toolName, count]) => (
                <div key={toolName} className="logs__item">
                  <div className="logs__item-content">
                    <div className="logs__details" style={{ flex: 1 }}>
                      <label className="login-page__label">{toolName}</label>
                      <p style={{ fontSize: '0.875rem', color: '#666' }}>Used {count} times</p>
                    </div>
                    <div style={{ fontSize: '1.5rem', fontWeight: '600', color: '#3b82f6' }}>{count}</div>
                  </div>
                </div>
              ))
            ) : (
              <div className="logs__item">
                <p style={{ textAlign: 'center', color: '#666' }}>No tool usage data yet</p>
              </div>
            )}
          </div>
        </div>

        <div className="logs__container">
          <h2 style={{ marginBottom: '1rem' }}>Recent Calls</h2>
          <div className="logs__list">
            {displayStats.recent_calls.length > 0 ? (
              displayStats.recent_calls.map((call, index) => (
                <div key={index} className="logs__item">
                  <div className="logs__item-content">
                    <div className="dashboard__icon" style={{ color: call.success ? '#22c55e' : '#ef4444' }}>
                      {call.success ? <CheckCircle /> : <XCircle />}
                    </div>
                    <div className="logs__details" style={{ flex: 1 }}>
                      <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', marginBottom: '0.25rem' }}>
                        <span style={{ fontWeight: '600', fontSize: '0.875rem' }}>{call.model}</span>
                        <span style={{ fontSize: '0.75rem', color: '#666' }}>
                          {new Date(call.timestamp).toLocaleString()}
                        </span>
                      </div>
                      <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', fontSize: '0.75rem', color: '#666' }}>
                        <Clock style={{ width: '0.75rem', height: '0.75rem' }} />
                        <span>{call.duration_ms}ms</span>
                      </div>
                      {call.tools_used.length > 0 && (
                        <div style={{ marginTop: '0.25rem', display: 'flex', gap: '0.25rem', flexWrap: 'wrap' }}>
                          {call.tools_used.map((tool, i) => (
                            <span
                              key={i}
                              style={{
                                fontSize: '0.75rem',
                                padding: '0.125rem 0.375rem',
                                background: '#f3f4f6',
                                borderRadius: '0.25rem',
                              }}
                            >
                              {tool}
                            </span>
                          ))}
                        </div>
                      )}
                    </div>
                  </div>
                </div>
              ))
            ) : (
              <div className="logs__item">
                <p style={{ textAlign: 'center', color: '#666' }}>No recent calls yet</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
