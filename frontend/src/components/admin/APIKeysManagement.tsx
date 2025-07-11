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
import { Checkbox } from '@/components/ui/checkbox';
import { Plus, Edit, Trash2, Eye, EyeOff } from 'lucide-react';
import { CreateAPIKeyRequest, APIKey } from '@/types/admin';
import { formatDate } from '@/lib/utils';

export default function APIKeysManagement() {
  const {
    apiKeys,
    isLoading,
    error,
    loadAPIKeys,
    createAPIKey,
    updateAPIKey,
    deleteAPIKey,
    clearError,
  } = useAdminStore();

  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [isEditOpen, setIsEditOpen] = useState(false);
  const [editingApiKey, setEditingApiKey] = useState<string | null>(null);
  const [showValues, setShowValues] = useState<Record<string, boolean>>({});
  const [showEditValue, setShowEditValue] = useState(false);
  const [formData, setFormData] = useState<CreateAPIKeyRequest>({
    name: '',
    value: '',
    type: 'chat',
    provider: '',
    api_url: '',
    proxy_url: '',
    enabled: true,
  });

  useEffect(() => {
    loadAPIKeys();
  }, [loadAPIKeys]);

  useEffect(() => {
    if (error) {
      setTimeout(() => clearError(), 5000);
    }
  }, [error, clearError]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      if (editingApiKey) {
        await updateAPIKey(editingApiKey, formData);
      } else {
        await createAPIKey(formData);
      }
      resetForm();
    } catch (error) {
      // Error is handled by the store
    }
  };

  const resetForm = () => {
    setFormData({
      name: '',
      value: '',
      type: 'chat',
      provider: '',
      api_url: '',
      proxy_url: '',
      enabled: true,
    });
    setEditingApiKey(null);
    setIsCreateOpen(false);
    setIsEditOpen(false);
    setShowEditValue(false);
  };

  const handleEdit = (apiKey: APIKey) => {
    setFormData({
      name: apiKey.name,
      value: apiKey.value,
      type: apiKey.type,
      provider: apiKey.provider,
      api_url: apiKey.api_url || '',
      proxy_url: apiKey.proxy_url || '',
      enabled: apiKey.enabled,
    });
    setEditingApiKey(apiKey.id);
    setIsEditOpen(true);
  };

  const handleDelete = async (id: string, name: string) => {
    if (confirm(`Are you sure you want to delete "${name}"?`)) {
      try {
        await deleteAPIKey(id);
      } catch (error) {
        // Error is handled by the store
      }
    }
  };

  const toggleShowValue = (id: string) => {
    setShowValues(prev => ({
      ...prev,
      [id]: !prev[id]
    }));
  };

  const handleToggleEnabled = async (id: string, enabled: boolean) => {
    try {
      await updateAPIKey(id, { enabled: !enabled });
    } catch (error) {
      // Error is handled by the store
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">API Keys 管理</h2>
          <p className="text-muted-foreground">管理 AI 服务提供商的 API 密钥</p>
        </div>
        
        <Dialog open={isCreateOpen} onOpenChange={setIsCreateOpen}>
          <DialogTrigger asChild>
            <Button>
              <Plus className="h-4 w-4 mr-2" />
              添加 API Key
            </Button>
          </DialogTrigger>
          <DialogContent className="max-w-md">
            <DialogHeader>
              <DialogTitle>添加新的 API Key</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <Label htmlFor="name">名称</Label>
                <Input
                  id="name"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  placeholder="例如: OpenAI Production"
                  required
                />
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
                    <SelectItem value="chat">聊天</SelectItem>
                    <SelectItem value="img">图片</SelectItem>
                    <SelectItem value="embedding">嵌入</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div>
                <Label htmlFor="value">API Key</Label>
                <Input
                  id="value"
                  type="password"
                  value={formData.value}
                  onChange={(e) => setFormData({ ...formData, value: e.target.value })}
                  placeholder="输入 API Key"
                  required
                />
              </div>

              <div>
                <Label htmlFor="api_url">API 地址 (可选)</Label>
                <Input
                  id="api_url"
                  value={formData.api_url}
                  onChange={(e) => setFormData({ ...formData, api_url: e.target.value })}
                  placeholder="例如: https://api.openai.com/v1"
                />
              </div>

              <div>
                <Label htmlFor="proxy_url">代理地址 (可选)</Label>
                <Input
                  id="proxy_url"
                  value={formData.proxy_url}
                  onChange={(e) => setFormData({ ...formData, proxy_url: e.target.value })}
                  placeholder="例如: http://proxy.example.com:8080"
                />
              </div>

              <div className="flex items-center space-x-2">
                <Checkbox
                  id="enabled"
                  checked={formData.enabled}
                  onCheckedChange={(checked) => setFormData({ ...formData, enabled: checked as boolean })}
                />
                <Label htmlFor="enabled">启用</Label>
              </div>

              <div className="flex justify-end gap-2">
                <Button type="button" variant="outline" onClick={resetForm}>
                  取消
                </Button>
                <Button type="submit" disabled={isLoading}>
                  {isLoading ? '创建中...' : '创建'}
                </Button>
              </div>
            </form>
          </DialogContent>
        </Dialog>

        {/* Edit Dialog */}
        <Dialog open={isEditOpen} onOpenChange={setIsEditOpen}>
          <DialogContent className="max-w-md">
            <DialogHeader>
              <DialogTitle>编辑 API Key</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <Label htmlFor="edit-name">名称</Label>
                <Input
                  id="edit-name"
                  value={formData.name}
                  onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                  placeholder="例如: OpenAI Production"
                  required
                />
              </div>

              <div>
                <Label htmlFor="edit-provider">提供商</Label>
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

              <div>
                <Label htmlFor="edit-type">类型</Label>
                <Select 
                  value={formData.type} 
                  onValueChange={(value) => setFormData({ ...formData, type: value })}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="chat">聊天</SelectItem>
                    <SelectItem value="img">图片</SelectItem>
                    <SelectItem value="embedding">嵌入</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div>
                <Label htmlFor="edit-value">API Key</Label>
                <div className="relative">
                  <Input
                    id="edit-value"
                    type={showEditValue ? "text" : "password"}
                    value={formData.value}
                    onChange={(e) => setFormData({ ...formData, value: e.target.value })}
                    placeholder="输入 API Key"
                    required
                    className="pr-10"
                  />
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    className="absolute right-1 top-1/2 -translate-y-1/2 h-6 w-6 p-0"
                    onClick={() => setShowEditValue(!showEditValue)}
                  >
                    {showEditValue ? (
                      <EyeOff className="h-3 w-3" />
                    ) : (
                      <Eye className="h-3 w-3" />
                    )}
                  </Button>
                </div>
              </div>

              <div>
                <Label htmlFor="edit-api_url">API 地址 (可选)</Label>
                <Input
                  id="edit-api_url"
                  value={formData.api_url}
                  onChange={(e) => setFormData({ ...formData, api_url: e.target.value })}
                  placeholder="例如: https://api.openai.com/v1"
                />
              </div>

              <div>
                <Label htmlFor="edit-proxy_url">代理地址 (可选)</Label>
                <Input
                  id="edit-proxy_url"
                  value={formData.proxy_url}
                  onChange={(e) => setFormData({ ...formData, proxy_url: e.target.value })}
                  placeholder="例如: http://proxy.example.com:8080"
                />
              </div>

              <div className="flex items-center space-x-2">
                <Checkbox
                  id="edit-enabled"
                  checked={formData.enabled}
                  onCheckedChange={(checked) => setFormData({ ...formData, enabled: checked as boolean })}
                />
                <Label htmlFor="edit-enabled">启用</Label>
              </div>

              <div className="flex justify-end gap-2">
                <Button type="button" variant="outline" onClick={resetForm}>
                  取消
                </Button>
                <Button type="submit" disabled={isLoading}>
                  {isLoading ? '更新中...' : '更新'}
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

      <Card>
        <CardHeader>
          <CardTitle>API Keys 列表</CardTitle>
        </CardHeader>
        <CardContent>
          {isLoading && apiKeys.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">
              加载中...
            </div>
          ) : apiKeys.length === 0 ? (
            <div className="text-center py-8 text-muted-foreground">
              暂无 API Keys
            </div>
          ) : (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>名称</TableHead>
                  <TableHead>提供商</TableHead>
                  <TableHead>类型</TableHead>
                  <TableHead>API Key</TableHead>
                  <TableHead>状态</TableHead>
                  <TableHead>最后使用</TableHead>
                  <TableHead>操作</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {apiKeys.map((apiKey) => (
                  <TableRow key={apiKey.id}>
                    <TableCell className="font-medium">{apiKey.name}</TableCell>
                    <TableCell>
                      <Badge variant="outline">{apiKey.provider}</Badge>
                    </TableCell>
                    <TableCell>
                      <Badge variant="secondary">{apiKey.type}</Badge>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <span className="font-mono text-sm">
                          {showValues[apiKey.id] ? apiKey.value : '***'}
                        </span>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => toggleShowValue(apiKey.id)}
                        >
                          {showValues[apiKey.id] ? (
                            <EyeOff className="h-3 w-3" />
                          ) : (
                            <Eye className="h-3 w-3" />
                          )}
                        </Button>
                      </div>
                    </TableCell>
                    <TableCell>
                      <Button
                        variant="ghost"
                        size="sm"
                        onClick={() => handleToggleEnabled(apiKey.id, apiKey.enabled)}
                      >
                        <Badge variant={apiKey.enabled ? "default" : "secondary"}>
                          {apiKey.enabled ? '启用' : '禁用'}
                        </Badge>
                      </Button>
                    </TableCell>
                    <TableCell className="text-sm text-muted-foreground">
                      {apiKey.last_used_at 
                        ? formatDate(apiKey.last_used_at)
                        : '从未使用'
                      }
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-1">
                        <Button 
                          variant="ghost" 
                          size="sm"
                          onClick={() => handleEdit(apiKey)}
                        >
                          <Edit className="h-3 w-3" />
                        </Button>
                        <Button 
                          variant="ghost" 
                          size="sm"
                          onClick={() => handleDelete(apiKey.id, apiKey.name)}
                        >
                          <Trash2 className="h-3 w-3" />
                        </Button>
                      </div>
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