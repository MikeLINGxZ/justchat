
import 'package:flutter/material.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<StatefulWidget> createState() => _HomePage();
}

class _HomePage extends State<HomePage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(body: Column(
      children: [
        Divider(
          height: 1,
          thickness: 1,
        ),
        Expanded(
          child: Text("data"),
        ),
      ],
    ));
  }
}