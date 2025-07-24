import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:markdown_widget/markdown_widget.dart';

class MessageViewThinking extends ConsumerStatefulWidget {
  final int messageIndex;
  final String reasoningContent;
  final VoidCallback? onExpansionChanged;

  const MessageViewThinking({
    super.key,
    required this.messageIndex,
    required this.reasoningContent,
    this.onExpansionChanged,
  });

  @override
  ConsumerState<MessageViewThinking> createState() => _MessageViewThinkingState();
}

class _MessageViewThinkingState extends ConsumerState<MessageViewThinking> {
  bool _isExpanded = false;

  void _toggleExpansion() {
    setState(() {
      _isExpanded = !_isExpanded;
    });
    
    // 如果展开思考过程，通知父组件
    if (_isExpanded && widget.onExpansionChanged != null) {
      Future.delayed(const Duration(milliseconds: 300), () {
        widget.onExpansionChanged!();
      });
    }
  }

  // 为思考过程定制的Markdown配置
  MarkdownConfig _buildReasoningMarkdownConfig() {
    return MarkdownConfig(
      configs: [
        PConfig(
          textStyle: TextStyle(
            fontSize: FontSizeUtils.getXSmallSize(ref),
            color: Colors.orange.shade800,
            height: 1.5,
          ),
        ),
        H1Config(
          style: TextStyle(
            fontSize: FontSizeUtils.getTitleSize(ref),
            height: 1.3,
            color: Colors.orange.shade900,
            fontWeight: FontWeight.bold,
          ),
        ),
        H2Config(
          style: TextStyle(
            fontSize: FontSizeUtils.getSubheadingSize(ref),
            height: 1.3,
            color: Colors.orange.shade900,
            fontWeight: FontWeight.w600,
          ),
        ),
        H3Config(
          style: TextStyle(
            fontSize: FontSizeUtils.getBodySize(ref),
            height: 1.3,
            color: Colors.orange.shade800,
            fontWeight: FontWeight.w600,
          ),
        ),
        H4Config(
          style: TextStyle(
            fontSize: FontSizeUtils.getSmallSize(ref),
            height: 1.3,
            color: Colors.orange.shade800,
            fontWeight: FontWeight.w500,
          ),
        ),
        H5Config(
          style: TextStyle(
            fontSize: FontSizeUtils.getXSmallSize(ref),
            height: 1.3,
            color: Colors.orange.shade800,
            fontWeight: FontWeight.w500,
          ),
        ),
        H6Config(
          style: TextStyle(
            fontSize: FontSizeUtils.getXSmallSize(ref),
            height: 1.3,
            color: Colors.orange.shade800,
            fontWeight: FontWeight.w500,
          ),
        ),
       ],
    );
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.only(top: 12, bottom: 8),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          GestureDetector(
            onTap: _toggleExpansion,
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  colors: [
                    Colors.orange.withOpacity(0.15),
                    Colors.orange.withOpacity(0.08),
                  ],
                  begin: Alignment.centerLeft,
                  end: Alignment.centerRight,
                ),
                borderRadius: BorderRadius.circular(12),
                border: Border.all(
                  color: Colors.orange.withOpacity(0.4),
                  width: 1.5,
                ),
                boxShadow: [
                  BoxShadow(
                    color: Colors.orange.withOpacity(0.1),
                    blurRadius: 4,
                    offset: const Offset(0, 2),
                  ),
                ],
              ),
              child: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Container(
                    padding: const EdgeInsets.all(4),
                    decoration: BoxDecoration(
                      color: Colors.orange.withOpacity(0.2),
                      borderRadius: BorderRadius.circular(6),
                    ),
                    child: Icon(
                      Icons.psychology_outlined,
                      size: 16,
                      color: Colors.orange.shade800,
                    ),
                  ),
                  const SizedBox(width: 8),
                  Text(
                    '思考过程',
                    style: TextStyle(
                      fontSize: FontSizeUtils.getSmallSize(ref),
                      color: Colors.orange.shade800,
                      fontWeight: FontWeight.w600,
                      letterSpacing: 0.5,
                    ),
                  ),
                  const SizedBox(width: 8),
                  AnimatedRotation(
                    turns: _isExpanded ? 0.5 : 0,
                    duration: const Duration(milliseconds: 200),
                    child: Icon(
                      Icons.keyboard_arrow_down,
                      size: 18,
                      color: Colors.orange.shade700,
                    ),
                  ),
                ],
              ),
            ),
          ),
          AnimatedContainer(
            duration: const Duration(milliseconds: 300),
            curve: Curves.easeInOut,
            height: _isExpanded ? null : 0,
            child: _isExpanded ? Column(
              children: [
                const SizedBox(height: 12),
                Container(
                  width: double.infinity,
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    gradient: LinearGradient(
                      colors: [
                        Colors.orange.withOpacity(0.03),
                        Colors.orange.withOpacity(0.08),
                      ],
                      begin: Alignment.topLeft,
                      end: Alignment.bottomRight,
                    ),
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(
                      color: Colors.orange.withOpacity(0.3),
                      width: 1,
                    ),
                    boxShadow: [
                      BoxShadow(
                        color: Colors.orange.withOpacity(0.05),
                        blurRadius: 8,
                        offset: const Offset(0, 2),
                        spreadRadius: 1,
                      ),
                    ],
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      // 思考过程标题栏
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                        decoration: BoxDecoration(
                          color: Colors.orange.withOpacity(0.15),
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Icon(
                              Icons.lightbulb_outline,
                              size: 14,
                              color: Colors.orange.shade700,
                            ),
                            const SizedBox(width: 6),
                            Text(
                              'AI 思考过程',
                              style: TextStyle(
                                fontSize: FontSizeUtils.getXSmallSize(ref),
                                color: Colors.orange.shade700,
                                fontWeight: FontWeight.w500,
                                letterSpacing: 0.3,
                              ),
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(height: 12),
                      // 思考过程内容
                      MarkdownBlock(
                        data: widget.reasoningContent,
                        config: _buildReasoningMarkdownConfig(),
                      ),
                    ],
                  ),
                ),
              ],
            ) : const SizedBox.shrink(),
          ),
        ],
      ),
    );
  }
} 