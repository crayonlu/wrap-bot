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
  password: string
}

export interface LoginResponse {
  token: string
}

export interface BotStatus {
  online: boolean
  account: {
    user_id: number
    nickname: string
  }
  stats: {
    message_sent: number
    message_received: number
  }
}

export interface Plugin {
  name: string
  description: string
  enabled: boolean
  commands?: string[]
}

export interface Task {
  id: string
  name: string
  description: string
  schedule: string
  last_run?: string
  next_run?: string
  status: 'running' | 'idle' | 'failed'
}

export interface ConfigItem {
  key: string
  value: string
  description?: string
}

export interface LogEntry {
  timestamp: string
  level: 'info' | 'warn' | 'error' | 'debug'
  message: string
  context?: Record<string, unknown>
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

export default api
