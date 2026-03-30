import enUSAntd from 'antd/locale/en_US';
import zhCNAntd from 'antd/locale/zh_CN';
import zhCN from './resources/zh-CN';
import enUS from './resources/en-US';

export const resources = {
  'zh-CN': {
    translation: zhCN,
  },
  'en-US': {
    translation: enUS,
  },
} as const;

export type TranslationSchema = typeof zhCN;
export type AppLanguage = keyof typeof resources;
export type AppRegion =
  | 'asia'
  | 'europe'
  | 'north-america'
  | 'south-america'
  | 'africa'
  | 'oceania'
  | 'antarctica';

export interface LanguageOption {
  value: AppLanguage;
  nativeLabel: string;
  localeLabelKey: string;
}

export interface RegionOption {
  value: AppRegion;
  labelKey: string;
}

export const DEFAULT_LANGUAGE: AppLanguage = 'zh-CN';
export const DEFAULT_REGION: AppRegion = 'asia';

export const localeRegistry = {
  'zh-CN': {
    nativeLabel: '简体中文',
    localeLabelKey: 'common.chineseSimplified',
    antd: zhCNAntd,
  },
  'en-US': {
    nativeLabel: 'English',
    localeLabelKey: 'common.english',
    antd: enUSAntd,
  },
} as const satisfies Record<AppLanguage, { nativeLabel: string; localeLabelKey: string; antd: typeof zhCNAntd }>;

export const LANGUAGE_OPTIONS: LanguageOption[] = Object.entries(localeRegistry).map(([value, meta]) => ({
  value: value as AppLanguage,
  nativeLabel: meta.nativeLabel,
  localeLabelKey: meta.localeLabelKey,
}));

export const REGION_OPTIONS: RegionOption[] = [
  { value: 'asia', labelKey: 'settings.languageRegion.regions.asia' },
  { value: 'europe', labelKey: 'settings.languageRegion.regions.europe' },
  { value: 'north-america', labelKey: 'settings.languageRegion.regions.northAmerica' },
  { value: 'south-america', labelKey: 'settings.languageRegion.regions.southAmerica' },
  { value: 'africa', labelKey: 'settings.languageRegion.regions.africa' },
  { value: 'oceania', labelKey: 'settings.languageRegion.regions.oceania' },
  { value: 'antarctica', labelKey: 'settings.languageRegion.regions.antarctica' },
];

export const antdLocales = Object.fromEntries(
  Object.entries(localeRegistry).map(([key, meta]) => [key, meta.antd]),
) as Record<AppLanguage, (typeof localeRegistry)[AppLanguage]['antd']>;
