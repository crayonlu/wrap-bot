import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from './stores/auth'
import LoginPage from './pages/LoginPage'
import DashboardPage from './pages/DashboardPage'
import PluginsPage from './pages/PluginsPage'
import TasksPage from './pages/TasksPage'
import ConfigPage from './pages/ConfigPage'
import LogsPage from './pages/LogsPage'
import MainLayout from './components/layout/MainLayout'
import NotificationContainer from './components/ui/NotificationContainer'

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated)
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />
}

function App() {
  return (
    <BrowserRouter>
      <NotificationContainer />
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
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App
