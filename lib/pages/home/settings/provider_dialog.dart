import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/models/llm_provider_v0.dart';
import 'package:lemon_tea/models/model_v0.dart';
import 'package:lemon_tea/utils/setting/provider_manager.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:lemon_tea/generated/l10n.dart';

class ProviderDialog extends ConsumerStatefulWidget {
  final LlmProvider_v0? provider; // 如果为null，则为添加模式；否则为编辑模式

  const ProviderDialog({super.key, this.provider});

  @override
  ConsumerState<ProviderDialog> createState() => _ProviderDialogState();
}

class _ProviderDialogState extends ConsumerState<ProviderDialog> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _baseUrlController = TextEditingController();
  final _apiKeyController = TextEditingController();
  final _aliasController = TextEditingController();
  final _descriptionController = TextEditingController();
  
  bool _isLoading = false;
  bool _showApiKey = false;

  @override
  void initState() {
    super.initState();
    if (widget.provider != null) {
      // 编辑模式，填充现有数据
      _nameController.text = widget.provider!.name;
      _baseUrlController.text = widget.provider!.baseUrl;
      _apiKeyController.text = widget.provider!.apiKey ?? '';
      _aliasController.text = widget.provider!.alias ?? '';
      _descriptionController.text = widget.provider!.description ?? '';
    }
  }

  @override
  void dispose() {
    _nameController.dispose();
    _baseUrlController.dispose();
    _apiKeyController.dispose();
    _aliasController.dispose();
    _descriptionController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final isEditMode = widget.provider != null;
    
    return AlertDialog(
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(4),
      ),
      title: Text(
        isEditMode ? '编辑模型供应商' : '添加模型供应商',
        style: TextStyle(fontSize: FontSizeUtils.getHeadingSize(ref)),
      ),
      content: SizedBox(
        width: 500,
        child: Form(
          key: _formKey,
          child: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _buildTextField(
                  controller: _nameController,
                  label: '供应商名称',
                  hint: '例如：OpenAI',
                  validator: (value) {
                    if (value == null || value.trim().isEmpty) {
                      return '请输入供应商名称';
                    }
                    return null;
                  },
                ),
                const SizedBox(height: 16),
                
                _buildTextField(
                  controller: _baseUrlController,
                  label: 'API基础URL',
                  hint: '例如：https://api.openai.com/v1',
                  validator: (value) {
                    if (value == null || value.trim().isEmpty) {
                      return '请输入API基础URL';
                    }
                    if (!Uri.tryParse(value)!.hasScheme == true) {
                      return '请输入有效的URL';
                    }
                    return null;
                  },
                ),
                const SizedBox(height: 16),
                
                _buildTextField(
                  controller: _apiKeyController,
                  label: 'API密钥',
                  hint: '输入您的API密钥',
                  obscureText: !_showApiKey,
                  suffixIcon: IconButton(
                    icon: Icon(_showApiKey ? Icons.visibility_off : Icons.visibility),
                    onPressed: () {
                      setState(() {
                        _showApiKey = !_showApiKey;
                      });
                    },
                  ),
                ),
                const SizedBox(height: 16),
                
                _buildTextField(
                  controller: _aliasController,
                  label: '显示名称（可选）',
                  hint: '例如：OpenAI',
                ),
                const SizedBox(height: 16),
                
                _buildTextField(
                  controller: _descriptionController,
                  label: '描述（可选）',
                  hint: '供应商的简要描述',
                  maxLines: 3,
                ),
              ],
            ),
          ),
        ),
      ),
      actions: [
        TextButton(
          onPressed: _isLoading ? null : () => Navigator.of(context).pop(),
          child: Text(S.of(context).cancel),
        ),
        if (isEditMode)
          TextButton(
            onPressed: _isLoading ? null : _testConnection,
            child: const Text('测试连接'),
          ),
        ElevatedButton(
          onPressed: _isLoading ? null : _saveProvider,
          child: _isLoading
              ? const SizedBox(
                  width: 16,
                  height: 16,
                  child: CircularProgressIndicator(strokeWidth: 2),
                )
              : Text(isEditMode ? '保存' : '添加'),
        ),
      ],
    );
  }

  Widget _buildTextField({
    required TextEditingController controller,
    required String label,
    required String hint,
    bool obscureText = false,
    Widget? suffixIcon,
    int maxLines = 1,
    String? Function(String?)? validator,
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: TextStyle(
            fontSize: FontSizeUtils.getBodySize(ref),
            fontWeight: FontWeight.w500,
          ),
        ),
        const SizedBox(height: 8),
        TextFormField(
          controller: controller,
          obscureText: obscureText,
          maxLines: maxLines,
          validator: validator,
          decoration: InputDecoration(
            hintText: hint,
            suffixIcon: suffixIcon,
            border: const OutlineInputBorder(),
            contentPadding: const EdgeInsets.symmetric(
              horizontal: 12,
              vertical: 12,
            ),
          ),
        ),
      ],
    );
  }

  Future<void> _saveProvider() async {
    if (!_formKey.currentState!.validate()) {
      return;
    }

    setState(() {
      _isLoading = true;
    });

    try {
      final provider = LlmProvider_v0(
        name: _nameController.text.trim(),
        baseUrl: _baseUrlController.text.trim(),
        apiKey: _apiKeyController.text.trim().isEmpty ? null : _apiKeyController.text.trim(),
        alias: _aliasController.text.trim().isEmpty ? null : _aliasController.text.trim(),
        description: _descriptionController.text.trim().isEmpty ? null : _descriptionController.text.trim(),
      );

      final providerManager = ref.read(providerManagerProvider.notifier);
      
      if (widget.provider != null) {
        // 编辑模式
        await providerManager.updateProvider(widget.provider!.name, provider);
      } else {
        // 添加模式
        await providerManager.addProvider(provider);
      }

      if (mounted) {
        Navigator.of(context).pop(true); // 返回true表示操作成功
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(widget.provider != null ? '供应商更新成功' : '供应商添加成功'),
          ),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('操作失败：${e.toString()}'),
            backgroundColor: Colors.red,
          ),
        );
      }
    } finally {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    }
  }

  Future<void> _testConnection() async {
    if (!_formKey.currentState!.validate()) {
      return;
    }

    setState(() {
      _isLoading = true;
    });

    try {
      final provider = LlmProvider_v0(
        name: _nameController.text.trim(),
        baseUrl: _baseUrlController.text.trim(),
        apiKey: _apiKeyController.text.trim().isEmpty ? null : _apiKeyController.text.trim(),
        alias: _aliasController.text.trim().isEmpty ? null : _aliasController.text.trim(),
        description: _descriptionController.text.trim().isEmpty ? null : _descriptionController.text.trim(),
      );

      final providerManager = ref.read(providerManagerProvider.notifier);
      final result = await providerManager.testProviderConnection(provider);

      if (mounted) {
        if (result['success']) {
          final models = result['models'] as List<Model_v0>;
          
          // 显示成功消息
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text('连接测试成功，获取到 ${models.length} 个模型'),
              backgroundColor: Colors.green,
              duration: const Duration(seconds: 3),
            ),
          );
          
          // 如果获取到了模型，显示对话框询问是否保存
          if (models.isNotEmpty) {
            _showSaveModelsDialog(provider, models);
          }
        } else {
          final error = result['error'] as String;
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text('连接测试失败：$error'),
              backgroundColor: Colors.red,
              duration: const Duration(seconds: 5),
            ),
          );
        }
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('连接测试失败：${e.toString()}'),
            backgroundColor: Colors.red,
          ),
        );
      }
    } finally {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    }
  }
  
  /// 显示保存模型对话框
  void _showSaveModelsDialog(LlmProvider_v0 provider, List<Model_v0> models) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(4),
        ),
        title: const Text('保存模型信息'),
        content: SizedBox(
          width: 400,
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text('是否将获取到的模型信息保存到此供应商？'),
              const SizedBox(height: 16),
              SizedBox(
                height: 200,
                child: ListView.builder(
                  shrinkWrap: true,
                  itemCount: models.length > 5 ? 5 : models.length,
                  itemBuilder: (context, index) {
                    final model = models[index];
                    return ListTile(
                      dense: true,
                      title: Text(model.displayName),
                      subtitle: Text('类型: ${model.object}'),
                    );
                  },
                ),
              ),
              if (models.length > 5)
                Padding(
                  padding: const EdgeInsets.only(top: 8.0),
                  child: Text('... 还有 ${models.length - 5} 个模型未显示'),
                ),
            ],
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('取消'),
          ),
          ElevatedButton(
            onPressed: () async {
              Navigator.of(context).pop();
              
              // 更新供应商的模型列表
              final updatedProvider = provider.copyWith(models: models);
              
              // 保存供应商信息
              try {
                final providerManager = ref.read(providerManagerProvider.notifier);
                
                if (widget.provider != null) {
                  // 编辑模式
                  await providerManager.updateProvider(widget.provider!.name, updatedProvider);
                } else {
                  // 添加模式
                  await providerManager.addProvider(updatedProvider);
                }
                
                ScaffoldMessenger.of(context).showSnackBar(
                  const SnackBar(
                    content: Text('模型信息已保存'),
                    backgroundColor: Colors.green,
                  ),
                );
              } catch (e) {
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(
                    content: Text('保存模型信息失败：${e.toString()}'),
                    backgroundColor: Colors.red,
                  ),
                );
              }
            },
            child: const Text('保存'),
          ),
        ],
      ),
    );
  }
} 