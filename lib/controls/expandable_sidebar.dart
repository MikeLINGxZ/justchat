import 'package:flutter/material.dart';
import 'package:flutter/material.dart';
import 'package:lemon_tea/utils/style.dart';

class ExpandableSidebar extends StatefulWidget {
  final int selectedIndex;
  final Function(int) onItemSelected;
  final List<SidebarItem> items;

  const ExpandableSidebar({
    super.key,
    required this.selectedIndex,
    required this.onItemSelected,
    required this.items,
  });

  @override
  State<ExpandableSidebar> createState() => _ExpandableSidebarState();
}

class _ExpandableSidebarState extends State<ExpandableSidebar>
    with TickerProviderStateMixin {
  bool _isExpanded = false;
  late AnimationController _animationController;
  late Animation<double> _widthAnimation;

  @override
  void initState() {
    super.initState();
    _animationController = AnimationController(
      duration: const Duration(milliseconds: 200),
      vsync: this,
    );
    _widthAnimation = Tween<double>(
      begin: 64.0, // 收缩时的宽度
      end: 200.0,  // 展开时的宽度
    ).animate(CurvedAnimation(
      parent: _animationController,
      curve: Curves.easeInOut,
    ));
  }

  @override
  void dispose() {
    _animationController.dispose();
    super.dispose();
  }

  void _toggleExpanded() {
    setState(() {
      _isExpanded = !_isExpanded;
      if (_isExpanded) {
        _animationController.forward();
      } else {
        _animationController.reverse();
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return AnimatedBuilder(
      animation: _widthAnimation,
      builder: (context, child) {
        return Container(
          width: _widthAnimation.value,
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(Style.radiusLv1),
            color: Style.sidebarBackground(context),
          ),
          padding: const EdgeInsets.symmetric(vertical: 20.0, horizontal: 12.0),
          child: Column(
            children: [
              // 顶部柠檬图标按钮
              _buildToggleButton(),
              const SizedBox(height: 14),
              
              // 菜单项
              ...widget.items.map((item) => _buildMenuItem(item)).toList(),
            ],
          ),
        );
      },
    );
  }

  Widget _buildToggleButton() {
    return MouseRegion(
      cursor: SystemMouseCursors.click,
      child: Tooltip(
        message: _isExpanded ? '收起菜单' : '展开菜单',
        child: GestureDetector(
          onTap: _toggleExpanded,
          child: Container(
            width: double.infinity,
            height: 40,
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(8),
              color: Colors.transparent,
            ),
            child: _isExpanded
                ? Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 12),
                    child: Row(
                      children: [
                        const Text('🍋', style: TextStyle(fontSize: 20)),
                        const Spacer(),
                        Icon(
                          Icons.keyboard_arrow_left,
                          size: 16,
                          color: Style.secondaryText(context),
                        ),
                      ],
                    ),
                  )
                : Center(
                    child: const Text('🍋', style: TextStyle(fontSize: 20)),
                  ),
          ),
        ),
      ),
    );
  }

  Widget _buildMenuItem(SidebarItem item) {
    final isSelected = widget.selectedIndex == item.index;
    
    return Padding(
      padding: const EdgeInsets.only(bottom: 14),
      child: Tooltip(
        message: _isExpanded ? '' : item.title,
        child: MouseRegion(
          cursor: SystemMouseCursors.click,
          child: GestureDetector(
            onTap: () => widget.onItemSelected(item.index),
            child: Container(
              width: double.infinity,
              height: 40,
              decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(8),
                color: isSelected ? Style.secondaryColor(context) : Colors.transparent,
              ),
              child: _isExpanded
                  ? Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 12),
                      child: Row(
                        children: [
                          Icon(
                            item.icon,
                            size: 22,
                            color: isSelected ? Style.primaryColor(context) : Colors.grey,
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: Text(
                              item.title,
                              style: TextStyle(
                                color: isSelected ? Style.primaryColor(context) : Style.primaryText(context),
                                fontSize: 14,
                                fontWeight: isSelected ? FontWeight.w600 : FontWeight.normal,
                              ),
                              overflow: TextOverflow.ellipsis,
                            ),
                          ),
                        ],
                      ),
                    )
                  : Center(
                      child: Icon(
                        item.icon,
                        size: 22,
                        color: isSelected ? Style.primaryColor(context) : Colors.grey,
                      ),
                    ),
            ),
          ),
        ),
      ),
    );
  }
}

class SidebarItem {
  final IconData icon;
  final String title;
  final int index;

  const SidebarItem({
    required this.icon,
    required this.title,
    required this.index,
  });
}