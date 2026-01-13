import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { useAuth } from '@/hooks/useAuth';
import { useWebSocket } from '@/hooks/useWebSocket';
import { Toaster } from '@/components/ui/sonner';
import { Login } from '@/pages/Login';
import { Main } from '@/components/layout/Main';
import { Dashboard } from '@/pages/Dashboard';
import { Plugins } from '@/pages/Plugins';
import { Tasks } from '@/pages/Tasks';
import { Config } from '@/pages/Config';
import { Logs } from '@/pages/Logs';
import { AI } from '@/pages/AI';
import { Presets } from '@/pages/Presets';

function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuth();
  
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }
  
  return <>{children}</>;
}

function App() {
  useWebSocket();

  return (
    <BrowserRouter>
      <Toaster />
      <Routes>
        <Route path="/login" element={<Login />} />
        
        <Route
          path="/"
          element={
            <ProtectedRoute>
              <Main />
            </ProtectedRoute>
          }
        >
          <Route index element={<Dashboard />} />
          <Route path="plugins" element={<Plugins />} />
          <Route path="tasks" element={<Tasks />} />
          <Route path="config" element={<Config />} />
          <Route path="logs" element={<Logs />} />
          <Route path="ai" element={<AI />} />
          <Route path="presets" element={<Presets />} />
        </Route>
        
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
