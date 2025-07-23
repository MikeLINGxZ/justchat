

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/utils/llm/models/message.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';

class MessageToolbar extends ConsumerWidget {
  final void Function(Message)? onCopy;
  final void Function(Message)? onCopyPlainText;
  final void Function(Message)? onRegenerate;
  final void Function(Message)? onDelete;
  final Message? message;
  final bool isVisible;

  const MessageToolbar({
    super.key,
    this.message,
    this.onCopy,
    this.onCopyPlainText,
    this.onRegenerate,
    this.onDelete,
    this.isVisible = true,
  });

  MessageToolbar setMessage(Message msg, {bool? visible}) {
    return MessageToolbar(
      message: msg,
      onCopy: onCopy,
      onCopyPlainText: onCopyPlainText,
      onDelete: onDelete,
      onRegenerate: onRegenerate,
      isVisible: visible ?? isVisible,
    );
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // 如果没有消息，返回空容器（模板状态）
    if (message == null) {
      return const SizedBox.shrink();
    }

    final isDark = Theme.of(context).brightness == Brightness.dark;
    
    return AnimatedOpacity(
      opacity: isVisible ? 1.0 : 0.0,
      duration: const Duration(milliseconds: 200),
      child: Container(
        // margin: const EdgeInsets.only(top: 12),
        padding: const EdgeInsets.symmetric(horizontal: 0, vertical: 4),

        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            _ToolbarButton(
              icon: Icons.content_copy_rounded,
              tooltip: '复制内容',
              onPressed: isVisible ? () => onCopy?.call(message!) : null,
              ref: ref,
            ),
            _ToolbarButton(
              icon: Icons.text_fields_rounded,
              tooltip: '复制纯文本',
              onPressed: isVisible ? () => onCopyPlainText?.call(message!) : null,
              ref: ref,
            ),
            _ToolbarButton(
              icon: Icons.refresh_rounded,
              tooltip: '重新生成',
              onPressed: isVisible ? () => onRegenerate?.call(message!) : null,
              ref: ref,
            ),
            _ToolbarButton(
              icon: Icons.delete_outline_rounded,
              tooltip: '删除消息',
              onPressed: isVisible ? () => onDelete?.call(message!) : null,
              ref: ref,
              isDestructive: true,
            ),
          ],
        ),
      ),
    );
  }
}

class _ToolbarButton extends StatefulWidget {
  final IconData icon;
  final String tooltip;
  final VoidCallback? onPressed;
  final WidgetRef ref;
  final bool isDestructive;

  const _ToolbarButton({
    required this.icon,
    required this.tooltip,
    required this.onPressed,
    required this.ref,
    this.isDestructive = false,
  });

  @override
  State<_ToolbarButton> createState() => _ToolbarButtonState();
}

class _ToolbarButtonState extends State<_ToolbarButton>
    with SingleTickerProviderStateMixin {
  late AnimationController _animationController;
  late Animation<double> _scaleAnimation;
  bool _isHovered = false;

  @override
  void initState() {
    super.initState();
    _animationController = AnimationController(
      duration: const Duration(milliseconds: 150),
      vsync: this,
    );
    _scaleAnimation = Tween<double>(
      begin: 1.0,
      end: 0.95,
    ).animate(CurvedAnimation(
      parent: _animationController,
      curve: Curves.easeInOut,
    ));
  }

  @override
  void dispose() {
    _animationController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final iconColor = widget.isDestructive
        ? Colors.red.shade400
        : isDark
            ? Colors.white.withOpacity(0.8)
            : Colors.black.withOpacity(0.7);

    return Tooltip(
      message: widget.tooltip,
      textStyle: TextStyle(
        fontSize: FontSizeUtils.getCaptionSize(widget.ref),
        color: Colors.white,
      ),
      decoration: BoxDecoration(
        color: Colors.black.withOpacity(0.8),
        borderRadius: BorderRadius.circular(6),
      ),
      child: MouseRegion(
        onEnter: (_) => setState(() => _isHovered = true),
        onExit: (_) => setState(() => _isHovered = false),
        child: GestureDetector(
          onTapDown: (_) => _animationController.forward(),
          onTapUp: (_) => _animationController.reverse(),
          onTapCancel: () => _animationController.reverse(),
          child: AnimatedBuilder(
            animation: _scaleAnimation,
            builder: (context, child) {
              return Transform.scale(
                scale: _scaleAnimation.value,
                child: AnimatedContainer(
                  duration: const Duration(milliseconds: 200),
                  width: FontSizeUtils.getSubheadingSize(widget.ref) * 1.8,
                  height: FontSizeUtils.getSubheadingSize(widget.ref) * 1.8,
                  decoration: BoxDecoration(
                    color: Colors.transparent,
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: IconButton(
                    icon: Icon(
                      widget.icon,
                      size: FontSizeUtils.getSubheadingSize(widget.ref),
                      color: _isHovered
                          ? (widget.isDestructive
                              ? Colors.red.shade500
                              : isDark
                                  ? Colors.white.withOpacity(0.9)
                                  : Colors.black.withOpacity(0.8))
                          : iconColor,
                    ),
                    onPressed: widget.onPressed,
                    padding: EdgeInsets.zero,
                    constraints: const BoxConstraints(),
                    splashRadius: FontSizeUtils.getSubheadingSize(widget.ref),
                    hoverColor: Colors.transparent,
                    highlightColor: Colors.transparent,
                    splashColor: Colors.transparent,
                  ),
                ),
              );
            },
          ),
        ),
      ),
    );
  }
}