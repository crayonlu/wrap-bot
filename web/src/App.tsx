import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { Toaster } from 'react-hot-toast'
import { useAuthStore } from './stores/auth'
import { useWebSocket } from './lib/hooks/useWebSocket'
import LoginPage from './pages/LoginPage'
import DashboardPage from './pages/DashboardPage'
import PluginsPage from './pages/PluginsPage'
import TasksPage from './pages/TasksPage'
import ConfigPage from './pages/ConfigPage'
import LogsPage from './pages/LogsPage'
import PresetsPage from './pages/PresetsPage'
import AIToolsPage from './pages/AIToolsPage'
import AIChatPage from './pages/AIChatPage'
import AIStatsPage from './pages/AIStatsPage'
import MainLayout from './components/layout/MainLayout'

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)
  useWebSocket()
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />
}

function App() {
  return (
    <BrowserRouter>
      <Toaster position="top-right" />
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route
          path="/"
          element={
            <ProtectedRoute>
              <MainLayout />
            </ProtectedRoute>
          }
        >
          <Route index element={<Navigate to="/dashboard" replace />} />
          <Route path="dashboard" element={<DashboardPage />} />
          <Route path="plugins" element={<PluginsPage />} />
          <Route path="tasks" element={<TasksPage />} />
          <Route path="config" element={<ConfigPage />} />
          <Route path="logs" element={<LogsPage />} />
          <Route path="presets" element={<PresetsPage />} />
          <Route path="ai/tools" element={<AIToolsPage />} />
          <Route path="ai/chat" element={<AIChatPage />} />
          <Route path="ai/stats" element={<AIStatsPage />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App
