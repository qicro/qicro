import { create } from 'zustand';
import { ChatMessage, Conversation, Model, Provider } from '@/types/chat';
import { chatAPI } from '@/lib/chat-api';

interface ChatState {
  conversations: Conversation[];
  currentConversation: Conversation | null;
  messages: ChatMessage[];
  models: Model[];
  providers: Provider[];
  isLoading: boolean;
  isStreaming: boolean;
  error: string | null;
  
  // Actions
  loadConversations: () => Promise<void>;
  loadModels: () => Promise<void>;
  loadProviders: () => Promise<void>;
  createConversation: (title: string, model: string) => Promise<Conversation>;
  generateTitle: (conversationId: string, firstMessage: string) => Promise<void>;
  selectConversation: (id: string) => Promise<void>;
  sendMessage: (content: string, stream?: boolean) => Promise<void>;
  sendMessageStream: (content: string) => Promise<void>;
  updateConversation: (id: string, updates: Partial<Conversation>) => Promise<void>;
  deleteConversation: (id: string) => Promise<void>;
  clearError: () => void;
  setLoading: (loading: boolean) => void;
}

export const useChatStore = create<ChatState>((set, get) => ({
  conversations: [],
  currentConversation: null,
  messages: [],
  models: [],
  providers: [],
  isLoading: false,
  isStreaming: false,
  error: null,

  loadConversations: async () => {
    try {
      set({ isLoading: true, error: null });
      const { conversations } = await chatAPI.getConversations();
      
      // If no conversations exist, create a default one
      if (conversations.length === 0) {
        const { models } = await chatAPI.getModels();
        const defaultModel = models.length > 0 ? models[0].id : 'gpt-3.5-turbo';
        const defaultConversation = await chatAPI.createConversation('New Chat', defaultModel);
        set({ 
          conversations: [defaultConversation], 
          currentConversation: defaultConversation,
          messages: [],
          isLoading: false 
        });
      } else {
        set({ conversations, isLoading: false });
      }
    } catch (error) {
      set({ 
        error: error instanceof Error ? error.message : 'Failed to load conversations',
        isLoading: false 
      });
    }
  },

  loadModels: async () => {
    try {
      const { models } = await chatAPI.getModels();
      set({ models });
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Failed to load models' });
    }
  },

  loadProviders: async () => {
    try {
      const { providers } = await chatAPI.getProviders();
      set({ providers });
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Failed to load providers' });
    }
  },

  createConversation: async (title: string, model: string) => {
    try {
      set({ isLoading: true, error: null });
      const conversation = await chatAPI.createConversation(title, model);
      const { conversations } = get();
      set({ 
        conversations: [conversation, ...conversations],
        currentConversation: conversation,
        messages: [],
        isLoading: false 
      });
      return conversation;
    } catch (error) {
      set({ 
        error: error instanceof Error ? error.message : 'Failed to create conversation',
        isLoading: false 
      });
      throw error;
    }
  },

  generateTitle: async (conversationId: string, firstMessage: string) => {
    try {
      // Generate a title based on the first message
      const titlePrompt = `Generate a concise, descriptive title (max 5 words) for a conversation that starts with: "${firstMessage.substring(0, 100)}..."`;
      
      // Use a simple approach for now - extract key words or use first few words
      const words = firstMessage.split(' ').slice(0, 4).join(' ');
      const title = words.length > 30 ? words.substring(0, 30) + '...' : words;
      
      await get().updateConversation(conversationId, { title });
    } catch (error) {
      console.error('Failed to generate title:', error);
      // If title generation fails, use a default
      await get().updateConversation(conversationId, { title: 'New Chat' });
    }
  },

  selectConversation: async (id: string) => {
    try {
      set({ isLoading: true, error: null });
      const { conversation, messages } = await chatAPI.getConversation(id);
      set({ 
        currentConversation: conversation,
        messages,
        isLoading: false 
      });
    } catch (error) {
      set({ 
        error: error instanceof Error ? error.message : 'Failed to load conversation',
        isLoading: false 
      });
    }
  },

  sendMessage: async (content: string, stream = false) => {
    const { currentConversation, messages } = get();
    if (!currentConversation) {
      set({ error: 'No conversation selected' });
      return;
    }

    // Check if this is the first message and title is still "New Chat"
    const isFirstMessage = messages.length === 0 && currentConversation.title === 'New Chat';

    try {
      set({ isLoading: true, error: null });
      
      if (stream) {
        await get().sendMessageStream(content);
      } else {
        const { user_message, assistant_message } = await chatAPI.sendMessage(
          currentConversation.id,
          { content }
        );
        
        const { messages: currentMessages } = get();
        set({
          messages: [...currentMessages, user_message, assistant_message],
          isLoading: false
        });
      }

      // Generate title if this is the first message
      if (isFirstMessage) {
        await get().generateTitle(currentConversation.id, content);
      }
    } catch (error) {
      set({ 
        error: error instanceof Error ? error.message : 'Failed to send message',
        isLoading: false 
      });
    }
  },

  sendMessageStream: async (content: string) => {
    const { currentConversation, messages } = get();
    if (!currentConversation) {
      set({ error: 'No conversation selected' });
      return;
    }

    // Check if this is the first message and title is still "New Chat"
    const isFirstMessage = messages.length === 0 && currentConversation.title === 'New Chat';

    try {
      set({ isStreaming: true, error: null });
      
      const eventSource = await chatAPI.sendMessageStream(
        currentConversation.id,
        { content, stream: true }
      );

      let assistantMessageContent = '';

      eventSource.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          
          if (event.lastEventId === 'user_message') {
            set({ messages: [...messages, data] });
          } else if (event.lastEventId === 'assistant_message') {
            assistantMessageContent += data.message.content;
            
            // 更新流式消息
            const { messages: currentMessages } = get();
            const lastMessage = currentMessages[currentMessages.length - 1];
            
            if (lastMessage && lastMessage.role === 'assistant') {
              // 更新现有消息
              const updatedMessages = [...currentMessages];
              updatedMessages[updatedMessages.length - 1] = {
                ...lastMessage,
                content: assistantMessageContent,
              };
              set({ messages: updatedMessages });
            } else {
              // 创建新的助手消息
              const assistantMessage: ChatMessage = {
                id: data.id,
                conversation_id: currentConversation.id,
                role: 'assistant',
                content: assistantMessageContent,
                created_at: new Date().toISOString(),
              };
              set({ messages: [...currentMessages, assistantMessage] });
            }
          } else if (event.lastEventId === 'done') {
            eventSource.close();
            set({ isStreaming: false });
            
            // Generate title if this is the first message
            if (isFirstMessage) {
              get().generateTitle(currentConversation.id, content);
            }
          }
        } catch (error) {
          console.error('Error parsing stream data:', error);
        }
      };

      eventSource.onerror = (error) => {
        console.error('EventSource error:', error);
        eventSource.close();
        set({ 
          error: 'Connection error during streaming',
          isStreaming: false 
        });
      };

    } catch (error) {
      set({ 
        error: error instanceof Error ? error.message : 'Failed to send stream message',
        isStreaming: false 
      });
    }
  },

  updateConversation: async (id: string, updates: Partial<Conversation>) => {
    try {
      const updatedConversation = await chatAPI.updateConversation(id, updates);
      const { conversations, currentConversation } = get();
      
      set({
        conversations: conversations.map(conv => 
          conv.id === id ? updatedConversation : conv
        ),
        currentConversation: currentConversation?.id === id ? updatedConversation : currentConversation,
      });
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Failed to update conversation' });
    }
  },

  deleteConversation: async (id: string) => {
    try {
      await chatAPI.deleteConversation(id);
      const { conversations, currentConversation } = get();
      
      set({
        conversations: conversations.filter(conv => conv.id !== id),
        currentConversation: currentConversation?.id === id ? null : currentConversation,
        messages: currentConversation?.id === id ? [] : get().messages,
      });
    } catch (error) {
      set({ error: error instanceof Error ? error.message : 'Failed to delete conversation' });
    }
  },

  clearError: () => set({ error: null }),
  setLoading: (loading: boolean) => set({ isLoading: loading }),
}));