import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/storage/llm_storage.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:lemon_tea/generated/l10n.dart';
import 'package:lemon_tea/models/llm_provider.dart';
import 'package:lemon_tea/models/model.dart';
import 'dialogs/add_provider_dialog.dart';
import 'dialogs/edit_provider_dialog.dart';
import 'dialogs/models_dialog.dart';

// 创建提供商列表的Provider
final providersProvider = FutureProvider<List<LlmProvider>>((ref) async {
  return LlmStorage.getAllProviders();
});

// 创建模型列表的Provider，接受提供商ID作为参数
final modelsProvider = FutureProvider.family<List<Model>, String>((ref, providerId) async {
  try {
    final models = await LlmStorage.getModelsByProviderId(providerId);
    return models; // 直接返回数据库中的模型列表，可能为空
  } catch (e) {
    debugPrint('获取模型列表出错: $e');
    return []; // 出错时返回空列表
  }
});

class ModelSettings extends ConsumerStatefulWidget {
  const ModelSettings({super.key});

  @override
  ConsumerState<ModelSettings> createState() => _ModelSettingsState();
}

class _ModelSettingsState extends ConsumerState<ModelSettings>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  // 存储本地修改的模型状态
  final Map<String, bool> _modelEnabledStates = {};
  // 存储本地修改的提供商状态
  final Map<String, bool> _providerEnabledStates = {};
  // 用于跟踪已预加载的提供商ID
  final Set<String> _preloadedProviderIds = {};

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
    
    // 预加载提供商数据
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _preloadProviders();
    });
  }
  
  // 预加载提供商和模型数据
  void _preloadProviders() async {
    final providers = await LlmStorage.getAllProviders();
    if (providers.isNotEmpty && mounted) {
      // 预加载第一个提供商的模型
      final firstProvider = providers.first;
      _preloadModels(firstProvider.id);
    }
  }
  
  // 预加载指定提供商的模型
  void _preloadModels(String providerId) {
    if (_preloadedProviderIds.contains(providerId)) return;
    
    _preloadedProviderIds.add(providerId);
    ref.read(modelsProvider(providerId));
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
          SnackBar(
            content: Text(
              '更新模型状态失败，请稍后重试',
              style: TextStyle(
                fontSize: FontSizeUtils.getBodySize(ref),
              ),
            ),
          ),
        );
      }
    }).catchError((e) {
      debugPrint('更新模型出错: $e');
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
              '更新模型状态出错: $e',
              style: TextStyle(
                fontSize: FontSizeUtils.getBodySize(ref),
              ),
            ),
          ),
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
          SnackBar(
            content: Text(
              '更新供应商状态失败，请稍后重试',
              style: TextStyle(
                fontSize: FontSizeUtils.getBodySize(ref),
              ),
            ),
          ),
        );
      }
    }).catchError((e) {
      debugPrint('更新供应商出错: $e');
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
              '更新供应商状态出错: $e',
              style: TextStyle(
                fontSize: FontSizeUtils.getBodySize(ref),
              ),
            ),
          ),
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
                child: Text(
                  '模型供应商',
                  style: TextStyle(
                    fontSize: FontSizeUtils.getBodySize(ref),
                  ),
                ),
                iconMargin: const EdgeInsets.only(bottom: 4),
                height: 64,
              ),
              Tab(
                icon: const Icon(Icons.text_fields),
                child: Text(
                  '提示词',
                  style: TextStyle(
                    fontSize: FontSizeUtils.getBodySize(ref),
                  ),
                ),
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
    return Stack(
      children: [
        ref.watch(providersProvider).when(
          data: (providers) {
            return SingleChildScrollView(
              // 添加顶部边距，为添加按钮留出空间
              padding: const EdgeInsets.fromLTRB(24, 70, 24, 24),
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
          error: (error, stack) => Center(
            child: Text(
              '加载失败: $error',
              style: TextStyle(
                fontSize: FontSizeUtils.getBodySize(ref),
                color: Theme.of(context).colorScheme.error,
              ),
            ),
          ),
        ),
        // 添加按钮放在左上角
        Align(
          alignment: Alignment.topLeft,
          child: Padding(
            padding: const EdgeInsets.only(top: 16, left: 24),
            child: ElevatedButton.icon(
              onPressed: () => showAddProviderDialog(context, ref),
              icon: const Icon(Icons.add),
              label: Text(
                '添加供应商',
                style: TextStyle(
                  fontSize: FontSizeUtils.getBodySize(ref),
                ),
              ),
              style: ElevatedButton.styleFrom(
                elevation: 0, // 去掉阴影
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(8),
                ),
              ),
            ),
          ),
        ),
      ],
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
                showModelsDialog(context, ref, provider, _preloadModels, getModelEnabledState, updateModelEnabledState);
              },
            ),
            SizedBox(
              height: 40,
              width: 40,
              child: PopupMenuButton<String>(
                padding: EdgeInsets.zero,
                icon: const Icon(Icons.more_vert, size: 20),
                tooltip: '更多操作',
                offset: const Offset(0, 10),
                position: PopupMenuPosition.under,
                onSelected: (value) {
                  if (value == 'edit') {
                    showEditProviderDialog(context, ref, provider);
                  } else if (value == 'delete') {
                    _showDeleteDialog(provider);
                  }
                },
                itemBuilder: (context) => [
                  PopupMenuItem<String>(
                    value: 'edit',
                    child: Row(
                      children: [
                        const Icon(Icons.edit, size: 18),
                        const SizedBox(width: 8),
                        Text(
                          '编辑',
                          style: TextStyle(
                            fontSize: FontSizeUtils.getBodySize(ref),
                          ),
                        ),
                      ],
                    ),
                  ),
                  PopupMenuItem<String>(
                    value: 'delete',
                    child: Row(
                      children: [
                        const Icon(Icons.delete, color: Colors.red, size: 18),
                        const SizedBox(width: 8),
                        Text(
                          '删除', 
                          style: TextStyle(
                            fontSize: FontSizeUtils.getBodySize(ref),
                            color: Colors.red,
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
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

  void _showDeleteDialog(LlmProvider provider) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(
          '删除供应商',
          style: TextStyle(
            fontSize: FontSizeUtils.getSubheadingSize(ref),
            fontWeight: FontWeight.bold,
          ),
        ),
        content: Text(
          '确定要删除 ${provider.name} 吗？此操作不可恢复。',
          style: TextStyle(
            fontSize: FontSizeUtils.getBodySize(ref),
          ),
        ),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            style: TextButton.styleFrom(
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(8),
              ),
            ),
            child: Text(
              '取消',
              style: TextStyle(
                fontSize: FontSizeUtils.getBodySize(ref),
              ),
            ),
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
                    SnackBar(
                      content: Text(
                        '删除失败',
                        style: TextStyle(
                          fontSize: FontSizeUtils.getBodySize(ref),
                        ),
                      ),
                    ),
                  );
                }
              }
            },
            style: TextButton.styleFrom(
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(8),
              ),
            ),
            child: Text(
              '删除',
              style: TextStyle(
                fontSize: FontSizeUtils.getBodySize(ref),
                color: Colors.red,
              ),
            ),
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