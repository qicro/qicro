'use client';

import { ChatMessage } from '@/types/chat';
import { Card, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { User, Bot, Copy, Check } from 'lucide-react';
import { cn, formatMessageTime } from '@/lib/utils';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import rehypeHighlight from 'rehype-highlight';
import rehypeRaw from 'rehype-raw';
import { useState } from 'react';

interface MessageBubbleProps {
  message: ChatMessage;
}

interface CodeBlockProps {
  children: React.ReactNode;
  className?: string;
  language?: string;
}

function CodeBlock({ children, className, language }: CodeBlockProps) {
  const [copied, setCopied] = useState(false);
  
  const copyToClipboard = async () => {
    try {
      await navigator.clipboard.writeText(String(children));
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error('Failed to copy:', err);
    }
  };

  return (
    <div className="relative group">
      <pre className="bg-muted p-3 rounded-md overflow-x-auto pr-12">
        <code className={className}>
          {children}
        </code>
      </pre>
      <Button
        variant="ghost"
        size="sm"
        onClick={copyToClipboard}
        className="absolute top-2 right-2 h-8 w-8 p-0 opacity-0 group-hover:opacity-100 transition-opacity"
      >
        {copied ? (
          <Check className="h-4 w-4 text-green-500" />
        ) : (
          <Copy className="h-4 w-4" />
        )}
      </Button>
    </div>
  );
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
          'relative rounded-2xl shadow-sm',
          isUser ? 'bg-primary text-primary-foreground' : 'bg-muted/70 backdrop-blur-sm'
        )}>
          <CardContent className="p-4">
            {isUser ? (
              <div className="whitespace-pre-wrap text-sm">
                {message.content}
              </div>
            ) : (
              <div className="prose prose-sm dark:prose-invert max-w-none">
                <ReactMarkdown
                  remarkPlugins={[remarkGfm]}
                  rehypePlugins={[rehypeHighlight, rehypeRaw]}
                  components={{
                    code: ({ className, children, ...props }: { className?: string; children?: React.ReactNode; [key: string]: any }) => {
                      const match = /language-(\w+)/.exec(className || '');
                      const language = match ? match[1] : undefined;
                      const isInline = !props['data-inline'] && !match;
                      
                      return !isInline && match ? (
                        <CodeBlock className={className} language={language}>
                          {children}
                        </CodeBlock>
                      ) : (
                        <code className="bg-muted px-1 py-0.5 rounded text-sm" {...props}>
                          {children}
                        </code>
                      );
                    },
                    blockquote: ({ children }) => (
                      <blockquote className="border-l-4 border-primary pl-4 italic">
                        {children}
                      </blockquote>
                    ),
                    table: ({ children }) => (
                      <div className="overflow-x-auto">
                        <table className="min-w-full border-collapse border border-border">
                          {children}
                        </table>
                      </div>
                    ),
                    th: ({ children }) => (
                      <th className="border border-border px-3 py-2 bg-muted font-semibold text-left">
                        {children}
                      </th>
                    ),
                    td: ({ children }) => (
                      <td className="border border-border px-3 py-2">
                        {children}
                      </td>
                    ),
                  }}
                >
                  {message.content}
                </ReactMarkdown>
              </div>
            )}
          </CardContent>
        </Card>
        
        <div className={cn(
          'text-xs text-muted-foreground',
          isUser ? 'text-right' : 'text-left'
        )}>
          {formatMessageTime(message.created_at)}
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