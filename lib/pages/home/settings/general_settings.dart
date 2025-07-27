import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:lemon_tea/generated/l10n.dart';
import 'package:lemon_tea/utils/setting/manager.dart' as app_theme;
import 'package:lemon_tea/utils/setting/provider_manager.dart';
import 'package:lemon_tea/utils/setting/storage.dart';
import 'package:lemon_tea/utils/style.dart';

class GeneralSettings extends ConsumerWidget {
  const GeneralSettings({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final themeMode = ref.watch(app_theme.themeManagerProvider);
    final themeManager = ref.read(app_theme.themeManagerProvider.notifier);
    final fontSizeMode = ref.watch(app_theme.fontSizeModeProvider);
    final fontSizeManager = ref.read(app_theme.fontSizeModeProvider.notifier);
    final settings = ref.watch(settingsManagerProvider);
    final settingsManager = ref.read(settingsManagerProvider.notifier);

    // 基础字体大小为14
    final double baseFontSize = 14.0;
    final double currentFontSize = app_theme.calculateFontSize(
      baseFontSize,
      fontSizeMode,
    );

    return Scaffold(
      backgroundColor: Style.primaryBackground(context),
      body: SingleChildScrollView(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // 页面标题
            Container(
              margin: const EdgeInsets.only(bottom: 20),
              child: Row(
                children: [
                  Container(
                    padding: const EdgeInsets.all(8),
                    decoration: BoxDecoration(
                      gradient: LinearGradient(
                        colors: [
                          Style.primaryColor(context).withOpacity(0.1),
                          Style.primaryColor(context).withOpacity(0.05),
                        ],
                      ),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Icon(
                      Icons.settings,
                      size: 20,
                      color: Style.primaryColor(context),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Text(
                    S.of(context).generalSettings,
                    style: TextStyle(
                      fontSize: FontSizeUtils.getHeadingSize(ref),
                      fontWeight: FontWeight.bold,
                      color: Style.primaryText(context),
                    ),
                  ),
                ],
              ),
            ),

            // 主题设置
            _buildModernSection(
              context: context,
              ref: ref,
              title: S.of(context).theme,
              icon: Icons.palette_outlined,
              children: [
                _buildModernListTile(
                  context: context,
                  ref: ref,
                  icon: Icons.brightness_6_outlined,
                  title: S.of(context).themeMode,
                  subtitle: themeManager.getLocalizedThemeModeName(context),
                  trailing: Container(
                    padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 2),
                    decoration: BoxDecoration(
                      color: Style.primaryColor(context).withOpacity(0.1),
                      borderRadius: BorderRadius.circular(16),
                      border: Border.all(
                        color: Style.primaryColor(context).withOpacity(0.2),
                      ),
                    ),
                    child: DropdownButtonHideUnderline(
                      child: DropdownButton<app_theme.ThemeMode>(
                        value: themeMode,
                        isDense: true,
                        style: TextStyle(
                          fontSize: FontSizeUtils.getBodySize(ref) - 1,
                          color: Theme.of(context).primaryColor,
                          fontWeight: FontWeight.w500,
                        ),
                        dropdownColor: Theme.of(context).cardColor,
                        borderRadius: BorderRadius.circular(8),
                        onChanged: (app_theme.ThemeMode? newValue) {
                          if (newValue != null) {
                            themeManager.setThemeMode(newValue);
                          }
                        },
                        items: app_theme.ThemeMode.values
                            .map<DropdownMenuItem<app_theme.ThemeMode>>((
                              app_theme.ThemeMode mode,
                            ) {
                              String modeName = '';
                              IconData modeIcon = Icons.brightness_auto;
                              switch (mode) {
                                case app_theme.ThemeMode.light:
                                  modeName = S.of(context).lightMode;
                                  modeIcon = Icons.light_mode;
                                  break;
                                case app_theme.ThemeMode.dark:
                                  modeName = S.of(context).darkMode;
                                  modeIcon = Icons.dark_mode;
                                  break;
                                case app_theme.ThemeMode.system:
                                  modeName = S.of(context).systemMode;
                                  modeIcon = Icons.brightness_auto;
                                  break;
                              }
                              return DropdownMenuItem<app_theme.ThemeMode>(
                                value: mode,
                                child: Row(
                                  mainAxisSize: MainAxisSize.min,
                                  children: [
                                    Icon(modeIcon, size: 14),
                                    const SizedBox(width: 6),
                                    Text(modeName),
                                  ],
                                ),
                              );
                            })
                            .toList(),
                      ),
                    ),
                  ),
                ),
              ],
            ),

            const SizedBox(height: 16),

            // 字体设置
            _buildModernSection(
              context: context,
              ref: ref,
              title: S.of(context).fontSize,
              icon: Icons.text_fields_outlined,
              children: [
                _buildModernListTile(
                  context: context,
                  ref: ref,
                  icon: Icons.format_size_outlined,
                  title: S.of(context).interfaceFont,
                  subtitle: '${app_theme.getLocalizedFontSizeModeName(context, fontSizeMode)} (${currentFontSize.toInt()}px)',
                  trailing: Container(
                    padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 2),
                    decoration: BoxDecoration(
                      color: Theme.of(context).primaryColor.withOpacity(0.1),
                      borderRadius: BorderRadius.circular(16),
                      border: Border.all(
                        color: Theme.of(context).primaryColor.withOpacity(0.2),
                      ),
                    ),
                    child: DropdownButtonHideUnderline(
                      child: DropdownButton<app_theme.FontSizeMode>(
                        value: fontSizeMode,
                        isDense: true,
                        style: TextStyle(
                          fontSize: FontSizeUtils.getBodySize(ref) - 1,
                          color: Theme.of(context).primaryColor,
                          fontWeight: FontWeight.w500,
                        ),
                        dropdownColor: Theme.of(context).cardColor,
                        borderRadius: BorderRadius.circular(8),
                        onChanged: (app_theme.FontSizeMode? newValue) {
                          if (newValue != null) {
                            fontSizeManager.setFontSizeMode(newValue);
                          }
                        },
                        items: app_theme.FontSizeMode.values
                            .map<DropdownMenuItem<app_theme.FontSizeMode>>((
                              app_theme.FontSizeMode mode,
                            ) {
                              return DropdownMenuItem<app_theme.FontSizeMode>(
                                value: mode,
                                child: Text(
                                  app_theme.getLocalizedFontSizeModeName(
                                    context,
                                    mode,
                                  ),
                                ),
                              );
                            })
                            .toList(),
                      ),
                    ),
                  ),
                ),
              ],
            ),

            const SizedBox(height: 16),

            // 语言设置
            _buildModernSection(
              context: context,
              ref: ref,
              title: S.of(context).language,
              icon: Icons.language_outlined,
              children: [
                _buildModernListTile(
                  context: context,
                  ref: ref,
                  icon: Icons.translate_outlined,
                  title: S.of(context).interfaceLanguage,
                  subtitle: settings.language,
                  trailing: Container(
                    padding: const EdgeInsets.all(6),
                    decoration: BoxDecoration(
                      color: Theme.of(context).primaryColor.withOpacity(0.1),
                      borderRadius: BorderRadius.circular(6),
                    ),
                    child: Icon(
                      Icons.arrow_forward_ios,
                      size: 14,
                      color: Theme.of(context).primaryColor,
                    ),
                  ),
                  onTap: () => _showLanguageDialog(context, ref),
                ),
              ],
            ),

            // 底部间距
            const SizedBox(height: 20),
          ],
        ),
      ),
    );
  }

