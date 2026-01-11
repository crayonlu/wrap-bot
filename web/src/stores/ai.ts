import { create } from 'zustand'
import type { AITool, AIStats, AIChatMessage } from '../lib/api'

interface AIState {
  tools: AITool[]
  stats: AIStats | null
  chatHistory: AIChatMessage[]
  setTools: (tools: AITool[]) => void
  setStats: (stats: AIStats) => void
  addChatMessage: (message: AIChatMessage) => void
  clearChatHistory: () => void
}

export const useAIStore = create<AIState>((set) => ({
  tools: [],
  stats: null,
  chatHistory: [],
  setTools: (tools) => set({ tools }),
  setStats: (stats) => set({ stats }),
  addChatMessage: (message) => set((state) => ({
    chatHistory: [...state.chatHistory, message],
  })),
  clearChatHistory: () => set({ chatHistory: [] }),
}))
