import { useState } from 'react'
import { useConfig, useUpdateConfig } from '../lib/hooks/useQuery'
import { useNotificationStore } from '../stores/notification'
import { Save, Settings } from 'lucide-react'
import type { ConfigItem } from '../lib/api'

export default function ConfigPage() {
  const { data: configData, isLoading } = useConfig()
  const updateConfig = useUpdateConfig()
  const addNotification = useNotificationStore((state) => state.addNotification)
  const [editedConfig, setEditedConfig] = useState<ConfigItem[]>([])

  const handleChange = (key: string, value: string) => {
    const updated = [...(editedConfig.length ? editedConfig : configData || [])]
    const index = updated.findIndex((item) => item.key === key)
    if (index !== -1) {
      updated[index] = { ...updated[index], value }
    }
    setEditedConfig(updated)
  }

  const handleSave = async () => {
    try {
      await updateConfig.mutateAsync(editedConfig)
      addNotification({
        type: 'success',
        message: 'Configuration updated successfully',
      })
      setEditedConfig([])
    } catch (error) {
      addNotification({
        type: 'error',
        message: 'Failed to update configuration',
      })
    }
  }

  const displayConfig = editedConfig.length ? editedConfig : configData || []

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-[#8B7355]"></div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-[#8B7355]">Configuration</h1>
          <p className="text-gray-600 mt-1">Manage bot environment variables</p>
        </div>
        {editedConfig.length > 0 && (
          <button
            onClick={handleSave}
            disabled={updateConfig.isPending}
            className="px-4 py-2 bg-[#8B7355] text-white rounded-lg hover:bg-[#6d5940] transition-colors flex items-center gap-2 disabled:opacity-50"
          >
            <Save className="w-4 h-4" />
            Save Changes
          </button>
        )}
      </div>

      <div className="bg-white rounded-xl shadow-sm border border-[#EBE6DF] overflow-hidden">
        <div className="divide-y divide-[#EBE6DF]">
          {displayConfig.map((item) => (
            <div key={item.key} className="p-6 hover:bg-[#FAF8F5] transition-colors">
              <div className="flex items-start gap-4">
                <div className="p-2 bg-[#FAF8F5] rounded-lg">
                  <Settings className="w-5 h-5 text-[#8B7355]" />
                </div>
                <div className="flex-1">
                  <label className="block text-sm font-medium text-[#8B7355] mb-1">
                    {item.key}
                  </label>
                  {item.description && (
                    <p className="text-xs text-gray-500 mb-2">{item.description}</p>
                  )}
                  <input
                    type="text"
                    value={item.value}
                    onChange={(e) => handleChange(item.key, e.target.value)}
                    className="w-full px-3 py-2 border border-[#D4C4B0] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#8B7355] focus:border-transparent"
                  />
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
