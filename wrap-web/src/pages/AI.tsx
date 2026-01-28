import { useState, useEffect } from 'react';
import { Image as ImageIcon } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { Label } from '@/components/ui/label';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Badge } from '@/components/ui/badge';
import { apiClient } from '@/lib/api';
import { toast } from 'sonner';
import type { AITool } from '@/types/api';

export function AI() {
  const [tools, setTools] = useState<AITool[]>([]);
  const [textMessage, setTextMessage] = useState('');
  const [imageMessage, setImageMessage] = useState('');
  const [imageUrls, setImageUrls] = useState('');
  const [textResponse, setTextResponse] = useState('');
  const [imageResponse, setImageResponse] = useState('');
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    fetchTools();
  }, []);

  const fetchTools = async () => {
    try {
      const toolsData = await apiClient.getAITools();
      setTools(toolsData);
    } catch (error: any) {
      toast.error(error.response?.data?.error || '获取AI工具失败');
    }
  };

  const handleTextChat = async () => {
    if (!textMessage.trim()) {
      toast.error('请输入消息');
      return;
    }

    try {
      setLoading(true);
      const response = await apiClient.testAIChat({ message: textMessage });
      setTextResponse(response.response);
      toast.success('对话成功');
    } catch (error: any) {
      toast.error(error.response?.data?.error || '对话失败');
    } finally {
      setLoading(false);
    }
  };

  const handleImageChat = async () => {
    if (!imageMessage.trim() || !imageUrls.trim()) {
      toast.error('请输入消息和图片URL');
      return;
    }

    const urls = imageUrls.split('\n').filter(url => url.trim());

    try {
      setLoading(true);
      const response = await apiClient.testAIImageChat({
        message: imageMessage,
        images: urls,
      });
      setImageResponse(response.response);
      toast.success('对话成功');
    } catch (error: any) {
      toast.error(error.response?.data?.error || '对话失败');
    } finally {
      setLoading(false);
    }
  };

  const getToolStatusBadge = (enabled: boolean) => {
    return enabled ? <Badge>已启用</Badge> : <Badge variant="secondary">未启用</Badge>;
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl lg:text-3xl font-bold">AI功能</h1>
        <p className="text-muted-foreground">AI工具和对话测试</p>
      </div>

      <Tabs defaultValue="tools" className="space-y-4">
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="tools">工具列表</TabsTrigger>
          <TabsTrigger value="text-chat">文本对话</TabsTrigger>
          <TabsTrigger value="image-chat">图像对话</TabsTrigger>
        </TabsList>

        <TabsContent value="tools">
          <Card>
            <CardHeader>
              <CardTitle>AI工具</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-3">
                {tools.map((tool) => (
                  <div
                    key={tool.name}
                    className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 rounded-lg border p-4"
                  >
                    <div className="space-y-1">
                      <div className="flex items-center gap-2">
                        <h3 className="font-semibold">{tool.name}</h3>
                        {getToolStatusBadge(tool.enabled)}
                      </div>
                      <p className="text-sm text-muted-foreground">
                        {tool.description}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="text-chat">
          <Card>
            <CardHeader>
              <CardTitle>文本对话测试</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="text-message">消息</Label>
                <Textarea
                  id="text-message"
                  placeholder="输入要发送给AI的消息..."
                  value={textMessage}
                  onChange={(e) => setTextMessage(e.target.value)}
                  rows={4}
                />
              </div>
              <Button onClick={handleTextChat} disabled={loading} className="min-h-[44px] w-full sm:w-auto">
                {loading ? '发送中...' : '发送'}
              </Button>
              {textResponse && (
                <div className="space-y-2">
                  <Label>AI响应</Label>
                  <div className="rounded-lg border p-4">
                    <p className="whitespace-pre-wrap">{textResponse}</p>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="image-chat">
          <Card>
            <CardHeader>
              <CardTitle>图像对话测试</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="image-message">消息</Label>
                <Textarea
                  id="image-message"
                  placeholder="输入要发送给AI的消息..."
                  value={imageMessage}
                  onChange={(e) => setImageMessage(e.target.value)}
                  rows={3}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="image-urls">图片URL (每行一个)</Label>
                <Textarea
                  id="image-urls"
                  placeholder="https://example.com/image1.jpg&#10;https://example.com/image2.jpg"
                  value={imageUrls}
                  onChange={(e) => setImageUrls(e.target.value)}
                  rows={3}
                />
              </div>
              <Button onClick={handleImageChat} disabled={loading} className="min-h-[44px] w-full sm:w-auto">
                <ImageIcon className="mr-2 h-4 w-4" />
                {loading ? '发送中...' : '发送'}
              </Button>
              {imageResponse && (
                <div className="space-y-2">
                  <Label>AI响应</Label>
                  <div className="rounded-lg border p-4">
                    <p className="whitespace-pre-wrap">{imageResponse}</p>
                  </div>
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}
