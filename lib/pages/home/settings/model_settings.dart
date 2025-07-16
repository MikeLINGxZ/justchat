import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/storage/llm_storage.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:lemon_tea/generated/l10n.dart';
import 'package:lemon_tea/models/llm_provider.dart';
import 'package:lemon_tea/models/model.dart';

class ModelSettings extends ConsumerStatefulWidget {
  const ModelSettings({super.key});

  @override
  ConsumerState<ModelSettings> createState() => _ModelSettingsState();
}

class _ModelSettingsState extends ConsumerState<ModelSettings>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  Map<String, bool> _expandedProviders = {};

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
              FutureBuilder<List<LlmProvider>>(
                future: LlmStorage.getAllProviders(),
                builder: (context, snapshot) {
                  if (snapshot.connectionState == ConnectionState.waiting) {
                    return const Center(child: CircularProgressIndicator());
                  }
                  
                  if (snapshot.hasError) {
                    return Center(child: Text('加载失败: ${snapshot.error}'));
                  }
                  
                  final providers = snapshot.data ?? [];
                  
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
              ),
              _buildPromptsTab()
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildProviderCard(LlmProvider provider) {
    final theme = Theme.of(context);
    final isExpanded = _expandedProviders[provider.id] ?? false;

    return Card(
      margin: const EdgeInsets.only(bottom: 16),
      elevation: 0, // 去除阴影
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(8), // 调整圆角与tab一致
        side: BorderSide(color: theme.colorScheme.outlineVariant.withOpacity(0.5)), // 添加边框替代阴影
      ),
      child: Column(
        children: [
          Padding(
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
                    setState(() {
                      _expandedProviders[provider.id] = !isExpanded;
                    });
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
                    value: provider.enable,
                    onChanged: (value) async {
                      final updatedProvider = LlmProvider(
                        id: provider.id,
                        name: provider.name,
                        baseUrl: provider.baseUrl,
                        apiKey: provider.apiKey,
                        alias: provider.alias,
                        description: provider.description,
                        enable: value,
                        checked: provider.checked,
                      );
                      
                      final success = await LlmStorage.updateProvider(updatedProvider);
                      if (success) {
                        setState(() {});
                      } else {
                        if (mounted) {
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(content: Text('更新状态失败')),
                          );
                        }
                      }
                    },
                  ),
                ),
              ],
            ),
          ),
          if (isExpanded)
            FutureBuilder<List<Model>>(
              future: LlmStorage.getModelsByProviderId(provider.id),
              builder: (context, snapshot) {
                if (snapshot.connectionState == ConnectionState.waiting) {
                  return const Padding(
                    padding: EdgeInsets.all(16.0),
                    child: Center(child: CircularProgressIndicator()),
                  );
                }
                
                if (snapshot.hasError) {
                  return Padding(
                    padding: const EdgeInsets.all(16.0),
                    child: Center(child: Text('加载模型失败: ${snapshot.error}')),
                  );
                }
                
                final models = snapshot.data ?? [];
                
                if (models.isEmpty) {
                  return const Padding(
                    padding: EdgeInsets.all(16.0),
                    child: Center(child: Text('暂无模型')),
                  );
                }
                
                return Container(
                  decoration: BoxDecoration(
                    color: theme.colorScheme.surfaceContainerLowest,
                    borderRadius: const BorderRadius.only(
                      bottomLeft: Radius.circular(12),
                      bottomRight: Radius.circular(12),
                    ),
                  ),
                  child: ListView.builder(
                    shrinkWrap: true,
                    physics: const NeverScrollableScrollPhysics(),
                    itemCount: models.length,
                    itemBuilder: (context, index) {
                      final model = models[index];
                      return ListTile(
                        title: Text(model.id),
                        subtitle: Text('提供者: ${model.ownedBy}'),
                        trailing: Transform.scale(
                          scale: 0.8, // 调小开关大小
                          child: Switch(
                            value: model.enabled,
                            onChanged: (value) async {
                              final updatedModel = Model(
                                llm_provider_id: model.llm_provider_id,
                                id: model.id,
                                object: model.object,
                                ownedBy: model.ownedBy,
                                enabled: value,
                              );
                              
                              final success = await LlmStorage.updateModel(updatedModel);
                              if (success) {
                                setState(() {});
                              } else {
                                if (mounted) {
                                  ScaffoldMessenger.of(context).showSnackBar(
                                    const SnackBar(content: Text('更新模型状态失败')),
                                  );
                                }
                              }
                            },
                          ),
                        ),
                      );
                    },
                  ),
                );
              },
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
                setState(() {});
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
