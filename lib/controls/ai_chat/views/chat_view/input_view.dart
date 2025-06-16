import 'package:flutter/material.dart';
import 'package:lemon_tea/controls/input.dart';

class InputView extends StatefulWidget {
  @override
  State<StatefulWidget> createState() => _InputView();
}

class _InputView extends State<InputView> {
  final GlobalKey _inputViewKey = GlobalKey();
  double _inputViewHeight = 0;
  final MenuController _menuController = MenuController();

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      final RenderBox? renderBox = _inputViewKey.currentContext?.findRenderObject() as RenderBox?;
      if (renderBox != null) {
        setState(() {
          _inputViewHeight = renderBox.size.height;
          print('InputView height: $_inputViewHeight');
        });
      }
    });
  }

  void _handleFileSelection(String type) {
    // TODO: 实现文件选择逻辑
    print('Selected file type: $type');
  }

  void _handleSend() {
    // TODO: 实现发送逻辑
    print('Send message');
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8.0, horizontal: 8.0),
      child: Column(
        key: _inputViewKey,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Input(
            minLines: 2,
            maxLines: 4,
            hintText: '输入消息...',
          ),
          const SizedBox(height: 8),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              MenuAnchor(
                controller: _menuController,
                menuChildren: [
                  MenuItemButton(
                    onPressed: () {
                      _handleFileSelection('image');
                      _menuController.close();
                    },
                    child: const Row(
                      children: [
                        Icon(Icons.image),
                        SizedBox(width: 8),
                        Text('上传图片'),
                      ],
                    ),
                  ),
                  MenuItemButton(
                    onPressed: () {
                      _handleFileSelection('file');
                      _menuController.close();
                    },
                    child: const Row(
                      children: [
                        Icon(Icons.insert_drive_file),
                        SizedBox(width: 8),
                        Text('上传文件'),
                      ],
                    ),
                  ),
                ],
                builder: (context, controller, child) {
                  return InkWell(
                    onTap: () {
                      if (controller.isOpen) {
                        controller.close();
                      } else {
                        controller.open();
                      }
                    },
                    child: const Padding(
                      padding: EdgeInsets.all(4.0),
                      child: Icon(Icons.add, size: 20),
                    ),
                  );
                },
              ),
              InkWell(
                onTap: _handleSend,
                child: const Padding(
                  padding: EdgeInsets.all(4.0),
                  child: Icon(Icons.send, size: 20),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}