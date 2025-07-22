import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:lemon_tea/generated/l10n.dart';

class TitleBarView extends ConsumerWidget {
  final String tag;
  final String title;
  final VoidCallback? onAddTap;

  const TitleBarView({
    super.key,
    this.tag = '',
    this.title = '',
    this.onAddTap,
  });

  Widget _buildIconButton({
    required IconData icon,
    required VoidCallback onTap,
    Color? color,
    String? title,
  }) {
    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(4.0),
        child: Padding(
          padding: const EdgeInsets.all(4.0),
          child: title == null || title.isEmpty
              ? Icon(icon, size: 21, color: color)
              : Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(icon, size: 21, color: color),
              const SizedBox(width: 4.0),
              Text(title),
            ],
          ),
        ),
      ),
    );
  }


  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // 如果没有提供tag和title，使用默认值
    final displayTag = tag.isEmpty ? S.of(context).conversation : tag;
    final displayTitle = title.isEmpty ? S.of(context).aiAssistant : title;

    return Container(
      padding: const EdgeInsets.fromLTRB(16, 8, 16, 8),
      decoration: BoxDecoration(
        // color: Theme.of(context).colorScheme.surface,
        // border: Border(
        //   bottom: BorderSide(
        //     color: Theme.of(context).dividerColor,
        //     width: 0.8,
        //   ),
        // ),
      ),
      child: Row(
        children: [
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
            decoration: BoxDecoration(
              color: Theme.of(context).colorScheme.primary.withAlpha(26),
              borderRadius: BorderRadius.circular(4),
            ),
            child: Text(
              displayTag,
              style: TextStyle(
                fontSize: FontSizeUtils.getSmallSize(ref),
                color: Theme.of(context).colorScheme.primary,
              ),
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Text(
              displayTitle,
              style: TextStyle(
                fontSize: FontSizeUtils.getBodySize(ref),
                fontWeight: FontWeight.w900,
              ),
              overflow: TextOverflow.ellipsis,
            ),
          ),
          if (onAddTap != null) ...[
            _buildIconButton(
              title: "新建对话",
              icon: Icons.add,
              onTap: onAddTap!,
            ),
          ],
        ],
      ),
    );
  }
}
