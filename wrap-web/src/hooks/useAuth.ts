import { useAuthStore } from '@/stores/auth';
import { apiClient } from '@/lib/api';
import { useNavigate } from 'react-router-dom';

export function useAuth() {
  const { token, isAuthenticated, setToken, logout } = useAuthStore();
  const navigate = useNavigate();

  const login = async (username: string, password: string) => {
    try {
      const response = await apiClient.login({ username, password });
      setToken(response.token);
      localStorage.setItem('token', response.token);
      return { success: true };
    } catch (error: any) {
      return {
        success: false,
        error: error.response?.data?.error || '登录失败',
      };
    }
  };

  const handleLogout = () => {
    logout();
    localStorage.removeItem('token');
    navigate('/login');
  };

  return {
    token,
    isAuthenticated,
    login,
    logout: handleLogout,
  };
}
