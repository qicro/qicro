'use client';

import { useState } from 'react';
import { useAuthStore } from '@/store/auth';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Settings, Key, Bot, Tag, LogOut, ArrowLeft } from 'lucide-react';
import APIKeysManagement from './APIKeysManagement';
import ChatModelsManagement from './ChatModelsManagement';
import Link from 'next/link';

export default function AdminDashboard() {
  const { user, logout } = useAuthStore();
  const [activeTab, setActiveTab] = useState('api-keys');

  const handleLogout = () => {
    logout();
    window.location.href = '/auth';
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center py-4">
            <div className="flex items-center gap-4">
              <Link href="/chat">
                <Button variant="ghost" size="sm">
                  <ArrowLeft className="h-4 w-4 mr-2" />
                  返回聊天
                </Button>
              </Link>
              <div>
                <h1 className="text-2xl font-bold text-gray-900 flex items-center gap-2">
                  <Settings className="h-6 w-6" />
                  系统管理
                </h1>
                <p className="text-sm text-gray-500">Qicro 配置管理中心</p>
              </div>
            </div>
            <div className="flex items-center space-x-4">
              <div className="flex items-center space-x-2">
                <span className="text-gray-700">{user?.email}</span>
                <span className="text-xs bg-blue-100 text-blue-800 px-2 py-1 rounded">
                  管理员
                </span>
              </div>
              <Button
                variant="outline"
                size="sm"
                onClick={handleLogout}
              >
                <LogOut className="h-4 w-4 mr-2" />
                退出
              </Button>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto py-8 px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
          {/* Sidebar */}
          <div className="lg:col-span-1">
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">管理功能</CardTitle>
              </CardHeader>
              <CardContent className="space-y-2">
                <Button
                  variant={activeTab === 'api-keys' ? 'default' : 'ghost'}
                  className="w-full justify-start"
                  onClick={() => setActiveTab('api-keys')}
                >
                  <Key className="h-4 w-4 mr-2" />
                  API Keys
                </Button>
                <Button
                  variant={activeTab === 'models' ? 'default' : 'ghost'}
                  className="w-full justify-start"
                  onClick={() => setActiveTab('models')}
                >
                  <Bot className="h-4 w-4 mr-2" />
                  模型管理
                </Button>
                <Button
                  variant={activeTab === 'app-types' ? 'default' : 'ghost'}
                  className="w-full justify-start"
                  onClick={() => setActiveTab('app-types')}
                >
                  <Tag className="h-4 w-4 mr-2" />
                  应用分类
                </Button>
              </CardContent>
            </Card>

            {/* Quick Stats */}
            <Card className="mt-6">
              <CardHeader>
                <CardTitle className="text-lg">系统概览</CardTitle>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-600">活跃 API Keys</span>
                  <span className="font-semibold">5</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-600">可用模型</span>
                  <span className="font-semibold">12</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-gray-600">应用分类</span>
                  <span className="font-semibold">4</span>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Content Area */}
          <div className="lg:col-span-3">
            <Tabs value={activeTab} onValueChange={setActiveTab}>
              <TabsList className="hidden">
                <TabsTrigger value="api-keys">API Keys</TabsTrigger>
                <TabsTrigger value="models">模型管理</TabsTrigger>
                <TabsTrigger value="app-types">应用分类</TabsTrigger>
              </TabsList>

              <TabsContent value="api-keys">
                <APIKeysManagement />
              </TabsContent>

              <TabsContent value="models">
                <ChatModelsManagement />
              </TabsContent>

              <TabsContent value="app-types">
                <Card>
                  <CardHeader>
                    <CardTitle>应用分类管理</CardTitle>
                  </CardHeader>
                  <CardContent>
                    <div className="text-center py-8 text-muted-foreground">
                      应用分类管理功能开发中...
                    </div>
                  </CardContent>
                </Card>
              </TabsContent>
            </Tabs>
          </div>
        </div>
      </main>
    </div>
  );
}