import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:lemon_tea/models/model.dart';
import 'package:lemon_tea/utils/cli/client/client.dart';
import 'package:lemon_tea/utils/storage/local_storage.dart';
import 'package:lemon_tea/rpc/common.pb.dart' as $1;
import 'package:lemon_tea/rpc/service.pb.dart';

/// 更新LLM配置到服务器
Future<void> updateLlmConfig() async {
  try {
    // 获取本地存储的LLM提供商信息
    final LocalStorage localStorage = LocalStorage();
    final providerManagerJson = await localStorage.getString('llm_providers');
    
    if (providerManagerJson != null && Client().stub != null) {
      final List<dynamic> jsonList = jsonDecode(providerManagerJson);
      final providers = jsonList
          .map((json) => json as Map<String, dynamic>)
          .toList();
          
      // 转换为gRPC请求对象
      final llmProviders = providers.map((providerJson) {
        // 从JSON创建LlmProvider对象
        final provider = $1.LlmProvider(
          id: providerJson['name'] ?? '',
          name: providerJson['name'] ?? '',
          baseUrl: providerJson['baseUrl'] ?? '',
          apiKey: providerJson['apiKey'] ?? '',
          alias: providerJson['alias'],
          description: providerJson['description'],
        );
        
        // 添加模型信息
        if (providerJson['models'] != null) {
          final modelsList = (providerJson['models'] as List).map((modelJson) {
            return $1.Model(
              id: modelJson['id'] ?? '',
              object: modelJson['object'] ?? '',
              ownedBy: modelJson['ownedBy'] ?? '',
              enabled: modelJson['enabled'] ?? true,
            );
          }).toList();
          provider.models.addAll(modelsList);
        }
        
        return provider;
      }).toList();
      
      // 创建请求并发送
      final request = UpdateLlmConfigRequest(llmProviders: llmProviders);
      await Client().stub!.updateLlmConfig(request);
      debugPrint('LLM配置已更新到服务器');
    }
  } catch (e) {
    debugPrint('更新LLM配置失败: $e');
  }
}

/// 获取LLM模型列表
///
/// [baseUrl] LLM提供商的基础URL
/// [apiKey] 可选的API密钥
/// 返回模型列表
Future<List<Model>> llmModels(String baseUrl, String? apiKey) async {
  try {
    // 检查客户端是否初始化
    if (Client().stub == null) {
      await Client().init();
      if (Client().stub == null) {
        debugPrint('无法初始化gRPC客户端');
        return [];
      }
    }
    
    // 创建请求对象
    final request = ModelsRequest(
      baseUrl: baseUrl,
      apiKey: apiKey ?? '',
    );
    
    // 调用Models方法获取模型列表
    final response = await Client().stub!.models(request);
    
    // 将gRPC模型对象转换为应用程序模型对象
    return response.models.map((pbModel) => Model(
      id: pbModel.id,
      object: pbModel.object,
      ownedBy: pbModel.ownedBy,
      enabled: pbModel.enabled,
    )).toList();
  } catch (e) {
    debugPrint('获取LLM模型列表失败: $e');
    return [];
  }
}