import axios, { type AxiosInstance, type AxiosError } from 'axios';
import type {
  LoginRequest,
  LoginResponse,
  BotStatus,
  Plugin,
  Task,
  ConfigItem,
  LogEntry,
  Preset,
  AITool,
  AIStats,
  ChatRequest,
  ChatResponse,
  ImageChatRequest,
  UpdateConfigRequest,
  UpdateConfigResponse,
  UpdatePresetRequest,
  UpdatePresetResponse,
  TriggerTaskResponse,
  TogglePluginResponse,
  ErrorResponse,
} from '@/types/api';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

class ApiClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    this.client.interceptors.request.use((config) => {
      const token = localStorage.getItem('token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });

    this.client.interceptors.response.use(
      (response) => response,
      (error: AxiosError<ErrorResponse>) => {
        if (error.response?.status === 401) {
          localStorage.removeItem('token');
          window.location.href = '/login';
        }
        return Promise.reject(error);
      }
    );
  }

  async login(data: LoginRequest): Promise<LoginResponse> {
    const response = await this.client.post<LoginResponse>('/api/auth/login', data);
    return response.data;
  }

  async getStatus(): Promise<BotStatus> {
    const response = await this.client.get<BotStatus>('/api/status');
    return response.data;
  }

  async getPlugins(): Promise<Plugin[]> {
    const response = await this.client.get<Plugin[]>('/api/plugins');
    return response.data;
  }

  async togglePlugin(name: string): Promise<TogglePluginResponse> {
    const response = await this.client.post<TogglePluginResponse>(`/api/plugins/${name}/toggle`);
    return response.data;
  }

  async getTasks(): Promise<Task[]> {
    const response = await this.client.get<Task[]>('/api/tasks');
    return response.data;
  }

  async triggerTask(id: string): Promise<TriggerTaskResponse> {
    const response = await this.client.post<TriggerTaskResponse>(`/api/tasks/${id}/trigger`);
    return response.data;
  }

  async getConfig(): Promise<ConfigItem[]> {
    const response = await this.client.get<ConfigItem[]>('/api/config');
    return response.data;
  }

  async updateConfig(items: UpdateConfigRequest[]): Promise<UpdateConfigResponse> {
    const response = await this.client.post<UpdateConfigResponse>('/api/config', items);
    return response.data;
  }

  async getLogs(level?: string, limit?: number): Promise<LogEntry[]> {
    const params: Record<string, string | number> = {};
    if (level) params.level = level;
    if (limit) params.limit = limit;

    const response = await this.client.get<LogEntry[]>('/api/logs', { params });
    return response.data;
  }

  async getPresets(): Promise<Preset[]> {
    const response = await this.client.get<Preset[]>('/api/presets');
    return response.data;
  }

  async getPreset(filename: string): Promise<Preset> {
    const response = await this.client.get<Preset>(`/api/presets/${filename}`);
    return response.data;
  }

  async updatePreset(filename: string, data: UpdatePresetRequest): Promise<UpdatePresetResponse> {
    const response = await this.client.put<UpdatePresetResponse>(`/api/presets/${filename}`, data);
    return response.data;
  }

  async getAITools(): Promise<AITool[]> {
    const response = await this.client.get<AITool[]>('/api/ai/tools');
    return response.data;
  }

  async getAIStats(): Promise<AIStats> {
    const response = await this.client.get<AIStats>('/api/ai/stats');
    return response.data;
  }

  async testAIChat(data: ChatRequest): Promise<ChatResponse> {
    const response = await this.client.post<ChatResponse>('/api/ai/chat', data);
    return response.data;
  }

  async testAIImageChat(data: ImageChatRequest): Promise<ChatResponse> {
    const response = await this.client.post<ChatResponse>('/api/ai/chat/image', data);
    return response.data;
  }
}

export const apiClient = new ApiClient();
