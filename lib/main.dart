import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/pages/home/home.dart';
import 'package:lemon_tea/utils/setting/manager.dart' as app_theme;
import 'package:lemon_tea/utils/setting/storage.dart';
import 'package:lemon_tea/utils/system.dart';
import 'package:lemon_tea/utils/cli/server/server.dart';
import 'package:lemon_tea/utils/cli/client/client.dart';
import 'package:lemon_tea/storage/sqlite_util.dart';
import 'package:window_manager/window_manager.dart';
import 'generated/l10n.dart';

/// 初始化应用设置
Future<void> _initializeAppSettings(ProviderContainer container) async {
  // 等待设置加载完成
  await container.read(settingsManagerProvider.notifier).loadSettings();
  
  // 等待主题和字体大小设置加载完成
  await container.read(app_theme.themeManagerProvider.notifier).loadThemeMode();
  await container.read(app_theme.fontSizeModeProvider.notifier).loadFontSizeMode();
  
  // 根据设置初始化语言
  final settings = container.read(settingsManagerProvider);
  if (settings.language == 'English') {
    S.load(const Locale('en', 'US'));
  } else {
    S.load(const Locale('zh', 'CN'));
  }
}

/// 初始化数据库
Future<void> _initializeDatabase() async {
  try {
    // 获取数据库实例，这将触发数据库初始化
    await SqliteUtil.instance.database;
    debugPrint('数据库初始化成功');
  } catch (e) {
    debugPrint('数据库初始化失败: $e');
  }
}

/// 启动CLI服务
Future<void> _startLemonTeaService(ProviderContainer container) async {
  // 启动CLI服务（非阻塞模式）
  final localServer = Server();
  
  // 使用Future.microtask在下一个事件循环中启动服务，避免阻塞UI
  Future.microtask(() async {
    try {
      debugPrint('正在启动CLI服务...');
      final port = await localServer.startService();
      
      if (port != null) {
        debugPrint('CLI服务已启动，端口: $port');
        
        // 初始化RPC客户端
        final client = Client();
        final initialized = await client.init();
        if (initialized) {
          debugPrint('RPC客户端初始化成功');
        } else {
          debugPrint('RPC客户端初始化失败');
        }
      } else {
        debugPrint('CLI服务启动失败 - 请检查端口是否被占用');
      }
    } catch (e) {
      debugPrint('启动CLI服务时发生错误: $e');
    }
  });
}

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // 创建ProviderContainer来初始化设置和服务
  final container = ProviderContainer();

  // 初始化数据库
  await _initializeDatabase();

  if (System.isDesktop) {
    // 初始化 WindowManager
    await windowManager.ensureInitialized();

    // 设置窗口大小和位置
    WindowOptions windowOptions = const WindowOptions(
      size: Size(1300, 800),
      center: true,
      skipTaskbar: false,
      titleBarStyle: TitleBarStyle.hidden,
      minimumSize: Size(520, 520),
    );
    windowManager.waitUntilReadyToShow(windowOptions, () async {
      await windowManager.show(); // 显示窗口
      await windowManager.focus(); // 聚焦窗口
    });
    
    // 启动CLI服务
    await _startLemonTeaService(container);
  }
  
  // 初始化应用设置
  await _initializeAppSettings(container);
  
  runApp(ProviderScope(
    parent: container,
    child: const LemonTea(),
  ));
}

class LemonTea extends ConsumerWidget {
  const LemonTea({super.key});

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    // 获取当前主题模式
    final themeMode = ref.watch(app_theme.themeManagerProvider);
    // 获取用户设置
    final settings = ref.watch(settingsManagerProvider);
    
    // 转换为 Flutter 的 ThemeMode
    ThemeMode flutterThemeMode;
    switch (themeMode) {
      case app_theme.ThemeMode.light:
        flutterThemeMode = ThemeMode.light;
        break;
      case app_theme.ThemeMode.dark:
        flutterThemeMode = ThemeMode.dark;
        break;
      case app_theme.ThemeMode.system:
        flutterThemeMode = ThemeMode.system;
        break;
    }
    
    return MaterialApp(
      localizationsDelegates: [
        S.delegate,
        GlobalMaterialLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
      ],
      supportedLocales: S.delegate.supportedLocales,
      title: "Lemon Tea",
      locale: settings.language == 'English' ? const Locale('en', 'US') : const Locale('zh', 'CN'),
      theme: ThemeData.light(),
      darkTheme: ThemeData.dark(),
      themeMode: flutterThemeMode,
      // home: Scaffold(backgroundColor: Colors.white, body: ViewWidget(LoginPage())),
      home: _content(),
    );
  }

  Widget _content() {
    return System.isDesktop
        ? Navigator(
          onGenerateRoute: (settings) {
            return MaterialPageRoute(
              builder: (_) => HomePage(),
              settings: settings,
            );
          },
        )
        : HomePage();
  }
}
