import { useState, useEffect } from 'react';
import { FileCode, Save, Eye } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { apiClient } from '@/lib/api';
import { toast } from 'sonner';
import type { Preset } from '@/types/api';

export function Presets() {
  const [presets, setPresets] = useState<Preset[]>([]);
  const [selectedPreset, setSelectedPreset] = useState<Preset | null>(null);
  const [content, setContent] = useState('');
  const [previewContent, setPreviewContent] = useState('');
  const [saving, setSaving] = useState(false);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetchPresets();
  }, []);

  const fetchPresets = async () => {
    try {
      setLoading(true);
      const data = await apiClient.getPresets();
      setPresets(data);
    } catch (error: any) {
      toast.error(error.response?.data?.error || '获取预设列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleSelect = async (preset: Preset) => {
    try {
      const data = await apiClient.getPreset(preset.name);
      setSelectedPreset(data);
      setContent(data.content);
    } catch (error: any) {
      toast.error(error.response?.data?.error || '获取预设内容失败');
    }
  };

  const handleSave = async () => {
    if (!selectedPreset) {
      toast.error('请先选择一个预设');
      return;
    }

    try {
      setSaving(true);
      await apiClient.updatePreset(selectedPreset.name, { content });
      toast.success('预设已保存');
    } catch (error: any) {
      toast.error(error.response?.data?.error || '保存预设失败');
    } finally {
      setSaving(false);
    }
  };

  const handlePreview = () => {
    setPreviewContent(content);
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl lg:text-3xl font-bold">预设管理</h1>
        <p className="text-muted-foreground">管理系统预设文件</p>
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        <Card className="lg:col-span-1">
          <CardHeader>
            <CardTitle>预设文件</CardTitle>
          </CardHeader>
          <CardContent>
            {loading ? (
              <div className="text-center py-4 text-muted-foreground">加载中...</div>
            ) : presets.length === 0 ? (
              <div className="text-center py-4 text-muted-foreground">暂无预设</div>
            ) : (
              <div className="space-y-2">
                {presets.map((preset) => (
                  <div
                    key={preset.name}
                    className={`flex cursor-pointer items-center justify-between rounded-lg border p-3 transition-colors hover:bg-accent ${
                      selectedPreset?.name === preset.name ? 'bg-accent' : ''
                    }`}
                    onClick={() => handleSelect(preset)}
                  >
                    <div className="flex items-center gap-2">
                      <FileCode className="h-4 w-4 text-muted-foreground" />
                      <span className="font-medium">{preset.name}</span>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>

        <Card className="lg:col-span-2">
          <CardHeader className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
            <CardTitle>
              {selectedPreset ? selectedPreset.name : '选择预设文件'}
            </CardTitle>
            <div className="flex gap-2">
              <Dialog>
                <DialogTrigger asChild>
                  <Button variant="outline" onClick={handlePreview} disabled={!selectedPreset} className="min-h-[44px]">
                    <Eye className="mr-2 h-4 w-4" />
                    预览
                  </Button>
                </DialogTrigger>
                <DialogContent className="max-w-2xl">
                  <DialogHeader>
                    <DialogTitle>预览: {selectedPreset?.name}</DialogTitle>
                  </DialogHeader>
                  <div className="max-h-[500px] overflow-auto rounded-lg border p-4">
                    <pre className="whitespace-pre-wrap font-mono text-sm">
                      {previewContent}
                    </pre>
                  </div>
                </DialogContent>
              </Dialog>
              <Button onClick={handleSave} disabled={!selectedPreset || saving} className="min-h-[44px]">
                <Save className="mr-2 h-4 w-4" />
                {saving ? '保存中...' : '保存'}
              </Button>
            </div>
          </CardHeader>
          <CardContent>
            {selectedPreset ? (
              <div className="space-y-2">
                <Textarea
                  value={content}
                  onChange={(e) => setContent(e.target.value)}
                  placeholder="编辑预设内容..."
                  className="min-h-[400px] font-mono"
                />
                <p className="text-xs text-muted-foreground">
                  路径: {selectedPreset.path}
                </p>
              </div>
            ) : (
              <div className="flex h-[400px] items-center justify-center text-muted-foreground">
                请从左侧选择一个预设文件进行编辑
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
