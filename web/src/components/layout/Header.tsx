import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '../../stores/auth'
import { useUIStore } from '../../stores/ui'
import { useNotificationStore } from '../../stores/notification'
import { Menu, LogOut } from 'lucide-react'

export default function Header() {
  const navigate = useNavigate()
  const logout = useAuthStore((state) => state.logout)
  const toggleSidebar = useUIStore((state) => state.toggleSidebar)
  const addNotification = useNotificationStore((state) => state.addNotification)

  const handleLogout = () => {
    logout()
    addNotification({
      type: 'info',
      message: 'Logged out successfully',
    })
    navigate('/login')
  }

  return (
    <header className="header">
      <div className="header__content">
        <button onClick={toggleSidebar} className="header__menu-btn">
          <Menu />
        </button>

        <button onClick={handleLogout} className="header__logout-btn">
          <LogOut />
          Logout
        </button>
      </div>
    </header>
  )
}
