import { create } from 'zustand';
import { 
  APIKey, 
  AppType, 
  ChatModel, 
  CreateAPIKeyRequest, 
  UpdateAPIKeyRequest, 
  CreateAppTypeRequest, 
  CreateChatModelRequest 
} from '@/types/admin';
import { adminAPI } from '@/lib/admin-api';

interface AdminState {
  apiKeys: APIKey[];
  appTypes: AppType[];
  chatModels: ChatModel[];
  isLoading: boolean;
  error: string | null;

  // API Keys
  loadAPIKeys: () => Promise<void>;
  createAPIKey: (data: CreateAPIKeyRequest) => Promise<void>;
  updateAPIKey: (id: string, data: UpdateAPIKeyRequest) => Promise<void>;
  deleteAPIKey: (id: string) => Promise<void>;

  // App Types
  loadAppTypes: () => Promise<void>;
  createAppType: (data: CreateAppTypeRequest) => Promise<void>;

  // Chat Models
  loadChatModels: (type?: string) => Promise<void>;
  createChatModel: (data: CreateChatModelRequest) => Promise<void>;

  // Utility
  clearError: () => void;
  setLoading: (loading: boolean) => void;
}

export const useAdminStore = create<AdminState>((set, get) => ({
  apiKeys: [],
  appTypes: [],
  chatModels: [],
  isLoading: false,
  error: null,

  // API Keys
  loadAPIKeys: async () => {
    try {
      set({ isLoading: true, error: null });
      const { api_keys } = await adminAPI.getAPIKeys();
      set({ apiKeys: api_keys, isLoading: false });
    } catch (error) {
      set({ 
        error: error instanceof Error ? error.message : 'Failed to load API keys',
        isLoading: false 
      });
    }
  },

  createAPIKey: async (data: CreateAPIKeyRequest) => {
    try {
      set({ isLoading: true, error: null });
      const newAPIKey = await adminAPI.createAPIKey(data);
      const { apiKeys } = get();
      set({ 
        apiKeys: [newAPIKey, ...apiKeys],
        isLoading: false 
      });
    } catch (error) {
      set({ 
        error: error instanceof Error ? error.message : 'Failed to create API key',
        isLoading: false 
      });
      throw error;
    }
  },

  updateAPIKey: async (id: string, data: UpdateAPIKeyRequest) => {
    try {
      set({ isLoading: true, error: null });
      const updatedAPIKey = await adminAPI.updateAPIKey(id, data);
      const { apiKeys } = get();
      set({ 
        apiKeys: apiKeys.map(key => key.id === id ? updatedAPIKey : key),
        isLoading: false 
      });
    } catch (error) {
      set({ 
        error: error instanceof Error ? error.message : 'Failed to update API key',
        isLoading: false 
      });
      throw error;
    }
  },

  deleteAPIKey: async (id: string) => {
    try {
      set({ isLoading: true, error: null });
      await adminAPI.deleteAPIKey(id);
      const { apiKeys } = get();
      set({ 
        apiKeys: apiKeys.filter(key => key.id !== id),
        isLoading: false 
      });
    } catch (error) {
      set({ 
        error: error instanceof Error ? error.message : 'Failed to delete API key',
        isLoading: false 
      });
      throw error;
    }
  },

  // App Types
  loadAppTypes: async () => {
    try {
      set({ isLoading: true, error: null });
      const { app_types } = await adminAPI.getAppTypes();
      set({ appTypes: app_types, isLoading: false });
    } catch (error) {
      set({ 
        error: error instanceof Error ? error.message : 'Failed to load app types',
        isLoading: false 
      });
    }
  },

  createAppType: async (data: CreateAppTypeRequest) => {
    try {
      set({ isLoading: true, error: null });
      const newAppType = await adminAPI.createAppType(data);
      const { appTypes } = get();
      set({ 
        appTypes: [...appTypes, newAppType],
        isLoading: false 
      });
    } catch (error) {
      set({ 
        error: error instanceof Error ? error.message : 'Failed to create app type',
        isLoading: false 
      });
      throw error;
    }
  },

  // Chat Models
  loadChatModels: async (type?: string) => {
    try {
      set({ isLoading: true, error: null });
      const { models } = await adminAPI.getChatModels(type);
      set({ chatModels: models, isLoading: false });
    } catch (error) {
      set({ 
        error: error instanceof Error ? error.message : 'Failed to load chat models',
        isLoading: false 
      });
    }
  },

  createChatModel: async (data: CreateChatModelRequest) => {
    try {
      set({ isLoading: true, error: null });
      const newModel = await adminAPI.createChatModel(data);
      const { chatModels } = get();
      set({ 
        chatModels: [...chatModels, newModel],
        isLoading: false 
      });
    } catch (error) {
      set({ 
        error: error instanceof Error ? error.message : 'Failed to create chat model',
        isLoading: false 
      });
      throw error;
    }
  },

  clearError: () => set({ error: null }),
  setLoading: (loading: boolean) => set({ isLoading: loading }),
}));