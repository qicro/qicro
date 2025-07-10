import { create } from 'zustand';
import { ChatMessage, Conversation, Model, Provider } from '@/types/chat';
import { chatAPI } from '@/lib/chat-api';

interface ChatState {
  conversations: Conversation[];
  currentConversation: Conversation | null;
  messages: ChatMessage[];
  models: Model[];
  providers: Provider[];
  selectedModel: string;
  isLoading: boolean;
  isStreaming: boolean;
  error: string | null;
  
  // Actions
  loadConversations: () => Promise<void>;
  loadModels: () => Promise<void>;
  loadProviders: () => Promise<void>;
  setSelectedModel: (model: string) => void;
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
  selectedModel: '',
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
          selectedModel: defaultConversation.model, // Update selected model to match default conversation
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
      
      // Set default selected model if none is selected
      const { selectedModel } = get();
      if (!selectedModel && models.length > 0) {
        const validModels = models.filter(model => model.id && model.id.trim() !== '');
        if (validModels.length > 0) {
          set({ selectedModel: validModels[0].id });
        }
      }
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

  setSelectedModel: (model: string) => {
    set({ selectedModel: model });
    
    // Update current conversation model if there is one
    const { currentConversation } = get();
    if (currentConversation && currentConversation.model !== model) {
      // Don't await this to avoid blocking the UI
      get().updateConversation(currentConversation.id, { model }).catch(error => {
        console.error('Failed to update conversation model:', error);
      });
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
        selectedModel: conversation.model, // Update selected model to match new conversation
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
      // Use a simple approach - extract key words or use first few words
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
        selectedModel: conversation.model, // Update selected model to match conversation
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

    console.log('Sending message:', { content, currentConversation: currentConversation.id, model: currentConversation.model, stream });

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
        // Don't await title generation to avoid blocking the UI
        get().generateTitle(currentConversation.id, content).catch(error => {
          console.error('Failed to generate title:', error);
        });
      }
    } catch (error) {
      console.error('Error sending message:', error);
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

    console.log('Sending stream message:', { content, currentConversation: currentConversation.id, model: currentConversation.model });

    // Check if this is the first message and title is still "New Chat"
    const isFirstMessage = messages.length === 0 && currentConversation.title === 'New Chat';

    try {
      set({ isStreaming: true, error: null });
      
      const response = await chatAPI.sendMessageStream(
        currentConversation.id,
        { content, stream: true }
      );

      console.log('Stream response received:', response.status);

      if (!response.body) {
        throw new Error('No response body');
      }

      const reader = response.body.getReader();
      const decoder = new TextDecoder();
      let buffer = '';
      let currentEventType = '';

      console.log('Debug: Starting to read SSE stream...');

      while (true) {
        console.log('Debug: Reading next chunk...');
        const { done, value } = await reader.read();
        if (done) {
          console.log('Debug: Stream reading completed');
          break;
        }

        buffer += decoder.decode(value, { stream: true });
        console.log('Debug: Received chunk, buffer length:', buffer.length);
        const lines = buffer.split('\n');
        buffer = lines.pop() || '';

        console.log('Debug: Processing', lines.length, 'lines');
        for (const line of lines) {
          console.log('Debug: Processing line:', `"${line}"`);
          if (line.trim() === '') {
            // Empty line indicates end of event - but don't reset currentEventType yet
            continue;
          }
          
          if (line.startsWith('event:')) {
            currentEventType = line.slice(6).trim();
            console.log('Debug: SSE event type:', currentEventType);
          } else if (line.startsWith('data:')) {
            const data = line.slice(5).trim();
            console.log('Debug: SSE data line:', `"${data}"`);
            
            try {
              const eventData = JSON.parse(data);
              console.log('Debug: SSE data for event type:', currentEventType, eventData);
              
              if (currentEventType === 'user_message') {
                // 添加用户消息到列表
                const { messages: currentMessages } = get();
                set({ messages: [...currentMessages, eventData] });
                console.log('Debug: Added user message');
              } else if (currentEventType === 'assistant_message') {
                const { messages: currentMessages } = get();
                const lastMessage = currentMessages[currentMessages.length - 1];
                
                // 检查eventData的结构 - 应该是ChatResponse格式
                console.log('Debug: Assistant message data structure:', eventData);
                
                // 从ChatResponse.message.content获取内容
                const newContent = eventData.message?.content || '';
                console.log('Debug: Extracted content:', `"${newContent}"`);
                
                // 检查是否stream结束
                if (eventData.finish_reason === 'stop') {
                  console.log('Debug: Stream finished with stop reason');
                  set({ isStreaming: false });
                }
                
                // 只有当有内容时才处理，避免创建空消息
                if (newContent || eventData.finish_reason === 'stop') {
                  if (lastMessage && lastMessage.role === 'assistant') {
                    // 更新现有的助手消息 - 追加新内容实现打字效果
                    const updatedMessages = [...currentMessages];
                    updatedMessages[updatedMessages.length - 1] = {
                      ...lastMessage,
                      content: lastMessage.content + newContent,
                    };
                    set({ messages: updatedMessages });
                    console.log('Debug: Updated existing assistant message, total content length:', updatedMessages[updatedMessages.length - 1].content.length);
                  } else if (newContent) {
                    // 只有当有内容时才创建新的助手消息
                    const assistantMessage: ChatMessage = {
                      id: eventData.id || Date.now().toString(),
                      conversation_id: eventData.conversation_id || currentConversation.id,
                      role: 'assistant',
                      content: newContent,
                      created_at: new Date().toISOString(),
                    };
                    set({ messages: [...currentMessages, assistantMessage] });
                    console.log('Debug: Created new assistant message with content:', `"${newContent}"`);
                  }
                } else {
                  console.log('Debug: Skipping empty content message');
                }
              } else if (currentEventType === 'done') {
                console.log('Debug: Stream completed');
                set({ isStreaming: false, isLoading: false });
                break;
              } else if (currentEventType === 'error') {
                console.error('Debug: Stream error:', eventData);
                set({ error: eventData.error, isStreaming: false, isLoading: false });
                break;
              }
              
              // Reset event type after processing data
              currentEventType = '';
            } catch (parseError) {
              console.error('Error parsing SSE data:', parseError, 'Raw data:', data);
            }
          }
        }
      }

      // Ensure streaming is stopped when loop ends
      console.log('Debug: Stream reading loop ended, stopping streaming');
      set({ isStreaming: false, isLoading: false });

      // Generate title if this is the first message
      if (isFirstMessage) {
        // Don't await title generation to avoid blocking the UI
        get().generateTitle(currentConversation.id, content).catch(error => {
          console.error('Failed to generate title:', error);
        });
      }
    } catch (error) {
      console.error('Error in stream message:', error);
      set({ 
        error: error instanceof Error ? error.message : 'Failed to send stream message',
        isStreaming: false,
        isLoading: false 
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