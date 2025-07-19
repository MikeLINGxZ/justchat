import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/rpc/service.pb.dart';
import 'package:lemon_tea/storage/llm_storage.dart';
import 'package:lemon_tea/utils/cli/client/client.dart';
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
  Map<String, bool> _expandedProviders = {};
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
              onPressed: () => _showAddProviderDialog(),
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
                        _showEditDialog(provider);
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

  // 显示模型列表对话框
  void _showModelsDialog(LlmProvider provider) {
    // 确保模型数据已预加载
    _preloadModels(provider.id);
    
    // 创建本地状态副本，用于UI显示
    final Map<String, bool> localModelStates = {};
    
    // 预先加载模型状态
    ref.read(modelsProvider(provider.id)).whenData((models) {
      for (final model in models) {
        final key = '${model.llmProviderId}_${model.id}';
        localModelStates[key] = getModelEnabledState(model);
      }
    });
    
    showDialog(
      context: context,
      builder: (context) => StatefulBuilder(
        builder: (context, dialogSetState) {
          // 使用Consumer直接访问模型数据，避免loading状态
          return AlertDialog(
            title: Text(
              '${provider.name} 模型列表',
              style: TextStyle(
                fontSize: FontSizeUtils.getSubheadingSize(ref),
                fontWeight: FontWeight.bold,
              ),
            ),
            content: SizedBox(
              width: 460, // 设置固定宽度，使对话框更窄
              height: 400, // 设置固定高度，确保对话框不会过大
              child: Consumer(
                builder: (context, ref, child) {
                  // 强制刷新模型数据
                  final modelsAsync = ref.watch(modelsProvider(provider.id));
                  
                  return modelsAsync.when(
                    data: (models) {
                      if (models.isEmpty) {
                        return Center(
                          child: Text(
                            '暂无模型',
                            style: TextStyle(
                              fontSize: FontSizeUtils.getBodySize(ref),
                            ),
                          ),
                        );
                      }
                      
                      return ListView.builder(
                        shrinkWrap: false, // 不收缩，允许滚动
                        physics: const AlwaysScrollableScrollPhysics(), // 始终可滚动
                        itemCount: models.length,
                        itemBuilder: (context, index) {
                          final model = models[index];
                          final key = '${model.llmProviderId}_${model.id}';
                          // 确保本地状态存在
                          if (!localModelStates.containsKey(key)) {
                            localModelStates[key] = getModelEnabledState(model);
                          }
                          
                          return ListTile(
                            title: Text(
                              model.id,
                              style: TextStyle(
                                fontSize: FontSizeUtils.getBodySize(ref),
                              ),
                            ),
                            subtitle: Text(
                              '提供者: ${model.ownedBy}',
                              style: TextStyle(
                                fontSize: FontSizeUtils.getSmallSize(ref),
                                color: Theme.of(context).colorScheme.onSurfaceVariant,
                              ),
                            ),
                            trailing: SizedBox(
                              height: 48, // 固定高度确保垂直居中
                              child: Row(
                                mainAxisSize: MainAxisSize.min,
                                mainAxisAlignment: MainAxisAlignment.center, // 水平居中
                                crossAxisAlignment: CrossAxisAlignment.center, // 垂直居中
                                children: [
                                  // 自定义标签放在编辑按钮前面
                                  if (model.isCustom)
                                    Container(
                                      margin: const EdgeInsets.only(right: 8),
                                      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                                      decoration: BoxDecoration(
                                        color: Theme.of(context).colorScheme.primaryContainer,
                                        borderRadius: BorderRadius.circular(4),
                                      ),
                                      child: Text(
                                        '自定义',
                                        style: TextStyle(
                                          fontSize: FontSizeUtils.getSmallSize(ref) - 1,
                                          color: Theme.of(context).colorScheme.onPrimaryContainer,
                                        ),
                                      ),
                                    ),
                                  Transform.scale(
                                    scale: 0.8,
                                    child: Switch(
                                      value: localModelStates[key]!,
                                      onChanged: (value) {
                                        // 更新本地状态和UI
                                        dialogSetState(() {
                                          localModelStates[key] = value;
                                        });
                                        // 更新数据库
                                        updateModelEnabledState(model, value);
                                      },
                                    ),
                                  ),
                                  // 所有模型都显示更多操作按钮
                                  SizedBox(
                                    height: 40,
                                    width: 40,
                                    child: PopupMenuButton<String>(
                                      padding: EdgeInsets.zero,
                                      icon: const Icon(Icons.more_vert, size: 20),
                                      tooltip: '更多操作',
                                      offset: const Offset(0, 10),
                                      position: PopupMenuPosition.under,
                                      itemBuilder: (context) => [
                                        PopupMenuItem<String>(
                                          value: 'edit',
                                          enabled: model.isCustom,
                                          child: Row(
                                            children: [
                                              Icon(Icons.edit, 
                                                size: 18, 
                                                color: model.isCustom 
                                                    ? null 
                                                    : Theme.of(context).colorScheme.onSurface.withOpacity(0.38),
                                              ),
                                              const SizedBox(width: 8),
                                              Text(
                                                '编辑',
                                                style: TextStyle(
                                                  fontSize: FontSizeUtils.getBodySize(ref),
                                                  color: model.isCustom 
                                                      ? null 
                                                      : Theme.of(context).colorScheme.onSurface.withOpacity(0.38),
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
                                      onSelected: (value) {
                                        if (value == 'edit' && model.isCustom) {
                                          Navigator.of(context).pop();
                                          _showEditModelDialog(model);
                                        } else if (value == 'delete') {
                                          _showDeleteModelDialog(model, context);
                                        }
                                      },
                                    ),
                                  ),
                                ],
                              ),
                            ),
                            contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
                          );
                        },
                      );
                    },
                    loading: () => const Center(child: CircularProgressIndicator()),
                    error: (error, stack) => Center(
                      child: Text(
                        '加载模型失败: $error',
                        style: TextStyle(
                          fontSize: FontSizeUtils.getBodySize(ref),
                          color: Theme.of(context).colorScheme.error,
                        ),
                      ),
                    ),
                  );
                },
              ),
            ),
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(8), // 设置弹窗圆角为8
            ),
            actions: [
              TextButton(
                onPressed: () {
                  // 不再关闭当前对话框，直接打开添加模型对话框
                  _showAddModelDialog(provider.id);
                },
                style: TextButton.styleFrom(
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(8),
                  ),
                ),
                child: Text(
                  '添加模型',
                  style: TextStyle(
                    fontSize: FontSizeUtils.getBodySize(ref),
                  ),
                ),
              ),
              TextButton(
                onPressed: () => Navigator.of(context).pop(),
                style: TextButton.styleFrom(
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(8),
                  ),
                ),
                child: Text(
                  '关闭',
                  style: TextStyle(
                    fontSize: FontSizeUtils.getBodySize(ref),
                  ),
                ),
              ),
            ],
          );
        }
      ),
    );
  }

  // 显示编辑模型对话框
  void _showEditModelDialog(Model model) {
    // 由于在UI中已经禁用了非自定义模型的编辑选项，这里不再需要额外检查
    
    final TextEditingController modelIdController = TextEditingController(text: model.id);
    final TextEditingController ownedByController = TextEditingController(text: model.ownedBy);
    bool isEnabled = model.enabled;

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(
          '编辑模型',
          style: TextStyle(
            fontSize: FontSizeUtils.getSubheadingSize(ref),
            fontWeight: FontWeight.bold,
          ),
        ),
        content: SizedBox(
          width: 400,
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              TextField(
                controller: modelIdController,
                decoration: InputDecoration(
                  labelText: '模型ID',
                  labelStyle: TextStyle(
                    fontSize: FontSizeUtils.getBodySize(ref),
                  ),
                ),
                style: TextStyle(
                  fontSize: FontSizeUtils.getBodySize(ref),
                ),
              ),
              const SizedBox(height: 16),
              TextField(
                controller: ownedByController,
                decoration: InputDecoration(
                  labelText: '提供者',
                  labelStyle: TextStyle(
                    fontSize: FontSizeUtils.getBodySize(ref),
                  ),
                ),
                style: TextStyle(
                  fontSize: FontSizeUtils.getBodySize(ref),
                ),
              ),
              const SizedBox(height: 16),
              Row(
                children: [
                  Text(
                    '启用状态',
                    style: TextStyle(
                      fontSize: FontSizeUtils.getBodySize(ref),
                    ),
                  ),
                  const Spacer(),
                  StatefulBuilder(
                    builder: (BuildContext context, StateSetter setState) {
                      return Switch(
                        value: isEnabled,
                        onChanged: (value) {
                          setState(() {
                            isEnabled = value;
                          });
                        },
                      );
                    },
                  ),
                ],
              ),
            ],
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
              final updatedModel = Model(
                llmProviderId: model.llmProviderId,
                id: modelIdController.text.trim(),
                object: model.object,
                ownedBy: ownedByController.text.trim(),
                enabled: isEnabled,
                isCustom: model.isCustom,
              );
              
              final success = await LlmStorage.updateModel(updatedModel);
              Navigator.of(context).pop();
              
              if (success) {
                // 刷新模型列表
                ref.refresh(modelsProvider(model.llmProviderId));
              } else {
                if (mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      content: Text(
                        '更新模型失败',
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
              '保存',
              style: TextStyle(
                fontSize: FontSizeUtils.getBodySize(ref),
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
        title: Text(
          '编辑供应商',
          style: TextStyle(
            fontSize: FontSizeUtils.getSubheadingSize(ref),
            fontWeight: FontWeight.bold,
          ),
        ),
        content: Text(
          '此功能尚未实现',
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
              '关闭',
              style: TextStyle(
                fontSize: FontSizeUtils.getBodySize(ref),
              ),
            ),
          ),
        ],
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

  // 显示添加模型对话框
  void _showAddModelDialog(String providerId) {
    final TextEditingController modelIdController = TextEditingController();
    final TextEditingController modelNameController = TextEditingController();
    final TextEditingController ownedByController = TextEditingController();
    bool isEnabled = true;

    showDialog(
      context: context,
      builder: (dialogContext) => AlertDialog(  // 使用不同的context变量名
        title: Text(
          '添加自定义模型',
          style: TextStyle(
            fontSize: FontSizeUtils.getSubheadingSize(ref),
            fontWeight: FontWeight.bold,
          ),
        ),
        content: SizedBox(
          width: 400,
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              TextField(
                controller: modelIdController,
                decoration: InputDecoration(
                  labelText: '模型ID *',
                  hintText: '例如: gpt-4-turbo',
                  labelStyle: TextStyle(
                    fontSize: FontSizeUtils.getBodySize(ref),
                  ),
                ),
                style: TextStyle(
                  fontSize: FontSizeUtils.getBodySize(ref),
                ),
              ),
              const SizedBox(height: 16),
              TextField(
                controller: modelNameController,
                decoration: InputDecoration(
                  labelText: '模型名称 *',
                  hintText: '例如: GPT-4 Turbo',
                  labelStyle: TextStyle(
                    fontSize: FontSizeUtils.getBodySize(ref),
                  ),
                ),
                style: TextStyle(
                  fontSize: FontSizeUtils.getBodySize(ref),
                ),
              ),
              const SizedBox(height: 16),
              TextField(
                controller: ownedByController,
                decoration: InputDecoration(
                  labelText: '提供者 *',
                  hintText: '例如: OpenAI',
                  labelStyle: TextStyle(
                    fontSize: FontSizeUtils.getBodySize(ref),
                  ),
                ),
                style: TextStyle(
                  fontSize: FontSizeUtils.getBodySize(ref),
                ),
              ),
              const SizedBox(height: 16),
              Row(
                children: [
                  Text(
                    '启用状态',
                    style: TextStyle(
                      fontSize: FontSizeUtils.getBodySize(ref),
                    ),
                  ),
                  const Spacer(),
                  StatefulBuilder(
                    builder: (BuildContext context, StateSetter setState) {
                      return Switch(
                        value: isEnabled,
                        onChanged: (value) {
                          setState(() {
                            isEnabled = value;
                          });
                        },
                      );
                    },
                  ),
                ],
              ),
            ],
          ),
        ),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(dialogContext).pop(),  // 使用dialogContext
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
              // 验证输入
              final modelId = modelIdController.text.trim();
              final modelName = modelNameController.text.trim();
              final ownedBy = ownedByController.text.trim();
              
              if (modelId.isEmpty || modelName.isEmpty || ownedBy.isEmpty) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text(
                      '请填写所有必填字段',
                      style: TextStyle(
                        fontSize: FontSizeUtils.getBodySize(ref),
                      ),
                    ),
                  ),
                );
                return;
              }
              
              // 创建模型对象
              final newModel = Model(
                llmProviderId: providerId,
                id: modelId,
                ownedBy: ownedBy,
                enabled: isEnabled,
                isCustom: true,
              );
              
              // 添加模型到数据库，使用自定义方法处理name字段
              final success = await _addModelWithName(newModel, modelName);
              Navigator.of(dialogContext).pop();  // 使用dialogContext
              
              if (success) {
                // 刷新模型列表
                ref.refresh(modelsProvider(providerId));
                
                // 延迟一点时间后重新打开模型列表对话框以显示更新后的列表
                Future.delayed(const Duration(milliseconds: 300), () async {
                  // 获取提供商对象
                  final provider = await LlmStorage.getProviderById(providerId);
                  if (provider != null && mounted) {
                    _showModelsDialog(provider);
                  }
                });
              } else {
                if (mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      content: Text(
                        '添加模型失败，可能模型ID已存在',
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
              '添加',
              style: TextStyle(
                fontSize: FontSizeUtils.getBodySize(ref),
              ),
            ),
          ),
        ],
      ),
    );
  }

  // 添加模型到数据库，处理name字段
  Future<bool> _addModelWithName(Model model, String name) async {
    try {
      // 获取模型的Map数据
      final modelMap = model.toMap();
      // 添加name字段
      modelMap['name'] = name;
      
      // 获取当前最大seq_id并设置新模型的seq_id为最大值+1
      final maxSeqId = await LlmStorage.getMaxModelSeqId(model.llmProviderId);
      modelMap['seq_id'] = maxSeqId + 1;
      
      // 插入数据库
      final result = await LlmStorage.addModelWithCustomFields(modelMap);
      return result;
    } catch (e) {
      debugPrint('添加模型失败: $e');
      return false;
    }
  }

  // 显示删除模型对话框
  void _showDeleteModelDialog(Model model, BuildContext parentContext) {
    showDialog(
      context: parentContext, // 使用父对话框的context，而不是全局context
      builder: (dialogContext) => AlertDialog(
        title: Text(
          '删除模型',
          style: TextStyle(
            fontSize: FontSizeUtils.getSubheadingSize(ref),
            fontWeight: FontWeight.bold,
          ),
        ),
        content: Text(
          '确定要删除模型 ${model.id} 吗？此操作不可恢复。',
          style: TextStyle(
            fontSize: FontSizeUtils.getBodySize(ref),
          ),
        ),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(dialogContext).pop(),
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
              final success = await LlmStorage.deleteModel(model.id);
              Navigator.of(dialogContext).pop();
              
              if (success) {
                // 刷新模型列表
                ref.refresh(modelsProvider(model.llmProviderId));
                
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text(
                      '模型已删除',
                      style: TextStyle(
                        fontSize: FontSizeUtils.getBodySize(ref),
                      ),
                    ),
                    duration: const Duration(seconds: 2),
                  ),
                );
              } else {
                if (mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      content: Text(
                        '删除模型失败',
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
              foregroundColor: Colors.red,
            ),
            child: Text(
              '删除',
              style: TextStyle(
                fontSize: FontSizeUtils.getBodySize(ref),
              ),
            ),
          ),
        ],
      ),
    );
  }
  
  // 显示添加供应商对话框
  void _showAddProviderDialog() {
    final TextEditingController nameController = TextEditingController();
    final TextEditingController baseUrlController = TextEditingController();
    final TextEditingController apiKeyController = TextEditingController();
    final TextEditingController aliasController = TextEditingController();
    final TextEditingController descriptionController = TextEditingController();
    bool isEnabled = true;

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(
          '添加模型供应商',
          style: TextStyle(
            fontSize: FontSizeUtils.getSubheadingSize(ref),
            fontWeight: FontWeight.bold,
          ),
        ),
        content: SizedBox(
          width: 400,
          child: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                TextField(
                  controller: nameController,
                  decoration: InputDecoration(
                    labelText: '供应商名称 *',
                    hintText: '例如: OpenAI',
                    labelStyle: TextStyle(
                      fontSize: FontSizeUtils.getBodySize(ref),
                    ),
                  ),
                  style: TextStyle(
                    fontSize: FontSizeUtils.getBodySize(ref),
                  ),
                ),
                const SizedBox(height: 16),
                TextField(
                  controller: baseUrlController,
                  decoration: InputDecoration(
                    labelText: '基础URL *',
                    hintText: '例如: https://api.openai.com/v1',
                    labelStyle: TextStyle(
                      fontSize: FontSizeUtils.getBodySize(ref),
                    ),
                  ),
                  style: TextStyle(
                    fontSize: FontSizeUtils.getBodySize(ref),
                  ),
                ),
                const SizedBox(height: 16),
                TextField(
                  controller: apiKeyController,
                  decoration: InputDecoration(
                    labelText: 'API密钥 *',
                    hintText: '您的API密钥',
                    labelStyle: TextStyle(
                      fontSize: FontSizeUtils.getBodySize(ref),
                    ),
                  ),
                  style: TextStyle(
                    fontSize: FontSizeUtils.getBodySize(ref),
                  ),
                  obscureText: true,
                ),
                const SizedBox(height: 16),
                TextField(
                  controller: aliasController,
                  decoration: InputDecoration(
                    labelText: '别名',
                    hintText: '可选',
                    labelStyle: TextStyle(
                      fontSize: FontSizeUtils.getBodySize(ref),
                    ),
                  ),
                  style: TextStyle(
                    fontSize: FontSizeUtils.getBodySize(ref),
                  ),
                ),
                const SizedBox(height: 16),
                TextField(
                  controller: descriptionController,
                  decoration: InputDecoration(
                    labelText: '描述',
                    hintText: '可选',
                    labelStyle: TextStyle(
                      fontSize: FontSizeUtils.getBodySize(ref),
                    ),
                  ),
                  style: TextStyle(
                    fontSize: FontSizeUtils.getBodySize(ref),
                  ),
                  maxLines: 2,
                ),
                const SizedBox(height: 16),
                Row(
                  children: [
                    Text(
                      '启用状态',
                      style: TextStyle(
                        fontSize: FontSizeUtils.getBodySize(ref),
                      ),
                    ),
                    const Spacer(),
                    StatefulBuilder(
                      builder: (BuildContext context, StateSetter setState) {
                        return Switch(
                          value: isEnabled,
                          onChanged: (value) {
                            setState(() {
                              isEnabled = value;
                            });
                          },
                        );
                      },
                    ),
                  ],
                ),
              ],
            ),
          ),
        ),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8),
        ),
        actions: [
          TextButton(
            onPressed:  () async {
              debugPrint("xxx");
              try {
                final request = ModelsRequest(name: nameController.text,apiKey: apiKeyController.text,baseUrl: baseUrlController.text);
                ModelsResponse response = await Client().stub!.models(request);
                debugPrint(response.models.toString());
              }catch (e) {
                debugPrint(e.toString());
              }
            },
            style: TextButton.styleFrom(
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(8),
              ),
            ),
            child: Text(
              '验证',
              style: TextStyle(
                fontSize: FontSizeUtils.getBodySize(ref),
              ),
            ),
          ),
          TextButton(
            onPressed: () async {
              // 验证输入
              final name = nameController.text.trim();
              final baseUrl = baseUrlController.text.trim();
              final apiKey = apiKeyController.text.trim();
              final alias = aliasController.text.trim();
              final description = descriptionController.text.trim();

              if (name.isEmpty || baseUrl.isEmpty || apiKey.isEmpty) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text(
                      '请填写所有必填字段',
                      style: TextStyle(
                        fontSize: FontSizeUtils.getBodySize(ref),
                      ),
                    ),
                  ),
                );
                return;
              }

              // 创建供应商对象
              final newProvider = LlmProvider(
                id: '${name.toLowerCase().replaceAll(' ', '_')}_${DateTime.now().millisecondsSinceEpoch}',
                name: name,
                baseUrl: baseUrl,
                apiKey: apiKey,
                alias: alias.isEmpty ? null : alias,
                description: description.isEmpty ? null : description,
                enable: isEnabled,
                checked: false, // 新添加的供应商默认未验证
              );

              // 添加供应商到数据库
              final success = await LlmStorage.addProvider(newProvider);
              Navigator.of(context).pop();

              if (success) {
                // 刷新供应商列表
                ref.refresh(providersProvider);

                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text(
                      '供应商已添加',
                      style: TextStyle(
                        fontSize: FontSizeUtils.getBodySize(ref),
                      ),
                    ),
                    duration: const Duration(seconds: 2),
                  ),
                );
              } else {
                if (mounted) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      content: Text(
                        '添加供应商失败',
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
              '添加',
              style: TextStyle(
                fontSize: FontSizeUtils.getBodySize(ref),
              ),
            ),
          ),
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
        ],
      ),
    );
  }
}
