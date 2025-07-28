import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/controls/ai_chat/views/chat_view/message_toolbar.dart';
import 'package:lemon_tea/controls/ai_chat/views/chat_view/message_view_thinking.dart';
import 'package:lemon_tea/controls/ai_chat/views/chat_view/file_preview.dart';
import 'package:lemon_tea/models/message_role.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:lemon_tea/utils/llm/models/message.dart';
import 'package:lemon_tea/utils/style.dart';
import 'package:markdown_widget/markdown_widget.dart';
import 'package:flutter/services.dart';

// 导出MessageView State类型以便其他文件可以访问
typedef MessageViewState = _MessageViewState;

class MessageView extends ConsumerStatefulWidget {
  final List<Message> historyMessages;
  final bool isStreaming;
  final double visibleWidth;
  final MessageToolbar? messageToolBar;
  final Function(bool)? onUserScrollChanged; // 添加用户滚动状态变化回调

  const MessageView(
    this.historyMessages, {
    super.key,
    this.isStreaming = false,
    this.visibleWidth = double.infinity,
    this.messageToolBar,
    this.onUserScrollChanged, // 添加回调参数
  });

  @override
  ConsumerState<MessageView> createState() => _MessageViewState();
}

class _MessageViewState extends ConsumerState<MessageView> {
  final ScrollController _scrollController = ScrollController();
  int _lastMessageCount = 0;
  String _lastMessageContent = '';
  String _lastReasoningContent = ''; // 添加思维链内容跟踪
  final Map<int, bool> _messageHovered = {}; // 跟踪每个消息的悬停状态
  
  // 用户滚动状态管理
  bool _userHasScrolled = false; // 用户是否手动滚动过
  bool _wasStreaming = false; // 上次的流式状态
  bool _isInitializing = true; // 是否正在初始化，防止初始化时误触发用户滚动检测
  bool _isAutoScrolling = false; // 是否正在执行自动滚动，防止误触发用户滚动检测

  @override
  void initState() {
    super.initState();
    _lastMessageCount = widget.historyMessages.length;
    _wasStreaming = widget.isStreaming;
    if (widget.historyMessages.isNotEmpty) {
      final lastMessage = widget.historyMessages.last;
      _lastMessageContent = lastMessage.content;
      _lastReasoningContent = lastMessage.reasoningContent ?? '';
    }

    // 添加滚动监听器来检测用户滚动
    _scrollController.addListener(_onScrollChanged);

    WidgetsBinding.instance.addPostFrameCallback((_) {
      _scrollToBottom();
      // 初始滚动完成后，启用用户滚动检测
      WidgetsBinding.instance.addPostFrameCallback((_) {
        if (mounted) {
          setState(() {
            _isInitializing = false;
          });
        }
      });
    });
  }

  // 滚动变化监听器
  void _onScrollChanged() {
    if (!_scrollController.hasClients || _isInitializing || _isAutoScrolling) return;
    
    // 检查是否是用户手动滚动（不在自动滚动过程中）
    // 如果当前位置不是最底部，说明用户向上滚动了
    final maxScrollExtent = _scrollController.position.maxScrollExtent;
    final currentPosition = _scrollController.position.pixels;
    
    // 给一个小的容差值（10像素），避免因为精度问题误判
    final tolerance = 10.0;
    final isAtBottom = maxScrollExtent - currentPosition <= tolerance;
    
    if (!isAtBottom) {
      // 用户离开了底部
      if (!_userHasScrolled) {
        debugPrint('检测到用户手动滚动，暂停自动滚动');
        setState(() {
          _userHasScrolled = true;
        });
        // 通知父组件用户滚动状态变化
        widget.onUserScrollChanged?.call(true);
      }
    } else {
      // 用户回到了底部
      if (_userHasScrolled) {
        debugPrint('用户滚动回到底部，恢复自动滚动并隐藏按钮');
        setState(() {
          _userHasScrolled = false;
        });
        // 通知父组件用户滚动状态变化
        widget.onUserScrollChanged?.call(false);
      }
    }
  }

