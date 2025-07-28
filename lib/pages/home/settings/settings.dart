import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:lemon_tea/generated/l10n.dart';
import 'package:lemon_tea/utils/setting/manager.dart' as app_theme;
import 'package:lemon_tea/utils/setting/storage.dart';
import 'package:lemon_tea/pages/home/settings/general_settings.dart';
import 'package:lemon_tea/pages/home/settings/model_settings/model_settings.dart';
import 'package:lemon_tea/pages/home/settings/data_settings.dart';
import 'package:lemon_tea/pages/home/settings/about_settings.dart';

class SettingsPage extends ConsumerStatefulWidget {
  const SettingsPage({super.key});

  @override
  ConsumerState<SettingsPage> createState() => _SettingsPageState();
}

class _SettingsPageState extends ConsumerState<SettingsPage> {
  @override
  void initState() {
    super.initState();
    // 初始化时加载设置
    WidgetsBinding.instance.addPostFrameCallback((_) {
      final settings = ref.read(settingsManagerProvider);
      // 如果语言设置与当前不同，需要重新加载语言
      if (settings.language != '中文' && settings.language == 'English') {
        S.load(const Locale('en', 'US'));
      }
    });
  }

  final List<Map<String, dynamic>> _menuItems = [
    {'title': 'general', 'icon': Icons.settings_outlined},
    {'title': 'model', 'icon': Icons.model_training_outlined},
    {'title': 'data', 'icon': Icons.storage_outlined},
    {'title': 'about', 'icon': Icons.info_outline},
  ];

  @override
  Widget build(BuildContext context) {
    final themeMode = ref.watch(app_theme.themeManagerProvider);
    final themeManager = ref.read(app_theme.themeManagerProvider.notifier);
    final settings = ref.watch(settingsManagerProvider);
    final settingsManager = ref.read(settingsManagerProvider.notifier);

    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        // 左侧菜单
        Container(
          width: 200,
          decoration: BoxDecoration(
            border: Border(
              right: BorderSide(
                color: Theme.of(context).dividerColor.withOpacity(0.2),
                width: 1,
              ),
            ),
          ),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Padding(
                padding: const EdgeInsets.all(16.0),
                child: Text(
                  S.of(context).settings,
                  style: TextStyle(
                    fontSize: FontSizeUtils.getHeadingSize(ref),
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
              Expanded(
                child: ListView.builder(
                  itemCount: _menuItems.length,
                  itemBuilder: (context, index) {
                    final item = _menuItems[index];
                    final isSelected = settings.selectedMenuIndex == index;

                    // 根据菜单项的title获取对应的多语言文本
                    String title;
                    switch (item['title']) {
                      case 'general':
                        title = S.of(context).general;
                        break;
                      case 'model':
                        title = S.of(context).model;
                        break;
                      case 'data':
                        title = S.of(context).data;
                        break;
                      case 'about':
                        title = S.of(context).about;
                        break;
                      default:
                        title = '';
                    }

                    return ListTile(
                      leading: Icon(
                        item['icon'],
                        color:
                            isSelected ? Theme.of(context).colorScheme.primary : null,
                      ),
                      title: Text(
                        title,
                        style: TextStyle(
                          fontWeight:
                              isSelected ? FontWeight.bold : FontWeight.normal,
                          color:
                              isSelected
                                  ? Theme.of(context).colorScheme.primary
                                  : null,
                          fontSize: FontSizeUtils.getBodySize(ref),
                        ),
                      ),
                      selected: isSelected,
                      onTap: () {
                        settingsManager.setSelectedMenuIndex(index);
                      },
                    );
                  },
                ),
              ),
            ],
          ),
        ),

        // 右侧内容区
        Expanded(child: _buildContent()),
      ],
    );
  }

  Widget _buildContent() {
    final settings = ref.watch(settingsManagerProvider);
    switch (settings.selectedMenuIndex) {
      case 0:
        return const GeneralSettings();
      case 1:
        return const ModelSettings();
      case 2:
        return const DataSettings();
      case 3:
        return const AboutSettings();
      default:
        return const GeneralSettings();
    }
  }

  Widget _buildGeneralSettings() {
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
            title: S.of(context).language,
            children: [
              ListTile(
                title: Text(S.of(context).interfaceLanguage),
                subtitle: Text(settings.language),
                trailing: const Icon(Icons.arrow_forward_ios, size: 16),
                onTap: () {
                  _showLanguageDialog();
                },
              ),
            ],
          ),
        ],
      ),
    );
  }


  Widget _buildSection({
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
                width: 1.0, // 可以自定义宽度
              ),
            ),
          ),
          child: Column(children: children),
        ),
      ],
    );
  }

  void _showLanguageDialog() {
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

  void _showClearDataDialog() {
    final settingsManager = ref.read(settingsManagerProvider.notifier);

    showDialog(
      context: context,
      builder:
          (context) => AlertDialog(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(4),
            ),
            title: Text(S.of(context).confirmDelete),
            content: const Text('确定要清除所有数据吗？此操作无法撤销。'),
            actions: [
              TextButton(
                onPressed: () => Navigator.of(context).pop(),
                child: Text(S.of(context).cancel),
              ),
              TextButton(
                onPressed: () async {
                  // 清除设置数据
                  await settingsManager.clearAllSettings();
                  Navigator.of(context).pop();
                  ScaffoldMessenger.of(
                    context,
                  ).showSnackBar(const SnackBar(content: Text('数据已清除')));
                },
                style: TextButton.styleFrom(foregroundColor: Colors.red),
                child: Text(S.of(context).delete),
              ),
            ],
          ),
    );
  }
}
