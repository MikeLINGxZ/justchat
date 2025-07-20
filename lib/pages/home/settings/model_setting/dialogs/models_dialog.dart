import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/storage/llm_storage.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:lemon_tea/models/llm_provider.dart';
import 'package:lemon_tea/models/model.dart';
import '../model_settings.dart';

void showModelsDialog(
  BuildContext context, 
  WidgetRef ref, 
  LlmProvider provider,
  Function(String) preloadModels,
  bool Function(Model) getModelEnabledState,
  void Function(Model, bool) updateModelEnabledState,
) {
  // 确保模型数据已预加载
  preloadModels(provider.id);
  
  // 创建本地状态副本，用于UI显示
  final Map<String, bool> localModelStates = {};
  
  // 预先加载模型状态
  ref.read(modelsProvider(provider.id)).whenData((models) {
    for (final model in models) {
      final key = '${model.llmProviderId}_${model.id}';
      localModelStates[key] = getModelEnabledState(model);
    }
  });
  
  showDialog(
    context: context,
    builder: (context) => StatefulBuilder(
      builder: (context, dialogSetState) {
        // 使用Consumer直接访问模型数据，避免loading状态
        return AlertDialog(
          title: Text(
            '${provider.name} 模型列表',
            style: TextStyle(
              fontSize: FontSizeUtils.getSubheadingSize(ref),
              fontWeight: FontWeight.bold,
            ),
          ),
          content: SizedBox(
            width: 460, // 设置固定宽度，使对话框更窄
            height: 400, // 设置固定高度，确保对话框不会过大
            child: Consumer(
              builder: (context, ref, child) {
                // 强制刷新模型数据
                final modelsAsync = ref.watch(modelsProvider(provider.id));
                
                return modelsAsync.when(
                  data: (models) {
                    if (models.isEmpty) {
                      return Center(
                        child: Text(
                          '暂无模型',
                          style: TextStyle(
                            fontSize: FontSizeUtils.getBodySize(ref),
                          ),
                        ),
                      );
                    }
                    
                    return ListView.builder(
                      shrinkWrap: false, // 不收缩，允许滚动
                      physics: const AlwaysScrollableScrollPhysics(), // 始终可滚动
                      itemCount: models.length,
                      itemBuilder: (context, index) {
                        final model = models[index];
                        final key = '${model.llmProviderId}_${model.id}';
                        // 确保本地状态存在
                        if (!localModelStates.containsKey(key)) {
                          localModelStates[key] = getModelEnabledState(model);
                        }
                        
                        return ListTile(
                          title: Text(
                            model.id,
                            style: TextStyle(
                              fontSize: FontSizeUtils.getBodySize(ref),
                            ),
                          ),
                          subtitle: Text(
                            '提供者: ${model.ownedBy}',
                            style: TextStyle(
                              fontSize: FontSizeUtils.getSmallSize(ref),
                              color: Theme.of(context).colorScheme.onSurfaceVariant,
                            ),
                          ),
                          trailing: SizedBox(
                            height: 48, // 固定高度确保垂直居中
                            child: Row(
                              mainAxisSize: MainAxisSize.min,
                              mainAxisAlignment: MainAxisAlignment.center, // 水平居中
                              crossAxisAlignment: CrossAxisAlignment.center, // 垂直居中
                              children: [
                                // 自定义标签放在编辑按钮前面
                                if (model.isCustom)
                                  Container(
                                    margin: const EdgeInsets.only(right: 8),
                                    padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                                    decoration: BoxDecoration(
                                      color: Theme.of(context).colorScheme.primaryContainer,
                                      borderRadius: BorderRadius.circular(4),
                                    ),
                                    child: Text(
                                      '自定义',
                                      style: TextStyle(
                                        fontSize: FontSizeUtils.getSmallSize(ref) - 1,
                                        color: Theme.of(context).colorScheme.onPrimaryContainer,
                                      ),
                                    ),
                                  ),
                                Transform.scale(
                                  scale: 0.8,
                                  child: Switch(
                                    value: localModelStates[key]!,
                                    onChanged: (value) {
                                      // 更新本地状态和UI
                                      dialogSetState(() {
                                        localModelStates[key] = value;
                                      });
                                      // 更新数据库
                                      updateModelEnabledState(model, value);
                                    },
                                  ),
                                ),
                                // 所有模型都显示更多操作按钮
                                SizedBox(
                                  height: 40,
                                  width: 40,
                                  child: PopupMenuButton<String>(
                                    padding: EdgeInsets.zero,
                                    icon: const Icon(Icons.more_vert, size: 20),
                                    tooltip: '更多操作',
                                    offset: const Offset(0, 10),
                                    position: PopupMenuPosition.under,
                                    itemBuilder: (context) => [
                                      PopupMenuItem<String>(
                                        value: 'edit',
                                        enabled: model.isCustom,
                                        child: Row(
                                          children: [
                                            Icon(Icons.edit, 
                                              size: 18, 
                                              color: model.isCustom 
                                                  ? null 
                                                  : Theme.of(context).colorScheme.onSurface.withOpacity(0.38),
                                            ),
                                            const SizedBox(width: 8),
                                            Text(
                                              '编辑',
                                              style: TextStyle(
                                                fontSize: FontSizeUtils.getBodySize(ref),
                                                color: model.isCustom 
                                                    ? null 
                                                    : Theme.of(context).colorScheme.onSurface.withOpacity(0.38),
                                              ),
                                            ),
                                          ],
                                        ),
                                      ),
                                      PopupMenuItem<String>(
                                        value: 'delete',
                                        child: Row(
                                          children: [
                                            const Icon(Icons.delete, color: Colors.red, size: 18),
                                            const SizedBox(width: 8),
                                            Text(
                                              '删除', 
                                              style: TextStyle(
                                                fontSize: FontSizeUtils.getBodySize(ref),
                                                color: Colors.red,
                                              ),
                                            ),
                                          ],
                                        ),
                                      ),
                                    ],
                                    onSelected: (value) {
                                      if (value == 'edit' && model.isCustom) {
                                        Navigator.of(context).pop();
                                        _showEditModelDialog(context, ref, model, provider);
                                      } else if (value == 'delete') {
                                        _showDeleteModelDialog(model, context, ref);
                                      }
                                    },
                                  ),
                                ),
                              ],
                            ),
                          ),
                          contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
                        );
                      },
                    );
                  },
                  loading: () => const Center(child: CircularProgressIndicator()),
                  error: (error, stack) => Center(
                    child: Text(
                      '加载模型失败: $error',
                      style: TextStyle(
                        fontSize: FontSizeUtils.getBodySize(ref),
                        color: Theme.of(context).colorScheme.error,
                      ),
                    ),
                  ),
                );
              },
            ),
          ),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(8), // 设置弹窗圆角为8
          ),
          actions: [
            TextButton(
              onPressed: () {
                // 不再关闭当前对话框，直接打开添加模型对话框
                _showAddModelDialog(context, ref, provider.id, provider);
              },
              style: TextButton.styleFrom(
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(8),
                ),
              ),
              child: Text(
                '添加模型',
                style: TextStyle(
                  fontSize: FontSizeUtils.getBodySize(ref),
                ),
              ),
            ),
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              style: TextButton.styleFrom(
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(8),
                ),
              ),
              child: Text(
                '关闭',
                style: TextStyle(
                  fontSize: FontSizeUtils.getBodySize(ref),
                ),
              ),
            ),
          ],
        );
      }
    ),
  );
}

