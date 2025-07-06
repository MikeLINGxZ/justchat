import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:lemon_tea/generated/l10n.dart';
import 'package:lemon_tea/utils/setting/storage.dart';

class DataSettings extends ConsumerWidget {
  const DataSettings({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
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
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 24),

          _buildSection(
            context: context,
            ref: ref,
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
            context: context,
            ref: ref,
            title: '数据管理',
            children: [
              ListTile(
                title: const Text('清空所有数据'),
                subtitle: const Text('删除所有对话和设置'),
                trailing: const Icon(Icons.delete_outline, color: Colors.red),
                onTap: () {
                  _showClearDataDialog(context, ref);
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

  void _showClearDataDialog(BuildContext context, WidgetRef ref) {
    final settingsManager = ref.read(settingsManagerProvider.notifier);

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text(S.of(context).confirmDelete),
        content: const Text('确定要清除所有数据吗？此操作无法撤销。'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: Text(S.of(context).cancel),
          ),
          TextButton(
            onPressed: () async {
              await settingsManager.clearAllSettings();
              Navigator.of(context).pop();
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(content: Text('数据已清除')),
              );
            },
            style: TextButton.styleFrom(foregroundColor: Colors.red),
            child: Text(S.of(context).delete),
          ),
        ],
      ),
    );
  }
} 