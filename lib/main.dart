import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/pages/home/home.dart';
import 'package:lemon_tea/utils/local_server/local_service_provider.dart';
import 'package:lemon_tea/utils/setting/manager.dart' as app_theme;
import 'package:lemon_tea/utils/setting/storage.dart';
import 'package:lemon_tea/utils/system.dart';
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

/// 启动CLI服务
Future<void> _startCliService(ProviderContainer container) async {
  try {
    // 使用Provider启动CLI服务
    final port = await container.read(cliServiceProvider.notifier).startService();
    
    if (port != null) {
      debugPrint('CLI服务已成功启动，端口: $port');
    } else {
      debugPrint('CLI服务启动失败');
    }
  } catch (e) {
    debugPrint('启动CLI服务时发生错误: $e');
  }
}

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // 创建ProviderContainer来初始化设置和服务
  final container = ProviderContainer();

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
    await _startCliService(container);
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

  // 定义浅色主题
  ThemeData _lightTheme() {
    return ThemeData(
      useMaterial3: true,
      colorScheme: ColorScheme.fromSeed(
        seedColor: const Color(0xFFFFD700), // 柠檬黄色
        brightness: Brightness.light,
      ),
      scaffoldBackgroundColor: const Color(0xFFF5F5F5),
      appBarTheme: const AppBarTheme(
        backgroundColor: Colors.white,
        elevation: 0,
      ),
      cardTheme: CardTheme(
        elevation: 0.5,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
        ),
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          elevation: 0,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(8),
          ),
        ),
      ),
      inputDecorationTheme: InputDecorationTheme(
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(8),
          borderSide: BorderSide.none,
        ),
        filled: true,
        fillColor: Colors.grey[100],
        contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      ),
    );
  }

  // 定义深色主题
  ThemeData _darkTheme() {
    return ThemeData(
      useMaterial3: true,
      colorScheme: ColorScheme.fromSeed(
        seedColor: const Color(0xFFFFD700), // 柠檬黄色
        brightness: Brightness.dark,
      ),
      scaffoldBackgroundColor: const Color(0xFF1E1E1E),
      appBarTheme: const AppBarTheme(
        backgroundColor: Color(0xFF2D2D2D),
        elevation: 0,
      ),
      cardTheme: CardTheme(
        elevation: 0.5,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
        ),
        color: const Color(0xFF2D2D2D),
      ),
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ElevatedButton.styleFrom(
          elevation: 0,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(8),
          ),
        ),
      ),
      inputDecorationTheme: InputDecorationTheme(
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(8),
          borderSide: BorderSide.none,
        ),
        filled: true,
        fillColor: const Color(0xFF3D3D3D),
        contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      ),
    );
  }

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
