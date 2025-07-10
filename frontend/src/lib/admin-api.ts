import { 
  APIKey, 
  AppType, 
  ChatModel, 
  CreateAPIKeyRequest, 
  UpdateAPIKeyRequest, 
  CreateAppTypeRequest, 
  CreateChatModelRequest,
  UpdateChatModelRequest
} from '@/types/admin';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

class AdminAPI {
  private getHeaders() {
    const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
    return {
      'Content-Type': 'application/json',
      ...(token && { Authorization: `Bearer ${token}` }),
    };
  }

  // API Keys
  async createAPIKey(data: CreateAPIKeyRequest): Promise<APIKey> {
    const response = await fetch(`${API_BASE}/api/admin/api-keys`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to create API key');
    }

    return response.json();
  }

  async getAPIKeys(): Promise<{ api_keys: APIKey[] }> {
    const response = await fetch(`${API_BASE}/api/admin/api-keys`, {
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to get API keys');
    }

    return response.json();
  }

  async getAPIKey(id: string): Promise<APIKey> {
    const response = await fetch(`${API_BASE}/api/admin/api-keys/${id}`, {
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to get API key');
    }

    return response.json();
  }

  async updateAPIKey(id: string, data: UpdateAPIKeyRequest): Promise<APIKey> {
    const response = await fetch(`${API_BASE}/api/admin/api-keys/${id}`, {
      method: 'PUT',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to update API key');
    }

    return response.json();
  }

  async deleteAPIKey(id: string): Promise<void> {
    const response = await fetch(`${API_BASE}/api/admin/api-keys/${id}`, {
      method: 'DELETE',
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to delete API key');
    }
  }

  // App Types
  async createAppType(data: CreateAppTypeRequest): Promise<AppType> {
    const response = await fetch(`${API_BASE}/api/admin/app-types`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to create app type');
    }

    return response.json();
  }

  async getAppTypes(): Promise<{ app_types: AppType[] }> {
    const response = await fetch(`${API_BASE}/api/admin/app-types`, {
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to get app types');
    }

    return response.json();
  }

  // Chat Models
  async createChatModel(data: CreateChatModelRequest): Promise<ChatModel> {
    const response = await fetch(`${API_BASE}/api/admin/chat-models`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to create chat model');
    }

    return response.json();
  }

  async getChatModels(type?: string): Promise<{ models: ChatModel[] }> {
    const url = new URL(`${API_BASE}/api/admin/chat-models`);
    if (type) {
      url.searchParams.set('type', type);
    }

    const response = await fetch(url.toString(), {
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to get chat models');
    }

    return response.json();
  }

  async getChatModel(id: string): Promise<ChatModel> {
    const response = await fetch(`${API_BASE}/api/admin/chat-models/${id}`, {
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to get chat model');
    }

    return response.json();
  }

  async updateChatModel(id: string, data: UpdateChatModelRequest): Promise<ChatModel> {
    const response = await fetch(`${API_BASE}/api/admin/chat-models/${id}`, {
      method: 'PUT',
      headers: this.getHeaders(),
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to update chat model');
    }

    return response.json();
  }

  async deleteChatModel(id: string): Promise<void> {
    const response = await fetch(`${API_BASE}/api/admin/chat-models/${id}`, {
      method: 'DELETE',
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to delete chat model');
    }
  }
}

export const adminAPI = new AdminAPI();