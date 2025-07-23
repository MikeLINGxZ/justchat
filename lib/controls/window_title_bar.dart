import 'package:flutter/material.dart' hide ThemeMode;
import 'package:lemon_tea/utils/style.dart';
import 'package:window_manager/window_manager.dart';
import 'package:lemon_tea/utils/cli/server/server.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/utils/setting/theme_manager.dart';

class WindowTitleBar extends ConsumerStatefulWidget {
  final String title;
  const WindowTitleBar({super.key, required this.title});

  @override
  ConsumerState<WindowTitleBar> createState() => _WindowTitleBar();
}

class _WindowTitleBar extends ConsumerState<WindowTitleBar> with WindowListener {
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

  /// 检查是否处于调试模式
  bool _isDebugMode() {
    bool inDebugMode = false;
    assert(() {
      inDebugMode = true;
      return true;
    }());
    return inDebugMode;
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

  /// 构建主题切换按钮
  Widget _buildThemeButton() {
    final themeMode = ref.watch(themeManagerProvider);
    final themeManager = ref.read(themeManagerProvider.notifier);
    
    // 根据主题模式显示不同的图标
    IconData getThemeIcon() {
      switch (themeMode) {
        case ThemeMode.light:
          return Icons.light_mode;
        case ThemeMode.dark:
          return Icons.dark_mode;
        case ThemeMode.system:
          return Icons.settings_brightness;
      }
    }
    
    // 获取主题模式提示文本
    String getThemeTooltip() {
      switch (themeMode) {
        case ThemeMode.light:
          return '当前：浅色模式\n点击切换到深色模式';
        case ThemeMode.dark:
          return '当前：深色模式\n点击切换到跟随系统';
        case ThemeMode.system:
          return '当前：跟随系统\n点击切换到浅色模式';
      }
    }
    
    return Tooltip(
      message: getThemeTooltip(),
      child: GestureDetector(
        onTap: () {
          themeManager.toggleTheme();
        },
        child: Container(
          width: 24,
          height: 24,
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(4),
            color: Colors.transparent,
          ),
          child: Icon(
            getThemeIcon(),
            size: 16,
            color: Theme.of(context).iconTheme.color?.withOpacity(0.8),
          ),
        ),
      ),
    );
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
            // 右侧控件区域
            Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                // 调试模式下显示主题切换按钮
                if (_isDebugMode()) ...[
                  _buildThemeButton(),
                  const SizedBox(width: 8),
                ],
                // 服务状态指示器
                Padding(
                  padding: const EdgeInsets.only(left: 12.0, right: 12.0),
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
