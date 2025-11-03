import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '../stores/auth'
import { authAPI } from '../lib/api'
import { Lock, LogIn } from 'lucide-react'
import toast from 'react-hot-toast'

export default function LoginPage() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()
  const login = useAuthStore((state) => state.login)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)

    try {
      const response = await authAPI.login({ username, password })
      login(response.data.token)
      toast.success('Login successful!')
      navigate('/dashboard')
    } catch (error) {
      toast.error('Invalid credentials')
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
              <label className="login-page__label">Username</label>
              <input
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                required
                className="login-page__input"
                placeholder="Enter username"
              />
            </div>

            <div>
              <label className="login-page__label">Password</label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                className="login-page__input"
                placeholder="Enter password"
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
