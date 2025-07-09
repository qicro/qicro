export interface APIKey {
  id: string;
  name: string;
  value: string;
  type: string;
  provider: string;
  api_url?: string;
  proxy_url?: string;
  last_used_at?: string;
  enabled: boolean;
  created_at: string;
  updated_at: string;
}

export interface AppType {
  id: string;
  name: string;
  icon?: string;
  sort_num: number;
  enabled: boolean;
  created_at: string;
  updated_at: string;
}

export interface ChatModel {
  id: string;
  type: string;
  name: string;
  value: string;
  provider: string;
  sort_num: number;
  enabled: boolean;
  power: number;
  temperature: number;
  max_tokens: number;
  max_context: number;
  open: boolean;
  api_key_id?: string;
  created_at: string;
  updated_at: string;
}

export interface CreateAPIKeyRequest {
  name: string;
  value: string;
  type: string;
  provider: string;
  api_url?: string;
  proxy_url?: string;
  enabled?: boolean;
}

export interface UpdateAPIKeyRequest {
  name?: string;
  value?: string;
  type?: string;
  provider?: string;
  api_url?: string;
  proxy_url?: string;
  enabled?: boolean;
}

export interface CreateAppTypeRequest {
  name: string;
  icon?: string;
  sort_num?: number;
  enabled?: boolean;
}

export interface CreateChatModelRequest {
  type: string;
  name: string;
  value: string;
  provider: string;
  sort_num?: number;
  enabled?: boolean;
  power?: number;
  temperature?: number;
  max_tokens?: number;
  max_context?: number;
  open?: boolean;
  api_key_id?: string;
}