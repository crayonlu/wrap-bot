import { usePlugins, useTogglePlugin } from '../lib/hooks/useQuery'
import { useNotificationStore } from '../stores/notification'
import { Puzzle, Power } from 'lucide-react'

export default function PluginsPage() {
  const { data: plugins, isLoading } = usePlugins()
  const togglePlugin = useTogglePlugin()
  const addNotification = useNotificationStore((state) => state.addNotification)

  const handleToggle = async (name: string, currentState: boolean) => {
    try {
      await togglePlugin.mutateAsync(name)
      addNotification({
        type: 'success',
        message: `Plugin ${name} ${currentState ? 'disabled' : 'enabled'} successfully`,
      })
    } catch (error) {
      addNotification({
        type: 'error',
        message: `Failed to toggle plugin ${name}`,
      })
    }
  }

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
        <h1 className="text-3xl font-bold text-[#8B7355]">Plugins</h1>
        <p className="text-gray-600 mt-1">Manage bot plugins and features</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {plugins?.map((plugin) => (
          <div
            key={plugin.name}
            className="bg-white rounded-xl shadow-sm border border-[#EBE6DF] p-6 hover:shadow-md transition-shadow"
          >
            <div className="flex items-start justify-between">
              <div className="flex items-start gap-4 flex-1">
                <div className="p-3 bg-[#FAF8F5] rounded-lg">
                  <Puzzle className="w-6 h-6 text-[#8B7355]" />
                </div>
                <div className="flex-1">
                  <h3 className="font-semibold text-[#8B7355] text-lg">{plugin.name}</h3>
                  <p className="text-gray-600 text-sm mt-1">{plugin.description}</p>
                  {plugin.commands && plugin.commands.length > 0 && (
                    <div className="mt-3 flex flex-wrap gap-2">
                      {plugin.commands.map((cmd) => (
                        <span
                          key={cmd}
                          className="px-2 py-1 bg-[#FAF8F5] text-[#8B7355] text-xs rounded-md font-mono"
                        >
                          {cmd}
                        </span>
                      ))}
                    </div>
                  )}
                </div>
              </div>
              <button
                onClick={() => handleToggle(plugin.name, plugin.enabled)}
                disabled={togglePlugin.isPending}
                className={`ml-4 p-2 rounded-lg transition-colors ${
                  plugin.enabled
                    ? 'bg-[#8B7355] text-white hover:bg-[#6d5940]'
                    : 'bg-gray-200 text-gray-600 hover:bg-gray-300'
                }`}
              >
                <Power className="w-5 h-5" />
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
