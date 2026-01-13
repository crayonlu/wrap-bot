// API类型定义

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
}

export interface BotStatus {
  running: boolean;
  uptime: number;
  version: string;
  go_version: string;
}

export interface Plugin {
  name: string;
  enabled: boolean;
  description: string;
}

export interface Task {
  id: string;
  name: string;
  schedule: string;
  next_run: string;
  last_run: string;
  status: string;
  can_trigger: boolean;
  description: string;
}

export interface ConfigItem {
  key: string;
  value: string;
  description: string;
}

export interface LogEntry {
  timestamp: string;
  level: string;
  message: string;
  context?: Record<string, any>;
}

export interface Preset {
  name: string;
  path: string;
  content: string;
}

export interface AITool {
  name: string;
  description: string;
  category: string;
  text_enabled: boolean;
  vision_enabled: boolean;
}

export interface AIStats {
  total_calls: number;
  tool_usage: Record<string, number>;
  success_rate: number;
  recent_calls: AICall[];
}

export interface AICall {
  timestamp: string;
  model: string;
  tools_used: string[];
  success: boolean;
  duration_ms: number;
}

export interface ChatRequest {
  message: string;
  conversation_id?: string;
  model?: string;
}

export interface ChatResponse {
  response: string;
  conversation_id: string;
  tool_calls: ToolCall[];
}

export interface ToolCall {
  name: string;
  arguments: string;
}

export interface ImageChatRequest {
  message: string;
  images: string[];
  conversation_id?: string;
  model?: string;
}

export interface UpdateConfigRequest {
  key: string;
  value: string;
}

export interface UpdateConfigResponse {
  status: string;
  updated_count: number;
  updated_keys: Record<string, boolean>;
}

export interface UpdatePresetRequest {
  content: string;
}

export interface UpdatePresetResponse {
  message: string;
}

export interface TriggerTaskResponse {
  status: string;
  task_id: string;
}

export interface TogglePluginResponse {
  name: string;
  enabled: boolean;
}

export interface ErrorResponse {
  error: string;
}

// WebSocket事件类型
export type WebSocketEventType = 'status' | 'plugins' | 'tasks' | 'log';

export interface WebSocketEvent<T = any> {
  type: WebSocketEventType;
  data: T;
}
