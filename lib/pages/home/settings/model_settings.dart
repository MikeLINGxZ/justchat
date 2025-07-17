import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/storage/llm_storage.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:lemon_tea/generated/l10n.dart';
import 'package:lemon_tea/models/llm_provider.dart';
import 'package:lemon_tea/models/model.dart';

// 创建提供商列表的Provider
final providersProvider = FutureProvider<List<LlmProvider>>((ref) async {
  return LlmStorage.getAllProviders();
});

// 创建模型列表的Provider，接受提供商ID作为参数
final modelsProvider = FutureProvider.family<List<Model>, String>((ref, providerId) async {
  try {
    final models = await LlmStorage.getModelsByProviderId(providerId);
    if (models.isNotEmpty) {
      return models;
    }
    // 如果没有数据，返回模拟数据
    return _createMockModels(providerId);
  } catch (e) {
    // 出错时返回模拟数据
    return _createMockModels(providerId);
  }
});

// 创建模拟模型数据
List<Model> _createMockModels(String providerId) {
  if (providerId.contains('openai')) {
    return [
      Model(
        llmProviderId: providerId,
        id: 'gpt-4-turbo',
        ownedBy: 'OpenAI',
        enabled: true,
      ),
      Model(
        llmProviderId: providerId,
        id: 'gpt-4',
        ownedBy: 'OpenAI',
        enabled: true,
      ),
      Model(
        llmProviderId: providerId,
        id: 'gpt-3.5-turbo',
        ownedBy: 'OpenAI',
        enabled: true,
      ),
    ];
  } else if (providerId.contains('anthropic')) {
    return [
      Model(
        llmProviderId: providerId,
        id: 'claude-3-opus',
        ownedBy: 'Anthropic',
        enabled: true,
      ),
      Model(
        llmProviderId: providerId,
        id: 'claude-3-sonnet',
        ownedBy: 'Anthropic',
        enabled: true,
      ),
      Model(
        llmProviderId: providerId,
        id: 'claude-3-haiku',
        ownedBy: 'Anthropic',
        enabled: true,
      ),
    ];
  } else {
    return [
      Model(
        llmProviderId: providerId,
        id: 'default-model',
        ownedBy: '未知提供商',
        enabled: true,
      ),
    ];
  }
}

class ModelSettings extends ConsumerStatefulWidget {
  const ModelSettings({super.key});

  @override
  ConsumerState<ModelSettings> createState() => _ModelSettingsState();
}

