import { useEffect } from 'react';
import { Clock, Play } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { useBotStore } from '@/stores/bot';
import { apiClient } from '@/lib/api';
import { toast } from 'sonner';
import { format } from 'date-fns';
import { zhCN } from 'date-fns/locale';

export function Tasks() {
  const { tasks, setTasks, setLoading, setError } = useBotStore();

  useEffect(() => {
    const fetchTasks = async () => {
      try {
        setLoading(true);
        const tasksData = await apiClient.getTasks();
        setTasks(tasksData);
      } catch (error: any) {
        setError(error.response?.data?.error || '获取任务列表失败');
      } finally {
        setLoading(false);
      }
    };

    fetchTasks();
  }, [setTasks, setLoading, setError]);

  const handleTrigger = async (id: string, name: string) => {
    try {
      await apiClient.triggerTask(id);
      toast.success(`任务 ${name} 已触发`);
    } catch (error: any) {
      toast.error(error.response?.data?.error || '触发任务失败');
    }
  };

  const formatSchedule = (cron: string) => {
    const parts = cron.split(' ');
    if (parts.length !== 5) return cron;
    
    const [minute, hour] = parts;
    
    if (minute === '0' && hour.includes('*/')) {
      const hours = parseInt(hour.replace('*/', ''));
      return `每 ${hours} 小时`;
    }
    
    return cron;
  };

  const formatDate = (dateStr: string) => {
    if (!dateStr) return '从未运行';
    return format(new Date(dateStr), 'yyyy-MM-dd HH:mm:ss', { locale: zhCN });
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">任务管理</h1>
        <p className="text-muted-foreground">管理定时任务</p>
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        {tasks.map((task) => (
          <Card key={task.id}>
            <CardHeader className="flex flex-row items-start justify-between space-y-0 pb-2">
              <div className="space-y-1">
                <CardTitle className="text-base">{task.name}</CardTitle>
                <p className="text-sm text-muted-foreground">
                  {task.description}
                </p>
              </div>
              <Clock className="h-5 w-5 text-muted-foreground" />
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">调度:</span>
                  <Badge variant="outline">{formatSchedule(task.schedule)}</Badge>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">下次运行:</span>
                  <span>{formatDate(task.next_run)}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">上次运行:</span>
                  <span>{formatDate(task.last_run)}</span>
                </div>
              </div>
              {task.can_trigger && (
                <Button
                  className="w-full"
                  onClick={() => handleTrigger(task.id, task.name)}
                >
                  <Play className="mr-2 h-4 w-4" />
                  手动触发
                </Button>
              )}
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
