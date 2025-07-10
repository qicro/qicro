'use client';

import { useState, useRef } from 'react';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
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
    <div className="border-t bg-background p-4">
      <Card className="p-4">
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="flex items-end gap-2">
            <div className="flex-1">
              <textarea
                ref={textareaRef}
                value={message}
                onChange={handleTextareaChange}
                onKeyDown={handleKeyDown}
                placeholder={isProcessing ? "AI 正在思考中..." : "输入消息..."}
                disabled={disabled || isProcessing}
                className="w-full min-h-[40px] max-h-[200px] p-3 border rounded-lg resize-none focus:outline-none focus:ring-2 focus:ring-primary bg-background"
                rows={1}
              />
            </div>
            
            <Button
              type="submit"
              size="sm"
              disabled={!message.trim() || disabled || isProcessing}
              className="h-10 w-10 p-0"
            >
              {isProcessing ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                <Send className="h-4 w-4" />
              )}
            </Button>
          </div>
        </form>
      </Card>
    </div>
  );
}