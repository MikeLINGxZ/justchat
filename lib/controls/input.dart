import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

class Input extends StatefulWidget {
  const Input({
    super.key,
    this.hintText,
    this.hintStyle,
    this.contentPadding,
    this.style,
    this.isDense,
    this.border,
    this.floatingLabelBehavior,
    this.onChanged,
    this.controller,
    this.hasError = false,
    this.obscureText = false,
    this.obscuringCharacter = '•',
    this.constraints,
    this.enable,
    this.focusNode,
    this.keyboardType,
    this.inputFormatters,
    this.maxLines = 3,  // 默认2行
    this.minLines = 1,  // 默认最小2行
  });

  final TextStyle? style;
  final String? hintText;
  final TextStyle? hintStyle;
  final EdgeInsetsGeometry? contentPadding;
  final bool? isDense;
  final InputBorder? border;
  final FloatingLabelBehavior? floatingLabelBehavior;
  final Function? onChanged;
  final TextEditingController? controller;
  final bool hasError;
  final bool obscureText;
  final String obscuringCharacter;
  final BoxConstraints? constraints;
  final bool? enable;
  final FocusNode? focusNode;
  final TextInputType? keyboardType;
  final List<TextInputFormatter>? inputFormatters;
  final int? maxLines;
  final int? minLines;

  @override
  State<StatefulWidget> createState() => _Input();
}

class _Input extends State<Input> {
  late TextEditingController _controller;
  bool _showScrollbar = false;

  @override
  void initState() {
    super.initState();
    _controller = widget.controller ?? TextEditingController();
    _controller.addListener(_updateScrollbarVisibility);
  }

  @override
  void dispose() {
    if (widget.controller == null) {
      _controller.dispose();
    }
    super.dispose();
  }

  void _updateScrollbarVisibility() {
    final text = _controller.text;
    if (text.isEmpty) {
      setState(() {
        _showScrollbar = false;
      });
      return;
    }

    final textStyle = widget.style ?? const TextStyle(fontSize: 15.0);
    final lineHeight = textStyle.height ?? 1.0;
    final fontSize = textStyle.fontSize ?? 15.0;
    final lineCount = text.split('\n').length;
    final contentHeight = lineCount * (lineHeight * fontSize);
    final maxHeight = (widget.maxLines ?? 1) * (lineHeight * fontSize);

    setState(() {
      _showScrollbar = contentHeight > maxHeight;
    });
  }

  @override
  Widget build(BuildContext context) {
    // 计算默认高度（1行的高度）
    final textStyle = widget.style ?? const TextStyle(fontSize: 15.0);
    final lineHeight = textStyle.height ?? 1.0;
    final fontSize = textStyle.fontSize ?? 15.0;
    final defaultHeight = (lineHeight * fontSize) + (widget.contentPadding?.vertical ?? 24);

    return ConstrainedBox(
      constraints: BoxConstraints(
        minHeight: defaultHeight,
        maxHeight: widget.maxLines != null
            ? defaultHeight * widget.maxLines!
            : double.infinity,
      ),
      child: TextFormField(
        keyboardType: widget.keyboardType,
        inputFormatters: widget.inputFormatters,
        focusNode: widget.focusNode,
        obscureText: widget.obscureText,
        obscuringCharacter: widget.obscuringCharacter,
        controller: _controller,
        style: textStyle.copyWith(height: 1.5),
        cursorWidth: 1.0,  // 设置光标宽度为1.0
        enabled: widget.enable ?? true,
        maxLines: widget.maxLines,  // 控制最大行数
        minLines: widget.minLines,  // 控制最小行数
        decoration: InputDecoration(
          floatingLabelBehavior: widget.floatingLabelBehavior ?? FloatingLabelBehavior.never,
          border: widget.border ?? const OutlineInputBorder(),
          hintText: widget.hintText ?? '',
          hintStyle: widget.hintStyle ?? const TextStyle(fontSize: 13.0, textBaseline: TextBaseline.alphabetic),
          isDense: widget.isDense ?? true,
          contentPadding: widget.contentPadding ?? const EdgeInsets.all(12),
          constraints: widget.constraints,
          enabledBorder: OutlineInputBorder(
            borderSide: BorderSide(color: widget.hasError ? Colors.red : Colors.grey, width: 0.5),
          ),
          errorBorder: widget.hasError ? const OutlineInputBorder(borderSide: BorderSide(color: Colors.red, width: 0.5)) : null,
          focusedBorder: OutlineInputBorder(
            borderSide: BorderSide(color: widget.hasError ? Colors.red : Colors.blue, width: 0.5),
          ),
        ),
        onChanged: (value) {
          widget.onChanged?.call(value);
        },
      ),
    );
  }
}