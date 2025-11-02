import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface User {
  token: string
}

interface AuthState {
  user: User | null
  isAuthenticated: boolean
  login: (token: string) => void
  logout: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      isAuthenticated: false,
      login: (token: string) => {
        localStorage.setItem('token', token)
        set({ user: { token }, isAuthenticated: true })
      },
      logout: () => {
        localStorage.removeItem('token')
        set({ user: null, isAuthenticated: false })
      },
    }),
    {
      name: 'auth-storage',
    }
  )
)
