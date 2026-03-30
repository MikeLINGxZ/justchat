package i18n

import (
	"fmt"
	"strings"
	"sync/atomic"
)

const (
	LocaleZhCN    = "zh-CN"
	LocaleEnUS    = "en-US"
	DefaultLocale = LocaleZhCN
)

var currentLocale atomic.Value

func init() {
	currentLocale.Store(DefaultLocale)
}

var resources = map[string]map[string]string{
	LocaleZhCN: zhCN,
	LocaleEnUS: enUS,
}

func NormalizeLocale(locale string) string {
	trimmed := strings.TrimSpace(locale)
	if _, ok := resources[trimmed]; ok {
		return trimmed
	}
	return DefaultLocale
}

func SetCurrentLocale(locale string) {
	currentLocale.Store(NormalizeLocale(locale))
}

func CurrentLocale() string {
	if val, ok := currentLocale.Load().(string); ok && val != "" {
		return NormalizeLocale(val)
	}
	return DefaultLocale
}

func T(locale, key string, params map[string]string) string {
	dict := resources[NormalizeLocale(locale)]
	if text, ok := dict[key]; ok {
		return applyParams(text, params)
	}
	if fallback, ok := resources[DefaultLocale][key]; ok {
		return applyParams(fallback, params)
	}
	return key
}

func TCurrent(key string, params map[string]string) string {
	return T(CurrentLocale(), key, params)
}

func Sprintf(locale, key string, args ...interface{}) string {
	return fmt.Sprintf(T(locale, key, nil), args...)
}

func applyParams(text string, params map[string]string) string {
	if len(params) == 0 {
		return text
	}
	result := text
	for key, value := range params {
		result = strings.ReplaceAll(result, "{{"+key+"}}", value)
	}
	return result
}