// 显示添加模型对话框
void _showAddModelDialog(BuildContext context, WidgetRef ref, String providerId, LlmProvider provider) {
  final TextEditingController modelIdController = TextEditingController();
  final TextEditingController modelNameController = TextEditingController();
  final TextEditingController ownedByController = TextEditingController();
  bool isEnabled = true;

  showDialog(
    context: context,
    builder: (dialogContext) => AlertDialog(
      title: Text(
        '添加自定义模型',
        style: TextStyle(
          fontSize: FontSizeUtils.getSubheadingSize(ref),
          fontWeight: FontWeight.bold,
        ),
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
                hintText: '例如: gpt-4-turbo',
                labelStyle: TextStyle(
                  fontSize: FontSizeUtils.getBodySize(ref),
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
                labelText: '模型名称 *',
                hintText: '例如: GPT-4 Turbo',
                labelStyle: TextStyle(
                  fontSize: FontSizeUtils.getBodySize(ref),
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
                labelText: '提供者 *',
                hintText: '例如: OpenAI',
                labelStyle: TextStyle(
                  fontSize: FontSizeUtils.getBodySize(ref),
                ),
              ),
              style: TextStyle(
                fontSize: FontSizeUtils.getBodySize(ref),
              ),
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Text(
                  '启用状态',
                  style: TextStyle(
                    fontSize: FontSizeUtils.getBodySize(ref),
                  ),
                ),
                const Spacer(),
                StatefulBuilder(
                  builder: (BuildContext context, StateSetter setState) {
                    return Switch(
                      value: isEnabled,
                      onChanged: (value) {
                        setState(() {
                          isEnabled = value;
                        });
                      },
                    );
                  },
                ),
              ],
            ),
          ],
        ),
      ),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(8),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(dialogContext).pop(),
          style: TextButton.styleFrom(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(8),
            ),
          ),
          child: Text(
            '取消',
            style: TextStyle(
              fontSize: FontSizeUtils.getBodySize(ref),
            ),
          ),
        ),
        TextButton(
          onPressed: () async {
            // 验证输入
            final modelId = modelIdController.text.trim();
            final modelName = modelNameController.text.trim();
            final ownedBy = ownedByController.text.trim();
            
            if (modelId.isEmpty || modelName.isEmpty || ownedBy.isEmpty) {
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(
                  content: Text(
                    '请填写所有必填字段',
                    style: TextStyle(
                      fontSize: FontSizeUtils.getBodySize(ref),
                    ),
                  ),
                ),
              );
              return;
            }
            
            // 创建模型对象
            final newModel = Model(
              llmProviderId: providerId,
              id: modelId,
              ownedBy: ownedBy,
              enabled: isEnabled,
              isCustom: true,
            );
            
            // 添加模型到数据库，使用自定义方法处理name字段
            final success = await _addModelWithName(newModel, modelName);
            Navigator.of(dialogContext).pop();
            
            if (success) {
              // 刷新模型列表
              ref.refresh(modelsProvider(providerId));
              
              // 延迟一点时间后重新打开模型列表对话框以显示更新后的列表
              Future.delayed(const Duration(milliseconds: 300), () async {
                if (context.mounted) {
                  showModelsDialog(
                    context, 
                    ref, 
                    provider, 
                    (String id) {}, // 空实现，因为已经刷新了
                    (Model model) {
                      final key = '${model.llmProviderId}_${model.id}';
                      return model.enabled; // 直接返回模型状态
                    },
                    (Model model, bool value) {
                      // 这里可以实现状态更新逻辑
                    },
                  );
                }
              });
            } else {
              if (context.mounted) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text(
                      '添加模型失败，可能模型ID已存在',
                      style: TextStyle(
                        fontSize: FontSizeUtils.getBodySize(ref),
                      ),
                    ),
                  ),
                );
              }
            }
          },
          style: TextButton.styleFrom(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(8),
            ),
          ),
          child: Text(
            '添加',
            style: TextStyle(
              fontSize: FontSizeUtils.getBodySize(ref),
            ),
          ),
        ),
      ],
    ),
  );
}

