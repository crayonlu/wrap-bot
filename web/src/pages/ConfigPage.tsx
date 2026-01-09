import { useState } from 'react'
import { useConfig, useUpdateConfig } from '../lib/hooks/useQuery'
import toast from 'react-hot-toast'
import { Save, Settings, ToggleLeft, ToggleRight } from 'lucide-react'
import type { ConfigItem } from '../lib/api'

export default function ConfigPage() {
  const { data: configData, isLoading } = useConfig()
  const updateConfig = useUpdateConfig()
  const [editedConfig, setEditedConfig] = useState<ConfigItem[]>([])
  const [originalConfig, setOriginalConfig] = useState<ConfigItem[]>([])

  const handleChange = (key: string, value: string) => {
    const updated = [...(editedConfig.length ? editedConfig : configData || [])]
    const index = updated.findIndex((item) => item.key === key)
    if (index !== -1) {
      updated[index] = { ...updated[index], value }
    }
    setEditedConfig(updated)
  }

  const handleToggle = (key: string) => {
    const updated = [...(editedConfig.length ? editedConfig : configData || [])]
    const index = updated.findIndex((item) => item.key === key)
    if (index !== -1) {
      const currentValue = updated[index].value.toLowerCase()
      const newValue = currentValue === 'true' ? 'false' : 'true'
      updated[index] = { ...updated[index], value: newValue }
    }
    setEditedConfig(updated)
  }

  const handleSave = async () => {
    try {
      const changedItems = editedConfig.filter((item) => {
        const originalItem = originalConfig.find(o => o.key === item.key)
        return !originalItem || originalItem.value !== item.value
      })
      
      if (changedItems.length === 0) {
        toast.success('No changes to save')
        setEditedConfig([])
        return
      }

      await updateConfig.mutateAsync(changedItems)
      toast.success(`Configuration updated successfully (${changedItems.length} items)`)
      setEditedConfig([])
      setOriginalConfig(configData || [])
    } catch (error) {
      toast.error('Failed to update configuration')
    }
  }

  const displayConfig = editedConfig.length ? editedConfig : configData || []

  const isBooleanConfig = (key: string) => {
    return key === 'AI_ENABLED' || key === 'AI_TOOLS_ENABLED' || key === 'AI_VISION_ENABLED' || key === 'SERVER_ENABLED' || key === 'DEBUG'
  }

  const isSelectConfig = (key: string) => {
    return key === 'AI_IMAGE_DETAIL'
  }

  if (isLoading) {
    return (
      <div className="loading">
        <div className="loading__spinner"></div>
      </div>
    )
  }

  if (configData && configData.length > 0 && originalConfig.length === 0) {
    setOriginalConfig(configData)
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
          {displayConfig?.map((item) => (
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
                  {isBooleanConfig(item.key) ? (
                    <button
                      onClick={() => handleToggle(item.key)}
                      style={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: '0.5rem',
                        padding: '0.5rem 1rem',
                        border: '1px solid #ddd',
                        borderRadius: '0.375rem',
                        background: '#fff',
                        cursor: 'pointer',
                        fontSize: '0.875rem'
                      }}
                    >
                      {item.value.toLowerCase() === 'true' ? (
                        <>
                          <ToggleRight style={{width: '1.25rem', height: '1.25rem', color: '#22c55e'}} />
                          <span style={{color: '#22c55e', fontWeight: '500'}}>Enabled</span>
                        </>
                      ) : (
                        <>
                          <ToggleLeft style={{width: '1.25rem', height: '1.25rem', color: '#ef4444'}} />
                          <span style={{color: '#ef4444', fontWeight: '500'}}>Disabled</span>
                        </>
                      )}
                    </button>
                  ) : isSelectConfig(item.key) ? (
                    <select
                      value={item.value}
                      onChange={(e) => handleChange(item.key, e.target.value)}
                      className="login-page__input"
                      style={{cursor: 'pointer'}}
                    >
                      <option value="auto">Auto</option>
                      <option value="high">High</option>
                      <option value="low">Low</option>
                    </select>
                  ) : (
                    <input
                      type="text"
                      value={item.value}
                      onChange={(e) => handleChange(item.key, e.target.value)}
                      className="login-page__input"
                    />
                  )}
                </div>
              </div>
            </div>
          ))}
          {(!displayConfig || displayConfig.length === 0) && (
            <div className="logs__item">
              <p style={{textAlign: 'center', color: '#666'}}>No configuration items available</p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
