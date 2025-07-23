import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/controls/input.dart';
import 'package:lemon_tea/generated/l10n.dart';
import 'package:lemon_tea/controls/ai_chat/views/chat_view/model_selector.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:lemon_tea/utils/style.dart';

class InputView extends ConsumerStatefulWidget {
  final Function(String)? onFileSelected;
  final Function(String)? onSend;
  final String? selectedProviderId;
  final String? selectedModelId;
  final Function(String providerId, String modelId)? onModelSelected;
  final bool isStreaming;

  const InputView({
    super.key,
    this.onFileSelected,
    this.onSend,
    this.selectedProviderId,
    this.selectedModelId,
    this.onModelSelected,
    this.isStreaming = false,
  });

  @override
  ConsumerState<InputView> createState() => _InputView();
}

class _InputView extends ConsumerState<InputView> {
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
          child: Icon(icon, size: 21, color: color),
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
          TextFormField(
            controller: _textController,
            minLines: 2,
            maxLines: 4,
            cursorWidth: 1.5,
            style: TextStyle(
              fontSize: FontSizeUtils.getSmallSize(ref),
              height: 1.5
            ),
            decoration: InputDecoration(
              floatingLabelBehavior: FloatingLabelBehavior.never,
              hintText: S.of(context).inputMessage,
              hintStyle: TextStyle(
                fontSize: FontSizeUtils.getSmallSize(ref),
                textBaseline: TextBaseline.alphabetic,
                color: Style.hintText(context),
              ),
              border: OutlineInputBorder(
                borderSide: BorderSide(color: Style.primaryBorder(context), width: 0.5),
                borderRadius:  BorderRadius.all(Radius.circular(10.0)),
              ),
              focusedBorder: OutlineInputBorder(
                borderSide: BorderSide(color: Style.focusedBorder(context), width: 1.0),
                borderRadius: BorderRadius.all(Radius.circular(10.0)),
              ),
              fillColor: Style.inputBackground(context),
              filled: true,
            ),
            onChanged: (value) {
              if (value.endsWith('\n')) {
                final isShiftPressed = HardwareKeyboard.instance.isShiftPressed;

                if (!isShiftPressed) {
                  final text = value.substring(0, value.length - 1);
                  _textController.text = text;
                  _textController.selection = TextSelection.fromPosition(
                    TextPosition(offset: text.length),
                  );
                  _handleSend();
                }
              }
            },
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
                    child: Row(
                      children: [
                        const Icon(Icons.image),
                        const SizedBox(width: 8),
                        Text(S.of(context).uploadImage),
                      ],
                    ),
                  ),
                  MenuItemButton(
                    onPressed: () {
                      _handleFileSelection('file');
                      _menuController.close();
                    },
                    child: Row(
                      children: [
                        const Icon(Icons.insert_drive_file),
                        const SizedBox(width: 8),
                        Text(S.of(context).uploadFile),
                      ],
                    ),
                  ),
                ],
                builder: (context, controller, child) {
                  return _buildIconButton(
                    icon: Icons.file_upload,
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
              Row(
                children: [
                  ModelSelector(
                    selectedProviderId: widget.selectedProviderId,
                    selectedModelId: widget.selectedModelId,
                    onModelSelected: widget.onModelSelected,
                  ),
                  const SizedBox(width: 8),
                  _buildIconButton(
                    icon: Icons.arrow_forward,
                    onTap: _handleSend,
                  ),
                ],
              ),
            ],
          ),
        ],
      ),
    );
  }
}