  // 添加公共方法供父组件调用，用于滚动到底部并重新启用自动滚动
  void scrollToBottomAndResumeAutoScroll() {
    debugPrint('手动滚动到底部，重新启用自动滚动');
    setState(() {
      _userHasScrolled = false;
      _isInitializing = false; // 确保用户滚动检测已启用
    });
    
    // 通知父组件用户滚动状态变化
    widget.onUserScrollChanged?.call(false);
    
    if (_scrollController.hasClients) {
      _isAutoScrolling = true;
      
      // 使用 addPostFrameCallback 确保在UI渲染完成后再执行滚动
      WidgetsBinding.instance.addPostFrameCallback((_) {
        if (_scrollController.hasClients) {
          // 添加一个小延迟确保所有布局计算完成
          Future.delayed(const Duration(milliseconds: 50), () {
            if (_scrollController.hasClients && mounted) {
              _scrollController.animateTo(
                _scrollController.position.maxScrollExtent,
                duration: const Duration(milliseconds: 300),
                curve: Curves.easeOut,
              ).then((_) {
                // 动画完成后，重新启用滚动检测
                _isAutoScrolling = false;
                
                // 二次检查并确保真正滚动到底部
                WidgetsBinding.instance.addPostFrameCallback((_) {
                  if (_scrollController.hasClients && mounted) {
                    final maxScrollExtent = _scrollController.position.maxScrollExtent;
                    final currentPosition = _scrollController.position.pixels;
                    final tolerance = 5.0;
                    
                    // 如果还没有真正到底部，再次滚动
                    if (maxScrollExtent - currentPosition > tolerance) {
                      debugPrint('二次滚动修正，当前位置: $currentPosition, 最大位置: $maxScrollExtent');
                      _scrollController.jumpTo(maxScrollExtent);
                    }
                  }
                });
              });
            }
          });
        }
      });
    }
  }

  @override
  void dispose() {
    _scrollController.removeListener(_onScrollChanged);
    _scrollController.dispose();
    super.dispose();
  }

  void _scrollToBottom() {
    // 如果用户已经手动滚动过，则不自动滚动
    if (_userHasScrolled) {
      return;
    }
    
    if (_scrollController.hasClients) {
      // 使用更短的延迟和更快的滚动动画，提升流式更新时的滚动体验
      WidgetsBinding.instance.addPostFrameCallback((_) {
        if (_scrollController.hasClients && !_userHasScrolled && mounted) {
          _isAutoScrolling = true;
          
          // 添加一个小延迟确保布局稳定
          Future.delayed(const Duration(milliseconds: 10), () {
            if (_scrollController.hasClients && !_userHasScrolled && mounted) {
              _scrollController.animateTo(
                _scrollController.position.maxScrollExtent,
                duration: const Duration(milliseconds: 50), // 缩短动画时间
                curve: Curves.easeOut,
              ).then((_) {
                // 动画完成后，重新启用滚动检测
                _isAutoScrolling = false;
                
                // 对于自动滚动，也添加精确性检查
                WidgetsBinding.instance.addPostFrameCallback((_) {
                  if (_scrollController.hasClients && !_userHasScrolled && mounted) {
                    final maxScrollExtent = _scrollController.position.maxScrollExtent;
                    final currentPosition = _scrollController.position.pixels;
                    final tolerance = 2.0; // 更小的容差，因为是自动滚动
                    
                    // 如果还没有真正到底部，直接跳转
                    if (maxScrollExtent - currentPosition > tolerance) {
                      _scrollController.jumpTo(maxScrollExtent);
                    }
                  }
                });
              });
            }
          });
        }
      });
    }
  }

  @override
  void didUpdateWidget(MessageView oldWidget) {
    super.didUpdateWidget(oldWidget);

    // 检测流式状态的变化
    if (!_wasStreaming && widget.isStreaming) {
      // 从非流式状态变为流式状态，重新启用自动滚动
      debugPrint('开始新的流式生成，重新启用自动滚动');
      setState(() {
        _userHasScrolled = false;
        _isInitializing = false; // 确保用户滚动检测已启用
      });
      // 通知父组件用户滚动状态变化
      widget.onUserScrollChanged?.call(false);
    } else if (_wasStreaming && !widget.isStreaming) {
      // 从流式状态变为非流式状态（生成完成），恢复默认状态
      debugPrint('生成完成，恢复默认状态（自动滚动，无按钮）');
      // 生成完成时，无论用户之前是否滚动过，都重置为初始状态
      setState(() {
        _userHasScrolled = false;
        _isInitializing = false;
      });
      // 通知父组件用户滚动状态变化（隐藏按钮）
      widget.onUserScrollChanged?.call(false);
      // 生成完成后滚动到底部
      _scrollToBottom();
    }
    _wasStreaming = widget.isStreaming;

    final currentMessageCount = widget.historyMessages.length;
    String currentLastContent = '';
    String currentLastReasoningContent = '';
    
    if (widget.historyMessages.isNotEmpty) {
      final lastMessage = widget.historyMessages.last;
      currentLastContent = lastMessage.content;
      currentLastReasoningContent = lastMessage.reasoningContent ?? '';
    }

    // 检查是否有新消息、最后一条消息内容发生变化，或思维链内容发生变化（流式更新）
    final hasNewMessage = currentMessageCount > _lastMessageCount;
    final hasContentChanged = _lastMessageContent != currentLastContent && currentLastContent.isNotEmpty;
    final hasReasoningChanged = _lastReasoningContent != currentLastReasoningContent && currentLastReasoningContent.isNotEmpty;
    
    if (hasNewMessage || hasContentChanged || hasReasoningChanged) {
      _lastMessageCount = currentMessageCount;
      _lastMessageContent = currentLastContent;
      _lastReasoningContent = currentLastReasoningContent;
      
      // 滚动到底部（如果用户没有手动滚动过的话）
      _scrollToBottom();
    }
  }





