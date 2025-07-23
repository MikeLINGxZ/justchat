import 'package:flutter/material.dart';
import 'package:lemon_tea/utils/style.dart';

class SidebarIconButton extends StatelessWidget {
  final IconData icon;
  final bool isSelected;
  final VoidCallback onPressed;

  const SidebarIconButton({
    super.key,
    required this.icon,
    required this.isSelected,
    required this.onPressed,
  });

  @override
  Widget build(BuildContext context) {
    return IconButton(
      icon: Icon(
        icon,
        size: 22,
        color: isSelected ? Style.primaryColor(context) : Colors.grey,
      ),
      onPressed: onPressed,
      style: IconButton.styleFrom(
        backgroundColor: isSelected ? Style.secondaryColor(context) : Colors.transparent,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8),
        ),
        minimumSize: const Size(40, 40),
        padding: EdgeInsets.zero,
      ),
    );
  }
} 