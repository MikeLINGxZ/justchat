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
  }) {
    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(4.0),
        child: Padding(
          padding: const EdgeInsets.all(4.0),
          child: Icon(icon, size: 21, color: color),
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
      height: 50,
      padding: const EdgeInsets.symmetric(horizontal: 16.0),
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.surface,
        border: Border(
          bottom: BorderSide(
            color: Theme.of(context).dividerColor.withAlpha(77),
            width: 1.0,
          ),
        ),
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
                fontSize: FontSizeUtils.getSubheadingSize(ref),
                fontWeight: FontWeight.bold,
              ),
              overflow: TextOverflow.ellipsis,
            ),
          ),
          if (onAddTap != null) ...[
            _buildIconButton(
              icon: Icons.add_circle_outline,
              onTap: onAddTap!,
            ),
          ],
        ],
      ),
    );
  }
}
