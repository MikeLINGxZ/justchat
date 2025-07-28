import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/generated/l10n.dart';
import 'package:lemon_tea/controls/ai_chat/views/chat_view/model_selector.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:lemon_tea/utils/style.dart';
import 'package:lemon_tea/utils/file_service.dart';
import 'package:lemon_tea/utils/llm/models/message.dart';

// 导出InputView State类型以便其他文件可以访问
typedef InputViewState = _InputView;

class InputView extends ConsumerStatefulWidget {
  final Function(String)? onFileSelected; // 保留向后兼容性
  final Function(String text, List<FileContent> files)? onSend;
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
  final FocusNode _focusNode = FocusNode();
  bool _isFocused = false;
  List<FileContent> _selectedFiles = []; // 存储选择的文件

  // 添加公共方法来请求聚焦
  void requestFocus() {
    _focusNode.requestFocus();
  }

  @override
  void dispose() {
    _textController.dispose();
    _focusNode.dispose();
    super.dispose();
  }

  @override
  void initState() {
    super.initState();
    _focusNode.addListener(() {
      setState(() {
        _isFocused = _focusNode.hasFocus;
      });
    });
  }

  void _handleFileSelection(String type) async {
    List<FileContent>? selectedFiles;
    
    switch (type) {
      case 'image':
        selectedFiles = await FileService.pickImages();
        break;
      case 'file':
        selectedFiles = await FileService.pickDocuments();
        break;
      case 'any':
        selectedFiles = await FileService.pickAnyFiles();
        break;
    }
    
    if (selectedFiles != null && selectedFiles.isNotEmpty) {
      setState(() {
        _selectedFiles.addAll(selectedFiles!);
      });
      
      // 向后兼容：调用原始的onFileSelected回调
      widget.onFileSelected?.call(type);
    }
  }

  void _removeFile(int index) {
    setState(() {
      _selectedFiles.removeAt(index);
    });
  }

  void _handleSend() {
    final text = _textController.text.trim();
    
    // 如果既没有文本内容也没有文件，则不发送
    if (text.isEmpty && _selectedFiles.isEmpty) return;
    
    // 调用新的onSend回调，传递文本和文件
    widget.onSend?.call(text, List.from(_selectedFiles));
    
    // 清空输入
    _textController.clear();
    setState(() {
      _selectedFiles.clear();
    });
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

  Widget _buildFilePreview() {
    if (_selectedFiles.isEmpty) return const SizedBox.shrink();
    
    return Container(
      margin: const EdgeInsets.only(bottom: 8.0),
      child: Wrap(
        spacing: 8.0,
        runSpacing: 8.0,
        children: _selectedFiles.asMap().entries.map((entry) {
          final index = entry.key;
          final file = entry.value;
          
          return Container(
            padding: const EdgeInsets.all(8.0),
            decoration: BoxDecoration(
              color: Theme.of(context).colorScheme.surfaceVariant,
              borderRadius: BorderRadius.circular(8.0),
              border: Border.all(
                color: Style.primaryBorder(context),
                width: 0.5,
              ),
            ),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                // 文件图标或图片预览
                if (FileService.isImage(file)) ...[
                  // 图片预览
                  Container(
                    width: 40,
                    height: 40,
                                         decoration: BoxDecoration(
                       borderRadius: BorderRadius.circular(4.0),
                       color: Theme.of(context).colorScheme.surfaceVariant,
                     ),
                    child: ClipRRect(
                      borderRadius: BorderRadius.circular(4.0),
                      child: Builder(
                        builder: (context) {
                          final imageData = FileService.getImagePreviewData(file);
                          if (imageData != null) {
                            return Image.memory(
                              imageData,
                              fit: BoxFit.cover,
                              errorBuilder: (context, error, stackTrace) {
                                return Icon(
                                  FileService.getFileTypeIcon(file.type),
                                  size: 20,
                                  color: Theme.of(context).colorScheme.onSurfaceVariant,
                                );
                              },
                            );
                          }
                          return Icon(
                            FileService.getFileTypeIcon(file.type),
                            size: 20,
                            color: Theme.of(context).colorScheme.onSurfaceVariant,
                          );
                        },
                      ),
                    ),
                  ),
                ] else ...[
                  // 非图片文件图标
                  Container(
                    width: 40,
                    height: 40,
                    decoration: BoxDecoration(
                      borderRadius: BorderRadius.circular(4.0),
                      color: Theme.of(context).colorScheme.surface,
                    ),
                    child: Icon(
                      FileService.getFileTypeIcon(file.type),
                      size: 20,
                      color: Theme.of(context).colorScheme.onSurface,
                    ),
                  ),
                ],
                
                const SizedBox(width: 8.0),
                
                // 文件信息
                Flexible(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Text(
                        file.name,
                        style: TextStyle(
                          fontSize: FontSizeUtils.getSmallSize(ref),
                          fontWeight: FontWeight.w500,
                        ),
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                      ),
                      Text(
                        FileService.formatFileSize(file.size),
                        style: TextStyle(
                          fontSize: FontSizeUtils.getSmallSize(ref) - 2,
                          color: Theme.of(context).colorScheme.onSurfaceVariant,
                        ),
                      ),
                    ],
                  ),
                ),
                
                const SizedBox(width: 8.0),
                
                // 删除按钮
                GestureDetector(
                  onTap: () => _removeFile(index),
                  child: Container(
                    padding: const EdgeInsets.all(2.0),
                    decoration: BoxDecoration(
                      color: Style.error(context),
                      borderRadius: BorderRadius.circular(12.0),
                    ),
                    child: Icon(
                      Icons.close,
                      size: 16,
                      color: Colors.white,
                    ),
                  ),
                ),
              ],
            ),
          );
        }).toList(),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8.0, horizontal: 8.0),
      child: Container(
        padding: EdgeInsets.all(10.0),
        decoration: BoxDecoration(
          border: Border.all(
            color: _isFocused ? Style.focusedBorder(context) : Style.primaryBorder(context)
          ),
          borderRadius: BorderRadius.all(Radius.circular(Style.radiusLv1))
        ),
        child: Column(
          key: _inputViewKey,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // 文件预览区域
            _buildFilePreview(),
            
            // 文本输入区域
            TextFormField(
              focusNode: _focusNode,
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
                border: InputBorder.none,
                focusedBorder: InputBorder.none,
                fillColor: Colors.transparent,
                filled: true,
                hoverColor:Colors.transparent,
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
                    MenuItemButton(
                      onPressed: () {
                        _handleFileSelection('any');
                        _menuController.close();
                      },
                      child: Row(
                        children: [
                          const Icon(Icons.attach_file),
                          const SizedBox(width: 8),
                          const Text('上传任意文件'),
                        ],
                      ),
                    ),
                  ],
                  builder: (context, controller, child) {
                    return _buildIconButton(
                      icon: Icons.file_upload,
                      onTap: widget.isStreaming ? () {} : () {
                        if (controller.isOpen) {
                          controller.close();
                        } else {
                          controller.open();
                        }
                      },
                      color: widget.isStreaming 
                        ? Style.disabledText(context)
                        : null,
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
                      onTap: widget.isStreaming ? () {} : _handleSend,
                      color: widget.isStreaming 
                        ? Style.disabledText(context)
                        : null,
                    ),
                  ],
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}