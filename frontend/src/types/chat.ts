export interface ChatMessage {
  id: string;
  conversation_id: string;
  role: 'user' | 'assistant' | 'system';
  content: string;
  metadata?: Record<string, any>;
  created_at: string;
}

export interface Conversation {
  id: string;
  user_id: string;
  title: string;
  model: string;
  settings?: Record<string, any>;
  created_at: string;
  updated_at: string;
}

export interface ChatRequest {
  content: string;
  stream?: boolean;
}

export interface ChatResponse {
  id: string;
  conversation_id: string;
  message: ChatMessage;
  usage?: {
    prompt_tokens: number;
    completion_tokens: number;
    total_tokens: number;
  };
  finish_reason?: string;
  metadata?: Record<string, any>;
}

export interface Model {
  id: string;
  name: string;
  provider: string;
  capabilities: string[];
  max_tokens: number;
}

export interface Provider {
  name: string;
  models: Model[];
}