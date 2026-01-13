import { Link, useLocation } from 'react-router-dom';
import { 
  LayoutDashboard, 
  Puzzle, 
  Clock, 
  Settings, 
  FileText, 
  Brain, 
  FileCode,
  LogOut,
  X
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { useAuth } from '@/hooks/useAuth';
import { Button } from '@/components/ui/button';
import { Sheet, SheetContent, SheetHeader, SheetTitle } from '@/components/ui/sheet';

const navigation = [
  { name: '仪表盘', href: '/', icon: LayoutDashboard },
  { name: '插件管理', href: '/plugins', icon: Puzzle },
  { name: '任务管理', href: '/tasks', icon: Clock },
  { name: '配置管理', href: '/config', icon: Settings },
  { name: '日志查看', href: '/logs', icon: FileText },
  { name: 'AI功能', href: '/ai', icon: Brain },
  { name: '预设管理', href: '/presets', icon: FileCode },
];

interface MobileNavProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function MobileNav({ open, onOpenChange }: MobileNavProps) {
  const location = useLocation();
  const { logout } = useAuth();

  const handleNavigate = () => {
    onOpenChange(false);
  };

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent side="left" className="w-80 p-0">
        <SheetHeader className="border-b p-6">
          <div className="flex items-center justify-between">
            <SheetTitle className="text-xl font-bold">就你是管理员啊？</SheetTitle>
            <Button
              variant="ghost"
              size="icon"
              onClick={() => onOpenChange(false)}
            >
              <X className="h-5 w-5" />
            </Button>
          </div>
        </SheetHeader>

        <nav className="flex-1 space-y-1 p-4">
          {navigation.map((item) => {
            const isActive = location.pathname === item.href;
            return (
              <Link
                key={item.name}
                to={item.href}
                onClick={handleNavigate}
                className={cn(
                  'flex items-center gap-3 rounded-lg px-4 py-3 text-sm font-medium transition-colors min-h-[44px]',
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
            onClick={() => {
              logout();
              onOpenChange(false);
            }}
          >
            <LogOut className="h-5 w-5" />
            退出登录
          </Button>
        </div>
      </SheetContent>
    </Sheet>
  );
}
