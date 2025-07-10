'use client';

import { useEffect } from 'react';
import { useAuthStore } from '@/store/auth';
import { useChatStore } from '@/store/chat';
import { useWebSocket } from '@/hooks/useWebSocket';
import ConversationSidebar from './ConversationSidebar';
import MessageList from './MessageList';
import MessageInput from './MessageInput';
import { Card } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Wifi, WifiOff, Settings, Bot } from 'lucide-react';
import Link from 'next/link';
import { ThemeToggle } from '@/components/theme-toggle';

export default function Chat() {
  const { user } = useAuthStore();
  const { 
    currentConversation, 
    isLoading, 
    isStreaming, 
    error,
    sendMessage,
    clearError,
    createConversation,
    models,
    selectedModel,
    setSelectedModel,
    loadModels
  } = useChatStore();

  const { sendMessage: sendWebSocketMessage, isConnected } = useWebSocket({
    onMessage: (data: unknown) => {
      console.log('Received WebSocket message:', data);
      // Handle real-time updates here
      if (typeof data === 'object' && data !== null && 'type' in data) {
        const messageData = data as { type: string; [key: string]: unknown };
        if (messageData.type === 'message_update') {
          // Update message in store
        }
      }
    },
    onConnect: () => {
      console.log('Connected to WebSocket');
    },
    onDisconnect: () => {
      console.log('Disconnected from WebSocket');
    },
    onError: (error) => {
      console.warn('WebSocket connection error:', error);
    },
  });

  useEffect(() => {
    if (error) {
      console.warn('Chat error:', error);
      setTimeout(() => clearError(), 5000);
    }
  }, [error, clearError]);

  useEffect(() => {
    if (user) {
      loadModels();
    }
  }, [user, loadModels]);

  const handleSendMessage = async (content: string) => {
    console.log('Debug: handleSendMessage called with content:', content);
    
    let conversationToUse = currentConversation;
    
    // 如果没有当前对话，创建一个新对话
    if (!conversationToUse) {
      try {
        // 使用选中的模型或默认模型创建新对话
        const validModels = models.filter(model => model.id && model.id.trim() !== '');
        const modelToUse = selectedModel || (validModels.length > 0 ? validModels[0].id : 'gpt-3.5-turbo');
        console.log('Debug: Creating new conversation with model:', modelToUse);
        conversationToUse = await createConversation('New Chat', modelToUse);
        console.log('Debug: New conversation created:', conversationToUse);
      } catch (error) {
        console.error('Failed to create conversation:', error);
        return;
      }
    }
    
    console.log('Debug: Using conversation:', conversationToUse.id);
    
    // Send via WebSocket for real-time updates
    console.log('Debug: Sending WebSocket message');
    sendWebSocketMessage({
      type: 'new_message',
      conversation_id: conversationToUse.id,
      content: content,
      timestamp: new Date().toISOString(),
    });
    
    // Also send via HTTP API for persistence and LLM response
    console.log('Debug: Calling sendMessage API with stream=true');
    try {
      await sendMessage(content, true);
      console.log('Debug: sendMessage completed successfully');
    } catch (error) {
      console.error('Debug: sendMessage failed:', error);
    }
  };

  if (!user) {
    return (
      <div className="h-screen flex">
        <div className="w-80 border-r bg-muted/50 flex flex-col">
          <div className="p-4 border-b">
            <h1 className="text-xl font-bold text-foreground">Qicro</h1>
            <p className="text-sm text-muted-foreground">AI 聊天平台</p>
          </div>
          <div className="flex-1 flex items-center justify-center">
            <div className="text-center space-y-4">
              <div className="absolute top-4 right-4">
                <ThemeToggle />
              </div>
              <h3 className="text-lg font-semibold">请先登录</h3>
              <p className="text-sm text-muted-foreground">登录后即可开始聊天</p>
              <Link href="/auth">
                <Button>
                  登录
                </Button>
              </Link>
            </div>
          </div>
        </div>
        
        <div className="flex-1 flex flex-col">
          <div className="flex items-center justify-between bg-muted/30 px-4 py-2">
            <div className="flex items-center gap-2">
              <WifiOff className="h-3 w-3 text-red-500" />
              <span className="text-red-600 text-sm">未连接</span>
            </div>
            <div className="absolute top-4 right-4">
              <ThemeToggle />
            </div>
          </div>
          
          <div className="flex-1 flex items-center justify-center">
            <Card className="p-8 max-w-md text-center">
              <h2 className="text-2xl font-bold mb-4">欢迎使用 Qicro</h2>
              <p className="text-muted-foreground mb-6">
                智能 AI 聊天平台，支持多模型对话、实时通信和高级配置管理
              </p>
              <div className="space-y-3">
                <Link href="/auth">
                  <Button className="w-full">
                    开始聊天
                  </Button>
                </Link>
                <p className="text-sm text-muted-foreground">
                  需要登录后才能使用聊天功能
                </p>
              </div>
            </Card>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="h-screen flex">
      <ConversationSidebar />
      
      <div className="flex-1 flex flex-col">
        {/* Connection Status and Model Selector */}
        <div className="flex items-center justify-between bg-muted/30 px-4 py-2">
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2">
              {isConnected ? (
                <>
                  <Wifi className="h-3 w-3 text-green-500" />
                  <span className="text-green-600 text-xs">Connected</span>
                </>
              ) : (
                <>
                  <WifiOff className="h-3 w-3 text-red-500" />
                  <span className="text-red-600 text-xs">Disconnected</span>
                </>
              )}
            </div>
            
            {/* Model Selector */}
            <div className="flex items-center gap-2">
              <Bot className="h-4 w-4 text-muted-foreground" />
              <Select value={selectedModel} onValueChange={setSelectedModel}>
                <SelectTrigger className="w-48 h-8 text-xs">
                  <SelectValue placeholder="选择模型" />
                </SelectTrigger>
                <SelectContent>
                  {models.filter(model => model.id && model.id.trim() !== '').map((model) => (
                    <SelectItem key={model.id} value={model.id}>
                      {model.name}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>
          
          <div className="flex items-center gap-2">
            <span className="text-muted-foreground text-xs">Real-time Chat</span>
            <Link href="/admin">
              <Button variant="ghost" size="sm" className="h-6 px-2">
                <Settings className="h-3 w-3" />
              </Button>
            </Link>
          </div>
        </div>

        {error && (
          <div className="bg-destructive/10 text-destructive px-4 py-2 text-sm">
            {error}
          </div>
        )}
        
        {currentConversation ? (
          <>
            <div className="border-b px-6 py-4">
              <h1 className="text-xl font-semibold">{currentConversation.title}</h1>
              <p className="text-sm text-muted-foreground">
                Model: {currentConversation.model}
              </p>
            </div>
            
            <MessageList />
            
            <MessageInput
              onSendMessage={handleSendMessage}
              isLoading={isLoading}
              isStreaming={isStreaming}
              disabled={!currentConversation}
            />
          </>
        ) : (
          <div className="flex-1 flex flex-col">
            <div className="border-b px-6 py-4">
              <h1 className="text-xl font-semibold">开始新对话</h1>
              <p className="text-sm text-muted-foreground">
                选择模型并开始聊天，对话标题将由 AI 自动生成
              </p>
            </div>
            
            <div className="flex-1 flex items-center justify-center">
              <Card className="p-8 max-w-2xl w-full mx-4">
                <div className="text-center mb-6">
                  <h2 className="text-2xl font-bold mb-2">Qicro AI</h2>
                  <p className="text-muted-foreground">
                    我可以帮助您解答问题、创作内容、分析数据等
                  </p>
                </div>
                
                <div className="space-y-4">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <Card className="p-4 cursor-pointer hover:bg-muted/50">
                      <h3 className="font-semibold mb-2">💬 日常对话</h3>
                      <p className="text-sm text-muted-foreground">
                        聊天、问答、生活建议
                      </p>
                    </Card>
                    <Card className="p-4 cursor-pointer hover:bg-muted/50">
                      <h3 className="font-semibold mb-2">📝 内容创作</h3>
                      <p className="text-sm text-muted-foreground">
                        写作、翻译、创意设计
                      </p>
                    </Card>
                    <Card className="p-4 cursor-pointer hover:bg-muted/50">
                      <h3 className="font-semibold mb-2">🔍 数据分析</h3>
                      <p className="text-sm text-muted-foreground">
                        数据处理、图表分析
                      </p>
                    </Card>
                    <Card className="p-4 cursor-pointer hover:bg-muted/50">
                      <h3 className="font-semibold mb-2">💻 编程助手</h3>
                      <p className="text-sm text-muted-foreground">
                        代码编写、调试、优化
                      </p>
                    </Card>
                  </div>
                </div>
              </Card>
            </div>
            
            <MessageInput
              onSendMessage={handleSendMessage}
              isLoading={isLoading}
              isStreaming={isStreaming}
              disabled={false}
            />
          </div>
        )}
      </div>
    </div>
  );
}