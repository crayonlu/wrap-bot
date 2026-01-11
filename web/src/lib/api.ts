import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
}

export interface BotStatus {
  running: boolean
  uptime: number
  version: string
  go_version: string
}

export interface Plugin {
  name: string
  enabled: boolean
  description?: string
  commands?: string[]
}

export interface Task {
  id: string
  name: string
  schedule: string
  next_run: string
  last_run: string
  status: string
  can_trigger: boolean
  description?: string
}

export interface ConfigItem {
  key: string
  value: string
  description?: string
}

export interface LogEntry {
  timestamp: string
  level: string
  message: string
  context?: Record<string, unknown>
}

export interface AITool {
  name: string
  description: string
  category: 'text' | 'vision' | 'both'
  enabled: boolean
}

export interface AIStats {
  total_calls: number
  tool_usage: Record<string, number>
  success_rate: number
  recent_calls: AICallRecord[]
}

export interface AICallRecord {
  timestamp: string
  model: string
  tools_used: string[]
  success: boolean
  duration_ms: number
}

export interface AIChatRequest {
  message: string
  images?: string[]
  model: 'text' | 'vision'
  conversation_id?: string
}

export interface AIChatResponse {
  response: string
  tool_calls?: ToolCall[]
  conversation_id: string
}

export interface ToolCall {
  name: string
  arguments: string
}

export interface AIChatMessage {
  role: 'user' | 'assistant'
  content: string
  images?: string[]
  tool_calls?: ToolCall[]
  timestamp: string
}

export const authAPI = {
  login: (data: LoginRequest) => 
    api.post<LoginResponse>('/api/auth/login', data),
}

export const statusAPI = {
  get: () => api.get<BotStatus>('/api/status'),
}

export const pluginsAPI = {
  list: () => api.get<Plugin[]>('/api/plugins'),
  toggle: (name: string) => api.post(`/api/plugins/${name}/toggle`),
}

export const tasksAPI = {
  list: () => api.get<Task[]>('/api/tasks'),
  trigger: (id: string) => api.post(`/api/tasks/${id}/trigger`),
}

export const configAPI = {
  list: () => api.get<ConfigItem[]>('/api/config'),
  update: (data: ConfigItem[]) => api.post('/api/config', data),
}

export const logsAPI = {
  list: (params?: { level?: string; limit?: number }) =>
    api.get<LogEntry[]>('/api/logs', { params }),
}

export const aiAPI = {
  getTools: () => api.get<AITool[]>('/api/ai/tools'),
  getStats: () => api.get<AIStats>('/api/ai/stats'),
  chat: (data: AIChatRequest) => api.post<AIChatResponse>('/api/ai/chat', data),
  chatWithImage: (data: AIChatRequest) => api.post<AIChatResponse>('/api/ai/chat/image', data),
}

export default api