  // 动态创建 Markdown 配置
  MarkdownConfig _buildLightConfig() {
    return MarkdownConfig(
      configs: [
        PConfig(textStyle: TextStyle(fontSize: FontSizeUtils.getBodySize(ref))),
        H1Config(
          style: TextStyle(fontSize: FontSizeUtils.getTitleLargeSize(ref) + 6, height: 38 / 30),
        ),
        H2Config(
          style: TextStyle(fontSize: FontSizeUtils.getTitleSize(ref), height: 28 / 22),
        ),
        H3Config(
          style: TextStyle(fontSize: FontSizeUtils.getSubheadingSize(ref), height: 23 / 18),
        ),
        H4Config(
          style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref), height: 18 / 14),
        ),
        H5Config(
          style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref), height: 18 / 14),
        ),
        H6Config(
          style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref), height: 18 / 14),
        ),
      ],
    );
  }

  MarkdownConfig _buildDarkConfig() {
    return MarkdownConfig(
      configs: [
        PConfig(textStyle: TextStyle(fontSize: FontSizeUtils.getBodySize(ref))),
        H1Config(
          style: TextStyle(fontSize: FontSizeUtils.getTitleLargeSize(ref) + 6, height: 38 / 30),
        ),
        H2Config(
          style: TextStyle(fontSize: FontSizeUtils.getTitleSize(ref), height: 28 / 22),
        ),
        H3Config(
          style: TextStyle(fontSize: FontSizeUtils.getSubheadingSize(ref), height: 23 / 18),
        ),
        H4Config(
          style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref), height: 18 / 14),
        ),
        H5Config(
          style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref), height: 18 / 14),
        ),
        H6Config(
          style: TextStyle(fontSize: FontSizeUtils.getBodySize(ref), height: 18 / 14),
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    // 检查是否需要显示loading：正在流式处理且最后一条消息是空的AI消息（content和reasoningContent都为空）
    final shouldShowLoading = widget.isStreaming && 
        widget.historyMessages.isNotEmpty && 
        widget.historyMessages.last.role == MessageRole.assistant &&
        widget.historyMessages.last.content.trim().isEmpty &&
        (widget.historyMessages.last.reasoningContent?.trim().isEmpty ?? true);
    
    return SizedBox(
      width: double.infinity,
      child: ListView.builder(
        controller: _scrollController,
        itemCount: widget.historyMessages.length,
        itemBuilder: (context, index) {
          final message = widget.historyMessages[index];
          final isLastMessage = index == widget.historyMessages.length - 1;
          final isStreamingMessage = isLastMessage &&
              message.role == MessageRole.assistant &&
              widget.isStreaming;

          // 如果是最后一条空的AI消息且应该显示loading，则显示loading widget
          if (isLastMessage && shouldShowLoading) {
            return _buildLoadingMessage();
          }

          return Center(
            child: SizedBox(
              width: widget.visibleWidth,
              child: Padding(
                padding: const EdgeInsets.symmetric(vertical: 16.0, horizontal: 20.0),
                child: Row(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    if (message.role != MessageRole.user) ...[
                      CircleAvatar(
                        radius: 20,
                        backgroundColor: Colors.green,
                        child: Text(
                          'A',
                          style: TextStyle(
                            color: Colors.white,
                            fontSize: FontSizeUtils.getSubheadingSize(ref),
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ),
                      const SizedBox(width: 12),
                      Flexible(
                        child: MouseRegion(
                          onEnter: (_) => setState(() => _messageHovered[index] = true),
                          onExit: (_) => setState(() => _messageHovered[index] = false),
                          child: Container(
                            constraints: const BoxConstraints(
                              minWidth: 0,
                            ),
                            padding: const EdgeInsets.symmetric(
                              horizontal: 12,
                              vertical: 2,
                            ),
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                // 显示思考过程（如果有的话）
                                if (message.reasoningContent != null &&
                                    message.reasoningContent!.isNotEmpty)
                                  AnimatedSwitcher(
                                    duration: const Duration(milliseconds: 200),
                                    child: MessageViewThinking(
                                      messageIndex: index,
                                      reasoningContent: message.reasoningContent!,
                                      onExpansionChanged: _scrollToBottom,
                                    ),
                                  ),

                                // 显示文件内容（如果有的话）
                                if (message.hasFiles)
                                  FilePreview(
                                    files: message.files!,
                                    isUserMessage: false,
                                  ),

                                // 显示主要内容
                                Container(
                                  color: Style.assistantChatBubble(context),
                                  child: MarkdownBlock(
                                    data: message.content.isEmpty ? ' ' : message.content,
                                    config: Theme.of(context).brightness == Brightness.dark
                                        ? _buildDarkConfig()
                                        : _buildLightConfig(),
                                  ),
                                ),
                                if (message.role == MessageRole.assistant && widget.messageToolBar != null)
                                  widget.messageToolBar!.setMessage(message, visible: !widget.isStreaming && (_messageHovered[index] ?? false)),
                              ],
                            ),
                          ),
                        ),
                      ),
                    ],
                    if (message.role == MessageRole.user) ...[
                      const Spacer(),
                      Container(
                        constraints: const BoxConstraints(
                          minWidth: 0,
                        ),
                        padding: const EdgeInsets.symmetric(
                          horizontal: 12,
                          vertical: 2,
                        ),
                        decoration: BoxDecoration(
                          color:  Style.userChatBubble(context),
                          borderRadius: BorderRadius.circular(10),
                        ),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            // 显示文件内容（如果有的话）
                            if (message.hasFiles)
                              FilePreview(
                                files: message.files!,
                                isUserMessage: true,
                              ),
                            
                            // 显示文本内容
                            if (message.content.isNotEmpty)
                              MarkdownBlock(
                                data: message.content,
                                config: Theme.of(context).brightness == Brightness.dark
                                    ? _buildDarkConfig()
                                    : _buildLightConfig(),
                              ),
                          ],
                        ),
                      ),
                      const SizedBox(width: 12),
                      CircleAvatar(
                        radius: 20,
                        backgroundColor: Colors.blue,
                        child: Text(
                          'U',
                          style: TextStyle(
                            color: Colors.white,
                            fontSize: FontSizeUtils.getSubheadingSize(ref),
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ),
                    ],
                  ],
                ),
              ),
            ),
          );
        },
      ),
    );
  }

  // 构建loading消息的Widget
  Widget _buildLoadingMessage() {
    return Center(
      child: SizedBox(
        width: widget.visibleWidth,
        child: Padding(
          padding: const EdgeInsets.symmetric(vertical: 16.0, horizontal: 20.0),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              CircleAvatar(
                radius: 20,
                backgroundColor: Colors.green,
                child: Text(
                  'A',
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: FontSizeUtils.getSubheadingSize(ref),
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ),
              const SizedBox(width: 12),
              Flexible(
                child: Container(
                  constraints: const BoxConstraints(
                    minWidth: 0,
                  ),
                  padding: const EdgeInsets.symmetric(
                    horizontal: 12,
                    vertical: 12,
                  ),
                  decoration: BoxDecoration(
                    color: Style.assistantChatBubble(context),
                    borderRadius: BorderRadius.circular(10),
                  ),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      SizedBox(
                        width: 16,
                        height: 16,
                        child: CircularProgressIndicator(
                          strokeWidth: 2,
                          valueColor: AlwaysStoppedAnimation<Color>(
                            Theme.of(context).textTheme.bodyMedium?.color ?? Colors.grey,
                          ),
                        ),
                      ),
                      const SizedBox(width: 8),
                      Text(
                        '正在思考...',
                        style: TextStyle(
                          fontSize: FontSizeUtils.getBodySize(ref),
                          color: Theme.of(context).textTheme.bodyMedium?.color,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }


}

