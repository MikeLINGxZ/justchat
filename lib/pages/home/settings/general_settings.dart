import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:lemon_tea/generated/l10n.dart';
import 'package:lemon_tea/utils/setting/manager.dart' as app_theme;
import 'package:lemon_tea/utils/setting/provider_manager.dart';
import 'package:lemon_tea/utils/setting/storage.dart';

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

    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            S.of(context).generalSettings,
            style: TextStyle(
              fontSize: FontSizeUtils.getHeadingSize(ref),
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 24),

          _buildSection(
            context: context,
            ref: ref,
            title: S.of(context).theme,
            children: [
              ListTile(
                title: Text(S.of(context).themeMode),
                subtitle: Text(themeManager.getLocalizedThemeModeName(context)),
                trailing: DropdownButton<app_theme.ThemeMode>(
                  value: themeMode,
                  underline: Container(),
                  onChanged: (app_theme.ThemeMode? newValue) {
                    if (newValue != null) {
                      themeManager.setThemeMode(newValue);
                    }
                  },
                  items:
                      app_theme.ThemeMode.values
                          .map<DropdownMenuItem<app_theme.ThemeMode>>((
                            app_theme.ThemeMode mode,
                          ) {
                            String modeName = '';
                            switch (mode) {
                              case app_theme.ThemeMode.light:
                                modeName = S.of(context).lightMode;
                                break;
                              case app_theme.ThemeMode.dark:
                                modeName = S.of(context).darkMode;
                                break;
                              case app_theme.ThemeMode.system:
                                modeName = S.of(context).systemMode;
                                break;
                            }
                            return DropdownMenuItem<app_theme.ThemeMode>(
                              value: mode,
                              child: Text(
                                modeName,
                                style: TextStyle(
                                  fontSize: FontSizeUtils.getBodySize(ref),
                                ),
                              ),
                            );
                          })
                          .toList(),
                ),
              ),
            ],
          ),

          const SizedBox(height: 24),

          _buildSection(
            context: context,
            ref: ref,
            title: S.of(context).fontSize,
            children: [
              ListTile(
                title: Text(S.of(context).interfaceFont),
                subtitle: Text(
                  '${app_theme.getLocalizedFontSizeModeName(context, fontSizeMode)} (${currentFontSize.toInt()}px)',
                ),
                trailing: DropdownButton<app_theme.FontSizeMode>(
                  value: fontSizeMode,
                  underline: Container(),
                  onChanged: (app_theme.FontSizeMode? newValue) {
                    if (newValue != null) {
                      fontSizeManager.setFontSizeMode(newValue);
                    }
                  },
                  items:
                      app_theme.FontSizeMode.values
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
                                style: TextStyle(
                                  fontSize: FontSizeUtils.getBodySize(ref),
                                ),
                              ),
                            );
                          })
                          .toList(),
                ),
              ),
            ],
          ),

          const SizedBox(height: 24),

          _buildSection(
            context: context,
            ref: ref,
            title: S.of(context).language,
            children: [
              ListTile(
                title: Text(S.of(context).interfaceLanguage),
                subtitle: Text(settings.language),
                trailing: const Icon(Icons.arrow_forward_ios, size: 16),
                onTap: () {
                  _showLanguageDialog(context, ref);
                },
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildSection({
    required BuildContext context,
    required WidgetRef ref,
    required String title,
    required List<Widget> children,
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          title,
          style: TextStyle(
            fontSize: FontSizeUtils.getSubheadingSize(ref),
            fontWeight: FontWeight.w600,
          ),
        ),
        const SizedBox(height: 8),
        Container(
          decoration: BoxDecoration(
            color: Theme.of(context).cardColor,
            border: Border(
              bottom: BorderSide(
                color: Colors.grey.withOpacity(0.2),
                width: 1.0,
              ),
            ),
          ),
          child: Column(children: children),
        ),
      ],
    );
  }

  void _showLanguageDialog(BuildContext context, WidgetRef ref) {
    final settings = ref.read(settingsManagerProvider);
    final settingsManager = ref.read(settingsManagerProvider.notifier);

    showDialog(
      context: context,
      builder:
          (context) => AlertDialog(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(4),
            ),
            title: Text(S.of(context).language),
            content: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                RadioListTile<String>(
                  title: Text(S.of(context).chinese),
                  value: '中文',
                  groupValue: settings.language,
                  onChanged: (value) {
                    settingsManager.setLanguage(value!);
                    S.load(const Locale('zh', 'CN'));
                    Navigator.of(context).pop();
                  },
                ),
                RadioListTile<String>(
                  title: const Text('English'),
                  value: 'English',
                  groupValue: settings.language,
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
                child: Text(S.of(context).cancel),
              ),
            ],
          ),
    );
  }
} 