// 添加模型到数据库，处理name字段
Future<bool> _addModelWithName(Model model, String name) async {
  try {
    // 获取模型的Map数据
    final modelMap = model.toMap();
    // 添加name字段
    modelMap['name'] = name;
    
    // 获取当前最大seq_id并设置新模型的seq_id为最大值+1
    final maxSeqId = await LlmStorage.getMaxModelSeqId(model.llmProviderId);
    modelMap['seq_id'] = maxSeqId + 1;
    
    // 插入数据库
    final result = await LlmStorage.addModelWithCustomFields(modelMap);
    return result;
  } catch (e) {
    debugPrint('添加模型失败: $e');
    return false;
  }
}

// 显示编辑模型对话框
void _showEditModelDialog(BuildContext context, WidgetRef ref, Model model, LlmProvider provider) {
  final TextEditingController modelIdController = TextEditingController(text: model.id);
  final TextEditingController ownedByController = TextEditingController(text: model.ownedBy);
  bool isEnabled = model.enabled;

  showDialog(
    context: context,
    builder: (context) => AlertDialog(
      title: Text(
        '编辑模型',
        style: TextStyle(
          fontSize: FontSizeUtils.getSubheadingSize(ref),
          fontWeight: FontWeight.bold,
        ),
      ),
      content: SizedBox(
        width: 400,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: modelIdController,
              decoration: InputDecoration(
                labelText: '模型ID',
                labelStyle: TextStyle(
                  fontSize: FontSizeUtils.getBodySize(ref),
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
                labelText: '提供者',
                labelStyle: TextStyle(
                  fontSize: FontSizeUtils.getBodySize(ref),
                ),
              ),
              style: TextStyle(
                fontSize: FontSizeUtils.getBodySize(ref),
              ),
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Text(
                  '启用状态',
                  style: TextStyle(
                    fontSize: FontSizeUtils.getBodySize(ref),
                  ),
                ),
                const Spacer(),
                StatefulBuilder(
                  builder: (BuildContext context, StateSetter setState) {
                    return Switch(
                      value: isEnabled,
                      onChanged: (value) {
                        setState(() {
                          isEnabled = value;
                        });
                      },
                    );
                  },
                ),
              ],
            ),
          ],
        ),
      ),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(8),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(context).pop(),
          style: TextButton.styleFrom(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(8),
            ),
          ),
          child: Text(
            '取消',
            style: TextStyle(
              fontSize: FontSizeUtils.getBodySize(ref),
            ),
          ),
        ),
        TextButton(
          onPressed: () async {
            final updatedModel = Model(
              llmProviderId: model.llmProviderId,
              id: modelIdController.text.trim(),
              object: model.object,
              ownedBy: ownedByController.text.trim(),
              enabled: isEnabled,
              isCustom: model.isCustom,
            );
            
            final success = await LlmStorage.updateModel(updatedModel);
            Navigator.of(context).pop();
            
            if (success) {
              // 刷新模型列表
              ref.refresh(modelsProvider(model.llmProviderId));
              
              // 重新打开模型列表对话框
              Future.delayed(const Duration(milliseconds: 300), () async {
                if (context.mounted) {
                  showModelsDialog(
                    context, 
                    ref, 
                    provider, 
                    (String id) {}, // 空实现
                    (Model model) => model.enabled,
                    (Model model, bool value) {},
                  );
                }
              });
            } else {
              if (context.mounted) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text(
                      '更新模型失败',
                      style: TextStyle(
                        fontSize: FontSizeUtils.getBodySize(ref),
                      ),
                    ),
                  ),
                );
              }
            }
          },
          style: TextButton.styleFrom(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(8),
            ),
          ),
          child: Text(
            '保存',
            style: TextStyle(
              fontSize: FontSizeUtils.getBodySize(ref),
            ),
          ),
        ),
      ],
    ),
  );
}

