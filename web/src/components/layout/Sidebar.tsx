import { NavLink } from 'react-router-dom'
import { useUIStore } from '../../stores/ui'
import { LayoutDashboard, Puzzle, Clock, Settings, FileText, X } from 'lucide-react'

const navigation = [
  { name: 'Dashboard', href: '/dashboard', icon: LayoutDashboard },
  { name: 'Plugins', href: '/plugins', icon: Puzzle },
  { name: 'Tasks', href: '/tasks', icon: Clock },
  { name: 'Config', href: '/config', icon: Settings },
  { name: 'Logs', href: '/logs', icon: FileText },
]

export default function Sidebar() {
  const { sidebarOpen, setSidebarOpen } = useUIStore()

  return (
    <>
      {sidebarOpen && (
        <div className="sidebar__overlay" onClick={() => setSidebarOpen(false)} />
      )}
      <aside className={`sidebar__container ${sidebarOpen ? 'sidebar__container--open' : 'sidebar__container--closed'}`}>
        <div className="sidebar__header">
          <h2 className="sidebar__title">Bot Admin</h2>
          <button onClick={() => setSidebarOpen(false)} className="sidebar__close">
            <X />
          </button>
        </div>
        <nav className="sidebar__nav">
          {navigation.map((item) => (
            <NavLink
              key={item.name}
              to={item.href}
              className={({ isActive }) =>
                `sidebar__link ${isActive ? 'sidebar__link--active' : ''}`
              }
            >
              <item.icon />
              <span>{item.name}</span>
            </NavLink>
          ))}
        </nav>
      </aside>
    </>
  )
}
