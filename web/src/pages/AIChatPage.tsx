import { useState } from 'react'
import { useAIChat, useAIChatWithImage } from '../lib/hooks/useQuery'
import { useAIStore } from '../stores/ai'
import toast from 'react-hot-toast'
import { Send, Image as ImageIcon, Trash2, Bot, User } from 'lucide-react'
import type { AIChatMessage } from '../lib/api'

export default function AIChatPage() {
  const [message, setMessage] = useState('')
  const [images, setImages] = useState<string[]>([])
  const [model, setModel] = useState<'text' | 'vision'>('text')
  const [conversationId, setConversationId] = useState('')
  const chatMutation = useAIChat()
  const chatWithImageMutation = useAIChatWithImage()
  const { chatHistory, addChatMessage, clearChatHistory } = useAIStore()

  const handleSendMessage = async () => {
    if (!message.trim() && images.length === 0) return

    const userMessage: AIChatMessage = {
      role: 'user',
      content: message,
      images: images.length > 0 ? images : undefined,
      timestamp: new Date().toISOString(),
    }
    addChatMessage(userMessage)

    const currentConversationId = conversationId || `chat-${Date.now()}`
    setConversationId(currentConversationId)

    try {
      let response
      if (images.length > 0) {
        response = await chatWithImageMutation.mutateAsync({
          message,
          images,
          model,
          conversation_id: currentConversationId,
        })
      } else {
        response = await chatMutation.mutateAsync({
          message,
          model,
          conversation_id: currentConversationId,
        })
      }

      const assistantMessage: AIChatMessage = {
        role: 'assistant',
        content: response.data.response,
        tool_calls: response.data.tool_calls,
        timestamp: new Date().toISOString(),
      }
      addChatMessage(assistantMessage)

      setMessage('')
      setImages([])
    } catch (error) {
      toast.error('Failed to send message')
    }
  }

  const handleClearHistory = () => {
    clearChatHistory()
    setConversationId('')
    toast.success('Chat history cleared')
  }

  const handleImageUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files
    if (!files) return

    const newImages: string[] = []
    Array.from(files).forEach((file) => {
      const reader = new FileReader()
      reader.onloadend = () => {
        newImages.push(reader.result as string)
        if (newImages.length === files.length) {
          setImages((prev) => [...prev, ...newImages])
        }
      }
      reader.readAsDataURL(file)
    })
  }

  const handleRemoveImage = (index: number) => {
    setImages((prev) => prev.filter((_, i) => i !== index))
  }

  return (
    <div className="config">
      <div className="config__header" style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
        <div>
          <h1>AI Chat Test</h1>
          <p>Test AI chat functionality with tools</p>
        </div>
        <div style={{ display: 'flex', gap: '0.5rem' }}>
          <select
            value={model}
            onChange={(e) => setModel(e.target.value as 'text' | 'vision')}
            className="login-page__input"
            style={{ width: 'auto', cursor: 'pointer' }}
          >
            <option value="text">Text Model</option>
            <option value="vision">Vision Model</option>
          </select>
          <button
            onClick={handleClearHistory}
            className="login-page__button"
            style={{ width: 'auto', display: 'flex', alignItems: 'center', gap: '0.5rem', background: '#ef4444' }}
          >
            <Trash2 style={{ width: '1rem', height: '1rem' }} />
            Clear
          </button>
        </div>
      </div>

      <div className="logs__container">
        <div className="logs__list" style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
          {chatHistory.map((msg, index) => (
            <div
              key={index}
              className="logs__item"
              style={{
                background: msg.role === 'user' ? '#f3f4f6' : '#fff',
                borderLeft: msg.role === 'user' ? '4px solid #3b82f6' : '4px solid #22c55e',
              }}
            >
              <div className="logs__item-content">
                <div className="dashboard__icon" style={{ color: msg.role === 'user' ? '#3b82f6' : '#22c55e' }}>
                  {msg.role === 'user' ? <User /> : <Bot />}
                </div>
                <div className="logs__details" style={{ flex: 1 }}>
                  <div style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', marginBottom: '0.5rem' }}>
                    <span style={{ fontWeight: '600', fontSize: '0.875rem' }}>
                      {msg.role === 'user' ? 'You' : 'AI'}
                    </span>
                    <span style={{ fontSize: '0.75rem', color: '#666' }}>
                      {new Date(msg.timestamp).toLocaleString()}
                    </span>
                  </div>
                  {msg.images && msg.images.length > 0 && (
                    <div style={{ display: 'flex', gap: '0.5rem', marginBottom: '0.5rem' }}>
                      {msg.images.map((img, i) => (
                        <img
                          key={i}
                          src={img}
                          alt={`Uploaded ${i}`}
                          style={{ width: '100px', height: '100px', objectFit: 'cover', borderRadius: '0.375rem' }}
                        />
                      ))}
                    </div>
                  )}
                  <p style={{ whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}>{msg.content}</p>
                  {msg.tool_calls && msg.tool_calls.length > 0 && (
                    <div style={{ marginTop: '0.5rem', padding: '0.5rem', background: '#f9fafb', borderRadius: '0.375rem' }}>
                      <p style={{ fontSize: '0.75rem', fontWeight: '600', marginBottom: '0.25rem' }}>Tool Calls:</p>
                      {msg.tool_calls.map((tc, i) => (
                        <div key={i} style={{ fontSize: '0.75rem', color: '#666' }}>
                          <strong>{tc.name}</strong>: {tc.arguments}
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            </div>
          ))}
          {chatHistory.length === 0 && (
            <div className="logs__item">
              <p style={{ textAlign: 'center', color: '#666' }}>No messages yet. Start a conversation!</p>
            </div>
          )}
        </div>

        <div style={{ marginTop: '1rem', padding: '1rem', background: '#fff', borderRadius: '0.5rem', border: '1px solid #e5e7eb' }}>
          {images.length > 0 && (
            <div style={{ display: 'flex', gap: '0.5rem', marginBottom: '0.5rem', flexWrap: 'wrap' }}>
              {images.map((img, index) => (
                <div key={index} style={{ position: 'relative' }}>
                  <img
                    src={img}
                    alt={`Upload ${index}`}
                    style={{ width: '80px', height: '80px', objectFit: 'cover', borderRadius: '0.375rem' }}
                  />
                  <button
                    onClick={() => handleRemoveImage(index)}
                    style={{
                      position: 'absolute',
                      top: '-8px',
                      right: '-8px',
                      width: '20px',
                      height: '20px',
                      borderRadius: '50%',
                      background: '#ef4444',
                      color: '#fff',
                      border: 'none',
                      cursor: 'pointer',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      fontSize: '0.75rem',
                    }}
                  >
                    ×
                  </button>
                </div>
              ))}
            </div>
          )}
          <div style={{ display: 'flex', gap: '0.5rem' }}>
            <label
              className="login-page__button"
              style={{
                width: 'auto',
                display: 'flex',
                alignItems: 'center',
                gap: '0.5rem',
                cursor: 'pointer',
                padding: '0.5rem 1rem',
              }}
            >
              <ImageIcon style={{ width: '1rem', height: '1rem' }} />
              <input
                type="file"
                accept="image/*"
                multiple
                onChange={handleImageUpload}
                style={{ display: 'none' }}
              />
              Add Image
            </label>
            <input
              type="text"
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && handleSendMessage()}
              placeholder="Type a message..."
              className="login-page__input"
              style={{ flex: 1 }}
            />
            <button
              onClick={handleSendMessage}
              disabled={chatMutation.isPending || chatWithImageMutation.isPending}
              className="login-page__button"
              style={{ width: 'auto', display: 'flex', alignItems: 'center', gap: '0.5rem' }}
            >
              <Send style={{ width: '1rem', height: '1rem' }} />
              Send
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
