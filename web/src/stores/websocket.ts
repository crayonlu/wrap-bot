import { create } from 'zustand'
import type { BotStatus, Plugin, Task, LogEntry } from '../lib/api'

interface WebSocketState {
  connected: boolean
  status: BotStatus | null
  logs: LogEntry[]
  plugins: Plugin[]
  tasks: Task[]

  setConnected: (connected: boolean) => void
  setStatus: (status: BotStatus) => void
  addLog: (log: LogEntry) => void
  setLogs: (logs: LogEntry[]) => void
  setPlugins: (plugins: Plugin[]) => void
  setTasks: (tasks: Task[]) => void
  updatePlugin: (name: string, enabled: boolean) => void
  updateTask: (task: Task) => void
}

export const useWebSocketStore = create<WebSocketState>((set) => ({
  connected: false,
  status: null,
  logs: [],
  plugins: [],
  tasks: [],

  setConnected: (connected) => set({ connected }),

  setStatus: (status) => set({ status }),

  addLog: (log) =>
    set((state) => ({
      logs: [...state.logs, log].slice(-500), // Keep last 500 logs
    })),

  setLogs: (logs) => set({ logs }),

  setPlugins: (plugins) => set({ plugins }),

  setTasks: (tasks) => set({ tasks }),

  updatePlugin: (name, enabled) =>
    set((state) => ({
      plugins: state.plugins.map((p) =>
        p.name === name ? { ...p, enabled } : p
      ),
    })),

  updateTask: (task) =>
    set((state) => ({
      tasks: state.tasks.map((t) => (t.id === task.id ? task : t)),
    })),
}))
