import { Bell, Menu } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { useBotStore } from '@/stores/bot';

interface HeaderProps {
  onMenuClick: () => void;
  showMenuButton: boolean;
}

export function Header({ onMenuClick, showMenuButton }: HeaderProps) {
  const { status } = useBotStore();

  return (
    <header className="flex h-16 items-center justify-between border-b bg-card px-4 lg:px-6">
      <div className="flex items-center gap-4">
        {showMenuButton && (
          <Button variant="ghost" size="icon" onClick={onMenuClick} className="lg:hidden min-h-[44px] min-w-[44px]">
            <Menu className="h-5 w-5" />
          </Button>
        )}
        <h2 className="text-lg font-semibold">管理面板</h2>
        {status && (
          <Badge variant={status.running ? 'default' : 'destructive'}>
            {status.running ? '运行中' : '已停止'}
          </Badge>
        )}
      </div>

      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" className="min-h-[44px] min-w-[44px]">
          <Bell className="h-5 w-5" />
        </Button>
      </div>
    </header>
  );
}
