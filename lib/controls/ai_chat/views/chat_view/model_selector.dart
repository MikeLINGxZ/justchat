import 'package:flutter/material.dart';
import 'package:lemon_tea/models/llm_provider.dart';
import 'package:lemon_tea/models/model.dart';
import 'package:lemon_tea/storage/llm_storage.dart';

class ModelSelector extends StatefulWidget {
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
  State<ModelSelector> createState() => _ModelSelectorState();
}

class _ModelSelectorState extends State<ModelSelector> {
  final MenuController _menuController = MenuController();
  List<LlmProvider> _providers = [];
  Map<String, List<Model>> _providerModels = {};
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    _loadProvidersAndModels();
  }

  Future<void> _loadProvidersAndModels() async {
    setState(() {
      _isLoading = true;
    });

    try {
      // 获取所有供应商
      final providers = await LlmStorage.getAllProviders();
      final Map<String, List<Model>> providerModels = {};

      // 为每个供应商加载模型
      for (final provider in providers) {
        final models = await LlmStorage.getModelsByProviderId(provider.id);
        if (models.isNotEmpty) {
          providerModels[provider.id] = models;
        }
      }

      setState(() {
        _providers = providers.where((p) => providerModels.containsKey(p.id)).toList();
        _providerModels = providerModels;
        _isLoading = false;
      });
    } catch (e) {
      debugPrint('加载供应商和模型失败: $e');
      setState(() {
        _isLoading = false;
      });
    }
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

  Widget _buildModelSubmenu(LlmProvider provider, List<Model> models) {
    return SubmenuButton(
      menuChildren: models.map((model) {
        final isSelected = widget.selectedProviderId == provider.id && 
                          widget.selectedModelId == model.id;
        
        return MenuItemButton(
          onPressed: () {
            widget.onModelSelected?.call(provider.id, model.id);
            _menuController.close();
          },
          child: Row(
            children: [
              if (isSelected) ...[
                const Icon(Icons.check, size: 16),
                const SizedBox(width: 8),
              ] else
                const SizedBox(width: 24),
              Expanded(
                child: Text(
                  model.id,
                  style: TextStyle(
                    fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
                  ),
                ),
              ),
            ],
          ),
        );
      }).toList(),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 12.0, vertical: 8.0),
        child: Row(
          children: [
            Icon(Icons.business, size: 16, color: Theme.of(context).colorScheme.primary),
            const SizedBox(width: 8),
            Text(provider.name),
            const Spacer(),
            const Icon(Icons.arrow_right, size: 16),
          ],
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return MenuAnchor(
      controller: _menuController,
      menuChildren: _isLoading
          ? [
              const Padding(
                padding: EdgeInsets.all(16.0),
                child: Center(
                  child: CircularProgressIndicator(),
                ),
              ),
            ]
          : _providers.map((provider) {
              final models = _providerModels[provider.id] ?? [];
              return _buildModelSubmenu(provider, models);
            }).toList(),
      builder: (context, controller, child) {
        return Material(
          color: Colors.transparent,
          child: InkWell(
            onTap: () {
              if (controller.isOpen) {
                controller.close();
              } else {
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
                        fontSize: 12,
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