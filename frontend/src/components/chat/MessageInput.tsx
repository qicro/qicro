'use client';

import { useState, useRef } from 'react';
import { Button } from '@/components/ui/button';
import { Send, Loader2 } from 'lucide-react';

interface MessageInputProps {
  onSendMessage: (content: string) => void;
  isLoading?: boolean;
  isStreaming?: boolean;
  disabled?: boolean;
}

export default function MessageInput({
  onSendMessage,
  isLoading = false,
  isStreaming = false,
  disabled = false
}: MessageInputProps) {
  const [message, setMessage] = useState('');
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (message.trim() && !disabled && !isLoading && !isStreaming) {
      onSendMessage(message.trim());
      setMessage('');
      if (textareaRef.current) {
        textareaRef.current.style.height = '40px';
      }
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSubmit(e);
    }
  };

  const handleTextareaChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setMessage(e.target.value);
    
    // Auto-resize textarea
    if (textareaRef.current) {
      textareaRef.current.style.height = '40px';
      textareaRef.current.style.height = Math.min(textareaRef.current.scrollHeight, 200) + 'px';
    }
  };

  const isProcessing = isLoading || isStreaming;

  return (
    <div className="border-t bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 p-4">
      <div className="max-w-4xl mx-auto">
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="relative flex items-end gap-3 bg-muted/50 rounded-2xl p-3 border border-border/50 shadow-sm hover:shadow-md transition-all duration-200 focus-within:ring-2 focus-within:ring-primary/20 focus-within:border-primary/30">
            <div className="flex-1 min-h-0">
              <textarea
                ref={textareaRef}
                value={message}
                onChange={handleTextareaChange}
                onKeyDown={handleKeyDown}
                placeholder={isProcessing ? "AI 正在思考中..." : "输入消息... (按 Enter 发送，Shift+Enter 换行)"}
                disabled={disabled || isProcessing}
                className="w-full min-h-[44px] max-h-[200px] bg-transparent resize-none placeholder:text-muted-foreground/70 border-0 focus:outline-none focus:ring-0 text-sm leading-6 py-2 px-0"
                rows={1}
              />
            </div>
            
            <Button
              type="submit"
              size="sm"
              disabled={!message.trim() || disabled || isProcessing}
              className="h-10 w-10 p-0 rounded-xl shrink-0 transition-all duration-200 hover:scale-105 disabled:hover:scale-100"
            >
              {isProcessing ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                <Send className="h-4 w-4" />
              )}
            </Button>
          </div>
          
          {!disabled && !isProcessing && (
            <div className="text-xs text-muted-foreground text-center">
              按 Enter 发送消息，Shift+Enter 换行
            </div>
          )}
        </form>
      </div>
    </div>
  );
}