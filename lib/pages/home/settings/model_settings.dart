import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:lemon_tea/generated/l10n.dart';
import 'package:lemon_tea/utils/setting/provider_manager.dart';
import 'package:lemon_tea/pages/home/settings/provider_dialog.dart';
import 'package:lemon_tea/models/llm_provider.dart';
import 'package:lemon_tea/models/model.dart';

class ModelSettings extends ConsumerWidget {
  const ModelSettings({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
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
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 24),

          _buildSectionWithAction(
            context: context,
            ref: ref,
            title: '模型供应商',
            action: IconButton(
              icon: const Icon(Icons.add),
              onPressed: () {
                _showProviderDialog(context, ref);
              },
              tooltip: '添加供应商',
            ),
            children: [
              if (providers.isNotEmpty)
                ...providers.map((provider) => _buildProviderTile(context, ref, provider))
              else
                ListTile(
                  title: Text(
                    '暂无供应商',
                    style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref)),
                  ),
                  subtitle: Text(
                    '点击右上角添加按钮添加新的AI模型供应商',
                    style: TextStyle(fontSize: FontSizeUtils.getSmallSize(ref)),
                  ),
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

  Widget _buildSectionWithAction({
    required BuildContext context,
    required WidgetRef ref,
    required String title,
    required Widget action,
    required List<Widget> children,
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(
              title,
              style: TextStyle(
                fontSize: FontSizeUtils.getSubheadingSize(ref),
                fontWeight: FontWeight.w600,
              ),
            ),
            action,
          ],
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

  Widget _buildProviderTile(BuildContext context, WidgetRef ref, LlmProvider provider) {
    // 获取最新的供应商数据，确保模型列表是最新的
    final providers = ref.watch(providerManagerProvider);
    final currentProvider = providers.firstWhere(
      (p) => p.name == provider.name,
      orElse: () => provider,
    );
    
    // 计算启用的模型数量
    final enabledModelsCount = currentProvider.models?.where((m) => m.enabled).length ?? 0;
    final totalModelsCount = currentProvider.models?.length ?? 0;
    
    return ListTile(
      title: Text(
        currentProvider.displayName,
        style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref)),
      ),
      subtitle: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            currentProvider.baseUrl,
            style: TextStyle(fontSize: FontSizeUtils.getSmallSize(ref)),
          ),
          if (currentProvider.description != null) 
            Text(
              currentProvider.description!,
              style: TextStyle(fontSize: FontSizeUtils.getSmallSize(ref)),
            ),
          Text(
            '模型数量: $enabledModelsCount / $totalModelsCount',
            style: TextStyle(fontSize: FontSizeUtils.getSmallSize(ref)),
          ),
        ],
      ),
      trailing: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
            decoration: BoxDecoration(
              color: currentProvider.hasApiKey ? Colors.green.withOpacity(0.1) : Colors.orange.withOpacity(0.1),
              borderRadius: BorderRadius.circular(12),
              border: Border.all(
                color: currentProvider.hasApiKey ? Colors.green : Colors.orange,
                width: 1,
              ),
            ),
            child: Text(
              currentProvider.hasApiKey ? '已配置' : '未配置',
              style: TextStyle(
                color: currentProvider.hasApiKey ? Colors.green : Colors.orange,
                fontSize: FontSizeUtils.getSmallSize(ref),
                fontWeight: FontWeight.w500,
              ),
            ),
          ),
          const SizedBox(width: 8),
          IconButton(
            icon: const Icon(Icons.list),
            onPressed: () {
              _showModelsDialog(context, ref, currentProvider);
            },
            tooltip: '查看模型列表',
          ),
          const SizedBox(width: 8),
          PopupMenuButton<String>(
            onSelected: (value) {
              switch (value) {
                case 'edit':
                  _showProviderDialog(context, ref, provider: currentProvider);
                  break;
                case 'delete':
                  _showDeleteProviderDialog(context, ref, currentProvider);
                  break;
                case 'refresh':
                  _testConnection(context, ref, currentProvider);
                  break;
              }
            },
            itemBuilder: (context) => [
              PopupMenuItem(
                value: 'edit',
                child: Row(
                  children: [
                    const Icon(Icons.edit),
                    const SizedBox(width: 8),
                    Text(
                      '编辑',
                      style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref)),
                    ),
                  ],
                ),
              ),
              PopupMenuItem(
                value: 'refresh',
                child: Row(
                  children: [
                    const Icon(Icons.refresh),
                    const SizedBox(width: 8),
                    Text(
                      '测试连接',
                      style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref)),
                    ),
                  ],
                ),
              ),
              PopupMenuItem(
                value: 'delete',
                child: Row(
                  children: [
                    const Icon(Icons.delete, color: Colors.red),
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
        ],
      ),
      onTap: () {
        ref.read(selectedProviderProvider.notifier).state = currentProvider;
        ref.read(selectedModelProvider.notifier).state = null;
      },
    );
  }

  Widget _buildModelTile(BuildContext context, WidgetRef ref, Model model, LlmProvider provider) {
    final selectedModel = ref.watch(selectedModelProvider);
    final isSelected = selectedModel?.id == model.id;

    return ListTile(
      title: Text(
        model.displayName,
        style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref)),
      ),
      subtitle: Text(
        '类型: ${model.object}',
        style: TextStyle(fontSize: FontSizeUtils.getSmallSize(ref)),
      ),
      trailing: Icon(
        isSelected ? Icons.check_circle : Icons.circle_outlined,
        color: isSelected ? Colors.green : null,
      ),
      onTap: () {
        ref.read(selectedModelProvider.notifier).state = model;
      },
    );
  }

  void _showProviderDialog(BuildContext context, WidgetRef ref, {LlmProvider? provider}) {
    showDialog(
      context: context,
      builder: (context) => ProviderDialog(provider: provider),
    );
  }

  void _showModelsDialog(BuildContext context, WidgetRef ref, LlmProvider provider) {
    showDialog(
      context: context,
      builder: (context) => Consumer(
        builder: (context, ref, child) {
          // 获取最新的供应商数据，确保模型列表是最新的
          final currentProviders = ref.watch(providerManagerProvider);
          final currentProvider = currentProviders.firstWhere(
            (p) => p.name == provider.name,
            orElse: () => provider,
          );

          // 调试信息
          print('显示模型列表对话框');
          print('供应商名称: ${currentProvider.name}');
          print('模型数量: ${currentProvider.models?.length ?? 0}');
          if (currentProvider.models != null) {
            print('模型列表: ${currentProvider.models!.map((m) => m.displayName).toList()}');
          }

          return AlertDialog(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(4),
            ),
            title: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Expanded(
                  child: Text(
                    '${currentProvider.displayName} 的模型列表',
                    style: TextStyle(fontSize: FontSizeUtils.getHeadingSize(ref)),
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
                IconButton(
                  icon: const Icon(Icons.refresh),
                  tooltip: '刷新模型列表',
                  onPressed: () {
                    // 强制刷新对话框
                    Navigator.of(context).pop();
                    _showModelsDialog(context, ref, currentProvider);
                  },
                ),
              ],
            ),
            content: SizedBox(
              width: 400,
              height: 300,
              child: currentProvider.models != null && currentProvider.models!.isNotEmpty
                  ? ListView.builder(
                      itemCount: currentProvider.models!.length,
                      itemBuilder: (context, index) {
                        final model = currentProvider.models![index];
                        return ListTile(
                          title: Text(
                            model.displayName,
                            style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref)),
                          ),
                          trailing: Transform.scale(
                            scale: 0.8,
                            child: Switch(
                              value: model.enabled,
                              onChanged: (value) async {
                                // 更新模型的启用状态
                                final updatedModel = model.copyWith(enabled: value);
                                final updatedModels = List<Model>.from(currentProvider.models!);
                                updatedModels[index] = updatedModel;
                                
                                final updatedProvider = currentProvider.copyWith(models: updatedModels);
                                
                                // 保存更新后的供应商信息
                                final providerManager = ref.read(providerManagerProvider.notifier);
                                await providerManager.updateProvider(currentProvider.name, updatedProvider);
                                
                                // 显示提示信息
                                ScaffoldMessenger.of(context).showSnackBar(
                                  SnackBar(
                                    content: Text(
                                      '${model.displayName} ${value ? '已启用' : '已禁用'}',
                                      style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref)),
                                    ),
                                    duration: const Duration(seconds: 2),
                                  ),
                                );
                              },
                            ),
                          ),
                        );
                      },
                    )
                  : Center(
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          const Icon(Icons.info_outline, size: 48, color: Colors.grey),
                          const SizedBox(height: 16),
                          Text(
                            '该供应商暂无可用模型',
                            style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref)),
                          ),
                          const SizedBox(height: 8),
                          ElevatedButton.icon(
                            icon: const Icon(Icons.refresh),
                            label: const Text('测试连接获取模型'),
                            onPressed: () {
                              Navigator.of(context).pop();
                              _showProviderDialog(context, ref, provider: currentProvider);
                            },
                          ),
                        ],
                      ),
                    ),
            ),
            actions: [
              TextButton(
                onPressed: () => Navigator.of(context).pop(),
                child: Text(
                  S.of(context).cancel,
                  style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref)),
                ),
              ),
            ],
          );
        },
      ),
    );
  }

  void _showModelSelectionDialog(BuildContext context, WidgetRef ref) {
    final providers = ref.read(providerManagerProvider);
    final selectedProvider = ref.read(selectedProviderProvider);
    final selectedModel = ref.read(selectedModelProvider);

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(4),
        ),
        title: Text(
          '选择模型',
          style: TextStyle(fontSize: FontSizeUtils.getHeadingSize(ref)),
        ),
        content: SizedBox(
          width: 400,
          height: 300,
          child: Column(
            children: [
              Text(
                '选择供应商:',
                style: TextStyle(
                  fontSize: FontSizeUtils.getSubheadingSize(ref),
                  fontWeight: FontWeight.bold,
                ),
              ),
              const SizedBox(height: 8),
              Expanded(
                child: ListView.builder(
                  itemCount: providers.length,
                  itemBuilder: (context, index) {
                    final provider = providers[index];
                    final isSelected = selectedProvider?.name == provider.name;

                    return ListTile(
                      title: Text(
                        provider.displayName,
                        style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref)),
                      ),
                      subtitle: Text(
                        provider.baseUrl,
                        style: TextStyle(fontSize: FontSizeUtils.getSmallSize(ref)),
                      ),
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
                Text(
                  '选择模型:',
                  style: TextStyle(
                    fontSize: FontSizeUtils.getSubheadingSize(ref),
                    fontWeight: FontWeight.bold,
                  ),
                ),
                const SizedBox(height: 8),
                Expanded(
                  child: ListView.builder(
                    itemCount: selectedProvider.models?.length ?? 0,
                    itemBuilder: (context, index) {
                      final model = selectedProvider.models![index];
                      final isSelected = selectedModel?.id == model.id;

                      return ListTile(
                        title: Text(
                          model.displayName,
                          style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref)),
                        ),
                        subtitle: Text(
                          '类型: ${model.object}',
                          style: TextStyle(fontSize: FontSizeUtils.getSmallSize(ref)),
                        ),
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
            child: Text(
              S.of(context).cancel,
              style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref)),
            ),
          ),
        ],
      ),
    );
  }

  void _showDeleteProviderDialog(BuildContext context, WidgetRef ref, LlmProvider provider) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(4),
        ),
        title: Text(
          '确认删除',
          style: TextStyle(fontSize: FontSizeUtils.getHeadingSize(ref)),
        ),
        content: Text(
          '确定要删除供应商 "${provider.displayName}" 吗？此操作无法撤销。',
          style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref)),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: Text(
              S.of(context).cancel,
              style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref)),
            ),
          ),
          TextButton(
            onPressed: () async {
              try {
                final providerManager = ref.read(providerManagerProvider.notifier);
                await providerManager.deleteProvider(provider.name);
                Navigator.of(context).pop();
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text(
                      '供应商删除成功',
                      style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref)),
                    ),
                  ),
                );
              } catch (e) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text(
                      '删除失败：${e.toString()}',
                      style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref)),
                    ),
                    backgroundColor: Colors.red,
                  ),
                );
              }
            },
            style: TextButton.styleFrom(foregroundColor: Colors.red),
            child: Text(
              S.of(context).delete,
              style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref)),
            ),
          ),
        ],
      ),
    );
  }
  
  /// 直接在列表中测试连接
  Future<void> _testConnection(BuildContext context, WidgetRef ref, LlmProvider provider) async {
    final providerManager = ref.read(providerManagerProvider.notifier);
    
    // 显示加载对话框
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (context) => const AlertDialog(
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            CircularProgressIndicator(),
            SizedBox(height: 16),
            Text('正在测试连接并获取模型列表...'),
          ],
        ),
      ),
    );
    
    try {
      final result = await providerManager.testProviderConnection(provider);
      
      // 关闭加载对话框
      Navigator.of(context).pop();
      
      if (result['success']) {
        final models = result['models'] as List<Model>;
        
        // 显示保存模型对话框
        if (models.isNotEmpty) {
          _showSaveModelsDialog(context, ref, provider, models);
        } else {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text('连接测试成功，但未获取到模型'),
              backgroundColor: Colors.orange,
            ),
          );
        }
      } else {
        final error = result['error'] as String;
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('连接测试失败：$error'),
            backgroundColor: Colors.red,
          ),
        );
      }
    } catch (e) {
      // 关闭加载对话框
      Navigator.of(context).pop();
      
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('连接测试失败：${e.toString()}'),
          backgroundColor: Colors.red,
        ),
      );
    }
  }
  
  /// 显示保存模型对话框
  void _showSaveModelsDialog(BuildContext context, WidgetRef ref, LlmProvider provider, List<Model> models) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(4),
        ),
        title: const Text('保存模型信息'),
        content: SizedBox(
          width: 400,
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text('是否将获取到的模型信息保存到此供应商？'),
              const SizedBox(height: 16),
              SizedBox(
                height: 200,
                child: ListView.builder(
                  shrinkWrap: true,
                  itemCount: models.length > 5 ? 5 : models.length,
                  itemBuilder: (context, index) {
                    final model = models[index];
                    return ListTile(
                      dense: true,
                      title: Text(model.displayName),
                      subtitle: Text('类型: ${model.object}'),
                    );
                  },
                ),
              ),
              if (models.length > 5)
                Padding(
                  padding: const EdgeInsets.only(top: 8.0),
                  child: Text('... 还有 ${models.length - 5} 个模型未显示'),
                ),
            ],
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () async {
              Navigator.of(context).pop();
              
              // 更新供应商的模型列表
              final updatedProvider = provider.copyWith(models: models);
              
              // 保存供应商信息
              try {
                final providerManager = ref.read(providerManagerProvider.notifier);
                await providerManager.updateProvider(provider.name, updatedProvider);
                
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(
                    content: Text('模型信息已保存'),
                    backgroundColor: Colors.green,
                  ),
                );
              } catch (e) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text('保存模型信息失败：${e.toString()}'),
                    backgroundColor: Colors.red,
                  ),
                );
              }
            },
            child: const Text('保存'),
          ),
        ],
      ),
    );
  }
} 