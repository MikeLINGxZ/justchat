import { useLanguageStore } from '@/stores/languageStore';
import { antdLocales } from '@/i18n/types';

export function useAppLocale() {
  const language = useLanguageStore((state) => state.language);
  const region = useLanguageStore((state) => state.region);

  return {
    language,
    region,
    antdLocale: antdLocales[language],
  };
}
