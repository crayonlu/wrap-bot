import { useEffect } from 'react';
import { Puzzle } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Switch } from '@/components/ui/switch';
import { Badge } from '@/components/ui/badge';
import { useBotStore } from '@/stores/bot';
import { apiClient } from '@/lib/api';
import { toast } from 'sonner';

export function Plugins() {
  const { plugins, setPlugins, updatePlugin, setLoading, setError } = useBotStore();

  useEffect(() => {
    const fetchPlugins = async () => {
      try {
        setLoading(true);
        const pluginsData = await apiClient.getPlugins();
        setPlugins(pluginsData);
      } catch (error: any) {
        setError(error.response?.data?.error || '获取插件列表失败');
      } finally {
        setLoading(false);
      }
    };

    fetchPlugins();
  }, [setPlugins, setLoading, setError]);

  const handleToggle = async (name: string) => {
    try {
      const result = await apiClient.togglePlugin(name);
      updatePlugin(name, result.enabled);
      toast.success(`插件 ${name} 已${result.enabled ? '启用' : '禁用'}`);
    } catch (error: any) {
      toast.error(error.response?.data?.error || '操作失败');
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold">插件管理</h1>
        <p className="text-muted-foreground">管理Bot的插件</p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {plugins.map((plugin) => (
          <Card key={plugin.name}>
            <CardHeader className="flex flex-row items-start justify-between space-y-0 pb-2">
              <div className="space-y-1">
                <CardTitle className="text-base">{plugin.name}</CardTitle>
                <p className="text-sm text-muted-foreground">
                  {plugin.description}
                </p>
              </div>
              <Puzzle className="h-5 w-5 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="flex items-center justify-between">
                <Badge variant={plugin.enabled ? 'default' : 'secondary'}>
                  {plugin.enabled ? '已启用' : '已禁用'}
                </Badge>
                <Switch
                  checked={plugin.enabled}
                  onCheckedChange={() => handleToggle(plugin.name)}
                />
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
