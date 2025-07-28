import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/rpc/service.pb.dart';
import 'package:lemon_tea/storage/llm_storage.dart';
import 'package:lemon_tea/utils/cli/client/client.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:lemon_tea/models/llm_provider.dart';
import 'package:lemon_tea/models/model.dart';
import '../model_settings.dart';

void showModelsDialog(BuildContext context, WidgetRef ref, LlmProvider provider) {
  showDialog(
    context: context,
    builder: (context) => _ModelsDialog(provider: provider, ref: ref),
  );
}

class _ModelsDialog extends StatefulWidget {
  final LlmProvider provider;
  final WidgetRef ref;

  const _ModelsDialog({required this.provider, required this.ref});

  @override
  State<_ModelsDialog> createState() => _ModelsDialogState();
}

class _ModelsDialogState extends State<_ModelsDialog> {
  List<Model> availableModels = [];
  bool isVerifying = false;
  String searchKeyword = '';
  final TextEditingController searchController = TextEditingController();
  GlobalKey<State<StatefulWidget>> futureBuilderKey = GlobalKey();

  @override
  void dispose() {
    searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: Row(
        children: [
          Icon(
            Icons.psychology,
            color: Theme.of(context).colorScheme.primary,
            size: 28,
          ),
          const SizedBox(width: 12),
          Text(
            '${widget.provider.name} - 模型列表',
            style: TextStyle(
              fontSize: FontSizeUtils.getSubheadingSize(widget.ref),
              fontWeight: FontWeight.bold,
            ),
          ),
        ],
      ),
      content: SizedBox(
        width: 600,
        height: 500,
        child: FutureBuilder<List<Model>>(
          key: futureBuilderKey,
          future: _loadModels(widget.provider.id),
          builder: (context, snapshot) {
            if (snapshot.connectionState == ConnectionState.waiting) {
              return const Center(child: CircularProgressIndicator());
            }

            if (snapshot.hasError) {
              return Center(
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Icon(
                      Icons.error_outline,
                      size: 48,
                      color: Theme.of(context).colorScheme.error,
                    ),
                    const SizedBox(height: 12),
                    Text(
                      '加载模型失败',
                      style: TextStyle(
                        fontSize: FontSizeUtils.getBodySize(widget.ref),
                        fontWeight: FontWeight.w500,
                        color: Theme.of(context).colorScheme.error,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      '${snapshot.error}',
                      style: TextStyle(
                        fontSize: FontSizeUtils.getSmallSize(widget.ref),
                        color: Theme.of(context).colorScheme.onSurfaceVariant,
                      ),
                      textAlign: TextAlign.center,
                    ),
                  ],
                ),
              );
            }

            final models = snapshot.data ?? [];
            availableModels = models;

            // 根据搜索关键字过滤模型
            final filteredModels = searchKeyword.isEmpty 
                ? availableModels
                : availableModels.where((model) {
                    final keyword = searchKeyword.toLowerCase();
                    return model.id.toLowerCase().contains(keyword) ||
                           model.ownedBy.toLowerCase().contains(keyword);
                  }).toList();

            return Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // 标题栏
                Container(
                  padding: const EdgeInsets.fromLTRB(12, 10, 12, 6),
                  decoration: BoxDecoration(
                    color: Theme.of(context).colorScheme.primaryContainer.withOpacity(0.3),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Row(
                    children: [
                      Icon(
                        Icons.psychology,
                        size: 18,
                        color: Theme.of(context).colorScheme.primary,
                      ),
                      const SizedBox(width: 8),
                      Text(
                        searchKeyword.isNotEmpty
                            ? '搜索结果 (${filteredModels.length}/${availableModels.length})'
                            : availableModels.isNotEmpty
                                ? '可用模型列表 (${availableModels.length}个)'
                                : '模型列表',
                        style: TextStyle(
                          fontSize: FontSizeUtils.getBodySize(widget.ref),
                          fontWeight: FontWeight.w600,
                          color: Theme.of(context).colorScheme.primary,
                        ),
                      ),
                      const Spacer(),
                      if (availableModels.isNotEmpty && searchKeyword.isEmpty) ...[
                        Container(
                          padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                          decoration: BoxDecoration(
                            color: Theme.of(context).colorScheme.primary.withOpacity(0.1),
                            borderRadius: BorderRadius.circular(8),
                          ),
                          child: Text(
                            '${availableModels.where((m) => !m.isCustom).length}个官方',
                            style: TextStyle(
                              fontSize: FontSizeUtils.getSmallSize(widget.ref) - 1,
                              color: Theme.of(context).colorScheme.primary,
                              fontWeight: FontWeight.w500,
                            ),
                          ),
                        ),
                        if (availableModels.any((m) => m.isCustom)) ...[
                          const SizedBox(width: 6),
                          Container(
                            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                            decoration: BoxDecoration(
                              color: Theme.of(context).colorScheme.secondary.withOpacity(0.1),
                              borderRadius: BorderRadius.circular(8),
                            ),
                            child: Text(
                              '${availableModels.where((m) => m.isCustom).length}个自定义',
                              style: TextStyle(
                                fontSize: FontSizeUtils.getSmallSize(widget.ref) - 1,
                                color: Theme.of(context).colorScheme.secondary,
                                fontWeight: FontWeight.w500,
                              ),
                            ),
                          ),
                        ],
                        const SizedBox(width: 8),
                      ],
                      // 添加模型按钮
                      Material(
                        color: Colors.transparent,
                        child: InkWell(
                          onTap: () => _showAddModelDialog(context, widget.ref, widget.provider, setState, availableModels),
                          borderRadius: BorderRadius.circular(6),
                          child: Container(
                            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                            decoration: BoxDecoration(
                              color: Theme.of(context).colorScheme.primaryContainer.withOpacity(0.6),
                              borderRadius: BorderRadius.circular(6),
                              border: Border.all(
                                color: Theme.of(context).colorScheme.primary.withOpacity(0.3),
                                width: 1,
                              ),
                            ),
                            child: Row(
                              mainAxisSize: MainAxisSize.min,
                              children: [
                                Icon(
                                  Icons.add,
                                  size: 14,
                                  color: Theme.of(context).colorScheme.primary,
                                ),
                                const SizedBox(width: 4),
                                Text(
                                  '添加模型',
                                  style: TextStyle(
                                    fontSize: FontSizeUtils.getSmallSize(widget.ref) - 1,
                                    color: Theme.of(context).colorScheme.primary,
                                    fontWeight: FontWeight.w500,
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 12),
                // 搜索框
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 10),
                  child: TextField(
                    controller: searchController,
                    onChanged: (value) {
                      setState(() {
                        searchKeyword = value;
                      });
                    },
                    decoration: InputDecoration(
                      hintText: '搜索模型ID或拥有者...',
                      hintStyle: TextStyle(
                        fontSize: FontSizeUtils.getSmallSize(widget.ref),
                        color: Theme.of(context).colorScheme.onSurfaceVariant.withOpacity(0.6),
                      ),
                      prefixIcon: Icon(
                        Icons.search,
                        size: 20,
                        color: Theme.of(context).colorScheme.onSurfaceVariant,
                      ),
                      suffixIcon: searchKeyword.isNotEmpty
                          ? IconButton(
                              icon: Icon(
                                Icons.clear,
                                size: 18,
                                color: Theme.of(context).colorScheme.onSurfaceVariant,
                              ),
                              onPressed: () {
                                setState(() {
                                  searchKeyword = '';
                                  searchController.clear();
                                });
                              },
                              padding: EdgeInsets.zero,
                              constraints: const BoxConstraints(
                                minWidth: 32,
                                minHeight: 32,
                              ),
                            )
                          : null,
                      border: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(8),
                        borderSide: BorderSide(
                          color: Theme.of(context).colorScheme.outline.withOpacity(0.3),
                          width: 1,
                        ),
                      ),
                      enabledBorder: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(8),
                        borderSide: BorderSide(
                          color: Theme.of(context).colorScheme.outline.withOpacity(0.3),
                          width: 1,
                        ),
                      ),
                      focusedBorder: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(8),
                        borderSide: BorderSide(
                          color: Theme.of(context).colorScheme.primary,
                          width: 1.5,
                        ),
                      ),
                      filled: true,
                      fillColor: Theme.of(context).colorScheme.surface,
                      contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                      isDense: true,
                    ),
                    style: TextStyle(
                      fontSize: FontSizeUtils.getSmallSize(widget.ref),
                    ),
                  ),
                ),
                const SizedBox(height: 12),
                // 模型列表内容
                Expanded(
                  child: filteredModels.isNotEmpty
                      ? ScrollConfiguration(
                    behavior: ScrollConfiguration.of(context).copyWith(
                      scrollbars: false,
                      overscroll: false,
                      physics: const ClampingScrollPhysics(),
                    ),
                    child: ListView.separated(
                      padding: const EdgeInsets.all(10),
                      itemCount: filteredModels.length,
                      separatorBuilder: (context, index) => const SizedBox(height: 6),
                      physics: const ClampingScrollPhysics(),
                      itemBuilder: (context, index) {
                        final model = filteredModels[index];
                        // 找到原始列表中的索引，用于更新状态
                        final originalIndex = availableModels.indexOf(model);
                        return Container(
                          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                          decoration: BoxDecoration(
                            color: model.isCustom
                                ? Theme.of(context).colorScheme.secondaryContainer.withOpacity(0.3)
                                : Theme.of(context).colorScheme.primaryContainer.withOpacity(0.2),
                            borderRadius: BorderRadius.circular(8),
                            border: Border.all(
                              color: model.isCustom
                                  ? Theme.of(context).colorScheme.secondary.withOpacity(0.2)
                                  : Theme.of(context).colorScheme.primary.withOpacity(0.2),
                              width: 1,
                            ),
                          ),
                          child: Row(
                            children: [
                              Container(
                                width: 20,
                                height: 20,
                                decoration: BoxDecoration(
                                  color: model.isCustom
                                      ? Theme.of(context).colorScheme.secondary
                                      : Theme.of(context).colorScheme.primary,
                                  borderRadius: BorderRadius.circular(10),
                                ),
                                child: Icon(
                                  model.isCustom ? Icons.person : Icons.smart_toy,
                                  size: 12,
                                  color: Colors.white,
                                ),
                              ),
                              const SizedBox(width: 12),
                              Expanded(
                                child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    // 高亮搜索关键字
                                    _buildHighlightedText(
                                      model.id,
                                      searchKeyword,
                                      TextStyle(
                                        fontSize: FontSizeUtils.getSmallSize(widget.ref),
                                        fontWeight: FontWeight.w500,
                                        color: Theme.of(context).colorScheme.onSurface,
                                      ),
                                      Theme.of(context).colorScheme.secondary.withOpacity(0.3),
                                    ),
                                    if (model.ownedBy.isNotEmpty && model.ownedBy != 'unknown') ...[
                                      const SizedBox(height: 2),
                                      _buildHighlightedText(
                                        'by ${model.ownedBy}',
                                        searchKeyword,
                                        TextStyle(
                                          fontSize: FontSizeUtils.getSmallSize(widget.ref) - 2,
                                          color: Theme.of(context).colorScheme.onSurfaceVariant,
                                        ),
                                        Theme.of(context).colorScheme.secondary.withOpacity(0.3),
                                      ),
                                    ],
                                  ],
                                ),
                              ),
                              if (model.isCustom)
                                Container(
                                  padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                                  decoration: BoxDecoration(
                                    color: Theme.of(context).colorScheme.secondary.withOpacity(0.2),
                                    borderRadius: BorderRadius.circular(4),
                                  ),
                                  child: Text(
                                    '自定义',
                                    style: TextStyle(
                                      fontSize: FontSizeUtils.getSmallSize(widget.ref) - 2,
                                      color: Theme.of(context).colorScheme.secondary,
                                      fontWeight: FontWeight.bold,
                                    ),
                                  ),
                                ),
                              const SizedBox(width: 8),
                              // 启用开关
                              Transform.scale(
                                scale: 0.8,
                                child: Switch(
                                  value: model.enabled,
                                  onChanged: (value) async {
                                    // 更新模型启用状态
                                    final updatedModel = Model(
                                      llmProviderId: model.llmProviderId,
                                      id: model.id,
                                      object: model.object,
                                      ownedBy: model.ownedBy,
                                      enabled: value,
                                      isCustom: model.isCustom,
                                      seqId: model.seqId,
                                    );

                                    final success = await LlmStorage.updateModel(updatedModel);
                                    if (success) {
                                      setState(() {
                                        availableModels[originalIndex] = updatedModel;
                                      });
                                      // 刷新 provider 数据
                                      widget.ref.refresh(modelsProvider(widget.provider.id));
                                    }
                                  },
                                  materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
                                ),
                              ),
                              const SizedBox(width: 4),
                              // 更多按钮
                              PopupMenuButton<String>(
                                onSelected: (value) async {
                                  if (value == 'edit' && model.isCustom) {
                                    _showEditModelDialog(context, widget.ref, model, widget.provider, setState, availableModels, originalIndex);
                                  } else if (value == 'delete') {
                                    // 删除模型
                                    final confirmed = await showDialog<bool>(
                                      context: context,
                                      builder: (context) => AlertDialog(
                                        title: Text(
                                          '确认删除',
                                          style: TextStyle(
                                            fontSize: FontSizeUtils.getSubheadingSize(widget.ref),
                                          ),
                                        ),
                                        content: Text(
                                          '确定要删除模型 "${model.id}" 吗？',
                                          style: TextStyle(
                                            fontSize: FontSizeUtils.getBodySize(widget.ref),
                                          ),
                                        ),
                                        actions: [
                                          TextButton(
                                            onPressed: () => Navigator.of(context).pop(false),
                                            child: Text(
                                              '取消',
                                              style: TextStyle(
                                                fontSize: FontSizeUtils.getBodySize(widget.ref),
                                              ),
                                            ),
                                          ),
                                          FilledButton(
                                            onPressed: () => Navigator.of(context).pop(true),
                                            child: Text(
                                              '删除',
                                              style: TextStyle(
                                                fontSize: FontSizeUtils.getBodySize(widget.ref),
                                              ),
                                            ),
                                          ),
                                        ],
                                      ),
                                    );

                                    if (confirmed == true) {
                                      final success = await LlmStorage.deleteModel(model.id);
                                      if (success) {
                                        setState(() {
                                          availableModels.removeAt(originalIndex);
                                        });
                                        // 刷新 provider 数据
                                        widget.ref.refresh(modelsProvider(widget.provider.id));
                                        ScaffoldMessenger.of(context).showSnackBar(
                                          SnackBar(
                                            content: Text(
                                              '模型已删除',
                                              style: TextStyle(
                                                fontSize: FontSizeUtils.getBodySize(widget.ref),
                                              ),
                                            ),
                                            backgroundColor: Theme.of(context).colorScheme.primaryContainer,
                                          ),
                                        );
                                      }
                                    }
                                  }
                                },
                                itemBuilder: (context) => [
                                  PopupMenuItem<String>(
                                    value: 'edit',
                                    enabled: model.isCustom,
                                    child: Row(
                                      children: [
                                        Icon(
                                          Icons.edit,
                                          size: 16,
                                          color: model.isCustom
                                              ? Theme.of(context).colorScheme.onSurface
                                              : Theme.of(context).colorScheme.onSurface.withOpacity(0.3),
                                        ),
                                        const SizedBox(width: 8),
                                        Text(
                                          '编辑',
                                          style: TextStyle(
                                            fontSize: FontSizeUtils.getSmallSize(widget.ref),
                                            color: model.isCustom
                                                ? Theme.of(context).colorScheme.onSurface
                                                : Theme.of(context).colorScheme.onSurface.withOpacity(0.3),
                                          ),
                                        ),
                                      ],
                                    ),
                                  ),
                                  PopupMenuItem<String>(
                                    value: 'delete',
                                    child: Row(
                                      children: [
                                        Icon(
                                          Icons.delete,
                                          size: 16,
                                          color: Theme.of(context).colorScheme.error,
                                        ),
                                        const SizedBox(width: 8),
                                        Text(
                                          '删除',
                                          style: TextStyle(
                                            fontSize: FontSizeUtils.getSmallSize(widget.ref),
                                            color: Theme.of(context).colorScheme.error,
                                          ),
                                        ),
                                      ],
                                    ),
                                  ),
                                ],
                                icon: Icon(
                                  Icons.more_vert,
                                  size: 16,
                                  color: Theme.of(context).colorScheme.onSurface.withOpacity(0.6),
                                ),
                                padding: EdgeInsets.zero,
                                iconSize: 16,
                              ),
                            ],
                          ),
                        );
                      },
                    ),
                  )
                      : Container(
                    padding: const EdgeInsets.all(40),
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        Icon(
                          searchKeyword.isNotEmpty ? Icons.search_off : Icons.psychology_outlined,
                          size: 64,
                          color: Theme.of(context).colorScheme.onSurfaceVariant.withOpacity(0.5),
                        ),
                        const SizedBox(height: 16),
                        Text(
                          searchKeyword.isNotEmpty ? '未找到匹配的模型' : '暂无模型数据',
                          style: TextStyle(
                            fontSize: FontSizeUtils.getBodySize(widget.ref),
                            fontWeight: FontWeight.w500,
                            color: Theme.of(context).colorScheme.onSurfaceVariant,
                          ),
                        ),
                        const SizedBox(height: 8),
                        Text(
                          searchKeyword.isNotEmpty 
                              ? '尝试使用其他关键字搜索，或清空搜索条件'
                              : '点击右上角"添加模型"按钮可添加自定义模型',
                          style: TextStyle(
                            fontSize: FontSizeUtils.getSmallSize(widget.ref),
                            color: Theme.of(context).colorScheme.onSurfaceVariant.withOpacity(0.7),
                          ),
                          textAlign: TextAlign.center,
                        ),
                      ],
                    ),
                  ),
                ),
              ],
            );
          },
        ),
      ),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
      ),
      actions: [
        TextButton.icon(
          onPressed: isVerifying ? null : () async {
            setState(() {
              isVerifying = true;
            });

            try {
              final request = ModelsRequest(
                  name: widget.provider.name,
                  apiKey: widget.provider.apiKey,
                  baseUrl: widget.provider.baseUrl
              );
              final response = await Client().stub!.models(request);

              // 更新数据库中的模型列表
              await _updateProviderModels(widget.provider.id, response.models);

              // 强制重新构建FutureBuilder
              setState(() {
                futureBuilderKey = GlobalKey();
                isVerifying = false;
              });

              // 刷新 provider 数据
              widget.ref.refresh(modelsProvider(widget.provider.id));

              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(
                  content: Text(
                    '模型列表已更新! 发现 ${response.models.length} 个模型',
                    style: TextStyle(
                      fontSize: FontSizeUtils.getBodySize(widget.ref),
                    ),
                  ),
                  backgroundColor: Theme.of(context).colorScheme.primaryContainer,
                ),
              );
            } catch (e) {
              setState(() {
                isVerifying = false;
              });

              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(
                  content: Text(
                    '重新验证失败: ${e.toString()}',
                    style: TextStyle(
                      fontSize: FontSizeUtils.getBodySize(widget.ref),
                    ),
                  ),
                  backgroundColor: Theme.of(context).colorScheme.errorContainer,
                ),
              );
            }
          },
          icon: isVerifying
              ? SizedBox(
            width: 16,
            height: 16,
            child: CircularProgressIndicator(
              strokeWidth: 2,
              color: Theme.of(context).colorScheme.primary,
            ),
          )
              : const Icon(Icons.refresh),
          label: Text(
            isVerifying ? '验证中...' : '重新验证',
            style: TextStyle(
              fontSize: FontSizeUtils.getBodySize(widget.ref),
            ),
          ),
          style: TextButton.styleFrom(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(8),
            ),
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
          ),
        ),
        TextButton.icon(
          onPressed: () => Navigator.of(context).pop(),
          icon: const Icon(Icons.close),
          label: Text(
            '关闭',
            style: TextStyle(
              fontSize: FontSizeUtils.getBodySize(widget.ref),
            ),
          ),
          style: TextButton.styleFrom(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(8),
            ),
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
          ),
        ),
      ],
      actionsPadding: const EdgeInsets.fromLTRB(16, 0, 16, 16),
    );
  }
}

