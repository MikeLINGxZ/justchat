import 'package:flutter/material.dart';
import 'package:window_manager/window_manager.dart';

class WindowTitleBar extends StatefulWidget {
  final String title;
  const WindowTitleBar({super.key, required this.title});

  @override
  _WindowTitleBar createState() => _WindowTitleBar();
}

class _WindowTitleBar extends State<WindowTitleBar> with WindowListener {

  final double _titleBarHeight = 30;

  @override
  void initState() {
    super.initState();
    windowManager.addListener(this);
  }

  @override
  void dispose() {
    windowManager.removeListener(this);
    super.dispose();
  }

  Widget _buildTitleBar() {
    final colorScheme = Theme.of(context).colorScheme;
    return GestureDetector(
      behavior: HitTestBehavior.translucent,
      onPanStart: (details) {
        windowManager.startDragging();
      },
      child: Container(
        height: _titleBarHeight,
        decoration: BoxDecoration(
          border: Border(
            bottom: BorderSide(
              color: colorScheme.onSurface.withValues(alpha: 0.2),
              width: 1.0,
            ),
          ),
        ),
        child: Center(
          child: Row(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              Text(widget.title)
            ],
          ),
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Column(
        children: [
          _buildTitleBar(),
        ],
      ),
    );
  }
}