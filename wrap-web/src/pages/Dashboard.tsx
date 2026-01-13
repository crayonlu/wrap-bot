import { useEffect } from 'react';
import { Activity, Clock } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useBotStore } from '@/stores/bot';
import { apiClient } from '@/lib/api';

export function Dashboard() {
  const { status, setLoading, setError } = useBotStore();

  useEffect(() => {
    const fetchStatus = async () => {
      try {
        setLoading(true);
        const statusData = await apiClient.getStatus();
        useBotStore.getState().setStatus(statusData);
      } catch (error: any) {
        setError(error.response?.data?.error || '获取状态失败');
      } finally {
        setLoading(false);
      }
    };

    fetchStatus();
  }, [setLoading, setError]);

  const formatUptime = (seconds: number) => {
    const days = Math.floor(seconds / 86400);
    const hours = Math.floor((seconds % 86400) / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);

    if (days > 0) return `${days}天 ${hours}小时`;
    if (hours > 0) return `${hours}小时 ${minutes}分钟`;
    return `${minutes}分钟`;
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">仪表盘</h1>
        <p className="text-muted-foreground">Bot运行状态</p>
      </div>

      {/* 状态卡片 */}
      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">运行状态</CardTitle>
            <Activity className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {status?.running ? '运行中' : '已停止'}
            </div>
            <p className="text-xs text-muted-foreground">
              {status?.version || '未知版本'}
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">运行时间</CardTitle>
            <Clock className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {status ? formatUptime(status.uptime) : '--'}
            </div>
            <p className="text-xs text-muted-foreground">
              Go {status?.go_version || '未知'}
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