  Widget _buildModernSection({
    required BuildContext context,
    required WidgetRef ref,
    required String title,
    required IconData icon,
    required List<Widget> children,
  }) {
    return Container(
      decoration: BoxDecoration(
        color: Theme.of(context).cardColor,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.04),
            blurRadius: 8,
            offset: const Offset(0, 1),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 分组标题
          Container(
            padding: const EdgeInsets.all(14),
            decoration: BoxDecoration(
              color: Style.primaryText(context).withOpacity(0.04),
              borderRadius: const BorderRadius.only(
                topLeft: Radius.circular(12),
                topRight: Radius.circular(12),
              ),
            ),
            child: Row(
              children: [
                Container(
                  padding: const EdgeInsets.all(6),
                  decoration: BoxDecoration(
                    color: Style.primaryColor(context).withOpacity(0.1),
                    borderRadius: BorderRadius.circular(6),
                  ),
                  child: Icon(
                    icon,
                    size: 16,
                    color: Style.primaryColor(context),
                  ),
                ),
                const SizedBox(width: 8),
                Text(
                  title,
                  style: TextStyle(
                    fontSize: FontSizeUtils.getBodyLargeSize(ref) - 1,
                    fontWeight: FontWeight.w600,
                    color: Style.buttonText(context),
                  ),
                ),
              ],
            ),
          ),
          // 内容
          ...children,
        ],
      ),
    );
  }

  Widget _buildModernListTile({
    required BuildContext context,
    required WidgetRef ref,
    required IconData icon,
    required String title,
    required String subtitle,
    required Widget trailing,
    VoidCallback? onTap,
  }) {
    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: onTap,
        borderRadius: const BorderRadius.only(
          bottomLeft: Radius.circular(12),
          bottomRight: Radius.circular(12),
        ),
        child: Container(
          padding: const EdgeInsets.all(14),
          child: Row(
            children: [
              // 图标
              Container(
                padding: const EdgeInsets.all(6),
                decoration: BoxDecoration(
                  color: Theme.of(context).primaryColor.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(6),
                ),
                child: Icon(
                  icon,
                  size: 16,
                  color: Theme.of(context).primaryColor,
                ),
              ),
              const SizedBox(width: 12),
              // 文本内容
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      title,
                      style: TextStyle(
                        fontSize: FontSizeUtils.getBodySize(ref) - 1,
                        fontWeight: FontWeight.w500,
                        color: Theme.of(context).textTheme.bodyLarge?.color,
                      ),
                    ),
                    const SizedBox(height: 2),
                    Text(
                      subtitle,
                      style: TextStyle(
                        fontSize: FontSizeUtils.getBodySize(ref) - 3,
                        color: Theme.of(context).textTheme.bodyMedium?.color?.withOpacity(0.7),
                      ),
                    ),
                  ],
                ),
              ),
              // 尾部组件
              trailing,
            ],
          ),
        ),
      ),
    );
  }

  void _showLanguageDialog(BuildContext context, WidgetRef ref) {
    final settings = ref.read(settingsManagerProvider);
    final settingsManager = ref.read(settingsManagerProvider.notifier);

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(16),
        ),
        elevation: 8,
        backgroundColor: Theme.of(context).cardColor,
        contentPadding: const EdgeInsets.all(16),
        title: Row(
          children: [
            Container(
              padding: const EdgeInsets.all(6),
              decoration: BoxDecoration(
                color: Theme.of(context).primaryColor.withOpacity(0.1),
                borderRadius: BorderRadius.circular(6),
              ),
              child: Icon(
                Icons.language,
                color: Theme.of(context).primaryColor,
                size: 18,
              ),
            ),
            const SizedBox(width: 8),
            Text(
              S.of(context).language,
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.w600,
                color: Theme.of(context).textTheme.titleLarge?.color,
              ),
            ),
          ],
        ),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            _buildLanguageOption(
              context: context,
              title: S.of(context).chinese,
              value: '中文',
              groupValue: settings.language,
              flag: '🇨🇳',
              onChanged: (value) {
                settingsManager.setLanguage(value!);
                S.load(const Locale('zh', 'CN'));
                Navigator.of(context).pop();
              },
            ),
            const SizedBox(height: 4),
            _buildLanguageOption(
              context: context,
              title: 'English',
              value: 'English',
              groupValue: settings.language,
              flag: '🇺🇸',
              onChanged: (value) {
                settingsManager.setLanguage(value!);
                S.load(const Locale('en', 'US'));
                Navigator.of(context).pop();
              },
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            style: TextButton.styleFrom(
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(6),
              ),
            ),
            child: Text(
              S.of(context).cancel,
              style: TextStyle(
                color: Theme.of(context).textTheme.bodyLarge?.color?.withOpacity(0.7),
                fontWeight: FontWeight.w500,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildLanguageOption({
    required BuildContext context,
    required String title,
    required String value,
    required String groupValue,
    required String flag,
    required ValueChanged<String?> onChanged,
  }) {
    final isSelected = value == groupValue;
    
    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: () => onChanged(value),
        borderRadius: BorderRadius.circular(8),
        child: Container(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
          decoration: BoxDecoration(
            color: isSelected 
                ? Theme.of(context).primaryColor.withOpacity(0.1)
                : Colors.transparent,
            borderRadius: BorderRadius.circular(8),
            border: Border.all(
              color: isSelected 
                  ? Theme.of(context).primaryColor.withOpacity(0.3)
                  : Colors.transparent,
            ),
          ),
          child: Row(
            children: [
              Text(
                flag,
                style: const TextStyle(fontSize: 18),
              ),
              const SizedBox(width: 8),
              Expanded(
                child: Text(
                  title,
                  style: TextStyle(
                    fontWeight: isSelected ? FontWeight.w600 : FontWeight.w400,
                    color: isSelected 
                        ? Theme.of(context).primaryColor
                        : Theme.of(context).textTheme.bodyLarge?.color,
                  ),
                ),
              ),
              if (isSelected)
                Icon(
                  Icons.check_circle,
                  color: Theme.of(context).primaryColor,
                  size: 18,
                ),
            ],
          ),
        ),
      ),
    );
  }
} 