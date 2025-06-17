import 'package:flutter/material.dart';
import 'package:lemon_tea/controls/input.dart';

class InputView extends StatefulWidget {
  final Function(String)? onFileSelected;
  final Function(String)? onSend;

  const InputView({
    super.key,
    this.onFileSelected,
    this.onSend,
  });

  @override
  State<StatefulWidget> createState() => _InputView();
}

class _InputView extends State<InputView> {
  final GlobalKey _inputViewKey = GlobalKey();
  final MenuController _menuController = MenuController();
  final TextEditingController _textController = TextEditingController();

  @override
  void dispose() {
    _textController.dispose();
    super.dispose();
  }

  @override
  void initState() {
    super.initState();
  }

  void _handleFileSelection(String type) {
    widget.onFileSelected?.call(type);
  }

  void _handleSend() {
    final text = _textController.text.trim();
    if (text.isNotEmpty) {
      widget.onSend?.call(text);
      _textController.clear();
    }
  }

  Widget _buildIconButton({
    required IconData icon,
    required VoidCallback onTap,
    Color? color,
  }) {
    return Material(
      color: Colors.transparent,
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(4.0),
        child: Padding(
          padding: const EdgeInsets.all(4.0),
          child: Icon(icon, size: 20, color: color),
        ),
      ),
    );
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
            controller: _textController,
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
                  return _buildIconButton(
                    icon: Icons.add,
                    onTap: () {
                      if (controller.isOpen) {
                        controller.close();
                      } else {
                        controller.open();
                      }
                    },
                  );
                },
              ),
              _buildIconButton(
                icon: Icons.arrow_forward,
                onTap: _handleSend,
              ),
            ],
          ),
        ],
      ),
    );
  }
}