import { useEffect } from 'react';
import { Activity, Clock, Zap, AlertCircle } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { ScrollArea } from '@/components/ui/scroll-area';
import { useBotStore } from '@/stores/bot';
import { apiClient } from '@/lib/api';

export function Dashboard() {
  const { status, logs, setLoading, setError } = useBotStore();

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

  const getLogColor = (level: string) => {
    switch (level) {
      case 'ERROR':
        return 'text-destructive';
      case 'WARN':
        return 'text-yellow-500';
      case 'INFO':
        return 'text-blue-500';
      default:
        return 'text-muted-foreground';
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">仪表盘</h1>
        <p className="text-muted-foreground">Bot运行状态和实时日志</p>
      </div>

      {/* 状态卡片 */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
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

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">活跃插件</CardTitle>
            <Zap className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {useBotStore.getState().plugins.filter(p => p.enabled).length}
            </div>
            <p className="text-xs text-muted-foreground">
              共 {useBotStore.getState().plugins.length} 个插件
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">活跃任务</CardTitle>
            <AlertCircle className="h-4 w-4 text-muted-foreground" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">
              {useBotStore.getState().tasks.length}
            </div>
            <p className="text-xs text-muted-foreground">
              定时任务
            </p>
          </CardContent>
        </Card>
      </div>

      {/* 实时日志 */}
      <Card>
        <CardHeader>
          <CardTitle>实时日志</CardTitle>
        </CardHeader>
        <CardContent>
          <ScrollArea className="h-[400px] w-full rounded-md border p-4">
            <div className="space-y-2 font-mono text-sm">
              {logs.length === 0 ? (
                <p className="text-muted-foreground">暂无日志</p>
              ) : (
                logs.map((log, index) => (
                  <div key={index} className="flex gap-2">
                    <span className="text-muted-foreground shrink-0">
                      {new Date(log.timestamp).toLocaleTimeString('zh-CN')}
                    </span>
                    <Badge variant="outline" className="shrink-0">
                      {log.level}
                    </Badge>
                    <span className={getLogColor(log.level)}>
                      {log.message}
                    </span>
                  </div>
                ))
              )}
            </div>
          </ScrollArea>
        </CardContent>
      </Card>
    </div>
  );
}
