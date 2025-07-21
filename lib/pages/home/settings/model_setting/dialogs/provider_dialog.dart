import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/rpc/service.pb.dart';
import 'package:lemon_tea/storage/llm_storage.dart';
import 'package:lemon_tea/utils/cli/client/client.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:lemon_tea/models/llm_provider.dart';
import '../model_settings.dart';

void showProviderDialog(BuildContext context, WidgetRef ref, {LlmProvider? provider}) {
  final bool isEditMode = provider != null;
  
  final TextEditingController nameController = TextEditingController(text: provider?.name ?? '');
  final TextEditingController baseUrlController = TextEditingController(text: provider?.baseUrl ?? '');
  final TextEditingController apiKeyController = TextEditingController(text: provider?.apiKey ?? '');
  final TextEditingController aliasController = TextEditingController(text: provider?.alias ?? '');
  final TextEditingController descriptionController = TextEditingController(text: provider?.description ?? '');
  bool isEnabled = provider?.enable ?? true;
  bool isVerifying = false;
  bool verificationSuccess = provider?.checked ?? false;
  bool verificationFailed = false;
  String verificationMessage = provider?.checked ?? false ? '已验证' : '';

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
                  margin: const EdgeInsets.only(bottom: 24),
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
                      const SizedBox(height: 8),
                      Container(
                        padding: const EdgeInsets.all(16),
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
                            const SizedBox(height: 16),
                            TextField(
                              controller: aliasController,
                              decoration: InputDecoration(
                                labelText: '别名',
                                hintText: '可选，用于显示的友好名称',
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
                          ],
                        ),
                      ),
                    ],
                  ),
                ),
                
                // 分组：连接信息
                Container(
                  margin: const EdgeInsets.only(bottom: 24),
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
                      const SizedBox(height: 8),
                      Container(
                        padding: const EdgeInsets.all(16),
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
                            const SizedBox(height: 16),
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
                  margin: const EdgeInsets.only(bottom: 16),
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
                      const SizedBox(height: 8),
                      Container(
                        padding: const EdgeInsets.all(16),
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
                              maxLines: 3,
                            ),
                            const SizedBox(height: 16),
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
                if (verificationSuccess || verificationFailed)
                  Container(
                    padding: const EdgeInsets.all(12),
                    margin: const EdgeInsets.only(top: 8),
                    decoration: BoxDecoration(
                      color: verificationSuccess 
                          ? Theme.of(context).colorScheme.primaryContainer.withOpacity(0.5)
                          : Theme.of(context).colorScheme.errorContainer.withOpacity(0.5),
                      borderRadius: BorderRadius.circular(8),
                      border: Border.all(
                        color: verificationSuccess 
                            ? Theme.of(context).colorScheme.primary.withOpacity(0.5)
                            : Theme.of(context).colorScheme.error.withOpacity(0.5),
                        width: 1,
                      ),
                    ),
                    child: Row(
                      children: [
                        Icon(
                          verificationSuccess ? Icons.check_circle : Icons.error,
                          color: verificationSuccess 
                              ? Theme.of(context).colorScheme.primary
                              : Theme.of(context).colorScheme.error,
                        ),
                        const SizedBox(width: 8),
                        Expanded(
                          child: Text(
                            verificationMessage,
                            style: TextStyle(
                              fontSize: FontSizeUtils.getSmallSize(ref),
                              color: verificationSuccess 
                                  ? Theme.of(context).colorScheme.primary
                                  : Theme.of(context).colorScheme.error,
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
              isVerifying ? '验证中...' : '验证连接',
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
          FilledButton.icon(
            onPressed: () async {
              // 验证输入
              final name = nameController.text.trim();
              final baseUrl = baseUrlController.text.trim();
              final apiKey = apiKeyController.text.trim();
              final alias = aliasController.text.trim();
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
                  alias: alias.isEmpty ? null : alias,
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
                  alias: alias.isEmpty ? null : alias,
                  description: description.isEmpty ? null : description,
                  enable: isEnabled,
                  checked: verificationSuccess, // 根据验证结果设置checked状态
                );

                // 添加供应商到数据库
                success = await LlmStorage.addProvider(newProvider);
              }

              Navigator.of(context).pop();

              if (success) {
                // 刷新供应商列表
                ref.refresh(providersProvider);

                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text(
                      isEditMode ? '供应商已更新' : '供应商已添加',
                      style: TextStyle(
                        fontSize: FontSizeUtils.getBodySize(ref),
                      ),
                    ),
                    backgroundColor: Theme.of(context).colorScheme.primaryContainer,
                    duration: const Duration(seconds: 2),
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