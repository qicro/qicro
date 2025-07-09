import { ChatMessage, Conversation, ChatRequest, Model, Provider } from '@/types/chat';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

class ChatAPI {
  private getHeaders() {
    const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
    return {
      'Content-Type': 'application/json',
      ...(token && { Authorization: `Bearer ${token}` }),
    };
  }

  // 模型相关
  async getModels(): Promise<{ models: Model[] }> {
    const response = await fetch(`${API_BASE}/api/models`, {
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to get models');
    }

    return response.json();
  }

  async getProviders(): Promise<{ providers: Provider[] }> {
    const response = await fetch(`${API_BASE}/api/providers`, {
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to get providers');
    }

    return response.json();
  }

  // 对话相关
  async createConversation(title: string, model: string): Promise<Conversation> {
    const response = await fetch(`${API_BASE}/api/conversations`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify({ title, model }),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to create conversation');
    }

    return response.json();
  }

  async getConversations(): Promise<{ conversations: Conversation[] }> {
    const response = await fetch(`${API_BASE}/api/conversations`, {
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to get conversations');
    }

    return response.json();
  }

  async getConversation(id: string): Promise<{ conversation: Conversation; messages: ChatMessage[] }> {
    const response = await fetch(`${API_BASE}/api/conversations/${id}`, {
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to get conversation');
    }

    return response.json();
  }

  async updateConversation(id: string, updates: Partial<Conversation>): Promise<Conversation> {
    const response = await fetch(`${API_BASE}/api/conversations/${id}`, {
      method: 'PUT',
      headers: this.getHeaders(),
      body: JSON.stringify(updates),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to update conversation');
    }

    return response.json();
  }

  async deleteConversation(id: string): Promise<void> {
    const response = await fetch(`${API_BASE}/api/conversations/${id}`, {
      method: 'DELETE',
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to delete conversation');
    }
  }

  // 消息相关
  async sendMessage(conversationId: string, request: ChatRequest): Promise<{ user_message: ChatMessage; assistant_message: ChatMessage }> {
    const response = await fetch(`${API_BASE}/api/conversations/${conversationId}/messages`, {
      method: 'POST',
      headers: this.getHeaders(),
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to send message');
    }

    return response.json();
  }

  async sendMessageStream(conversationId: string, request: ChatRequest): Promise<EventSource> {
    const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
    const url = new URL(`${API_BASE}/api/conversations/${conversationId}/messages`);
    
    const eventSource = new EventSource(url.toString());
    
    // 发送消息
    fetch(url.toString(), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...(token && { Authorization: `Bearer ${token}` }),
      },
      body: JSON.stringify({ ...request, stream: true }),
    }).catch(error => {
      console.error('Failed to send stream message:', error);
    });

    return eventSource;
  }

  async getMessages(conversationId: string): Promise<{ messages: ChatMessage[] }> {
    const response = await fetch(`${API_BASE}/api/conversations/${conversationId}/messages`, {
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to get messages');
    }

    return response.json();
  }
}

export const chatAPI = new ChatAPI();