// 显示删除模型对话框
void _showDeleteModelDialog(Model model, BuildContext parentContext, WidgetRef ref) {
  showDialog(
    context: parentContext, // 使用父对话框的context，而不是全局context
    builder: (dialogContext) => AlertDialog(
      title: Text(
        '删除模型',
        style: TextStyle(
          fontSize: FontSizeUtils.getSubheadingSize(ref),
          fontWeight: FontWeight.bold,
        ),
      ),
      content: Text(
        '确定要删除模型 ${model.id} 吗？此操作不可恢复。',
        style: TextStyle(
          fontSize: FontSizeUtils.getBodySize(ref),
        ),
      ),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(8),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.of(dialogContext).pop(),
          style: TextButton.styleFrom(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(8),
            ),
          ),
          child: Text(
            '取消',
            style: TextStyle(
              fontSize: FontSizeUtils.getBodySize(ref),
            ),
          ),
        ),
        TextButton(
          onPressed: () async {
            final success = await LlmStorage.deleteModel(model.id);
            Navigator.of(dialogContext).pop();
            
            if (success) {
              // 刷新模型列表
              ref.refresh(modelsProvider(model.llmProviderId));
              
              ScaffoldMessenger.of(parentContext).showSnackBar(
                SnackBar(
                  content: Text(
                    '模型已删除',
                    style: TextStyle(
                      fontSize: FontSizeUtils.getBodySize(ref),
                    ),
                  ),
                  duration: const Duration(seconds: 2),
                ),
              );
            } else {
              if (parentContext.mounted) {
                ScaffoldMessenger.of(parentContext).showSnackBar(
                  SnackBar(
                    content: Text(
                      '删除模型失败',
                      style: TextStyle(
                        fontSize: FontSizeUtils.getBodySize(ref),
                      ),
                    ),
                  ),
                );
              }
            }
          },
          style: TextButton.styleFrom(
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(8),
            ),
            foregroundColor: Colors.red,
          ),
          child: Text(
            '删除',
            style: TextStyle(
              fontSize: FontSizeUtils.getBodySize(ref),
            ),
          ),
        ),
      ],
    ),
  );
} 