import { create } from 'zustand';
import type { BotStatus, Plugin, Task, LogEntry } from '@/types/api';

interface BotState {
  status: BotStatus | null;
  plugins: Plugin[];
  tasks: Task[];
  logs: LogEntry[];
  isLoading: boolean;
  error: string | null;
  
  setStatus: (status: BotStatus) => void;
  setPlugins: (plugins: Plugin[]) => void;
  updatePlugin: (name: string, enabled: boolean) => void;
  setTasks: (tasks: Task[]) => void;
  addLog: (log: LogEntry) => void;
  setLogs: (logs: LogEntry[]) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
}

export const useBotStore = create<BotState>((set) => ({
  status: null,
  plugins: [],
  tasks: [],
  logs: [],
  isLoading: false,
  error: null,

  setStatus: (status) => set({ status }),
  
  setPlugins: (plugins) => set({ plugins }),
  
  updatePlugin: (name, enabled) => set((state) => ({
    plugins: state.plugins.map((plugin) =>
      plugin.name === name ? { ...plugin, enabled } : plugin
    ),
  })),
  
  setTasks: (tasks) => set({ tasks }),
  
  addLog: (log) => set((state) => ({
    logs: [...state.logs.slice(-99), log],
  })),
  
  setLogs: (logs) => set({ logs }),
  
  setLoading: (isLoading) => set({ isLoading }),
  
  setError: (error) => set({ error }),
}));
