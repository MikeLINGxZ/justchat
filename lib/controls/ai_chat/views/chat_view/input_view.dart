
import 'package:flutter/material.dart';
import 'package:lemon_tea/controls/input.dart';

class InputView extends StatefulWidget {
  @override
  State<StatefulWidget> createState() => _InputView();
}

class _InputView extends State<InputView> {
  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Expanded(child: Input(maxLines: 3))
      ],
    );
  }

}