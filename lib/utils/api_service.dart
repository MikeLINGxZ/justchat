import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:lemon_tea/models/llm_provider.dart';
import 'package:lemon_tea/models/model.dart';
import 'package:lemon_tea/utils/cli/cli_utils/llm_config_utils.dart';

/// API服务类，负责与LLM供应商的API通信
class ApiService {
  /// 获取供应商的模型列表
  static Future<List<Model>> getModels(LlmProvider provider) async {
    if (!provider.hasApiKey) {
      throw Exception('API密钥未配置');
    }

    try {
      final url = Uri.parse('${provider.baseUrl}/models');
      
      final response = await http.get(
        url,
        headers: {
          'Authorization': 'Bearer ${provider.apiKey}',
          'Content-Type': 'application/json',
        },
      ).timeout(const Duration(seconds: 30));

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        final modelsData = data['data'] as List<dynamic>? ?? [];
        
        return modelsData.map((modelData) {
          final modelJson = modelData as Map<String, dynamic>;
          // 处理不同的JSON结构
          return Model(
            id: modelJson['id']?.toString() ?? '',
            object: modelJson['object']?.toString() ?? 'model',
            ownedBy: modelJson['owned_by']?.toString() ?? 
                     modelJson['OwnedBy']?.toString() ?? 
                     provider.name.toLowerCase(),
            enabled: true, // 默认启用
          );
        }).where((model) => model.id.isNotEmpty).toList();
      } else {
        throw Exception('获取模型列表失败: ${response.statusCode} - ${response.body}');
      }
    } catch (e) {
      if (e is Exception) {
        rethrow;
      }
      throw Exception('网络请求失败: $e');
    }
  }

  /// 测试供应商连接
  static Future<bool> testConnection(LlmProvider provider) async {
    if (!provider.hasApiKey) {
      return false;
    }

    try {
      // 尝试获取模型列表来测试连接
      await getModels(provider);
      return true;
    } catch (e) {
      return false;
    }
  }

  /// 测试连接并获取模型列表
  static Future<Map<String, dynamic>> testConnectionAndGetModels(LlmProvider provider) async {
    if (!provider.hasApiKey) {
      throw Exception('API密钥未配置');
    }

    try {
      // 使用llmModels函数获取模型列表
      final models = await llmModels(provider.baseUrl, provider.apiKey);
      
      if (models.isNotEmpty) {
        return {
          'success': true,
          'models': models,
        };
      } else {
        // 如果没有获取到模型，尝试使用HTTP请求获取
        final httpModels = await getModels(provider);
        return {
          'success': true,
          'models': httpModels,
        };
      }
    } catch (e) {
      return {
        'success': false,
        'error': e.toString(),
        'models': <Model>[],
      };
    }
  }
} 