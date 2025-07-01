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
  int _selectedMenuIndex = 0;

  final List<Map<String, dynamic>> _menuItems = [
    {'title': '通用', 'icon': Icons.settings_outlined},
    {'title': '模型', 'icon': Icons.model_training_outlined},
    {'title': '数据', 'icon': Icons.storage_outlined},
    {'title': '关于', 'icon': Icons.info_outline},
  ];

  @override
  Widget build(BuildContext context) {
    final themeMode = ref.watch(app_theme.themeManagerProvider);
    final themeManager = ref.read(app_theme.themeManagerProvider.notifier);

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
                  '设置',
                  style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
                ),
              ),
              Expanded(
                child: ListView.builder(
                  itemCount: _menuItems.length,
                  itemBuilder: (context, index) {
                    final item = _menuItems[index];
                    final isSelected = _selectedMenuIndex == index;

                    return ListTile(
                      leading: Icon(
                        item['icon'],
                        color:
                            isSelected ? Theme.of(context).primaryColor : null,
                      ),
                      title: Text(
                        item['title'],
                        style: TextStyle(
                          fontWeight:
                              isSelected ? FontWeight.bold : FontWeight.normal,
                          color:
                              isSelected
                                  ? Theme.of(context).primaryColor
                                  : null,
                          fontSize: 14,
                        ),
                      ),
                      selected: isSelected,
                      onTap: () {
                        setState(() {
                          _selectedMenuIndex = index;
                        });
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
    switch (_selectedMenuIndex) {
      case 0:
        return _buildGeneralSettings();
      case 1:
        return _buildModelSettings();
      case 2:
        return _buildDataSettings();
      case 3:
        return _buildAboutSettings();
      default:
        return _buildGeneralSettings();
    }
  }

  Widget _buildGeneralSettings() {
    final themeMode = ref.watch(app_theme.themeManagerProvider);
    final themeManager = ref.read(app_theme.themeManagerProvider.notifier);

    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '通用设置',
            style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 24),

          _buildSection(
            title: '主题',
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
            ],
          ),

          const SizedBox(height: 24),

          _buildSection(
            title: '字体大小',
            children: [
              ListTile(
                title: const Text('界面字体'),
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

          _buildSection(
            title: '语言',
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
        ],
      ),
    );
  }

  Widget _buildModelSettings() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '模型设置',
            style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 24),

          _buildSection(
            title: '添加模型',
            children: [
              ListTile(
                title: const Text('添加新模型'),
                trailing: const Icon(Icons.add),
                onTap: () {
                  // TODO: 实现添加模型功能
                },
              ),
            ],
          ),

          const SizedBox(height: 24),

          _buildSection(
            title: '模型列表',
            children: [
              ListTile(
                title: const Text('GPT-4'),
                subtitle: const Text('OpenAI'),
                trailing: const Icon(Icons.check_circle, color: Colors.green),
              ),
              const Divider(height: 1),
              ListTile(
                title: const Text('Claude 3'),
                subtitle: const Text('Anthropic'),
                trailing: const Icon(Icons.circle_outlined),
              ),
              const Divider(height: 1),
              ListTile(
                title: const Text('Gemini Pro'),
                subtitle: const Text('Google'),
                trailing: const Icon(Icons.circle_outlined),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildDataSettings() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '数据设置',
            style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 24),

          _buildSection(
            title: '数据存储',
            children: [
              SwitchListTile(
                title: const Text('自动保存数据'),
                subtitle: const Text('自动保存对话内容'),
                value: _autoSave,
                onChanged: (value) {
                  setState(() {
                    _autoSave = value;
                  });
                },
              ),
            ],
          ),

          const SizedBox(height: 24),

          _buildSection(
            title: '数据管理',
            children: [
              ListTile(
                title: const Text('清空所有数据'),
                subtitle: const Text('删除所有对话和设置'),
                trailing: const Icon(Icons.delete_outline, color: Colors.red),
                onTap: () {
                  _showClearDataDialog();
                },
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildAboutSettings() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '关于',
            style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 24),

          _buildSection(
            title: '应用信息',
            children: [
              ListTile(
                title: const Text('Lemon Tea'),
                subtitle: const Text('版本 1.0.0'),
              ),
              const Divider(height: 1),
              ListTile(
                title: const Text('开发者'),
                subtitle: const Text('Lemon Tea Team'),
              ),
            ],
          ),

          const SizedBox(height: 24),

          _buildSection(
            title: '帮助',
            children: [
              ListTile(
                title: const Text('帮助文档'),
                trailing: const Icon(Icons.arrow_forward_ios, size: 16),
                onTap: () {
                  // TODO: 打开帮助文档
                },
              ),
              const Divider(height: 1),
              ListTile(
                title: const Text('反馈问题'),
                trailing: const Icon(Icons.arrow_forward_ios, size: 16),
                onTap: () {
                  // TODO: 打开反馈页面
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
          style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w600),
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
}
