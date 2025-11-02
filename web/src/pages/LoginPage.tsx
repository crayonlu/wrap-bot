import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '../stores/auth'
import { useNotificationStore } from '../stores/notification'
import { authAPI } from '../lib/api'
import { Lock, LogIn } from 'lucide-react'

export default function LoginPage() {
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()
  const login = useAuthStore((state) => state.login)
  const addNotification = useNotificationStore((state) => state.addNotification)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)

    try {
      const response = await authAPI.login({ password })
      login(response.data.token)
      addNotification({
        type: 'success',
        message: 'Login successful!',
      })
      navigate('/dashboard')
    } catch (error) {
      addNotification({
        type: 'error',
        message: 'Invalid password',
      })
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="login-page">
      <div className="login-page__container">
        <div className="login-page__card">
          <div className="login-page__header">
            <div className="login-page__icon">
              <Lock />
            </div>
            <h1 className="login-page__title">Bot Admin</h1>
            <p className="login-page__subtitle">Enter your password to continue</p>
          </div>

          <form onSubmit={handleSubmit} className="login-page__form">
            <div>
              <label className="login-page__label">Password</label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                className="login-page__input"
                placeholder="Enter admin password"
              />
            </div>

            <button type="submit" disabled={loading} className="login-page__button">
              {loading ? (
                <div className="login-page__spinner"></div>
              ) : (
                <>
                  <LogIn />
                  Login
                </>
              )}
            </button>
          </form>
        </div>
      </div>
    </div>
  )
}
