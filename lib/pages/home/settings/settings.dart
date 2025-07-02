import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:lemon_tea/generated/l10n.dart';
import 'package:lemon_tea/utils/setting/manager.dart' as app_theme;
import 'package:lemon_tea/utils/setting/storage.dart';
import 'package:lemon_tea/utils/setting/provider_manager.dart';
import 'package:lemon_tea/pages/home/settings/provider_dialog.dart';
import 'package:lemon_tea/models/llm_provider.dart';
import 'package:lemon_tea/models/model.dart';

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
                    fontWeight: FontWeight.bold
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
                            isSelected ? Theme.of(context).primaryColor : null,
                      ),
                      title: Text(
                        title,
                        style: TextStyle(
                          fontWeight:
                              isSelected ? FontWeight.bold : FontWeight.normal,
                          color:
                              isSelected
                                  ? Theme.of(context).primaryColor
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
    final fontSizeMode = ref.watch(app_theme.fontSizeModeProvider);
    final fontSizeManager = ref.read(app_theme.fontSizeModeProvider.notifier);
    final settings = ref.watch(settingsManagerProvider);
    final settingsManager = ref.read(settingsManagerProvider.notifier);
    
    // 基础字体大小为14
    final double baseFontSize = 14.0;
    final double currentFontSize = app_theme.calculateFontSize(baseFontSize, fontSizeMode);

    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            S.of(context).generalSettings,
            style: TextStyle(
              fontSize: FontSizeUtils.getHeadingSize(ref),
              fontWeight: FontWeight.bold
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
                              child: Text(modeName,style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref)),),
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
                subtitle: Text('${app_theme.getLocalizedFontSizeModeName(context, fontSizeMode)} (${currentFontSize.toInt()}px)'),
                trailing: DropdownButton<app_theme.FontSizeMode>(
                  value: fontSizeMode,
                  underline: Container(),
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
                          child: Text(app_theme.getLocalizedFontSizeModeName(context, mode), style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref)),),
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

  Widget _buildModelSettings() {
    final providers = ref.watch(providerManagerProvider);
    final selectedProvider = ref.watch(selectedProviderProvider);
    final selectedModel = ref.watch(selectedModelProvider);
    
    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            S.of(context).modelSettings,
            style: TextStyle(
              fontSize: FontSizeUtils.getHeadingSize(ref),
              fontWeight: FontWeight.bold
            ),
          ),
          const SizedBox(height: 24),

          _buildSection(
            title: '模型供应商',
            children: [
              ListTile(
                title: const Text('添加供应商'),
                subtitle: const Text('添加新的AI模型供应商'),
                trailing: const Icon(Icons.add),
                onTap: () {
                  _showProviderDialog();
                },
              ),
              if (providers.isNotEmpty) ...[
                const Divider(height: 1),
                ...providers.map((provider) => _buildProviderTile(provider)),
              ],
            ],
          ),

          if (selectedProvider != null) ...[
            const SizedBox(height: 24),
            _buildSection(
              title: '${selectedProvider.displayName} 的模型',
              children: [
                if (selectedProvider.models != null && selectedProvider.models!.isNotEmpty)
                  ...selectedProvider.models!.map((model) => _buildModelTile(model, selectedProvider))
                else
                  const ListTile(
                    title: Text('暂无模型'),
                    subtitle: Text('该供应商暂无可用模型'),
                  ),
              ],
            ),
          ],

          const SizedBox(height: 24),
          _buildSection(
            title: '当前选择',
            children: [
              ListTile(
                title: Text(selectedProvider?.displayName ?? '未选择供应商'),
                subtitle: Text(selectedModel?.displayName ?? '未选择模型'),
                trailing: const Icon(Icons.settings),
                onTap: () {
                  _showModelSelectionDialog();
                },
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildProviderTile(LlmProvider provider) {
    return ListTile(
      title: Text(provider.displayName),
      subtitle: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(provider.baseUrl),
          if (provider.description != null) Text(provider.description!),
          Text('模型数量: ${provider.models?.length ?? 0}'),
        ],
      ),
      trailing: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            provider.hasApiKey ? Icons.check_circle : Icons.warning,
            color: provider.hasApiKey ? Colors.green : Colors.orange,
          ),
          const SizedBox(width: 8),
          PopupMenuButton<String>(
            onSelected: (value) {
              switch (value) {
                case 'edit':
                  _showProviderDialog(provider: provider);
                  break;
                case 'delete':
                  _showDeleteProviderDialog(provider);
                  break;
              }
            },
            itemBuilder: (context) => [
              const PopupMenuItem(
                value: 'edit',
                child: Row(
                  children: [
                    Icon(Icons.edit),
                    SizedBox(width: 8),
                    Text('编辑'),
                  ],
                ),
              ),
              const PopupMenuItem(
                value: 'delete',
                child: Row(
                  children: [
                    Icon(Icons.delete, color: Colors.red),
                    SizedBox(width: 8),
                    Text('删除', style: TextStyle(color: Colors.red)),
                  ],
                ),
              ),
            ],
          ),
        ],
      ),
      onTap: () {
        ref.read(selectedProviderProvider.notifier).state = provider;
        ref.read(selectedModelProvider.notifier).state = null;
      },
    );
  }

  Widget _buildModelTile(Model model, LlmProvider provider) {
    final selectedModel = ref.watch(selectedModelProvider);
    final isSelected = selectedModel?.id == model.id;
    
    return ListTile(
      title: Text(model.displayName),
      subtitle: Text('类型: ${model.object}'),
      trailing: Icon(
        isSelected ? Icons.check_circle : Icons.circle_outlined,
        color: isSelected ? Colors.green : null,
      ),
      onTap: () {
        ref.read(selectedModelProvider.notifier).state = model;
      },
    );
  }

  Widget _buildDataSettings() {
    final settings = ref.watch(settingsManagerProvider);
    final settingsManager = ref.read(settingsManagerProvider.notifier);
    
    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            S.of(context).dataSettings,
            style: TextStyle(
              fontSize: FontSizeUtils.getHeadingSize(ref),
              fontWeight: FontWeight.bold
            ),
          ),
          const SizedBox(height: 24),

          _buildSection(
            title: S.of(context).dataStorage,
            children: [
              SwitchListTile(
                title: Text(S.of(context).autoSaveData),
                subtitle: const Text('自动保存对话内容'),
                value: settings.autoSave,
                onChanged: (value) {
                  settingsManager.setAutoSave(value);
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
          Text(
            '关于',
            style: TextStyle(
              fontSize: FontSizeUtils.getHeadingSize(ref),
              fontWeight: FontWeight.bold
            ),
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
          style: TextStyle(
            fontSize: FontSizeUtils.getSubheadingSize(ref),
            fontWeight: FontWeight.w600
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

  void _showProviderDialog({LlmProvider? provider}) {
    showDialog(
      context: context,
      builder: (context) => ProviderDialog(provider: provider),
    ).then((result) {
      if (result == true) {
        // 刷新供应商列表
        setState(() {});
      }
    });
  }

  void _showModelSelectionDialog() {
    final providers = ref.read(providerManagerProvider);
    final selectedProvider = ref.read(selectedProviderProvider);
    final selectedModel = ref.read(selectedModelProvider);
    
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('选择模型'),
        content: SizedBox(
          width: 400,
          height: 300,
          child: Column(
            children: [
              // 供应商选择
              const Text('选择供应商:', style: TextStyle(fontWeight: FontWeight.bold)),
              const SizedBox(height: 8),
              Expanded(
                child: ListView.builder(
                  itemCount: providers.length,
                  itemBuilder: (context, index) {
                    final provider = providers[index];
                    final isSelected = selectedProvider?.name == provider.name;
                    
                    return ListTile(
                      title: Text(provider.displayName),
                      subtitle: Text(provider.baseUrl),
                      trailing: Icon(
                        isSelected ? Icons.check_circle : Icons.circle_outlined,
                        color: isSelected ? Colors.green : null,
                      ),
                      onTap: () {
                        ref.read(selectedProviderProvider.notifier).state = provider;
                        ref.read(selectedModelProvider.notifier).state = null;
                        Navigator.of(context).pop();
                      },
                    );
                  },
                ),
              ),
              
              if (selectedProvider != null) ...[
                const Divider(),
                const Text('选择模型:', style: TextStyle(fontWeight: FontWeight.bold)),
                const SizedBox(height: 8),
                Expanded(
                  child: ListView.builder(
                    itemCount: selectedProvider.models?.length ?? 0,
                    itemBuilder: (context, index) {
                      final model = selectedProvider.models![index];
                      final isSelected = selectedModel?.id == model.id;
                      
                      return ListTile(
                        title: Text(model.displayName),
                        subtitle: Text('类型: ${model.object}'),
                        trailing: Icon(
                          isSelected ? Icons.check_circle : Icons.circle_outlined,
                          color: isSelected ? Colors.green : null,
                        ),
                        onTap: () {
                          ref.read(selectedModelProvider.notifier).state = model;
                          Navigator.of(context).pop();
                        },
                      );
                    },
                  ),
                ),
              ],
            ],
          ),
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

  void _showDeleteProviderDialog(LlmProvider provider) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('确认删除'),
        content: Text('确定要删除供应商 "${provider.displayName}" 吗？此操作无法撤销。'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: Text(S.of(context).cancel),
          ),
          TextButton(
            onPressed: () async {
              try {
                final providerManager = ref.read(providerManagerProvider.notifier);
                await providerManager.deleteProvider(provider.name);
                Navigator.of(context).pop();
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(content: Text('供应商删除成功')),
                );
              } catch (e) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text('删除失败：${e.toString()}'),
                    backgroundColor: Colors.red,
                  ),
                );
              }
            },
            style: TextButton.styleFrom(foregroundColor: Colors.red),
            child: Text(S.of(context).delete),
          ),
        ],
      ),
    );
  }
}
