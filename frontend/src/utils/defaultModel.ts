import { Events } from '@wailsio/runtime';

export const DEFAULT_MODEL_KEY = 'chat_default_model';

export interface DefaultModelConfig {
  modelId: number;
  modelName: string;
}

export function getDefaultModelConfig(): DefaultModelConfig | null {
  try {
    const raw = localStorage.getItem(DEFAULT_MODEL_KEY);
    if (!raw) return null;
    return JSON.parse(raw) as DefaultModelConfig;
  } catch {
    return null;
  }
}

export function setDefaultModelConfig(config: DefaultModelConfig) {
  localStorage.setItem(DEFAULT_MODEL_KEY, JSON.stringify(config));
  window.dispatchEvent(new CustomEvent('chat-default-model-changed', { detail: config }));
  void Events.Emit('chat-default-model-changed', config);
}

export function clearDefaultModelConfig() {
  localStorage.removeItem(DEFAULT_MODEL_KEY);
  window.dispatchEvent(new CustomEvent('chat-default-model-changed', { detail: null }));
  void Events.Emit('chat-default-model-changed', null);
}
