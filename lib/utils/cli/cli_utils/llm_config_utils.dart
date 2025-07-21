import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:lemon_tea/models/model_v0.dart';
import 'package:lemon_tea/utils/cli/client/client.dart';
import 'package:lemon_tea/storage/llm_storage.dart';
import 'package:lemon_tea/rpc/common.pb.dart' as $1;
import 'package:lemon_tea/rpc/service.pb.dart';

/// 更新LLM配置到服务器
Future<void> updateLlmConfig() async {
  try {
    // 从SQLite数据库获取LLM提供商信息
    final providers = await LlmStorage.getAllProviders();
    
    if (providers.isNotEmpty && Client().stub != null) {
      // 转换为gRPC请求对象
      final llmProviders = <$1.LlmProvider>[];
      
      for (final provider in providers) {
        // 从数据库创建LlmProvider对象
        final grpcProvider = $1.LlmProvider(
          id: provider.id,
          name: provider.name,
          baseUrl: provider.baseUrl,
          apiKey: provider.apiKey ?? '',
          alias: provider.alias,
          description: provider.description,
        );
        
        // 获取该提供商的所有模型
        final models = await LlmStorage.getModelsByProviderId(provider.id);
        final modelsList = models.map((model) {
          return $1.Model(
            id: model.id,
            object: model.object,
            ownedBy: model.ownedBy,
            enabled: model.enabled,
          );
        }).toList();
        
        grpcProvider.models.addAll(modelsList);
        llmProviders.add(grpcProvider);
      }
      
      // 创建请求并发送
      final request = UpdateLlmConfigRequest(llmProviders: llmProviders);
      await Client().stub!.updateLlmConfig(request);
      debugPrint('LLM配置已更新到服务器，包含${providers.length}个提供商和${llmProviders.fold(0, (sum, p) => sum + p.models.length)}个模型');
    } else {
      debugPrint('未找到LLM提供商配置或gRPC客户端未初始化');
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
Future<List<Model_v0>> llmModels(String baseUrl, String? apiKey) async {
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
    return response.models.map((pbModel) => Model_v0(
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