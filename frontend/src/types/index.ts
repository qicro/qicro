export interface User {
  id: string;
  email: string;
  oauthProvider?: string;
  oauthId?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Conversation {
  id: string;
  userId: string;
  title: string;
  settings?: Record<string, any>;
  createdAt: string;
  updatedAt: string;
}

export interface Message {
  id: string;
  conversationId: string;
  role: 'user' | 'assistant' | 'system';
  content: MessageContent;
  artifacts?: Artifact[];
  createdAt: string;
}

export interface MessageContent {
  text?: string;
  images?: string[];
  audio?: string;
  video?: string;
}

export interface Artifact {
  id: string;
  type: 'html' | 'mermaid' | 'svg';
  content: string;
  title?: string;
}

export interface LLMProvider {
  id: string;
  name: string;
  models: LLMModel[];
  enabled: boolean;
}

export interface LLMModel {
  id: string;
  name: string;
  provider: string;
  capabilities: ModelCapabilities;
  maxTokens: number;
}

export interface ModelCapabilities {
  text: boolean;
  images: boolean;
  audio: boolean;
  video: boolean;
  tools: boolean;
}

export interface ChatRequest {
  conversationId?: string;
  message: MessageContent;
  model: string;
  stream?: boolean;
  tools?: Tool[];
}

export interface ChatResponse {
  id: string;
  content: MessageContent;
  artifacts?: Artifact[];
  usage?: TokenUsage;
}

export interface TokenUsage {
  promptTokens: number;
  completionTokens: number;
  totalTokens: number;
}

export interface Tool {
  id: string;
  name: string;
  description: string;
  schema: any;
  enabled: boolean;
}

export interface KnowledgeBase {
  id: string;
  userId: string;
  name: string;
  type: 'internal' | 'external';
  config: Record<string, any>;
  createdAt: string;
  updatedAt: string;
}