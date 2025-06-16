import 'package:flutter/material.dart';
import 'package:lemon_tea/controls/ai_chat/views/chat_view/chat_view.dart';
import 'package:lemon_tea/utils/system.dart';

import '../../controls/window_title_bar.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<StatefulWidget> createState() => _HomePage();
}

class _HomePage extends State<HomePage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      // appBar: ,
      body: Center(
        child: ChatView(),
      ),
    );
  }
}
