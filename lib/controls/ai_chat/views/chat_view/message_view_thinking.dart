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

class _MessageViewThinkingState extends ConsumerState<MessageViewThinking>
    with TickerProviderStateMixin {
  bool _isExpanded = false;
  bool _isHovered = false;
  late AnimationController _animationController;
  late AnimationController _pulseController;
  late Animation<double> _fadeAnimation;
  late Animation<double> _slideAnimation;
  late Animation<double> _pulseAnimation;

  @override
  void initState() {
    super.initState();
    _animationController = AnimationController(
      duration: const Duration(milliseconds: 350),
      vsync: this,
    );
    _pulseController = AnimationController(
      duration: const Duration(milliseconds: 1800),
      vsync: this,
    );
    
    _fadeAnimation = Tween<double>(
      begin: 0.0,
      end: 1.0,
    ).animate(CurvedAnimation(
      parent: _animationController,
      curve: Curves.easeOutCubic,
    ));
    
    _slideAnimation = Tween<double>(
      begin: -15.0,
      end: 0.0,
    ).animate(CurvedAnimation(
      parent: _animationController,
      curve: Curves.easeOutBack,
    ));

    _pulseAnimation = Tween<double>(
      begin: 1.0,
      end: 1.03,
    ).animate(CurvedAnimation(
      parent: _pulseController,
      curve: Curves.easeInOut,
    ));

    _pulseController.repeat(reverse: true);
  }

  @override
  void dispose() {
    _animationController.dispose();
    _pulseController.dispose();
    super.dispose();
  }

  void _toggleExpansion() {
    setState(() {
      _isExpanded = !_isExpanded;
    });
    
    if (_isExpanded) {
      _animationController.forward();
      if (widget.onExpansionChanged != null) {
        Future.delayed(const Duration(milliseconds: 350), () {
          widget.onExpansionChanged!();
        });
      }
    } else {
      _animationController.reverse();
    }
  }

  // 为思考过程定制的Markdown配置
  MarkdownConfig _buildReasoningMarkdownConfig() {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final primaryColor = isDark ? const Color(0xFF4FC3F7) : const Color(0xFF1976D2);
    final textColor = isDark ? const Color(0xFFE1F5FE) : const Color(0xFF263238);
    
    return MarkdownConfig(
      configs: [
        PConfig(
          textStyle: TextStyle(
            fontSize: FontSizeUtils.getCaptionSize(ref),
            color: textColor,
            height: 1.4,
            letterSpacing: 0.1,
          ),
        ),
        H1Config(
          style: TextStyle(
            fontSize: FontSizeUtils.getSubheadingSize(ref),
            height: 1.3,
            color: primaryColor,
            fontWeight: FontWeight.bold,
          ),
        ),
        H2Config(
          style: TextStyle(
            fontSize: FontSizeUtils.getBodySize(ref),
            height: 1.3,
            color: primaryColor,
            fontWeight: FontWeight.w600,
          ),
        ),
        H3Config(
          style: TextStyle(
            fontSize: FontSizeUtils.getSmallSize(ref),
            height: 1.3,
            color: primaryColor.withOpacity(0.9),
            fontWeight: FontWeight.w600,
          ),
        ),
        H4Config(
          style: TextStyle(
            fontSize: FontSizeUtils.getCaptionSize(ref),
            height: 1.3,
            color: textColor,
            fontWeight: FontWeight.w500,
          ),
        ),
        H5Config(
          style: TextStyle(
            fontSize: FontSizeUtils.getXSmallSize(ref),
            height: 1.3,
            color: textColor,
            fontWeight: FontWeight.w500,
          ),
        ),
        H6Config(
          style: TextStyle(
            fontSize: FontSizeUtils.getXSmallSize(ref),
            height: 1.3,
            color: textColor,
            fontWeight: FontWeight.w500,
          ),
        ),
       ],
    );
  }

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final primaryColor = isDark ? const Color(0xFF4FC3F7) : const Color(0xFF1976D2);
    final backgroundColor = isDark ? const Color(0xFF1A1A1A) : const Color(0xFFFAFAFA);
    final cardColor = isDark ? const Color(0xFF2D2D2D) : Colors.white;
    
    return Container(
      margin: const EdgeInsets.symmetric(vertical: 6),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 主标题栏
          MouseRegion(
            onEnter: (_) => setState(() => _isHovered = true),
            onExit: (_) => setState(() => _isHovered = false),
            child: GestureDetector(
              onTap: _toggleExpansion,
              child: AnimatedBuilder(
                animation: _pulseAnimation,
                builder: (context, child) {
                  return Transform.scale(
                    scale: _isHovered ? _pulseAnimation.value : 1.0,
                    child: AnimatedContainer(
                      duration: const Duration(milliseconds: 180),
                      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
                      decoration: BoxDecoration(
                        gradient: LinearGradient(
                          colors: _isHovered
                              ? [
                                  primaryColor.withOpacity(0.12),
                                  primaryColor.withOpacity(0.20),
                                ]
                              : [
                                  primaryColor.withOpacity(0.06),
                                  primaryColor.withOpacity(0.10),
                                ],
                          begin: Alignment.topLeft,
                          end: Alignment.bottomRight,
                        ),
                        borderRadius: BorderRadius.circular(10),
                        border: Border.all(
                          color: primaryColor.withOpacity(_isHovered ? 0.5 : 0.25),
                          width: _isHovered ? 1.5 : 1,
                        ),
                        boxShadow: [
                          BoxShadow(
                            color: primaryColor.withOpacity(_isHovered ? 0.15 : 0.08),
                            blurRadius: _isHovered ? 8 : 4,
                            offset: const Offset(0, 2),
                            spreadRadius: _isHovered ? 1 : 0,
                          ),
                        ],
                      ),
                      child: Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          // 图标容器
                          Container(
                            padding: const EdgeInsets.all(4),
                            decoration: BoxDecoration(
                              gradient: LinearGradient(
                                colors: [
                                  primaryColor.withOpacity(0.15),
                                  primaryColor.withOpacity(0.30),
                                ],
                                begin: Alignment.topLeft,
                                end: Alignment.bottomRight,
                              ),
                              borderRadius: BorderRadius.circular(6),
                              boxShadow: [
                                BoxShadow(
                                  color: primaryColor.withOpacity(0.2),
                                  blurRadius: 2,
                                  offset: const Offset(0, 1),
                                ),
                              ],
                            ),
                            child: Icon(
                              Icons.psychology_outlined,
                              size: 14,
                              color: primaryColor,
                            ),
                          ),
                          const SizedBox(width: 6),
                          // 文字
                          Text(
                            '🧠 思维链',
                            style: TextStyle(
                              fontSize: FontSizeUtils.getCaptionSize(ref),
                              color: primaryColor,
                              fontWeight: FontWeight.w600,
                              letterSpacing: 0.3,
                            ),
                          ),
                          const SizedBox(width: 6),
                          // 展开图标
                          AnimatedRotation(
                            turns: _isExpanded ? 0.5 : 0,
                            duration: const Duration(milliseconds: 250),
                            child: Container(
                              padding: const EdgeInsets.all(2),
                              decoration: BoxDecoration(
                                color: primaryColor.withOpacity(0.08),
                                borderRadius: BorderRadius.circular(4),
                              ),
                              child: Icon(
                                Icons.keyboard_arrow_down_rounded,
                                size: 14,
                                color: primaryColor,
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                  );
                }
              ),
            ),
          ),
          
          // 展开内容
          AnimatedContainer(
            duration: const Duration(milliseconds: 350),
            curve: Curves.easeInOutCubic,
            height: _isExpanded ? null : 0,
            child: _isExpanded
                ? AnimatedBuilder(
                    animation: _fadeAnimation,
                    builder: (context, child) {
                      return Transform.translate(
                        offset: Offset(0, _slideAnimation.value),
                        child: Opacity(
                          opacity: _fadeAnimation.value,
                          child: Column(
                            children: [
                              const SizedBox(height: 8),
                              Container(
                                width: double.infinity,
                                decoration: BoxDecoration(
                                  color: cardColor,
                                  borderRadius: BorderRadius.circular(12),
                                  border: Border.all(
                                    color: primaryColor.withOpacity(0.15),
                                    width: 1,
                                  ),
                                  boxShadow: [
                                    BoxShadow(
                                      color: Colors.black.withOpacity(isDark ? 0.2 : 0.06),
                                      blurRadius: 12,
                                      offset: const Offset(0, 4),
                                      spreadRadius: 1,
                                    ),
                                    BoxShadow(
                                      color: primaryColor.withOpacity(0.08),
                                      blurRadius: 20,
                                      offset: const Offset(0, 2),
                                      spreadRadius: -2,
                                    ),
                                  ],
                                ),
                                child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    // 内容标题栏
                                    Container(
                                      width: double.infinity,
                                      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                                      decoration: BoxDecoration(
                                        gradient: LinearGradient(
                                          colors: [
                                            primaryColor.withOpacity(0.03),
                                            primaryColor.withOpacity(0.08),
                                          ],
                                          begin: Alignment.topLeft,
                                          end: Alignment.bottomRight,
                                        ),
                                        borderRadius: const BorderRadius.only(
                                          topLeft: Radius.circular(12),
                                          topRight: Radius.circular(12),
                                        ),
                                      ),
                                      child: Row(
                                        children: [
                                          Container(
                                            padding: const EdgeInsets.all(4),
                                            decoration: BoxDecoration(
                                              gradient: LinearGradient(
                                                colors: [
                                                  primaryColor.withOpacity(0.15),
                                                  primaryColor.withOpacity(0.25),
                                                ],
                                              ),
                                              borderRadius: BorderRadius.circular(6),
                                            ),
                                            child: Icon(
                                              Icons.lightbulb_outline_rounded,
                                              size: 12,
                                              color: primaryColor,
                                            ),
                                          ),
                                          const SizedBox(width: 6),
                                          Text(
                                            '思考过程',
                                            style: TextStyle(
                                              fontSize: FontSizeUtils.getXSmallSize(ref),
                                              color: primaryColor,
                                              fontWeight: FontWeight.w500,
                                              letterSpacing: 0.2,
                                            ),
                                          ),
                                        ],
                                      ),
                                    ),
                                    
                                    // 思考内容
                                    Container(
                                      padding: const EdgeInsets.all(12),
                                      child: MarkdownBlock(
                                        data: widget.reasoningContent,
                                        config: _buildReasoningMarkdownConfig(),
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                            ],
                          ),
                        ),
                      );
                    },
                  )
                : const SizedBox.shrink(),
          ),
        ],
      ),
    );
  }
} 