import { useState } from 'react'
import { useConfig, useUpdateConfig } from '../lib/hooks/useQuery'
import toast from 'react-hot-toast'
import { Save, Settings } from 'lucide-react'
import type { ConfigItem } from '../lib/api'

export default function ConfigPage() {
  const { data: configData, isLoading } = useConfig()
  const updateConfig = useUpdateConfig()
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
      toast.success('Configuration updated successfully')
      setEditedConfig([])
    } catch (error) {
      toast.error('Failed to update configuration')
    }
  }

  const displayConfig = editedConfig.length ? editedConfig : configData || []

  if (isLoading) {
    return (
      <div className="loading">
        <div className="loading__spinner"></div>
      </div>
    )
  }

  return (
    <div className="config">
      <div className="config__header" style={{display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start'}}>
        <div>
          <h1>Configuration</h1>
          <p>Manage bot environment variables</p>
        </div>
        {editedConfig.length > 0 && (
          <button
            onClick={handleSave}
            disabled={updateConfig.isPending}
            className="login-page__button"
            style={{width: 'auto', display: 'flex', alignItems: 'center', gap: '0.5rem'}}
          >
            <Save style={{width: '1rem', height: '1rem'}} />
            Save Changes
          </button>
        )}
      </div>

      <div className="logs__container">
        <div className="logs__list">
          {displayConfig.map((item) => (
            <div key={item.key} className="logs__item">
              <div className="logs__item-content">
                <div className="dashboard__icon">
                  <Settings />
                </div>
                <div className="logs__details" style={{flex: 1}}>
                  <label className="login-page__label">
                    {item.key}
                  </label>
                  {item.description && (
                    <p style={{fontSize: '0.75rem', color: '#666', marginBottom: '0.5rem'}}>{item.description}</p>
                  )}
                  <input
                    type="text"
                    value={item.value}
                    onChange={(e) => handleChange(item.key, e.target.value)}
                    className="login-page__input"
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
