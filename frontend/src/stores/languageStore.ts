import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import i18n from '@/i18n';
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service';
import { AppPreferences } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models';
import { DEFAULT_LANGUAGE, DEFAULT_REGION } from '@/i18n/types';
import type { AppLanguage, AppRegion } from '@/i18n/types';

interface LanguageState {
  language: AppLanguage;
  region: AppRegion;
  isHydrated: boolean;
  setLanguage: (language: AppLanguage) => Promise<void>;
  setRegion: (region: AppRegion) => Promise<void>;
  setPreferences: (preferences: { language: AppLanguage; region: AppRegion }) => void;
  markHydrated: () => void;
}

function applyLanguage(language: AppLanguage) {
  document.documentElement.lang = language;
  void i18n.changeLanguage(language);
}

async function persistPreferences(preferences: { language: AppLanguage; region: AppRegion }) {
  const current = await Service.GetAppPreferences();
  await Service.UpdateAppPreferences(new AppPreferences({
    ...current,
    language: preferences.language,
    region: preferences.region,
  }));
}

export const useLanguageStore = create<LanguageState>()(
  persist(
    (set) => ({
      language: DEFAULT_LANGUAGE,
      region: DEFAULT_REGION,
      isHydrated: false,
      setLanguage: async (language: AppLanguage) => {
        const current = useLanguageStore.getState();
        set({ language });
        applyLanguage(language);
        await persistPreferences({ language, region: current.region });
      },
      setRegion: async (region: AppRegion) => {
        const current = useLanguageStore.getState();
        set({ region });
        await persistPreferences({ language: current.language, region });
      },
      setPreferences: ({ language, region }) => set({ language, region }),
      markHydrated: () => set({ isHydrated: true }),
    }),
    {
      name: 'lemon-tea-language',
      version: 3,
      migrate: (persistedState: unknown, version) => {
        const state = (persistedState ?? {}) as { language?: AppLanguage; region?: AppRegion };
        if (version < 3) {
          return {
            ...state,
            language: state.language ?? DEFAULT_LANGUAGE,
            region: state.region ?? DEFAULT_REGION,
            isHydrated: false,
          };
        }
        return {
          language: state.language ?? DEFAULT_LANGUAGE,
          region: state.region ?? DEFAULT_REGION,
          isHydrated: false,
        };
      },
      onRehydrateStorage: () => (state) => {
        if (state) {
          applyLanguage(state.language);
        }
      },
    },
  ),
);

if (typeof window !== 'undefined') {
  window.addEventListener('storage', (event) => {
    if (event.key !== 'lemon-tea-language' || !event.newValue) {
      return;
    }
    try {
      const nextState = JSON.parse(event.newValue);
      const language = nextState.state?.language as AppLanguage | undefined;
      const region = nextState.state?.region as AppRegion | undefined;
      if (!language) {
        return;
      }
        const currentState = useLanguageStore.getState();
        if (language !== currentState.language || region !== currentState.region) {
          useLanguageStore.setState({
            language,
            region: region ?? DEFAULT_REGION,
        });
        applyLanguage(language);
      }
    } catch (error) {
      console.error('Failed to sync language across tabs:', error);
    }
  });
}

export function initializeLanguage() {
  applyLanguage(useLanguageStore.getState().language);
}

export async function hydrateLanguagePreferences() {
  const state = useLanguageStore.getState();
  try {
    const preferences = await Service.GetAppPreferences();
    if (preferences) {
      state.setPreferences({
        language: preferences.language as AppLanguage,
        region: preferences.region as AppRegion,
      });
      applyLanguage(preferences.language as AppLanguage);
    } else {
      await persistPreferences({ language: state.language, region: state.region });
    }
  } catch (error) {
    console.error('Failed to hydrate language preferences from backend:', error);
  } finally {
    useLanguageStore.getState().markHydrated();
  }
}
