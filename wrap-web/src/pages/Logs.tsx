import { useState, useEffect, useRef } from 'react';
import { FileText, Filter, Trash2, Download } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { ScrollArea } from '@/components/ui/scroll-area';
import { Badge } from '@/components/ui/badge';
import { useBotStore } from '@/stores/bot';
import { apiClient } from '@/lib/api';
import { format } from 'date-fns';

export function Logs() {
  const { logs, setLogs, setLoading, setError } = useBotStore();
  const [levelFilter, setLevelFilter] = useState<string>('all');
  const [autoScroll, setAutoScroll] = useState(true);
  const scrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const fetchLogs = async () => {
      try {
        setLoading(true);
        const logsData = await apiClient.getLogs(levelFilter === 'all' ? undefined : levelFilter, 200);
        setLogs(logsData);
      } catch (error: any) {
        setError(error.response?.data?.error || '获取日志失败');
      } finally {
        setLoading(false);
      }
    };

    fetchLogs();
  }, [levelFilter, setLogs, setLoading, setError]);

  useEffect(() => {
    if (autoScroll && scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [logs, autoScroll]);

  const handleClear = () => {
    setLogs([]);
  };

  const handleExport = () => {
    const logText = logs
      .map(
        (log) =>
          `[${log.timestamp}] [${log.level}] ${log.message}${
            log.context ? ` ${JSON.stringify(log.context)}` : ''
          }`
      )
      .join('\n');

    const blob = new Blob([logText], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `logs-${format(new Date(), 'yyyy-MM-dd-HH-mm-ss')}.txt`;
    a.click();
    URL.revokeObjectURL(url);
  };

  const getLogColor = (level: string) => {
    switch (level) {
      case 'ERROR':
        return 'text-destructive';
      case 'WARN':
        return 'text-yellow-500';
      case 'INFO':
        return 'text-blue-500';
      case 'DEBUG':
        return 'text-muted-foreground';
      default:
        return 'text-muted-foreground';
    }
  };

  const filteredLogs =
    levelFilter === 'all' ? logs : logs.filter((log) => log.level === levelFilter);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">日志查看</h1>
          <p className="text-muted-foreground">查看Bot运行日志</p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" onClick={handleClear}>
            <Trash2 className="mr-2 h-4 w-4" />
            清空
          </Button>
          <Button variant="outline" onClick={handleExport}>
            <Download className="mr-2 h-4 w-4" />
            导出
          </Button>
        </div>
      </div>

      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <CardTitle>日志流</CardTitle>
            <div className="flex items-center gap-4">
              <div className="flex items-center gap-2">
                <Filter className="h-4 w-4 text-muted-foreground" />
                <Select value={levelFilter} onValueChange={setLevelFilter}>
                  <SelectTrigger className="w-[120px]">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">全部</SelectItem>
                    <SelectItem value="DEBUG">DEBUG</SelectItem>
                    <SelectItem value="INFO">INFO</SelectItem>
                    <SelectItem value="WARN">WARN</SelectItem>
                    <SelectItem value="ERROR">ERROR</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <label className="flex items-center gap-2 text-sm">
                <input
                  type="checkbox"
                  checked={autoScroll}
                  onChange={(e) => setAutoScroll(e.target.checked)}
                  className="rounded"
                />
                自动滚动
              </label>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <ScrollArea className="h-[600px] w-full rounded-md border p-4">
            <div ref={scrollRef} className="space-y-2 font-mono text-sm">
              {filteredLogs.length === 0 ? (
                <div className="flex h-full items-center justify-center text-muted-foreground">
                  <FileText className="mr-2 h-8 w-8" />
                  暂无日志
                </div>
              ) : (
                filteredLogs.map((log, index) => (
                  <div key={index} className="flex gap-3">
                    <span className="text-muted-foreground shrink-0">
                      {format(new Date(log.timestamp), 'HH:mm:ss')}
                    </span>
                    <Badge variant="outline" className="shrink-0">
                      {log.level}
                    </Badge>
                    <span className={getLogColor(log.level)}>
                      {log.message}
                    </span>
                    {log.context && (
                      <span className="text-muted-foreground">
                        {JSON.stringify(log.context)}
                      </span>
                    )}
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
