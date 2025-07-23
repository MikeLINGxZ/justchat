

import 'package:flutter/material.dart';

class MessageToolbar extends StatelessWidget {
  final VoidCallback? onCopy;
  final VoidCallback? onCopyPlainText;
  final VoidCallback? onRegenerate;
  final VoidCallback? onDelete;

  const MessageToolbar({
    super.key,
    this.onCopy,
    this.onCopyPlainText,
    this.onRegenerate,
    this.onDelete,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.only(top: 10),
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
      decoration: BoxDecoration(
        color: Theme.of(context).brightness == Brightness.dark
            ? Colors.white.withOpacity(0.06)
            : Colors.black.withOpacity(0.04),
        borderRadius: BorderRadius.circular(10),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.04),
            blurRadius: 6,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Tooltip(
            message: '复制',
            child: IconButton(
              icon: const Icon(Icons.copy_rounded, size: 20),
              onPressed: onCopy,
            ),
          ),
          Tooltip(
            message: '重新生成',
            child: IconButton(
              icon: const Icon(Icons.refresh_rounded, size: 20),
              onPressed: onRegenerate,
            ),
          ),
          Tooltip(
            message: '删除',
            child: IconButton(
              icon: const Icon(Icons.delete_outline_rounded, size: 20, color: Colors.redAccent),
              onPressed: onDelete,
            ),
          ),
        ],
      ),
    );
  }
}