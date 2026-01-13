import { Link, useLocation } from 'react-router-dom';
import { 
  LayoutDashboard, 
  Puzzle, 
  Clock, 
  Settings, 
  FileText, 
  Brain, 
  FileCode,
  LogOut
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { useAuth } from '@/hooks/useAuth';
import { Button } from '@/components/ui/button';

const navigation = [
  { name: '仪表盘', href: '/', icon: LayoutDashboard },
  { name: '插件管理', href: '/plugins', icon: Puzzle },
  { name: '任务管理', href: '/tasks', icon: Clock },
  { name: '配置管理', href: '/config', icon: Settings },
  { name: '日志查看', href: '/logs', icon: FileText },
  { name: 'AI功能', href: '/ai', icon: Brain },
  { name: '预设管理', href: '/presets', icon: FileCode },
];

export function Sidebar() {
  const location = useLocation();
  const { logout } = useAuth();

  return (
    <div className="hidden lg:flex h-full w-64 flex-col border-r bg-card">
      <div className="flex h-16 items-center border-b px-6">
        <h1 className="text-xl font-bold">就你是管理员啊？</h1>
      </div>
      
      <nav className="flex-1 space-y-1 p-4">
        {navigation.map((item) => {
          const isActive = location.pathname === item.href;
          return (
            <Link
              key={item.name}
              to={item.href}
              className={cn(
                'flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors min-h-[44px]',
                isActive
                  ? 'bg-primary text-primary-foreground'
                  : 'text-muted-foreground hover:bg-accent hover:text-accent-foreground'
              )}
            >
              <item.icon className="h-5 w-5" />
              {item.name}
            </Link>
          );
        })}
      </nav>

      <div className="border-t p-4">
        <Button
          variant="ghost"
          className="w-full justify-start gap-3 min-h-[44px]"
          onClick={logout}
        >
          <LogOut className="h-5 w-5" />
          退出登录
        </Button>
      </div>
    </div>
  );
}
