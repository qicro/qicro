'use client';

import { ChatMessage } from '@/types/chat';
import { Card, CardContent } from '@/components/ui/card';
import { User, Bot } from 'lucide-react';
import { cn } from '@/lib/utils';

interface MessageBubbleProps {
  message: ChatMessage;
}

export default function MessageBubble({ message }: MessageBubbleProps) {
  const isUser = message.role === 'user';
  
  return (
    <div className={cn(
      'flex gap-3 p-4',
      isUser ? 'justify-end' : 'justify-start'
    )}>
      {!isUser && (
        <div className="flex-shrink-0">
          <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center">
            <Bot className="h-4 w-4 text-primary" />
          </div>
        </div>
      )}
      
      <div className={cn(
        'max-w-[70%] space-y-2',
        isUser ? 'items-end' : 'items-start'
      )}>
        <Card className={cn(
          'relative',
          isUser ? 'bg-primary text-primary-foreground' : 'bg-muted'
        )}>
          <CardContent className="p-3">
            <div className="whitespace-pre-wrap text-sm">
              {message.content}
            </div>
          </CardContent>
        </Card>
        
        <div className={cn(
          'text-xs text-muted-foreground',
          isUser ? 'text-right' : 'text-left'
        )}>
          {new Date(message.created_at).toLocaleTimeString()}
        </div>
      </div>
      
      {isUser && (
        <div className="flex-shrink-0">
          <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center">
            <User className="h-4 w-4 text-primary" />
          </div>
        </div>
      )}
    </div>
  );
}