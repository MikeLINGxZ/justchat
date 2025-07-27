import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/models/llm_provider.dart';
import 'package:lemon_tea/models/model.dart';
import 'package:lemon_tea/storage/llm_storage.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';

class ModelSelector extends ConsumerStatefulWidget {
  final String? selectedProviderId;
  final String? selectedModelId;
  final Function(String providerId, String modelId)? onModelSelected;

  const ModelSelector({
    super.key,
    this.selectedProviderId,
    this.selectedModelId,
    this.onModelSelected,
  });

  @override
  ConsumerState<ModelSelector> createState() => _ModelSelectorState();
}

class _ModelSelectorState extends ConsumerState<ModelSelector> {
  final MenuController _menuController = MenuController();
  final GlobalKey _buttonKey = GlobalKey();
  List<LlmProvider> _providers = [];
  Map<String, List<Model>> _providerModels = {};
  List<_ModelItem> _allModels = [];
  List<_ModelItem> _filteredModels = [];
  List<Widget> _displayItems = [];
  bool _isLoading = false;
  String _searchQuery = '';
  final TextEditingController _searchController = TextEditingController();

  @override
  void initState() {
    super.initState();
    _loadProvidersAndModels();
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  Future<void> _loadProvidersAndModels() async {
    setState(() {
      _isLoading = true;
    });

    try {
      // 获取所有供应商
      final providers = await LlmStorage.getAllProviders();
      final Map<String, List<Model>> providerModels = {};
      final List<_ModelItem> allModels = [];

      // 为每个供应商加载模型
      for (final provider in providers) {
        final models = await LlmStorage.getModelsByProviderId(provider.id);
        if (models.isNotEmpty) {
          providerModels[provider.id] = models;
          // 将模型转换为 _ModelItem
          for (final model in models) {
            allModels.add(_ModelItem(
              provider: provider,
              model: model,
            ));
          }
        }
      }

      setState(() {
        _providers = providers.where((p) => providerModels.containsKey(p.id)).toList();
        _providerModels = providerModels;
        _allModels = allModels;
        _filteredModels = allModels;
        _buildDisplayItems();
        _isLoading = false;
      });
    } catch (e) {
      debugPrint('加载供应商和模型失败: $e');
      setState(() {
        _isLoading = false;
      });
    }
  }

  void _filterModels(String query) {
    setState(() {
      _searchQuery = query;
      if (query.isEmpty) {
        _filteredModels = _allModels;
      } else {
        _filteredModels = _allModels.where((item) {
          return item.model.id.toLowerCase().contains(query.toLowerCase()) ||
                 item.provider.name.toLowerCase().contains(query.toLowerCase());
        }).toList();
      }
      _buildDisplayItems();
    });
  }

  void _buildDisplayItems() {
    _displayItems.clear();
    
    if (_searchQuery.isEmpty) {
      // 没有搜索时，按供应商分组显示
      for (final provider in _providers) {
        final models = _providerModels[provider.id] ?? [];
        if (models.isNotEmpty) {
          // 添加供应商标题
          _displayItems.add(_buildProviderHeader(provider));
          // 添加该供应商的所有模型
          for (final model in models) {
            _displayItems.add(_buildModelItem(_ModelItem(provider: provider, model: model)));
          }
        }
      }
    } else {
      // 有搜索时，直接显示匹配的模型
      for (final item in _filteredModels) {
        _displayItems.add(_buildModelItem(item));
      }
    }
  }

  double _getButtonWidth() {
    final RenderBox? renderBox = _buttonKey.currentContext?.findRenderObject() as RenderBox?;
    return renderBox?.size.width ?? 200;
  }

  String _getCurrentDisplayText() {
    if (widget.selectedProviderId != null && widget.selectedModelId != null) {
      final provider = _providers.firstWhere(
        (p) => p.id == widget.selectedProviderId,
        orElse: () => LlmProvider(
          id: '',
          name: '未知供应商',
          baseUrl: '',
          apiKey: '',
          seqId: 0,
        ),
      );
      
      final models = _providerModels[widget.selectedProviderId] ?? [];
      final model = models.firstWhere(
        (m) => m.id == widget.selectedModelId,
        orElse: () => Model(
          id: '',
          object: '',
          ownedBy: '',
          enabled: true,
          llmProviderId: '',
          seqId: 0,
        ),
      );
      
      return '${provider.name} / ${model.id}';
    }
    return '选择模型';
  }

  Widget _buildProviderHeader(LlmProvider provider) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 6.0),
      child: Row(
        children: [
          Icon(
            Icons.business,
            size: 14,
            color: Theme.of(context).colorScheme.primary,
          ),
          const SizedBox(width: 6),
          Text(
            provider.name,
            style: TextStyle(
              fontSize: FontSizeUtils.getSmallSize(ref),
              fontWeight: FontWeight.w600,
              color: Theme.of(context).colorScheme.primary,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSearchBar() {
    return Container(
      padding: const EdgeInsets.all(8.0),
      decoration: BoxDecoration(
        border: Border(
          top: BorderSide(
            color: Theme.of(context).colorScheme.outline.withOpacity(0.2),
          ),
        ),
      ),
      child: TextField(
        controller: _searchController,
        decoration: InputDecoration(
          hintText: '搜索模型...',
          hintStyle: TextStyle(fontSize: FontSizeUtils.getXSmallSize(ref)),
          prefixIcon: const Icon(Icons.search, size: 18),
          suffixIcon: _searchQuery.isNotEmpty
              ? IconButton(
                  icon: const Icon(Icons.clear, size: 18),
                  onPressed: () {
                    _searchController.clear();
                    _filterModels('');
                  },
                )
              : null,
          isDense: true,
          contentPadding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
          border: OutlineInputBorder(
            borderRadius: BorderRadius.circular(6),
            borderSide: BorderSide(
              color: Theme.of(context).colorScheme.outline.withOpacity(0.3),
            ),
          ),
          enabledBorder: OutlineInputBorder(
            borderRadius: BorderRadius.circular(6),
            borderSide: BorderSide(
              color: Theme.of(context).colorScheme.outline.withOpacity(0.3),
            ),
          ),
          focusedBorder: OutlineInputBorder(
            borderRadius: BorderRadius.circular(6),
            borderSide: BorderSide(
              color: Theme.of(context).colorScheme.primary,
            ),
          ),
        ),
        onChanged: _filterModels,
        style: TextStyle(fontSize: FontSizeUtils.getSmallSize(ref)),
      ),
    );
  }

  Widget _buildModelList() {
    if (_isLoading) {
      return const SizedBox(
        height: 150,
        child: Center(
          child: CircularProgressIndicator(),
        ),
      );
    }

    if (_displayItems.isEmpty) {
      return Container(
        height: 80,
        alignment: Alignment.center,
        child: Text(
          _searchQuery.isNotEmpty ? '未找到匹配的模型' : '暂无可用模型',
          style: TextStyle(
            fontSize: FontSizeUtils.getSmallSize(ref),
            color: Theme.of(context).colorScheme.onSurface.withOpacity(0.6),
          ),
        ),
      );
    }

    return ConstrainedBox(
      constraints: const BoxConstraints(
        maxHeight: 250,
        minHeight: 80,
      ),
      child: SingleChildScrollView(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: _displayItems,
        ),
      ),
    );
  }

  Widget _buildModelItem(_ModelItem item) {
    final isSelected = widget.selectedProviderId == item.provider.id && 
                      widget.selectedModelId == item.model.id;
    
    return MenuItemButton(
      onPressed: () {
        widget.onModelSelected?.call(item.provider.id, item.model.id);
        _menuController.close();
      },
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 0.0, vertical: 4.0),
        child: Row(
          children: [
            SizedBox(
              width: 20,
              child: isSelected
                  ? Icon(
                      Icons.check,
                      size: 14,
                      color: Theme.of(context).colorScheme.primary,
                    )
                  : null,
            ),
            const SizedBox(width: 6),
            Expanded(
              child: _searchQuery.isNotEmpty
                  ? Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Text(
                          item.model.id,
                          style: TextStyle(
                            fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
                            fontSize: FontSizeUtils.getSmallSize(ref),
                          ),
                        ),
                        const SizedBox(height: 2),
                        Text(
                          item.provider.name,
                          style: TextStyle(
                            fontSize: FontSizeUtils.getXSmallSize(ref),
                            color: Theme.of(context).colorScheme.onSurface.withOpacity(0.6),
                          ),
                        ),
                      ],
                    )
                  : Text(
                      item.model.id,
                      style: TextStyle(
                        fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
                        fontSize: FontSizeUtils.getSmallSize(ref),
                      ),
                    ),
            ),
          ],
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return MenuAnchor(
      controller: _menuController,
      style: MenuStyle(
        elevation: WidgetStateProperty.all(8),
        maximumSize: WidgetStateProperty.all(Size(_getButtonWidth().clamp(250, 400), 450)),
        padding: WidgetStateProperty.all(EdgeInsets.zero),
      ),
      menuChildren: [
        // 创建一个容器来包含模型列表和搜索框
        Container(
          width: _getButtonWidth().clamp(230, 380),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              // 模型列表在上方
              _buildModelList(),
              // 搜索框在底部
              _buildSearchBar(),
            ],
          ),
        ),
      ],
      builder: (context, controller, child) {
        return Material(
          color: Colors.transparent,
          child: InkWell(
            key: _buttonKey,
            onTap: () {
              if (controller.isOpen) {
                controller.close();
              } else {
                // 打开菜单时清空搜索
                _searchController.clear();
                _filterModels('');
                controller.open();
              }
            },
            borderRadius: BorderRadius.circular(4.0),
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 8.0, vertical: 4.0),
              decoration: BoxDecoration(
                border: Border.all(
                  color: Theme.of(context).colorScheme.outline.withOpacity(0.5),
                ),
                borderRadius: BorderRadius.circular(4.0),
              ),
              child: Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Icon(
                    Icons.smart_toy,
                    size: 16,
                    color: Theme.of(context).colorScheme.primary,
                  ),
                  const SizedBox(width: 4),
                  ConstrainedBox(
                    constraints: const BoxConstraints(maxWidth: 150),
                    child: Text(
                      _getCurrentDisplayText(),
                      style: TextStyle(
                        fontSize: FontSizeUtils.getCaptionSize(ref),
                        color: Theme.of(context).colorScheme.onSurface,
                      ),
                      overflow: TextOverflow.ellipsis,
                    ),
                  ),
                  const SizedBox(width: 4),
                  Icon(
                    Icons.keyboard_arrow_down,
                    size: 16,
                    color: Theme.of(context).colorScheme.onSurface,
                  ),
                ],
              ),
            ),
          ),
        );
      },
    );
  }
}

class _ModelItem {
  final LlmProvider provider;
  final Model model;

  _ModelItem({
    required this.provider,
    required this.model,
  });
} 