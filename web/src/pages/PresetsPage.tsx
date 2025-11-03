import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import api from '../lib/api'
import { FileText, Save, Eye, Edit3 } from 'lucide-react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'

interface Preset {
  name: string
  path: string
  content: string
}

export default function PresetsPage() {
  const [selectedPreset, setSelectedPreset] = useState<string | null>(null)
  const [isEditing, setIsEditing] = useState(false)
  const [editContent, setEditContent] = useState('')
  const queryClient = useQueryClient()

  const { data: presets, isLoading } = useQuery<Preset[]>({
    queryKey: ['presets'],
    queryFn: async () => {
      const response = await api.get('/api/presets')
      return response.data
    },
  })

  const { data: currentPreset } = useQuery<Preset>({
    queryKey: ['preset', selectedPreset],
    queryFn: async () => {
      const response = await api.get(`/api/presets/${selectedPreset}`)
      return response.data
    },
    enabled: !!selectedPreset,
  })

  const updateMutation = useMutation({
    mutationFn: async ({ filename, content }: { filename: string; content: string }) => {
      await api.put(`/api/presets/${filename}`, { content })
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['presets'] })
      queryClient.invalidateQueries({ queryKey: ['preset', selectedPreset] })
      setIsEditing(false)
    },
  })

  const handleEdit = () => {
    if (currentPreset) {
      setEditContent(currentPreset.content)
      setIsEditing(true)
    }
  }

  const handleSave = () => {
    if (selectedPreset) {
      updateMutation.mutate({ filename: selectedPreset, content: editContent })
    }
  }

  const handleCancel = () => {
    setIsEditing(false)
    setEditContent('')
  }

  if (isLoading) {
    return (
      <div className="loading">
        <div className="loading__spinner"></div>
      </div>
    )
  }

  return (
    <div className="presets">
      <div className="presets__header">
        <h1>Presets</h1>
        <p>Manage AI prompt presets</p>
      </div>

      <div className="presets__layout">
        <div className="presets__sidebar">
          <h3>Preset Files</h3>
          <div className="presets__list">
            {presets?.map((preset) => (
              <button
                key={preset.name}
                className={`presets__item ${selectedPreset === preset.name ? 'presets__item--active' : ''}`}
                onClick={() => {
                  setSelectedPreset(preset.name)
                  setIsEditing(false)
                }}
              >
                <FileText size={16} />
                <span>{preset.name}</span>
              </button>
            ))}
          </div>
        </div>

        <div className="presets__content">
          {!selectedPreset ? (
            <div className="presets__empty">
              <FileText size={48} />
              <p>Select a preset to view or edit</p>
            </div>
          ) : (
            <>
              <div className="presets__content-header">
                <h2>{selectedPreset}</h2>
                <div className="presets__actions">
                  {!isEditing ? (
                    <>
                      <button onClick={handleEdit} className="presets__button presets__button--primary">
                        <Edit3 size={16} />
                        Edit
                      </button>
                    </>
                  ) : (
                    <>
                      <button onClick={handleCancel} className="presets__button">
                        Cancel
                      </button>
                      <button
                        onClick={handleSave}
                        className="presets__button presets__button--primary"
                        disabled={updateMutation.isPending}
                      >
                        <Save size={16} />
                        {updateMutation.isPending ? 'Saving...' : 'Save'}
                      </button>
                    </>
                  )}
                </div>
              </div>

              {isEditing ? (
                <textarea
                  className="presets__editor"
                  value={editContent}
                  onChange={(e) => setEditContent(e.target.value)}
                  placeholder="Enter preset content..."
                />
              ) : (
                <div className="presets__preview">
                  <div className="presets__preview-header">
                    <Eye size={16} />
                    <span>Preview</span>
                  </div>
                  <div className="presets__markdown">
                    <ReactMarkdown remarkPlugins={[remarkGfm]}>
                      {currentPreset?.content || ''}
                    </ReactMarkdown>
                  </div>
                </div>
              )}
            </>
          )}
        </div>
      </div>
    </div>
  )
}