// 加载模型数据
Future<List<Model>> _loadModels(String providerId) async {
  try {
    final models = await LlmStorage.getModelsByProviderId(providerId);
    // 排序：非自定义模型在前，自定义模型在后
    models.sort((a, b) {
      if (a.isCustom == b.isCustom) {
        return a.id.compareTo(b.id);
      }
      return a.isCustom ? 1 : -1;
    });
    return models;
  } catch (e) {
    debugPrint('加载模型失败: $e');
    rethrow;
  }
}

// 显示添加模型对话框
void _showAddModelDialog(BuildContext context, WidgetRef ref, LlmProvider provider,
    void Function(void Function()) setState, List<Model> availableModels) {
  final TextEditingController modelIdController = TextEditingController();
  final TextEditingController modelNameController = TextEditingController();
  final TextEditingController ownedByController = TextEditingController(text: 'custom');
  bool modelEnabled = true;

  showDialog(
    context: context,
    builder: (context) => StatefulBuilder(
      builder: (context, setDialogState) => AlertDialog(
        title: Row(
          children: [
            Icon(
              Icons.add_circle,
              color: Theme.of(context).colorScheme.secondary,
              size: 24,
            ),
            const SizedBox(width: 12),
            Text(
              '添加自定义模型',
              style: TextStyle(
                fontSize: FontSizeUtils.getSubheadingSize(ref),
                fontWeight: FontWeight.bold,
              ),
            ),
          ],
        ),
        content: SizedBox(
          width: 400,
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              TextField(
                controller: modelIdController,
                decoration: InputDecoration(
                  labelText: '模型ID *',
                  hintText: '例如: my-custom-model',
                  prefixIcon: const Icon(Icons.psychology),
                  labelStyle: TextStyle(
                    fontSize: FontSizeUtils.getBodySize(ref),
                  ),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(8),
                  ),
                ),
                style: TextStyle(
                  fontSize: FontSizeUtils.getBodySize(ref),
                ),
              ),
              const SizedBox(height: 16),
              TextField(
                controller: modelNameController,
                decoration: InputDecoration(
                  labelText: '模型名称',
                  hintText: '例如: My Custom Model (可选，默认使用ID)',
                  prefixIcon: const Icon(Icons.label),
                  labelStyle: TextStyle(
                    fontSize: FontSizeUtils.getBodySize(ref),
                  ),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(8),
                  ),
                ),
                style: TextStyle(
                  fontSize: FontSizeUtils.getBodySize(ref),
                ),
              ),
              const SizedBox(height: 16),
              TextField(
                controller: ownedByController,
                decoration: InputDecoration(
                  labelText: '拥有者',
                  hintText: '例如: custom, user, organization',
                  prefixIcon: const Icon(Icons.person),
                  labelStyle: TextStyle(
                    fontSize: FontSizeUtils.getBodySize(ref),
                  ),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(8),
                  ),
                ),
                style: TextStyle(
                  fontSize: FontSizeUtils.getBodySize(ref),
                ),
              ),
              const SizedBox(height: 16),
              SwitchListTile(
                title: Text(
                  '启用此模型',
                  style: TextStyle(
                    fontSize: FontSizeUtils.getBodySize(ref),
                  ),
                ),
                subtitle: Text(
                  '关闭后将不会在模型选择中显示',
                  style: TextStyle(
                    fontSize: FontSizeUtils.getSmallSize(ref),
                    color: Theme.of(context).colorScheme.onSurfaceVariant,
                  ),
                ),
                value: modelEnabled,
                onChanged: (value) {
                  setDialogState(() {
                    modelEnabled = value;
                  });
                },
                secondary: Icon(
                  modelEnabled ? Icons.toggle_on : Icons.toggle_off,
                  color: modelEnabled ? Theme.of(context).colorScheme.secondary : null,
                ),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(8),
                ),
              ),
            ],
          ),
        ),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: Text(
              '取消',
              style: TextStyle(
                fontSize: FontSizeUtils.getBodySize(ref),
              ),
            ),
          ),
          FilledButton(
            onPressed: () async {
              final modelId = modelIdController.text.trim();
              final modelName = modelNameController.text.trim();
              final ownedBy = ownedByController.text.trim();

              if (modelId.isEmpty) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text(
                      '请输入模型ID',
                      style: TextStyle(
                        fontSize: FontSizeUtils.getBodySize(ref),
                      ),
                    ),
                  ),
                );
                return;
              }

              // 检查模型ID是否已存在
              final existingModel = availableModels.firstWhere(
                    (model) => model.id == modelId,
                orElse: () => Model(
                  llmProviderId: '',
                  id: '',
                  object: '',
                  ownedBy: '',
                  enabled: false,
                  isCustom: false,
                  seqId: 0,
                ),
              );

              if (existingModel.id.isNotEmpty) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text(
                      '模型ID已存在，请使用其他ID',
                      style: TextStyle(
                        fontSize: FontSizeUtils.getBodySize(ref),
                      ),
                    ),
                  ),
                );
                return;
              }

              try {
                // 获取最大序号
                final maxSeqId = await LlmStorage.getMaxModelSeqId(provider.id);

                // 创建模型数据
                final modelMap = {
                  'llm_provider_id': provider.id,
                  'id': modelId,
                  'name': modelName.isEmpty ? modelId : modelName,
                  'object': 'model',
                  'owned_by': ownedBy.isEmpty ? 'custom' : ownedBy,
                  'enabled': modelEnabled ? 1 : 0,
                  'is_custom': 1, // 标记为自定义模型
                  'seq_id': maxSeqId + 1,
                };

                // 添加到数据库
                await LlmStorage.addModelWithCustomFields(modelMap);

                // 创建Model对象并添加到列表
                final newModel = Model(
                  llmProviderId: provider.id,
                  id: modelId,
                  object: 'model',
                  ownedBy: ownedBy.isEmpty ? 'custom' : ownedBy,
                  enabled: modelEnabled,
                  isCustom: true,
                  seqId: maxSeqId + 1,
                );

                // 更新UI状态
                setState(() {
                  availableModels.add(newModel);
                  // 重新排序：非自定义模型在前，自定义模型在后
                  availableModels.sort((a, b) {
                    if (a.isCustom == b.isCustom) {
                      return a.id.compareTo(b.id);
                    }
                    return a.isCustom ? 1 : -1;
                  });
                });

                // 刷新 provider 数据
                ref.refresh(modelsProvider(provider.id));

                Navigator.of(context).pop();

                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text(
                      '自定义模型已添加',
                      style: TextStyle(
                        fontSize: FontSizeUtils.getBodySize(ref),
                      ),
                    ),
                    backgroundColor: Theme.of(context).colorScheme.primaryContainer,
                  ),
                );
              } catch (e) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text(
                      '添加模型失败: $e',
                      style: TextStyle(
                        fontSize: FontSizeUtils.getBodySize(ref),
                      ),
                    ),
                    backgroundColor: Theme.of(context).colorScheme.errorContainer,
                  ),
                );
              }
            },
            child: Text(
              '添加',
              style: TextStyle(
                fontSize: FontSizeUtils.getBodySize(ref),
              ),
            ),
          ),
        ],
      ),
    ),
  );
}

