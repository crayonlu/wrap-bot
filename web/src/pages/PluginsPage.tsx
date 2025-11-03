import { usePlugins, useTogglePlugin } from '../lib/hooks/useQuery'
import toast from 'react-hot-toast'
import { Puzzle } from 'lucide-react'

export default function PluginsPage() {
  const { data: plugins, isLoading } = usePlugins()
  const togglePlugin = useTogglePlugin()

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
        <h1>Plugins</h1>
        <p>Manage bot plugins and features</p>
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
