import 'package:flutter/material.dart';

class ResizableDivider extends StatefulWidget {
  final Widget leftChild;
  final Widget rightChild;
  final double leftWidth; // 左侧组件的固定宽度
  final double minLeftWidth; // 左侧组件的最小宽度
  final double maxLeftWidth; // 左侧组件的最大宽度
  final double dividerWidth; // 分隔线的宽度
  final Color? dividerColor;
  final VoidCallback? onResize;

  const ResizableDivider({
    super.key,
    required this.leftChild,
    required this.rightChild,
    this.leftWidth = 500.0,
    this.minLeftWidth = 300.0,
    this.maxLeftWidth = 800.0,
    this.dividerWidth = 4.0,
    this.dividerColor,
    this.onResize,
  });

  @override
  State<ResizableDivider> createState() => _ResizableDividerState();
}

class _ResizableDividerState extends State<ResizableDivider> {
  late double _leftWidth;
  bool _isDragging = false;

  @override
  void initState() {
    super.initState();
    _leftWidth = widget.leftWidth;
  }

  @override
  Widget build(BuildContext context) {
    return LayoutBuilder(
      builder: (context, constraints) {
        final totalWidth = constraints.maxWidth;
        final rightWidth = totalWidth - _leftWidth - widget.dividerWidth;

        return Row(
          children: [
            // 左侧组件
            SizedBox(
              width: _leftWidth,
              child: widget.leftChild,
            ),
            // 可拖拽的分隔线
            GestureDetector(
              onPanStart: (details) {
                setState(() {
                  _isDragging = true;
                });
              },
              onPanUpdate: (details) {
                setState(() {
                  _leftWidth = (_leftWidth + details.delta.dx).clamp(
                    widget.minLeftWidth,
                    widget.maxLeftWidth,
                  );
                });
                widget.onResize?.call();
              },
              onPanEnd: (details) {
                setState(() {
                  _isDragging = false;
                });
              },
              child: MouseRegion(
                cursor: SystemMouseCursors.resizeLeftRight,
                child: Container(
                  width: widget.dividerWidth,
                  color: _isDragging
                      ? (widget.dividerColor ?? Theme.of(context).dividerColor).withOpacity(0.8)
                      : widget.dividerColor ?? Theme.of(context).dividerColor,
                  child: _isDragging
                      ? Container(
                          decoration: BoxDecoration(
                            color: Theme.of(context).colorScheme.primary.withOpacity(0.1),
                          ),
                        )
                      : null,
                ),
              ),
            ),
            // 右侧组件
            Expanded(
              child: widget.rightChild,
            ),
          ],
        );
      },
    );
  }
} 