// 更新供应商模型列表的私有函数
Future<void> _updateProviderModels(String providerId, List<dynamic> apiModels) async {
  try {
    // 1. 获取该供应商下的所有模型
    final existingModels = await LlmStorage.getModelsByProviderId(providerId);

    // 2. 删除所有非自定义模型（isCustom = false）
    final nonCustomModels = existingModels.where((model) => !model.isCustom).toList();
    for (final model in nonCustomModels) {
      await LlmStorage.deleteModel(model.id);
    }

    // 3. 获取当前最大序号，用于新模型的序号分配
    int maxSeqId = await LlmStorage.getMaxModelSeqId(providerId);

    // 4. 添加新获取的模型（标记为非自定义）
    for (int i = 0; i < apiModels.length; i++) {
      final apiModel = apiModels[i];

      // 直接构造包含name字段的Map，避免Model类缺少name字段的问题
      final modelMap = {
        'llm_provider_id': providerId,
        'id': apiModel.id,
        'name': apiModel.id, // 使用模型id作为name
        'object': apiModel.object ?? 'model',
        'owned_by': apiModel.ownedBy ?? 'unknown',
        'enabled': (apiModel.enabled ?? true) ? 1 : 0,
        'is_custom': 0, // 标记为非自定义模型
        'seq_id': maxSeqId + i + 1, // 分配序号
      };

      await LlmStorage.addModelWithCustomFields(modelMap);
    }

    debugPrint('成功更新供应商 $providerId 的模型列表，删除 ${nonCustomModels.length} 个旧模型，添加 ${apiModels.length} 个新模型');
  } catch (e) {
    debugPrint('更新供应商模型列表失败: $e');
    rethrow;
  }
}

