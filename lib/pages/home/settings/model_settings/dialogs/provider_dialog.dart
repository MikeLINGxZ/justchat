import 'package:flutter/material.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/rpc/service.pb.dart';
import 'package:lemon_tea/storage/llm_storage.dart';
import 'package:lemon_tea/utils/cli/client/client.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:lemon_tea/models/llm_provider.dart';
import 'package:lemon_tea/models/model.dart';
import '../model_settings.dart';
import 'models_dialog.dart';

void showProviderDialog(BuildContext context, WidgetRef ref, {LlmProvider? provider}) {
  final bool isEditMode = provider != null;

  final TextEditingController nameController = TextEditingController(text: provider?.name ?? '');
  final TextEditingController baseUrlController = TextEditingController(text: provider?.baseUrl ?? '');
  final TextEditingController apiKeyController = TextEditingController(text: provider?.apiKey ?? '');
  final TextEditingController descriptionController = TextEditingController(text: provider?.description ?? '');
  bool isEnabled = provider?.enable ?? true;
  bool isVerifying = false;
  bool verificationSuccess = provider?.checked ?? false;
  bool verificationFailed = false;
  String verificationMessage = provider?.checked ?? false ? '已验证' : '待验证';



  showDialog(
    context: context,
    builder: (context) => StatefulBuilder(
      builder: (context, setState) => AlertDialog(
        title: Row(
          children: [
            Icon(
              isEditMode ? Icons.edit : Icons.add_circle,
              color: Theme.of(context).colorScheme.primary,
              size: 28,
            ),
            const SizedBox(width: 12),
            Text(
              isEditMode ? '编辑模型供应商' : '添加模型供应商',
              style: TextStyle(
                fontSize: FontSizeUtils.getSubheadingSize(ref),
                fontWeight: FontWeight.bold,
              ),
            ),
          ],
        ),
        content: SizedBox(
          width: 500,
          child: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // 分组：基本信息
                Container(
                  margin: const EdgeInsets.only(bottom: 16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        '基本信息',
                        style: TextStyle(
                          fontSize: FontSizeUtils.getBodySize(ref),
                          fontWeight: FontWeight.bold,
                          color: Theme.of(context).colorScheme.primary,
                        ),
                      ),
                      const SizedBox(height: 6),
                      Container(
                        padding: const EdgeInsets.all(12),
                        decoration: BoxDecoration(
                          color: Theme.of(context).colorScheme.surfaceContainerLow,
                          borderRadius: BorderRadius.circular(8),
                          border: Border.all(
                            color: Theme.of(context).colorScheme.outlineVariant.withOpacity(0.5),
                            width: 1,
                          ),
                        ),
                        child: Column(
                          children: [
                            TextField(
                              controller: nameController,
                              decoration: InputDecoration(
                                labelText: '供应商名称 *',
                                hintText: '例如: OpenAI',
                                prefixIcon: const Icon(Icons.business),
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
                          ],
                        ),
                      ),
                    ],
                  ),
                ),

                // 分组：连接信息
                Container(
                  margin: const EdgeInsets.only(bottom: 16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        '连接信息',
                        style: TextStyle(
                          fontSize: FontSizeUtils.getBodySize(ref),
                          fontWeight: FontWeight.bold,
                          color: Theme.of(context).colorScheme.primary,
                        ),
                      ),
                      const SizedBox(height: 6),
                      Container(
                        padding: const EdgeInsets.all(12),
                        decoration: BoxDecoration(
                          color: Theme.of(context).colorScheme.surfaceContainerLow,
                          borderRadius: BorderRadius.circular(8),
                          border: Border.all(
                            color: Theme.of(context).colorScheme.outlineVariant.withOpacity(0.5),
                            width: 1,
                          ),
                        ),
                        child: Column(
                          children: [
                            TextField(
                              controller: baseUrlController,
                              decoration: InputDecoration(
                                labelText: '基础URL *',
                                hintText: '例如: https://api.openai.com/v1',
                                prefixIcon: const Icon(Icons.link),
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
                            const SizedBox(height: 12),
                            TextField(
                              controller: apiKeyController,
                              decoration: InputDecoration(
                                labelText: 'API密钥 *',
                                hintText: '您的API密钥',
                                prefixIcon: const Icon(Icons.key),
                                suffixIcon: IconButton(
                                  icon: const Icon(Icons.visibility_off),
                                  onPressed: () {
                                    // 切换密码可见性功能可以在这里实现
                                  },
                                ),
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
                              obscureText: true,
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),

                // 分组：附加信息
                Container(
                  margin: const EdgeInsets.only(bottom: 12),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        '附加信息',
                        style: TextStyle(
                          fontSize: FontSizeUtils.getBodySize(ref),
                          fontWeight: FontWeight.bold,
                          color: Theme.of(context).colorScheme.primary,
                        ),
                      ),
                      const SizedBox(height: 6),
                      Container(
                        padding: const EdgeInsets.all(12),
                        decoration: BoxDecoration(
                          color: Theme.of(context).colorScheme.surfaceContainerLow,
                          borderRadius: BorderRadius.circular(8),
                          border: Border.all(
                            color: Theme.of(context).colorScheme.outlineVariant.withOpacity(0.5),
                            width: 1,
                          ),
                        ),
                        child: Column(
                          children: [
                            TextField(
                              controller: descriptionController,
                              decoration: InputDecoration(
                                labelText: '描述',
                                hintText: '可选，添加关于此供应商的描述',
                                prefixIcon: const Icon(Icons.description),
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
                              maxLines: 1,
                            ),
                            const SizedBox(height: 12),
                            SwitchListTile(
                              title: Text(
                                '启用此供应商',
                                style: TextStyle(
                                  fontSize: FontSizeUtils.getBodySize(ref),
                                ),
                              ),
                              subtitle: Text(
                                '关闭后将不会在模型列表中显示',
                                style: TextStyle(
                                  fontSize: FontSizeUtils.getSmallSize(ref),
                                  color: Theme.of(context).colorScheme.onSurfaceVariant,
                                ),
                              ),
                              value: isEnabled,
                              onChanged: (value) {
                                setState(() {
                                  isEnabled = value;
                                });
                              },
                              secondary: Icon(
                                isEnabled ? Icons.toggle_on : Icons.toggle_off,
                                color: isEnabled ? Theme.of(context).colorScheme.primary : null,
                              ),
                              shape: RoundedRectangleBorder(
                                borderRadius: BorderRadius.circular(8),
                              ),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                ),

                // 验证结果显示区域
                if (verificationSuccess || verificationFailed || (!verificationSuccess && !verificationFailed))
                  Container(
                    padding: const EdgeInsets.all(10),
                    margin: const EdgeInsets.only(top: 6),
                    decoration: BoxDecoration(
                      color: verificationSuccess
                          ? Theme.of(context).colorScheme.primaryContainer.withOpacity(0.5)
                          : verificationFailed
                          ? Theme.of(context).colorScheme.errorContainer.withOpacity(0.5)
                          : Theme.of(context).colorScheme.surfaceContainerLow,
                      borderRadius: BorderRadius.circular(8),
                      border: Border.all(
                        color: verificationSuccess
                            ? Theme.of(context).colorScheme.primary.withOpacity(0.5)
                            : verificationFailed
                            ? Theme.of(context).colorScheme.error.withOpacity(0.5)
                            : Theme.of(context).colorScheme.outlineVariant.withOpacity(0.5),
                        width: 1,
                      ),
                    ),
                    child: Column(
                      children: [
                        InkWell(
                          onTap: null,
                          borderRadius: BorderRadius.circular(8),
                          child: Padding(
                            padding: const EdgeInsets.all(4),
                            child: Row(
                              children: [
                                Icon(
                                  verificationSuccess
                                      ? Icons.check_circle
                                      : verificationFailed
                                      ? Icons.error
                                      : Icons.help_outline,
                                  color: verificationSuccess
                                      ? Theme.of(context).colorScheme.primary
                                      : verificationFailed
                                      ? Theme.of(context).colorScheme.error
                                      : Theme.of(context).colorScheme.onSurfaceVariant,
                                ),
                                const SizedBox(width: 8),
                                Expanded(
                                  child: Text(
                                    verificationMessage,
                                    style: TextStyle(
                                      fontSize: FontSizeUtils.getSmallSize(ref),
                                      color: verificationSuccess
                                          ? Theme.of(context).colorScheme.primary
                                          : verificationFailed
                                          ? Theme.of(context).colorScheme.error
                                          : Theme.of(context).colorScheme.onSurfaceVariant,
                                    ),
                                  ),
                                ),

                              ],
                            ),
                          ),
                        ),
                        if (verificationSuccess)
                          Container(
                            margin: const EdgeInsets.only(top: 8),
                            decoration: BoxDecoration(
                              color: Theme.of(context).colorScheme.primaryContainer.withOpacity(0.3),
                              borderRadius: BorderRadius.circular(8),
                              border: Border.all(
                                color: Theme.of(context).colorScheme.primary.withOpacity(0.3),
                                width: 1,
                              ),
                            ),
                            child: Material(
                              color: Colors.transparent,
                              child: InkWell(
                                onTap: () {
                                  if (isEditMode && provider != null) {
                                    // 导入新的模型对话框
                                    showModelsDialog(context, ref, provider!);
                                  }
                                },
                                borderRadius: BorderRadius.circular(8),
                                child: Padding(
                                  padding: const EdgeInsets.all(12),
                                  child: Row(
                                    children: [
                                      Icon(
                                        Icons.psychology,
                                        size: 20,
                                        color: Theme.of(context).colorScheme.primary,
                                      ),
                                      const SizedBox(width: 12),
                                      Expanded(
                                        child: Column(
                                          crossAxisAlignment: CrossAxisAlignment.start,
                                          children: [
                                            Text(
                                              '查看模型列表',
                                              style: TextStyle(
                                                fontSize: FontSizeUtils.getBodySize(ref),
                                                fontWeight: FontWeight.w600,
                                                color: Theme.of(context).colorScheme.primary,
                                              ),
                                            ),
                                            const SizedBox(height: 2),
                                            Text(
                                              '点击查看和管理可用模型',
                                              style: TextStyle(
                                                fontSize: FontSizeUtils.getSmallSize(ref),
                                                color: Theme.of(context).colorScheme.onSurfaceVariant,
                                              ),
                                            ),
                                          ],
                                        ),
                                      ),
                                      Icon(
                                        Icons.arrow_forward_ios,
                                        size: 16,
                                        color: Theme.of(context).colorScheme.primary,
                                      ),
                                    ],
                                  ),
                                ),
                              ),
                            ),
                          ),
                      ],
                    ),
                  ),
              ],
            ),
          ),
        ),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
        ),
        actions: [
          TextButton.icon(
            onPressed: isVerifying ? null : () async {
              // 验证输入
              final name = nameController.text.trim();
              final baseUrl = baseUrlController.text.trim();
              final apiKey = apiKeyController.text.trim();

              if (name.isEmpty || baseUrl.isEmpty || apiKey.isEmpty) {
                setState(() {
                  verificationFailed = true;
                  verificationSuccess = false;
                  verificationMessage = '请填写所有必填字段后再验证';
                });
                return;
              }

              setState(() {
                isVerifying = true;
                verificationSuccess = false;
                verificationFailed = false;
                verificationMessage = '';
              });

              try {
                final request = ModelsRequest(
                    name: name,
                    apiKey: apiKey,
                    baseUrl: baseUrl
                );
                ModelsResponse response = await Client().stub!.models(request);

                // 如果是编辑模式且供应商存在，更新数据库中的模型列表
                if (isEditMode && provider != null) {
                  await _updateProviderModels(provider.id, response.models);
                }

                // 如果是编辑模式，从数据库重新查询模型列表
                List<Model> updatedModels = [];
                if (isEditMode && provider != null) {
                  try {
                    final dbModels = await LlmStorage.getModelsByProviderId(provider.id);
                    dbModels.sort((a, b) {
                      if (a.isCustom == b.isCustom) {
                        return a.id.compareTo(b.id);
                      }
                      return a.isCustom ? 1 : -1;
                    });
                    updatedModels = dbModels;
                  } catch (e) {
                    debugPrint('查询数据库模型失败: $e');
                  }
                } else {
                  // 新建模式，创建临时Model对象用于显示
                  updatedModels = response.models.map((apiModel) => Model(
                    llmProviderId: 'temp',
                    id: apiModel.id,
                    object: apiModel.object ?? 'model',
                    ownedBy: apiModel.ownedBy ?? 'unknown',
                    enabled: apiModel.enabled ?? true,
                    isCustom: false,
                    seqId: 0,
                  )).toList();
                }

                setState(() {
                  isVerifying = false;
                  verificationSuccess = true;
                  verificationFailed = false;
                  verificationMessage = '验证成功! 发现 ${response.models.length} 个模型';
                });
              } catch (e) {
                setState(() {
                  isVerifying = false;
                  verificationSuccess = false;
                  verificationFailed = true;
                  verificationMessage = '验证失败: ${e.toString()}';
                });
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
                : const Icon(Icons.verified),
            label: Text(
              isVerifying ? '验证中...' : (verificationSuccess ? '重新验证' : '验证连接'),
              style: TextStyle(
                fontSize: FontSizeUtils.getBodySize(ref),
              ),
            ),
            style: TextButton.styleFrom(
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(8),
              ),
              padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            ),
          ),
          Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              FilledButton.icon(
                onPressed: () async {
                  // 验证输入
                  final name = nameController.text.trim();
                  final baseUrl = baseUrlController.text.trim();
                  final apiKey = apiKeyController.text.trim();
                  final description = descriptionController.text.trim();

                  if (name.isEmpty || baseUrl.isEmpty || apiKey.isEmpty) {
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

                  bool success;
                  if (isEditMode) {
                    // 更新供应商对象
                    final updatedProvider = LlmProvider(
                      id: provider!.id, // 保持原ID不变
                      name: name,
                      baseUrl: baseUrl,
                      apiKey: apiKey,
                      description: description.isEmpty ? null : description,
                      enable: isEnabled,
                      checked: verificationSuccess, // 根据验证结果设置checked状态
                    );

                    // 更新供应商到数据库
                    success = await LlmStorage.updateProvider(updatedProvider);
                  } else {
                    // 创建供应商对象
                    final newProvider = LlmProvider(
                      id: '${name.toLowerCase().replaceAll(' ', '_')}_${DateTime.now().millisecondsSinceEpoch}',
                      name: name,
                      baseUrl: baseUrl,
                      apiKey: apiKey,
                      description: description.isEmpty ? null : description,
                      enable: isEnabled,
                      checked: verificationSuccess, // 根据验证结果设置checked状态
                    );

                    // 添加供应商到数据库
                    success = await LlmStorage.addProvider(newProvider);

                    // 如果验证成功，尝试添加模型到数据库
                    if (success && verificationSuccess) {
                      try {
                        // 获取API模型数据
                        final modelsRequest = ModelsRequest(
                            name: name,
                            apiKey: apiKey,
                            baseUrl: baseUrl
                        );
                        final modelsResponse = await Client().stub!.models(modelsRequest);
                        await _updateProviderModels(newProvider.id, modelsResponse.models);
                      } catch (e) {
                        debugPrint('添加模型失败，但供应商已创建: $e');
                      }
                    }
                  }

                  Navigator.of(context).pop();

                  if (success) {
                    // 刷新供应商列表
                    ref.refresh(providersProvider);

                    ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(
                        content: Text(
                          isEditMode
                              ? '供应商已更新${verificationSuccess ? '，模型列表已同步' : ''}'
                              : '供应商已添加${verificationSuccess ? '，模型列表已同步' : ''}',
                          style: TextStyle(
                            fontSize: FontSizeUtils.getBodySize(ref),
                          ),
                        ),
                        backgroundColor: Theme.of(context).colorScheme.primaryContainer,
                        duration: const Duration(seconds: 3),
                      ),
                    );
                  } else {
                    ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(
                        content: Text(
                          isEditMode ? '更新供应商失败' : '添加供应商失败',
                          style: TextStyle(
                            fontSize: FontSizeUtils.getBodySize(ref),
                          ),
                        ),
                        backgroundColor: Theme.of(context).colorScheme.errorContainer,
                      ),
                    );
                  }
                },
                icon: Icon(isEditMode ? Icons.save : Icons.add),
                label: Text(
                  isEditMode ? '保存修改' : '添加供应商',
                  style: TextStyle(
                    fontSize: FontSizeUtils.getBodySize(ref),
                  ),
                ),
                style: FilledButton.styleFrom(
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(8),
                  ),
                  padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                ),
              ),
              const SizedBox(width: 8),
              TextButton.icon(
                onPressed: () => Navigator.of(context).pop(),
                icon: const Icon(Icons.cancel),
                label: Text(
                  '取消',
                  style: TextStyle(
                    fontSize: FontSizeUtils.getBodySize(ref),
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
          ),
        ],
        actionsPadding: const EdgeInsets.fromLTRB(16, 0, 16, 16),
        actionsAlignment: MainAxisAlignment.spaceBetween,
      ),
    ),
  );
}

// 便捷函数：显示添加供应商对话框
void showAddProviderDialog(BuildContext context, WidgetRef ref) {
  showProviderDialog(context, ref);
}

// 便捷函数：显示编辑供应商对话框
void showEditProviderDialog(BuildContext context, WidgetRef ref, LlmProvider provider) {
  showProviderDialog(context, ref, provider: provider);
}

// 显示添加模型对话框的私有函数
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