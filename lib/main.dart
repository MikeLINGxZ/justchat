import 'package:flutter/material.dart';
import 'package:flutter/material.dart';
import 'package:flutter_localizations/flutter_localizations.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/controls/window_title_bar.dart';
import 'package:lemon_tea/pages/home/home.dart';
import 'package:lemon_tea/utils/system.dart';
import 'package:window_manager/window_manager.dart';
import 'generated/l10n.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

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
  }
  runApp(ProviderScope(child: const LemonTea()));
}

class LemonTea extends ConsumerWidget {
  const LemonTea({super.key});

  // This widget is the root of your application.
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return MaterialApp(
      localizationsDelegates: [
        S.delegate,
        GlobalMaterialLocalizations.delegate,
        GlobalWidgetsLocalizations.delegate,
        GlobalCupertinoLocalizations.delegate,
      ],
      supportedLocales: S.delegate.supportedLocales,
      title: "Lemon Tea",
      locale: Locale("zh"),
      theme: ThemeData.light(),
      darkTheme: ThemeData.dark(),
      themeMode: ThemeMode.system,
      // home: Scaffold(backgroundColor: Colors.white, body: ViewWidget(LoginPage())),
      home: _content(),
    );
  }

  Widget _content() {
    return System.isDesktop
        ? Column(
      children: [
        // if (System.isDesktop && Service().getAccountCount() == 0) ...[SizedBox(height: 30, child: WindowTitleBar(title: "Teamail"))],
        Expanded(
          child: Navigator(
            onGenerateRoute: (settings) {
              return MaterialPageRoute(
                builder: (_) => HomePage(),
                settings: settings,
              );
            },
          ),
        ),
      ],
    )
        : HomePage();
  }
}
