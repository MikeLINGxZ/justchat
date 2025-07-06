import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';

class AboutSettings extends ConsumerWidget {
  const AboutSettings({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            '关于',
            style: TextStyle(
              fontSize: FontSizeUtils.getHeadingSize(ref),
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 24),

          _buildSection(
            context: context,
            ref: ref,
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
            context: context,
            ref: ref,
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
} 