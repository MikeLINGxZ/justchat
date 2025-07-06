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

          _buildSection(
            context: context,
            ref: ref,
            title: '模型供应商',
            children: [
              ListTile(
                title: const Text('添加供应商'),
                subtitle: const Text('添加新的AI模型供应商'),
                trailing: const Icon(Icons.add),
                onTap: () {
                  _showProviderDialog(context, ref);
                },
              ),
              if (providers.isNotEmpty) ...[
                const Divider(height: 1),
                ...providers.map((provider) => _buildProviderTile(context, ref, provider)),
              ],
            ],
          ),

          if (selectedProvider != null) ...[
            const SizedBox(height: 24),
            _buildSection(
              context: context,
              ref: ref,
              title: '${selectedProvider.displayName} 的模型',
              children: [
                if (selectedProvider.models != null &&
                    selectedProvider.models!.isNotEmpty)
                  ...selectedProvider.models!.map(
                    (model) => _buildModelTile(context, ref, model, selectedProvider),
                  )
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
            context: context,
            ref: ref,
            title: '当前选择',
            children: [
              ListTile(
                title: Text(selectedProvider?.displayName ?? '未选择供应商'),
                subtitle: Text(selectedModel?.displayName ?? '未选择模型'),
                trailing: const Icon(Icons.settings),
                onTap: () {
                  _showModelSelectionDialog(context, ref);
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

  Widget _buildProviderTile(BuildContext context, WidgetRef ref, LlmProvider provider) {
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
                  _showProviderDialog(context, ref, provider: provider);
                  break;
                case 'delete':
                  _showDeleteProviderDialog(context, ref, provider);
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

  Widget _buildModelTile(BuildContext context, WidgetRef ref, Model model, LlmProvider provider) {
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

  void _showProviderDialog(BuildContext context, WidgetRef ref, {LlmProvider? provider}) {
    showDialog(
      context: context,
      builder: (context) => ProviderDialog(provider: provider),
    );
  }

  void _showModelSelectionDialog(BuildContext context, WidgetRef ref) {
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
              const Text(
                '选择供应商:',
                style: TextStyle(fontWeight: FontWeight.bold),
              ),
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
                const Text(
                  '选择模型:',
                  style: TextStyle(fontWeight: FontWeight.bold),
                ),
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

  void _showDeleteProviderDialog(BuildContext context, WidgetRef ref, LlmProvider provider) {
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