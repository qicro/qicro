'use client';

import { useEffect, useState } from 'react';
import { useAdminStore } from '@/store/admin';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Plus, Settings } from 'lucide-react';
import { CreateChatModelRequest } from '@/types/admin';

export default function ChatModelsManagement() {
  const {
    chatModels,
    apiKeys,
    isLoading,
    error,
    loadChatModels,
    loadAPIKeys,
    createChatModel,
    clearError,
  } = useAdminStore();

  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [selectedType, setSelectedType] = useState<string>('');
  const [formData, setFormData] = useState<CreateChatModelRequest>({
    type: 'chat',
    name: '',
    value: '',
    provider: '',
    sort_num: 0,
    enabled: true,
    power: 1,
    temperature: 1.0,
    max_tokens: 1024,
    max_context: 4096,
    open: true,
  });

  useEffect(() => {
    loadChatModels();
    loadAPIKeys();
  }, [loadChatModels, loadAPIKeys]);

  useEffect(() => {
    if (error) {
      setTimeout(() => clearError(), 5000);
    }
  }, [error, clearError]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await createChatModel(formData);
      setFormData({
        type: 'chat',
        name: '',
        value: '',
        provider: '',
        sort_num: 0,
        enabled: true,
        power: 1,
        temperature: 1.0,
        max_tokens: 1024,
        max_context: 4096,
        open: true,
      });
      setIsCreateOpen(false);
    } catch (error) {
      // Error is handled by the store
    }
  };

  const handleFilterByType = (type: string) => {
    setSelectedType(type);
    if (type) {
      loadChatModels(type);
    } else {
      loadChatModels();
    }
  };

  const availableAPIKeys = apiKeys.filter(key => 
    key.provider === formData.provider && key.enabled
  );

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">模型管理</h2>
          <p className="text-muted-foreground">管理 AI 聊天和图片生成模型</p>
        </div>
        
        <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="h-4 w-4 mr-2" />
              添加模型
            </Button>
          </DialogTrigger>
          <DialogContent className="max-w-2xl max-h-[80vh] overflow-y-auto">
            <DialogHeader>
              <DialogTitle>添加新模型</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="type">类型</Label>
                  <Select 
                    value={formData.type} 
                    onValueChange={(value) => setFormData({ ...formData, type: value })}
                  >
                    <SelectTrigger>
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="chat">聊天模型</SelectItem>
                      <SelectItem value="img">图片生成</SelectItem>
                      <SelectItem value="embedding">嵌入模型</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div>
                  <Label htmlFor="provider">提供商</Label>
                  <Select 
                    value={formData.provider} 
                    onValueChange={(value) => setFormData({ ...formData, provider: value })}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="选择提供商" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="openai">OpenAI</SelectItem>
                      <SelectItem value="anthropic">Anthropic</SelectItem>
                      <SelectItem value="google">Google</SelectItem>
                      <SelectItem value="cohere">Cohere</SelectItem>
                      <SelectItem value="deepseek">DeepSeek</SelectItem>
                      <SelectItem value="qwen">通义千问</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="name">显示名称</Label>
                  <Input
                    id="name"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    placeholder="例如: GPT-4 Turbo"
                    required
                  />
                </div>

                <div>
                  <Label htmlFor="value">模型值</Label>
                  <Input
                    id="value"
                    value={formData.value}
                    onChange={(e) => setFormData({ ...formData, value: e.target.value })}
                    placeholder="例如: gpt-4-turbo-preview"
                    required
                  />
                </div>
              </div>

              <div className="grid grid-cols-3 gap-4">
                <div>
                  <Label htmlFor="power">算力消耗</Label>
                  <Input
                    id="power"
                    type="number"
                    min="1"
                    value={formData.power}
                    onChange={(e) => setFormData({ ...formData, power: parseInt(e.target.value) })}
                  />
                </div>

                <div>
                  <Label htmlFor="sort_num">排序</Label>
                  <Input
                    id="sort_num"
                    type="number"
                    value={formData.sort_num}
                    onChange={(e) => setFormData({ ...formData, sort_num: parseInt(e.target.value) })}
                  />
                </div>

                <div>
                  <Label htmlFor="temperature">创意度</Label>
                  <Input
                    id="temperature"
                    type="number"
                    step="0.1"
                    min="0"
                    max="2"
                    value={formData.temperature}
                    onChange={(e) => setFormData({ ...formData, temperature: parseFloat(e.target.value) })}
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="max_tokens">最大输出长度</Label>
                  <Input
                    id="max_tokens"
                    type="number"
                    min="1"
                    value={formData.max_tokens}
                    onChange={(e) => setFormData({ ...formData, max_tokens: parseInt(e.target.value) })}
                  />
                </div>

                <div>
                  <Label htmlFor="max_context">最大上下文长度</Label>
                  <Input
                    id="max_context"
                    type="number"
                    min="1"
                    value={formData.max_context}
                    onChange={(e) => setFormData({ ...formData, max_context: parseInt(e.target.value) })}
                  />
                </div>
              </div>

              {formData.provider && (
                <div>
                  <Label htmlFor="api_key_id">关联 API Key</Label>
                  <Select 
                    value={formData.api_key_id || ''} 
                    onValueChange={(value) => setFormData({ ...formData, api_key_id: value || undefined })}
                  >
                    <SelectTrigger>
                      <SelectValue placeholder="选择 API Key (可选)" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="">不关联</SelectItem>
                      {availableAPIKeys.map((key) => (
                        <SelectItem key={key.id} value={key.id}>
                          {key.name}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              )}

              <div className="flex justify-end gap-2">
                <Button type="button" variant="outline" onClick={() => setIsCreateOpen(false)}>
                  取消
                </Button>
                <Button type="submit" disabled={isLoading}>
                  {isLoading ? '创建中...' : '创建'}
                </Button>
              </div>
            </form>
          </DialogContent>
        </Dialog>
      </div>

      {error && (
        <div className="bg-destructive/10 text-destructive px-4 py-2 rounded-md">
          {error}
        </div>
      )}

      <div className="flex items-center gap-2">
        <Label>筛选类型:</Label>
        <Select value={selectedType} onValueChange={handleFilterByType}>
          <SelectTrigger className="w-48">
            <SelectValue placeholder="所有类型" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="">所有类型</SelectItem>
            <SelectItem value="chat">聊天模型</SelectItem>
            <SelectItem value="img">图片生成</SelectItem>
            <SelectItem value="embedding">嵌入模型</SelectItem>
          </SelectContent>
        </Select>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>模型列表</CardTitle>
        </CardHeader>
        <CardContent>
          {isLoading && chatModels.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">
              加载中...
            </div>
          ) : chatModels.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">
              暂无模型
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>名称</TableHead>
                  <TableHead>类型</TableHead>
                  <TableHead>提供商</TableHead>
                  <TableHead>模型值</TableHead>
                  <TableHead>算力</TableHead>
                  <TableHead>状态</TableHead>
                  <TableHead>配置</TableHead>
                  <TableHead>操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {chatModels.map((model) => (
                  <TableRow key={model.id}>
                    <TableCell className="font-medium">{model.name}</TableCell>
                    <TableCell>
                      <Badge variant="outline">{model.type}</Badge>
                    </TableCell>
                    <TableCell>
                      <Badge variant="secondary">{model.provider}</Badge>
                    </TableCell>
                    <TableCell className="font-mono text-sm">{model.value}</TableCell>
                    <TableCell>{model.power} 点</TableCell>
                    <TableCell>
                      <div className="flex gap-1">
                        <Badge variant={model.enabled ? "default" : "secondary"}>
                          {model.enabled ? '启用' : '禁用'}
                        </Badge>
                        {model.open && (
                          <Badge variant="outline">开放</Badge>
                        )}
                      </div>
                    </TableCell>
                    <TableCell className="text-sm text-muted-foreground">
                      T: {model.temperature}, Tokens: {model.max_tokens}
                    </TableCell>
                    <TableCell>
                      <Button variant="ghost" size="sm">
                        <Settings className="h-3 w-3" />
                      </Button>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>
    </div>
  );
}