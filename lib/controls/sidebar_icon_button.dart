import 'package:flutter/material.dart';

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
        color: isSelected ? Colors.blue : Colors.grey,
      ),
      onPressed: onPressed,
      style: IconButton.styleFrom(
        backgroundColor: isSelected ? Colors.blue.withOpacity(0.1) : Colors.transparent,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(8),
        ),
        minimumSize: const Size(40, 40),
        padding: EdgeInsets.zero,
      ),
    );
  }
} 