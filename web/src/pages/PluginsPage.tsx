import { useEffect } from 'react'
import { usePlugins, useTogglePlugin } from '../lib/hooks/useQuery'
import { useWebSocket } from '../lib/hooks/useWebSocket'
import { useWebSocketStore } from '../stores/websocket'
import toast from 'react-hot-toast'
import { Puzzle, Wifi, WifiOff } from 'lucide-react'

export default function PluginsPage() {
  const { data: initialPlugins, isLoading } = usePlugins()
  const togglePlugin = useTogglePlugin()
  const { connected } = useWebSocketStore()
  const plugins = useWebSocketStore((state) => state.plugins)
  const setPlugins = useWebSocketStore((state) => state.setPlugins)

  // Initialize plugins from API response
  useEffect(() => {
    if (initialPlugins && initialPlugins.length > 0) {
      setPlugins(initialPlugins)
    }
  }, [initialPlugins, setPlugins])

  // WebSocket integration
  useWebSocket({
    enabled: true,
    onMessage: (message) => {
      if (message.type === 'plugins') {
        const updatedPlugins = message.data as any[]
        setPlugins(updatedPlugins)
      }
    },
  })

  const handleToggle = async (name: string, currentState: boolean) => {
    try {
      await togglePlugin.mutateAsync(name)
      toast.success(`Plugin ${name} ${currentState ? 'disabled' : 'enabled'}`)
    } catch (error) {
      toast.error(`Failed to toggle plugin ${name}`)
    }
  }

  if (isLoading) {
    return (
      <div className="loading">
        <div className="loading__spinner"></div>
      </div>
    )
  }

  return (
    <div className="plugins">
      <div className="plugins__header">
        <div>
          <h1>Plugins</h1>
          <p>Manage bot plugins and features</p>
        </div>
        <div className="plugins__connection-status">
          {connected ? (
            <span className="plugins__status plugins__status--connected">
              <Wifi size={16} />
              Connected
            </span>
          ) : (
            <span className="plugins__status plugins__status--disconnected">
              <WifiOff size={16} />
              Disconnected
            </span>
          )}
        </div>
      </div>

      <div className="plugins__grid">
        {plugins?.map((plugin) => (
          <div key={plugin.name} className="plugins__card">
            <div className="plugins__card-header">
              <div className="dashboard__icon">
                <Puzzle />
              </div>
              <button
                onClick={() => handleToggle(plugin.name, plugin.enabled)}
                disabled={togglePlugin.isPending}
                className={`plugins__toggle ${plugin.enabled ? 'plugins__toggle--active' : ''}`}
              />
            </div>
            <h3 className="plugins__card-title">{plugin.name}</h3>
            <p className="plugins__card-description">{plugin.description || 'No description'}</p>
          </div>
        ))}
      </div>
    </div>
  )
}
