import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/utils/theme_manager.dart' as app_theme;

class SettingsPage extends ConsumerStatefulWidget {
  const SettingsPage({super.key});

  @override
  ConsumerState<SettingsPage> createState() => _SettingsPageState();
}

class _SettingsPageState extends ConsumerState<SettingsPage> {
  bool _autoSave = true;
  String _selectedLanguage = '中文';
  double _fontSize = 14.0;
  double lrPadding = 16.0;

  @override
  Widget build(BuildContext context) {
    // 获取当前主题模式
    final themeMode = ref.watch(app_theme.themeManagerProvider);
    final themeManager = ref.read(app_theme.themeManagerProvider.notifier);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: EdgeInsets.fromLTRB(
            lrPadding,
            lrPadding,
            lrPadding,
            lrPadding / 3,
          ),
          child: const Text(
            '设置',
            style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
          ),
        ),
        const SizedBox(height: 20),

        Expanded(
          child: SingleChildScrollView(
            padding: EdgeInsets.fromLTRB(
              lrPadding,
              lrPadding,
              lrPadding,
              lrPadding / 3,
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // 外观设置
                _buildSection(
                  title: '外观',
                  icon: Icons.palette_outlined,
                  children: [
                    ListTile(
                      title: const Text('主题模式'),
                      subtitle: Text(themeManager.getThemeModeName()),
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
                                      modeName = '浅色模式';
                                      break;
                                    case app_theme.ThemeMode.dark:
                                      modeName = '深色模式';
                                      break;
                                    case app_theme.ThemeMode.system:
                                      modeName = '跟随系统';
                                      break;
                                  }
                                  return DropdownMenuItem<app_theme.ThemeMode>(
                                    value: mode,
                                    child: Text(modeName),
                                  );
                                })
                                .toList(),
                      ),
                    ),
                    const Divider(),
                    ListTile(
                      title: const Text('字体大小'),
                      subtitle: Text('${_fontSize.toInt()}px'),
                      trailing: SizedBox(
                        width: 200,
                        child: Slider(
                          value: _fontSize,
                          min: 12.0,
                          max: 20.0,
                          divisions: 8,
                          onChanged: (value) {
                            setState(() {
                              _fontSize = value;
                            });
                          },
                        ),
                      ),
                    ),
                  ],
                ),

                const SizedBox(height: 24),

                // 语言设置
                _buildSection(
                  title: '语言',
                  icon: Icons.language_outlined,
                  children: [
                    ListTile(
                      title: const Text('界面语言'),
                      subtitle: Text(_selectedLanguage),
                      trailing: const Icon(Icons.arrow_forward_ios, size: 16),
                      onTap: () {
                        _showLanguageDialog();
                      },
                    ),
                  ],
                ),

                const SizedBox(height: 24),

                // 数据设置
                _buildSection(
                  title: '数据',
                  icon: Icons.storage_outlined,
                  children: [
                    SwitchListTile(
                      title: const Text('自动保存'),
                      subtitle: const Text('自动保存对话内容'),
                      value: _autoSave,
                      onChanged: (value) {
                        setState(() {
                          _autoSave = value;
                        });
                      },
                    ),
                    const Divider(),
                    ListTile(
                      title: const Text('清除所有数据'),
                      subtitle: const Text('删除所有对话和设置'),
                      trailing: const Icon(
                        Icons.delete_outline,
                        color: Colors.red,
                      ),
                      onTap: () {
                        _showClearDataDialog();
                      },
                    ),
                  ],
                ),

                const SizedBox(height: 24),

                // 关于
                _buildSection(
                  title: '关于',
                  icon: Icons.info_outline,
                  children: [
                    ListTile(
                      title: const Text('版本信息'),
                      subtitle: const Text('Lemon Tea v1.0.0'),
                      trailing: const Icon(Icons.arrow_forward_ios, size: 16),
                      onTap: () {
                        _showAboutDialog();
                      },
                    ),
                    const Divider(),
                    ListTile(
                      title: const Text('帮助文档'),
                      subtitle: const Text('查看使用说明'),
                      trailing: const Icon(Icons.arrow_forward_ios, size: 16),
                      onTap: () {
                        // TODO: 打开帮助文档
                      },
                    ),
                  ],
                ),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildSection({
    required String title,
    required IconData icon,
    required List<Widget> children,
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            Icon(icon, size: 20, color: Colors.grey[600]),
            const SizedBox(width: 8),
            Text(
              title,
              style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w600),
            ),
          ],
        ),
        const SizedBox(height: 12),
        Container(
          decoration: BoxDecoration(
            color: Theme.of(context).cardColor,
            // borderRadius: BorderRadius.circular(3),
            border: Border.all(color: Colors.grey.withOpacity(0.2)),
          ),
          child: Column(children: children),
        ),
      ],
    );
  }

  void _showLanguageDialog() {
    showDialog(
      context: context,
      builder:
          (context) => AlertDialog(
            title: const Text('选择语言'),
            content: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                RadioListTile<String>(
                  title: const Text('中文'),
                  value: '中文',
                  groupValue: _selectedLanguage,
                  onChanged: (value) {
                    setState(() {
                      _selectedLanguage = value!;
                    });
                    Navigator.of(context).pop();
                  },
                ),
                RadioListTile<String>(
                  title: const Text('English'),
                  value: 'English',
                  groupValue: _selectedLanguage,
                  onChanged: (value) {
                    setState(() {
                      _selectedLanguage = value!;
                    });
                    Navigator.of(context).pop();
                  },
                ),
              ],
            ),
            actions: [
              TextButton(
                onPressed: () => Navigator.of(context).pop(),
                child: const Text('取消'),
              ),
            ],
          ),
    );
  }

  void _showClearDataDialog() {
    showDialog(
      context: context,
      builder:
          (context) => AlertDialog(
            title: const Text('确认清除'),
            content: const Text('确定要清除所有数据吗？此操作无法撤销。'),
            actions: [
              TextButton(
                onPressed: () => Navigator.of(context).pop(),
                child: const Text('取消'),
              ),
              TextButton(
                onPressed: () {
                  // TODO: 实现清除数据功能
                  Navigator.of(context).pop();
                  ScaffoldMessenger.of(
                    context,
                  ).showSnackBar(const SnackBar(content: Text('数据已清除')));
                },
                style: TextButton.styleFrom(foregroundColor: Colors.red),
                child: const Text('清除'),
              ),
            ],
          ),
    );
  }

  void _showAboutDialog() {
    showDialog(
      context: context,
      builder:
          (context) => AlertDialog(
            title: const Text('关于 Lemon Tea'),
            content: const Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Lemon Tea v1.0.0'),
                SizedBox(height: 8),
                Text('一个简洁的AI助手应用'),
                SizedBox(height: 16),
                Text('© 2024 Lemon Tea Team'),
              ],
            ),
            actions: [
              TextButton(
                onPressed: () => Navigator.of(context).pop(),
                child: const Text('确定'),
              ),
            ],
          ),
    );
  }
}