class _ModelSettingsState extends ConsumerState<ModelSettings>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  Map<String, bool> _expandedProviders = {};
  // 存储本地修改的模型状态
  final Map<String, bool> _modelEnabledStates = {};
  // 存储本地修改的提供商状态
  final Map<String, bool> _providerEnabledStates = {};

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  // 获取模型的启用状态，优先使用本地状态
  bool getModelEnabledState(Model model) {
    final key = '${model.llmProviderId}_${model.id}';
    return _modelEnabledStates.containsKey(key) 
        ? _modelEnabledStates[key]! 
        : model.enabled;
  }

  // 获取提供商的启用状态，优先使用本地状态
  bool getProviderEnabledState(LlmProvider provider) {
    return _providerEnabledStates.containsKey(provider.id) 
        ? _providerEnabledStates[provider.id]! 
        : provider.enable;
  }

  // 更新模型启用状态
  void updateModelEnabledState(Model model, bool value) {
    final key = '${model.llmProviderId}_${model.id}';
    setState(() {
      _modelEnabledStates[key] = value;
    });
    
    // 异步更新数据库，不影响UI响应
    LlmStorage.updateModel(Model(
      llmProviderId: model.llmProviderId,
      id: model.id,
      object: model.object,
      ownedBy: model.ownedBy,
      enabled: value,
    )).then((success) {
      if (!success && mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('更新模型状态失败，请稍后重试')),
        );
      }
    }).catchError((e) {
      debugPrint('更新模型出错: $e');
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('更新模型状态出错: $e')),
        );
      }
    });
  }

  // 更新提供商启用状态
  void updateProviderEnabledState(LlmProvider provider, bool value) {
    setState(() {
      _providerEnabledStates[provider.id] = value;
    });
    
    // 异步更新数据库，不影响UI响应
    LlmStorage.updateProvider(LlmProvider(
      id: provider.id,
      name: provider.name,
      baseUrl: provider.baseUrl,
      apiKey: provider.apiKey,
      alias: provider.alias,
      description: provider.description,
      enable: value,
      checked: provider.checked,
    )).then((success) {
      if (!success && mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('更新供应商状态失败，请稍后重试')),
        );
      }
    }).catchError((e) {
      debugPrint('更新供应商出错: $e');
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('更新供应商状态出错: $e')),
        );
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(24, 24, 24, 0),
          child: Text(
            S.of(context).modelSettings,
            style: TextStyle(
              fontSize: FontSizeUtils.getHeadingSize(ref),
              fontWeight: FontWeight.bold,
            ),
          ),
        ),
        const SizedBox(height: 24),
        Container(
          margin: const EdgeInsets.symmetric(horizontal: 24),
          decoration: BoxDecoration(
            color: theme.colorScheme.surfaceContainerHighest,
            borderRadius: BorderRadius.circular(8),
          ),
          child: TabBar(
            controller: _tabController,
            labelColor: theme.colorScheme.primary,
            unselectedLabelColor: theme.colorScheme.onSurface,
            indicatorSize: TabBarIndicatorSize.tab,
            dividerColor: Colors.transparent,
            indicator: BoxDecoration(
              color: theme.colorScheme.surface,
              borderRadius: BorderRadius.circular(8),
              boxShadow: [
                BoxShadow(
                  color: Colors.black.withAlpha(13),
                  blurRadius: 4,
                  offset: const Offset(0, 2),
                ),
              ],
            ),
            splashBorderRadius: BorderRadius.circular(8),
            padding: const EdgeInsets.all(4),
            labelStyle: TextStyle(
              fontSize: FontSizeUtils.getBodySize(ref),
              fontWeight: FontWeight.w600,
            ),
            unselectedLabelStyle: TextStyle(
              fontSize: FontSizeUtils.getBodySize(ref),
              fontWeight: FontWeight.normal,
            ),
            tabs: [
              Tab(
                icon: const Icon(Icons.cloud),
                text: '模型供应商',
                iconMargin: const EdgeInsets.only(bottom: 4),
                height: 64,
              ),
              Tab(
                icon: const Icon(Icons.text_fields),
                text: '提示词',
                iconMargin: const EdgeInsets.only(bottom: 4),
                height: 64,
              ),
            ],
          ),
        ),
        const SizedBox(height: 16),
        Expanded(
          child: TabBarView(
            controller: _tabController,
            children: [
              _buildProvidersTab(),
              _buildPromptsTab()
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildProvidersTab() {
    return ref.watch(providersProvider).when(
      data: (providers) {
        return SingleChildScrollView(
          padding: const EdgeInsets.all(24),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start, 
            children: [
              for (final provider in providers)
                _buildProviderCard(provider),
            ],
          ),
        );
      },
      loading: () => const Center(child: CircularProgressIndicator()),
      error: (error, stack) => Center(child: Text('加载失败: $error')),
    );
  }

  Widget _buildProviderCard(LlmProvider provider) {
    final theme = Theme.of(context);
    final isEnabled = getProviderEnabledState(provider);

    return Card(
      margin: const EdgeInsets.only(bottom: 16),
      elevation: 0, // 去除阴影
      color: theme.brightness == Brightness.light 
          ? theme.colorScheme.surfaceContainerHighest.withValues(alpha: 0.4)
          : theme.colorScheme.surfaceContainerHighest.withValues(alpha: 0.6),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(8), // 调整圆角与tab一致
        // side: BorderSide(color: theme.colorScheme.outlineVariant.withOpacity(0.5)), // 添加边框替代阴影
      ),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Row(
          children: [
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Text(
                        provider.name,
                        style: TextStyle(
                          fontSize: FontSizeUtils.getSubheadingSize(ref),
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(width: 8),
                      if (!provider.checked)
                        Container(
                          padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                          decoration: BoxDecoration(
                            color: theme.colorScheme.errorContainer,
                            borderRadius: BorderRadius.circular(4),
                          ),
                          child: Text(
                            '未验证',
                            style: TextStyle(
                              fontSize: FontSizeUtils.getSmallSize(ref),
                              color: theme.colorScheme.onErrorContainer,
                            ),
                          ),
                        ),
                    ],
                  ),
                  const SizedBox(height: 4),
                  Text(
                    provider.baseUrl,
                    style: TextStyle(
                      fontSize: FontSizeUtils.getBodySize(ref),
                      color: theme.colorScheme.onSurfaceVariant,
                    ),
                  ),
                  if (provider.description != null && provider.description!.isNotEmpty)
                    Padding(
                      padding: const EdgeInsets.only(top: 4),
                      child: Text(
                        provider.description!,
                        style: TextStyle(
                          fontSize: FontSizeUtils.getSmallSize(ref),
                          color: theme.colorScheme.onSurfaceVariant,
                        ),
                      ),
                    ),
                ],
              ),
            ),
                IconButton(
                  icon: const Icon(Icons.list),
                  tooltip: '查看模型列表',
                  onPressed: () {
                    _showModelsDialog(provider);
                  },
                ),
                PopupMenuButton<String>(
                  icon: const Icon(Icons.more_vert),
                  onSelected: (value) {
                    if (value == 'edit') {
                      _showEditDialog(provider);
                    } else if (value == 'delete') {
                      _showDeleteDialog(provider);
                    }
                  },
                  itemBuilder: (context) => [
                    const PopupMenuItem<String>(
                      value: 'edit',
                      child: Row(
                        children: [
                          Icon(Icons.edit),
                          SizedBox(width: 8),
                          Text('编辑'),
                        ],
                      ),
                    ),
                    const PopupMenuItem<String>(
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
                // 启用按钮放在最右边，并调小
                Transform.scale(
                  scale: 0.8, // 调小开关大小
                  child: Switch(
                    value: isEnabled,
                    onChanged: (value) => updateProviderEnabledState(provider, value),
                  ),
                ),
              ],
            ),
          ),
    );
  }

  // 显示模型列表对话框
  void _showModelsDialog(LlmProvider provider) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('${provider.name} 模型列表'),
        content: SizedBox(
          width: 460, // 设置固定宽度，使对话框更窄
          child: ref.watch(modelsProvider(provider.id)).when(
            data: (models) {
              if (models.isEmpty) {
                return const Center(child: Text('暂无模型'));
              }
              
              return ListView.builder(
                shrinkWrap: true,
                itemCount: models.length,
                itemBuilder: (context, index) {
                  final model = models[index];
                  final isEnabled = getModelEnabledState(model);
                  
                  return ListTile(
                    title: Text(model.id),
                    subtitle: Text('提供者: ${model.ownedBy}'),
                    trailing: Transform.scale(
                      scale: 0.8,
                      child: Switch(
                        value: isEnabled,
                        onChanged: (value) {
                          updateModelEnabledState(model, value);
                          // 强制对话框内容刷新
                          setState(() {});
                        },
                      ),
                    ),
                  );
                },
              );
            },
            loading: () => const Center(child: CircularProgressIndicator()),
            error: (error, stack) => Center(child: Text('加载模型失败: $error')),
          ),
        ),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8), // 设置弹窗圆角为8
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('关闭'),
            style: TextButton.styleFrom(
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(8), // 设置按钮圆角为8
              ),
            ),
          ),
        ],
      ),
    );
  }

  void _showEditDialog(LlmProvider provider) {
    // TODO: 实现编辑对话框
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('编辑供应商'),
        content: const Text('此功能尚未实现'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('关闭'),
          ),
        ],
      ),
    );
  }

  void _showDeleteDialog(LlmProvider provider) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('删除供应商'),
        content: Text('确定要删除 ${provider.name} 吗？此操作不可恢复。'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('取消'),
          ),
          TextButton(
            onPressed: () async {
              Navigator.of(context).pop();
              final success = await LlmStorage.deleteProvider(provider.id);
              if (success) {
                // 刷新提供商列表
                ref.refresh(providersProvider);
              } else {
                if (mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    const SnackBar(content: Text('删除失败')),
                  );
                }
              }
            },
            child: const Text('删除', style: TextStyle(color: Colors.red)),
          ),
        ],
      ),
    );
  }

  Widget _buildPromptsTab() {
    return Center(
      child: Text(
        '提示词',
        style: TextStyle(
          fontSize: FontSizeUtils.getHeadingSize(ref),
          fontWeight: FontWeight.bold,
        ),
      ),
    );
  }
}
