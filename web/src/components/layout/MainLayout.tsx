import { Outlet } from 'react-router-dom'
import Sidebar from './Sidebar'
import Header from './Header'
import { useUIStore } from '../../stores/ui'

export default function MainLayout() {
  const sidebarOpen = useUIStore((state) => state.sidebarOpen)

  return (
    <div className="main-layout">
      <Sidebar />
      <div className={`main-layout__content ${sidebarOpen ? 'main-layout__content--sidebar-open' : 'main-layout__content--sidebar-closed'}`}>
        <Header />
        <main className="main-layout__main">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
