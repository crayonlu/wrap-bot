import { useState, useEffect } from 'react';
import { Save, Search } from 'lucide-react';
import { Card, CardContent } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Button } from '@/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { apiClient } from '@/lib/api';
import { toast } from 'sonner';
import type { ConfigItem } from '@/types/api';

const configCategories = {
  'NapCat': ['NAPCAT_HTTP_URL', 'NAPCAT_WS_URL', 'NAPCAT_HTTP_TOKEN', 'NAPCAT_WS_TOKEN'],
  '服务器': ['SERVER_PORT', 'SERVER_ENABLED', 'DEBUG', 'COMMAND_PREFIX'],
  'AI': ['AI_ENABLED', 'AI_URL', 'AI_KEY', 'AI_USE_UNIFIED', 'AI_UNIFIED_MODEL', 'AI_TEXT_MODEL', 'AI_VISION_MODEL', 'AI_TEMPERATURE', 'AI_TOP_P', 'AI_MAX_TOKENS', 'AI_MAX_HISTORY', 'AI_IMAGE_DETAIL', 'AI_TEXT_MODEL_TOOLS', 'AI_VISION_MODEL_TOOLS'],
  '推送': ['TECH_PUSH_GROUPS', 'TECH_PUSH_USERS', 'RSS_PUSH_GROUPS', 'RSS_PUSH_USERS'],
  '权限': ['ALLOWED_USERS', 'ALLOWED_GROUPS', 'ADMIN_IDS'],
  'API': ['HOT_API_HOST', 'HOT_API_KEY', 'RSS_API_HOST', 'SERP_API_KEY', 'WEATHER_API_KEY'],
  '预设': ['SYSTEM_PROMPT_PATH', 'ANALYZER_PROMPT_PATH'],
};

export function Config() {
  const [configs, setConfigs] = useState<ConfigItem[]>([]);
  const [filteredConfigs, setFilteredConfigs] = useState<ConfigItem[]>([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [editedConfigs, setEditedConfigs] = useState<Record<string, string>>({});

  useEffect(() => {
    fetchConfigs();
  }, []);

  useEffect(() => {
    if (searchQuery) {
      const filtered = configs.filter(
        (config) =>
          config.key.toLowerCase().includes(searchQuery.toLowerCase()) ||
          config.description.toLowerCase().includes(searchQuery.toLowerCase())
      );
      setFilteredConfigs(filtered);
    } else {
      setFilteredConfigs(configs);
    }
  }, [searchQuery, configs]);

  const fetchConfigs = async () => {
    try {
      setLoading(true);
      const data = await apiClient.getConfig();
      setConfigs(data);
      setFilteredConfigs(data);
    } catch (error: any) {
      toast.error(error.response?.data?.error || '获取配置失败');
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    const updates = Object.entries(editedConfigs).map(([key, value]) => ({
      key,
      value,
    }));

    if (updates.length === 0) {
      toast.info('没有修改需要保存');
      return;
    }

    try {
      setSaving(true);
      const response = await apiClient.updateConfig(updates);
      
      if (response.status === 'restarting') {
        toast.success(response.message || '配置已保存，容器正在重启...');
        setEditedConfigs({});
        
        setTimeout(() => {
          window.location.reload();
        }, 5000);
      } else {
        toast.success('配置已保存');
        setEditedConfigs({});
        await fetchConfigs();
      }
    } catch (error: any) {
      toast.error(error.response?.data?.error || '保存配置失败');
    } finally {
      setSaving(false);
    }
  };

  const handleConfigChange = (key: string, value: string) => {
    setEditedConfigs((prev) => ({ ...prev, [key]: value }));
  };

  const getConfigValue = (key: string) => {
    return editedConfigs[key] ?? configs.find((c) => c.key === key)?.value ?? '';
  };

  const hasChanges = Object.keys(editedConfigs).length > 0;

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 className="text-2xl lg:text-3xl font-bold">配置管理</h1>
          <p className="text-muted-foreground">管理Bot配置</p>
        </div>
        <Button
          onClick={handleSave}
          disabled={!hasChanges || saving}
          className="min-h-[44px] min-w-[120px]"
        >
          <Save className="mr-2 h-4 w-4" />
          {saving ? '保存中...' : '保存更改'}
        </Button>
      </div>

      <div className="relative">
        <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
        <Input
          placeholder="搜索配置项..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="pl-9"
        />
      </div>

      {loading ? (
        <div className="text-center py-8 text-muted-foreground">加载中...</div>
      ) : (
        <Tabs defaultValue="NapCat" className="space-y-4">
          <TabsList className="grid w-full grid-cols-4 lg:grid-cols-8">
            {Object.keys(configCategories).map((category) => (
              <TabsTrigger key={category} value={category}>
                {category}
              </TabsTrigger>
            ))}
          </TabsList>

          {Object.entries(configCategories).map(([category, keys]) => (
            <TabsContent key={category} value={category}>
              <div className="space-y-4">
                {keys
                  .filter((key) => filteredConfigs.some((c) => c.key === key))
                  .map((key) => {
                    const config = configs.find((c) => c.key === key);
                    if (!config) return null;

                    return (
                      <Card key={key}>
                        <CardContent className="pt-6">
                          <div className="space-y-2">
                            <Label htmlFor={key}>{config.key}</Label>
                            <Input
                              id={key}
                              value={getConfigValue(key)}
                              onChange={(e) => handleConfigChange(key, e.target.value)}
                              placeholder={config.value}
                            />
                            <p className="text-xs text-muted-foreground">
                              {config.description}
                            </p>
                          </div>
                        </CardContent>
                      </Card>
                    );
                  })}
              </div>
            </TabsContent>
          ))}
        </Tabs>
      )}
    </div>
  );
}
