import 'package:flutter/material.dart';
import 'package:lemon_tea/controls/ai_chat/views/chat_view/input_view.dart';
import 'package:lemon_tea/controls/ai_chat/views/chat_view/message_toolbar.dart';
import 'package:lemon_tea/controls/ai_chat/views/chat_view/message_view.dart';
import 'package:lemon_tea/controls/ai_chat/views/chat_view/title_bar_view.dart';
import 'package:lemon_tea/utils/llm/models/message.dart';
import 'package:lemon_tea/generated/l10n.dart';

// 导出ChatView State类型以便其他文件可以访问
typedef ChatViewState = _ChatView;

class ChatView extends StatefulWidget {
  const ChatView({
    super.key, 
    this.onFileSelected, 
    this.onSend, 
    required this.historyMessages,
    this.onNewConversation,
    this.currentTitle = '',
    this.selectedProviderId,
    this.selectedModelId,
    this.onModelSelected,
    this.isStreaming = false,
    this.visibleWidth = double.infinity,
    this.messageToolBar
  });

  final Function(String)? onFileSelected;
  final Function(String)? onSend;
  final List<Message> historyMessages;
  final VoidCallback? onNewConversation;
  final String currentTitle;
  final String? selectedProviderId;
  final String? selectedModelId;
  final Function(String providerId, String modelId)? onModelSelected;
  final bool isStreaming;
  final double visibleWidth;
  final MessageToolbar? messageToolBar;

  @override
  State<StatefulWidget> createState() => _ChatView();
}

class _ChatView extends State<ChatView> {
  final GlobalKey<InputViewState> _inputViewKey = GlobalKey<InputViewState>();
  final GlobalKey<MessageViewState> _messageViewKey = GlobalKey<MessageViewState>(); // 添加MessageView的key
  bool _showScrollToBottomButton = false; // 控制是否显示滚动到底部按钮

  // 添加公共方法来请求输入框聚焦
  void requestInputFocus() {
    _inputViewKey.currentState?.requestFocus();
  }

  // 处理用户滚动状态变化
  void _onUserScrollChanged(bool userHasScrolled) {
    if (_showScrollToBottomButton != userHasScrolled) {
      // 延迟执行setState，避免在build阶段调用
      WidgetsBinding.instance.addPostFrameCallback((_) {
        if (mounted) {
          setState(() {
            _showScrollToBottomButton = userHasScrolled;
          });
        }
      });
    }
  }

  // 滚动到底部按钮点击处理
  void _onScrollToBottomTapped() {
    _messageViewKey.currentState?.scrollToBottomAndResumeAutoScroll();
  }

  @override
  Widget build(BuildContext context) {
    final displayTitle = widget.currentTitle.isEmpty ? S.of(context).aiAssistant : widget.currentTitle;
    
    return ConstrainedBox(
      constraints: const BoxConstraints(
        minWidth: 300.0,
      ),
      child: Column(
        children: [
          // 顶部部件
          SizedBox(
            child: TitleBarView(
              title: displayTitle,
              onAddTap: widget.onNewConversation,
              visibleWidth: widget.visibleWidth,
            ),
          ),
          Expanded(
            child: Stack(
              children: [
                // 聊天消息区域
                MessageView(
                  key: _messageViewKey, // 添加key
                  widget.historyMessages,
                  isStreaming: widget.isStreaming,
                  visibleWidth: widget.visibleWidth,
                  messageToolBar: widget.messageToolBar,
                  onUserScrollChanged: _onUserScrollChanged, // 传递回调函数
                ),
                
                // 悬浮的滚动到底部按钮
                if (_showScrollToBottomButton)
                  Positioned(
                    bottom: 16.0,
                    left: 0,
                    right: 0,
                    child: Center(
                      child: Material(
                        color: Theme.of(context).colorScheme.surface,
                        elevation: 6,
                        borderRadius: BorderRadius.circular(20),
                        child: InkWell(
                          onTap: _onScrollToBottomTapped,
                          borderRadius: BorderRadius.circular(20),
                          child: Container(
                            padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 8.0),
                            child: Row(
                              mainAxisSize: MainAxisSize.min,
                              children: [
                                Icon(
                                  Icons.keyboard_arrow_down,
                                  size: 20,
                                  color: Theme.of(context).colorScheme.primary,
                                ),
                                const SizedBox(width: 4),
                                Text(
                                  S.of(context).scrollToBottom,
                                  style: TextStyle(
                                    color: Theme.of(context).colorScheme.primary,
                                    fontSize: 14,
                                    fontWeight: FontWeight.w500,
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ),
                      ),
                    ),
                  ),
                ],
            ),
          ),

          // 底部部件
          SizedBox(
            width: widget.visibleWidth,
            child: InputView(
              key: _inputViewKey,
              onFileSelected: widget.onFileSelected,
              onSend: (msg) {
                widget.onSend?.call(msg);
              },
              selectedProviderId: widget.selectedProviderId,
              selectedModelId: widget.selectedModelId,
              onModelSelected: widget.onModelSelected,
              isStreaming: widget.isStreaming,
            ),
          ),
        ],
      ),
    );
  }
}