// 显示编辑模型对话框
void _showEditModelDialog(BuildContext context, WidgetRef ref, Model model, LlmProvider provider,
    void Function(void Function()) setState, List<Model> availableModels, int index) {
  final TextEditingController modelIdController = TextEditingController(text: model.id);
  final TextEditingController ownedByController = TextEditingController(text: model.ownedBy);
  bool modelEnabled = model.enabled;

  showDialog(
    context: context,
    builder: (context) => StatefulBuilder(
      builder: (context, setDialogState) => AlertDialog(
        title: Row(
          children: [
            Icon(
              Icons.edit,
              color: Theme.of(context).colorScheme.secondary,
              size: 24,
            ),
            const SizedBox(width: 12),
            Text(
              '编辑自定义模型',
              style: TextStyle(
                fontSize: FontSizeUtils.getSubheadingSize(ref),
                fontWeight: FontWeight.bold,
              ),
            ),
          ],
        ),
        content: SizedBox(
          width: 400,
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              TextField(
                controller: modelIdController,
                decoration: InputDecoration(
                  labelText: '模型ID *',
                  prefixIcon: const Icon(Icons.psychology),
                  labelStyle: TextStyle(
                    fontSize: FontSizeUtils.getBodySize(ref),
                  ),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(8),
                  ),
                ),
                style: TextStyle(
                  fontSize: FontSizeUtils.getBodySize(ref),
                ),
              ),
              const SizedBox(height: 16),
              TextField(
                controller: ownedByController,
                decoration: InputDecoration(
                  labelText: '拥有者',
                  prefixIcon: const Icon(Icons.person),
                  labelStyle: TextStyle(
                    fontSize: FontSizeUtils.getBodySize(ref),
                  ),
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(8),
                  ),
                ),
                style: TextStyle(
                  fontSize: FontSizeUtils.getBodySize(ref),
                ),
              ),
              const SizedBox(height: 16),
              SwitchListTile(
                title: Text(
                  '启用此模型',
                  style: TextStyle(
                    fontSize: FontSizeUtils.getBodySize(ref),
                  ),
                ),
                subtitle: Text(
                  '关闭后将不会在模型选择中显示',
                  style: TextStyle(
                    fontSize: FontSizeUtils.getSmallSize(ref),
                    color: Theme.of(context).colorScheme.onSurfaceVariant,
                  ),
                ),
                value: modelEnabled,
                onChanged: (value) {
                  setDialogState(() {
                    modelEnabled = value;
                  });
                },
                secondary: Icon(
                  modelEnabled ? Icons.toggle_on : Icons.toggle_off,
                  color: modelEnabled ? Theme.of(context).colorScheme.secondary : null,
                ),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(8),
                ),
              ),
            ],
          ),
        ),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: Text(
              '取消',
              style: TextStyle(
                fontSize: FontSizeUtils.getBodySize(ref),
              ),
            ),
          ),
          FilledButton(
            onPressed: () async {
              final modelId = modelIdController.text.trim();
              final ownedBy = ownedByController.text.trim();

              if (modelId.isEmpty) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text(
                      '请输入模型ID',
                      style: TextStyle(
                        fontSize: FontSizeUtils.getBodySize(ref),
                      ),
                    ),
                  ),
                );
                return;
              }

              try {
                // 更新模型
                final updatedModel = Model(
                  llmProviderId: model.llmProviderId,
                  id: modelId,
                  object: model.object,
                  ownedBy: ownedBy.isEmpty ? 'custom' : ownedBy,
                  enabled: modelEnabled,
                  isCustom: model.isCustom,
                  seqId: model.seqId,
                );

                final success = await LlmStorage.updateModel(updatedModel);
                if (success) {
                  setState(() {
                    availableModels[index] = updatedModel;
                  });
                  // 刷新 provider 数据
                  ref.refresh(modelsProvider(provider.id));
                }

                Navigator.of(context).pop();

                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text(
                      success ? '模型已更新' : '更新模型失败',
                      style: TextStyle(
                        fontSize: FontSizeUtils.getBodySize(ref),
                      ),
                    ),
                    backgroundColor: success
                        ? Theme.of(context).colorScheme.primaryContainer
                        : Theme.of(context).colorScheme.errorContainer,
                  ),
                );
              } catch (e) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text(
                      '更新模型失败: $e',
                      style: TextStyle(
                        fontSize: FontSizeUtils.getBodySize(ref),
                      ),
                    ),
                    backgroundColor: Theme.of(context).colorScheme.errorContainer,
                  ),
                );
              }
            },
            child: Text(
              '保存',
              style: TextStyle(
                fontSize: FontSizeUtils.getBodySize(ref),
              ),
            ),
          ),
        ],
      ),
    ),
  );
}

// 辅助函数：高亮文本
Widget _buildHighlightedText(String text, String keyword, TextStyle style, Color highlightColor) {
  if (keyword.isEmpty) {
    return Text(text, style: style);
  }
  
  final List<TextSpan> spans = [];
  final RegExp regExp = RegExp(RegExp.escape(keyword), caseSensitive: false);
  final matches = regExp.allMatches(text);

  int lastMatchEnd = 0;
  for (final match in matches) {
    if (match.start > lastMatchEnd) {
      spans.add(TextSpan(text: text.substring(lastMatchEnd, match.start)));
    }
    spans.add(TextSpan(
      text: text.substring(match.start, match.end),
      style: style.copyWith(backgroundColor: highlightColor),
    ));
    lastMatchEnd = match.end;
  }
  if (lastMatchEnd < text.length) {
    spans.add(TextSpan(text: text.substring(lastMatchEnd)));
  }
  
  return RichText(
    text: TextSpan(
      children: spans,
      style: style,
    ),
  );
}