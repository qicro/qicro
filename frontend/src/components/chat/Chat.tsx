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
import { Loader2, Wifi, WifiOff, Settings } from 'lucide-react';
import Link from 'next/link';

export default function Chat() {
  const { user } = useAuthStore();
  const { 
    currentConversation, 
    isLoading, 
    isStreaming, 
    error,
    sendMessage,
    clearError
  } = useChatStore();

  const { sendMessage: sendWebSocketMessage, isConnected } = useWebSocket({
    onMessage: (data) => {
      console.log('Received WebSocket message:', data);
      // Handle real-time updates here
      if (data.type === 'message_update') {
        // Update message in store
      }
    },
    onConnect: () => {
      console.log('Connected to WebSocket');
    },
    onDisconnect: () => {
      console.log('Disconnected from WebSocket');
    },
    onError: (error) => {
      console.error('WebSocket error:', error);
    },
  });

  useEffect(() => {
    if (error) {
      console.error('Chat error:', error);
      setTimeout(() => clearError(), 5000);
    }
  }, [error, clearError]);

  const handleSendMessage = async (content: string) => {
    if (!currentConversation) {
      console.error('No conversation selected');
      return;
    }
    
    // Send via WebSocket for real-time updates
    sendWebSocketMessage({
      type: 'new_message',
      conversation_id: currentConversation.id,
      content: content,
      timestamp: new Date().toISOString(),
    });
    
    // Also send via HTTP API for persistence and LLM response
    await sendMessage(content, true);
  };

  if (!user) {
    return (
      <div className="h-screen flex items-center justify-center">
        <div className="text-center">
          <Loader2 className="h-8 w-8 animate-spin mx-auto mb-4" />
          <p className="text-muted-foreground">Loading...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="h-screen flex">
      <ConversationSidebar />
      
      <div className="flex-1 flex flex-col">
        {/* Connection Status */}
        <div className="flex items-center justify-between bg-muted/30 px-4 py-1 text-xs">
          <div className="flex items-center gap-2">
            {isConnected ? (
              <>
                <Wifi className="h-3 w-3 text-green-500" />
                <span className="text-green-600">Connected</span>
              </>
            ) : (
              <>
                <WifiOff className="h-3 w-3 text-red-500" />
                <span className="text-red-600">Disconnected</span>
              </>
            )}
          </div>
          <div className="flex items-center gap-2">
            <span className="text-muted-foreground">Real-time Chat</span>
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
          <div className="flex-1 flex items-center justify-center">
            <Card className="p-8 max-w-md text-center">
              <h2 className="text-xl font-semibold mb-4">Welcome to Chat</h2>
              <p className="text-muted-foreground mb-4">
                Select an existing conversation or create a new one to get started.
              </p>
              <p className="text-sm text-muted-foreground">
                Choose a model and start chatting with AI assistants.
              </p>
            </Card>
          </div>
        )}
      </div>
    </div>
  );
}