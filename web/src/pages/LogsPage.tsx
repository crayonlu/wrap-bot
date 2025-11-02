import { useState } from 'react'
import { useLogs } from '../lib/hooks/useQuery'
import { AlertCircle, Info, AlertTriangle } from 'lucide-react'

const levelIcons = {
  error: AlertCircle,
  warn: AlertTriangle,
  info: Info,
  debug: Info,
}

const levelColors = {
  error: 'text-red-600 bg-red-50',
  warn: 'text-yellow-600 bg-yellow-50',
  info: 'text-blue-600 bg-blue-50',
  debug: 'text-gray-600 bg-gray-50',
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
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-[#8B7355]"></div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-[#8B7355]">Logs</h1>
        <p className="text-gray-600 mt-1">System logs and events</p>
      </div>

      <div className="flex items-center gap-4">
        <select
          value={selectedLevel}
          onChange={(e) => setSelectedLevel(e.target.value)}
          className="px-4 py-2 border border-[#D4C4B0] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#8B7355]"
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
          className="px-4 py-2 border border-[#D4C4B0] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#8B7355]"
        >
          <option value={50}>50 entries</option>
          <option value={100}>100 entries</option>
          <option value={200}>200 entries</option>
          <option value={500}>500 entries</option>
        </select>
      </div>

      <div className="bg-white rounded-xl shadow-sm border border-[#EBE6DF] overflow-hidden">
        <div className="divide-y divide-[#EBE6DF] max-h-[600px] overflow-y-auto">
          {logs?.map((log, index) => {
            const Icon = levelIcons[log.level]
            return (
              <div key={index} className="p-4 hover:bg-[#FAF8F5] transition-colors">
                <div className="flex items-start gap-3">
                  <div className={`p-2 rounded-lg ${levelColors[log.level]}`}>
                    <Icon className="w-4 h-4" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-3 mb-1">
                      <span className="text-xs text-gray-500">
                        {new Date(log.timestamp).toLocaleString()}
                      </span>
                      <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${levelColors[log.level]}`}>
                        {log.level.toUpperCase()}
                      </span>
                    </div>
                    <p className="text-sm text-gray-800 font-mono">{log.message}</p>
                    {log.context && (
                      <pre className="mt-2 text-xs text-gray-600 bg-gray-50 p-2 rounded overflow-x-auto">
                        {JSON.stringify(log.context, null, 2)}
                      </pre>
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
