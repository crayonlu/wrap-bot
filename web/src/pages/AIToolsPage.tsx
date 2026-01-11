import { useEffect, useState } from 'react'
import { useAITools, useUpdateConfig } from '../lib/hooks/useQuery'
import { useAIStore } from '../stores/ai'
import toast from 'react-hot-toast'
import { Settings, ToggleLeft, ToggleRight, Save } from 'lucide-react'
import type { AITool } from '../lib/api'

export default function AIToolsPage() {
  const { data: tools, isLoading } = useAITools()
  const updateConfig = useUpdateConfig()
  const { setTools } = useAIStore()
  const [editedTools, setEditedTools] = useState<AITool[]>([])

  useEffect(() => {
    if (tools) {
      setTools(tools)
    }
  }, [tools, setTools])

  const handleToggle = (toolName: string) => {
    const updated = editedTools.length ? [...editedTools] : [...(tools || [])]
    const index = updated.findIndex((t) => t.name === toolName)
    if (index !== -1) {
      updated[index] = { ...updated[index], enabled: !updated[index].enabled }
    }
    setEditedTools(updated)
  }

  const handleSave = async () => {
    try {
      const textTools = editedTools
        .filter((t) => t.enabled && (t.category === 'text' || t.category === 'both'))
        .map((t) => t.name)
      const visionTools = editedTools
        .filter((t) => t.enabled && (t.category === 'vision' || t.category === 'both'))
        .map((t) => t.name)

      const configItems = [
        { key: 'AI_TEXT_MODEL_TOOLS', value: textTools.join(',') },
        { key: 'AI_VISION_MODEL_TOOLS', value: visionTools.join(',') },
      ]

      await updateConfig.mutateAsync(configItems)
      toast.success('AI tools configuration updated')
      setEditedTools([])
      setTools(editedTools)
    } catch (error) {
      toast.error('Failed to update AI tools configuration')
    }
  }

  const displayTools = editedTools.length ? editedTools : tools || []

  const getCategoryLabel = (category: string) => {
    switch (category) {
      case 'text':
        return 'Text'
      case 'vision':
        return 'Vision'
      case 'both':
        return 'Both'
      default:
        return category
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
    <div className="config">
      <div className="config__header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
        <div>
          <h1>AI Tools Configuration</h1>
          <p>Manage AI tools for text and vision models</p>
        </div>
        {editedTools.length > 0 && (
          <button
            onClick={handleSave}
            disabled={updateConfig.isPending}
            className="login-page__button"
            style={{ width: 'auto', display: 'flex', alignItems: 'center', gap: '0.5rem' }}
          >
            <Save style={{ width: '1rem', height: '1rem' }} />
            Save Changes
          </button>
        )}
      </div>

      <div className="logs__container">
        <div className="logs__list">
          {displayTools.map((tool) => (
            <div key={tool.name} className="logs__item">
              <div className="logs__item-content">
                <div className="dashboard__icon">
                  <Settings />
                </div>
                <div className="logs__details" style={{ flex: 1 }}>
                  <label className="login-page__label">{tool.name}</label>
                  <p style={{ fontSize: '0.75rem', color: '#666', marginBottom: '0.5rem' }}>{tool.description}</p>
                  <div style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
                    <span
                      style={{
                        fontSize: '0.75rem',
                        padding: '0.25rem 0.5rem',
                        background: '#f3f4f6',
                        borderRadius: '0.25rem',
                      }}
                    >
                      {getCategoryLabel(tool.category)}
                    </span>
                    <button
                      onClick={() => handleToggle(tool.name)}
                      style={{
                        display: 'flex',
                        alignItems: 'center',
                        gap: '0.5rem',
                        padding: '0.5rem 1rem',
                        border: '1px solid #ddd',
                        borderRadius: '0.375rem',
                        background: '#fff',
                        cursor: 'pointer',
                        fontSize: '0.875rem',
                      }}
                    >
                      {tool.enabled ? (
                        <>
                          <ToggleRight style={{ width: '1.25rem', height: '1.25rem', color: '#22c55e' }} />
                          <span style={{ color: '#22c55e', fontWeight: '500' }}>Enabled</span>
                        </>
                      ) : (
                        <>
                          <ToggleLeft style={{ width: '1.25rem', height: '1.25rem', color: '#ef4444' }} />
                          <span style={{ color: '#ef4444', fontWeight: '500' }}>Disabled</span>
                        </>
                      )}
                    </button>
                  </div>
                </div>
              </div>
            </div>
          ))}
          {(!displayTools || displayTools.length === 0) && (
            <div className="logs__item">
              <p style={{ textAlign: 'center', color: '#666' }}>No AI tools available</p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
