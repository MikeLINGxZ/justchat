import 'package:flutter/material.dart';
import 'package:lemon_tea/utils/style.dart';
import 'package:window_manager/window_manager.dart';
import 'package:lemon_tea/utils/cli/server/server.dart';

class WindowTitleBar extends StatefulWidget {
  final String title;
  const WindowTitleBar({super.key, required this.title});

  @override
  _WindowTitleBar createState() => _WindowTitleBar();
}

class _WindowTitleBar extends State<WindowTitleBar> with WindowListener {
  final double _titleBarHeight = 30;
  final Server _server = Server();
  bool _isServerRunning = false;
  int? _serverPort;
  bool _isMaximized = false;

  @override
  void initState() {
    super.initState();
    windowManager.addListener(this);
    _initializeServerStatus();
    _checkWindowState();
  }

  @override
  void dispose() {
    windowManager.removeListener(this);
    super.dispose();
  }

  @override
  void onWindowMaximize() {
    setState(() {
      _isMaximized = true;
    });
  }

  @override
  void onWindowUnmaximize() {
    setState(() {
      _isMaximized = false;
    });
  }

  void _checkWindowState() async {
    final isMaximized = await windowManager.isMaximized();
    setState(() {
      _isMaximized = isMaximized;
    });
  }

  void _toggleMaximize() async {
    if (_isMaximized) {
      await windowManager.unmaximize();
    } else {
      await windowManager.maximize();
    }
  }

  void _initializeServerStatus() {
    // 初始化服务状态
    setState(() {
      _isServerRunning = _server.isRunning;
      _serverPort = _server.port;
    });

    // 监听服务端口变化
    _server.portStream.listen((port) {
      if (mounted) {
        setState(() {
          _isServerRunning = port != null;
          _serverPort = port;
        });
      }
    });
  }

  Widget _buildTitleBar() {
    return GestureDetector(
      behavior: HitTestBehavior.translucent,
      onPanStart: (details) {
        windowManager.startDragging();
      },
      onDoubleTap: _toggleMaximize,
      child: Container(
        color: Style.titleBarBackground(context),
        height: _titleBarHeight,
        child: Row(
          children: [
            // 右侧占位，保持对称
            const SizedBox(width: 32),
            // 中间标题
            Expanded(
              child: Center(
                child: Text(widget.title),
              ),
            ),
            // 左侧服务状态指示器
            Padding(
              padding: const EdgeInsets.only(left: 12.0,right: 12.0),
              child: Tooltip(
                message: _isServerRunning
                    ? '服务正在运行 (端口: $_serverPort)'
                    : '服务未运行',
                child: Container(
                  width: 8,
                  height: 8,
                  decoration: BoxDecoration(
                    shape: BoxShape.circle,
                    color: _isServerRunning ? Colors.green : Colors.red,
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: _titleBarHeight,
      child: Scaffold(body: Column(children: [_buildTitleBar()]),backgroundColor: Theme.of(context).appBarTheme.backgroundColor),
    );
  }
}
