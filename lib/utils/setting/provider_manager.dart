import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:lemon_tea/models/llm_provider.dart';
import 'package:lemon_tea/models/model.dart';

/// 模型供应商管理器Provider
final providerManagerProvider = StateNotifierProvider<ProviderManager, List<LlmProvider>>((ref) {
  return ProviderManager();
});

/// 当前选中的模型供应商Provider
final selectedProviderProvider = StateProvider<LlmProvider?>((ref) => null);

/// 当前选中的模型Provider
final selectedModelProvider = StateProvider<Model?>((ref) => null);

/// 模型供应商管理器
class ProviderManager extends StateNotifier<List<LlmProvider>> {
  static const String _providersKey = 'llm_providers';
  
  ProviderManager() : super([]) {
    loadProviders();
  }

  /// 加载保存的模型供应商
  Future<void> loadProviders() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final providersJson = prefs.getString(_providersKey);
      
      if (providersJson != null) {
        final List<dynamic> jsonList = jsonDecode(providersJson);
        final providers = jsonList
            .map((json) => LlmProvider.fromJson(json as Map<String, dynamic>))
            .toList();
        state = providers;
      } else {
        // 如果没有保存的数据，添加一些默认的供应商
        state = _getDefaultProviders();
      }
    } catch (e) {
      debugPrint('Failed to load providers: $e');
      // 使用默认供应商
      state = _getDefaultProviders();
    }
  }

  /// 获取默认的模型供应商
  List<LlmProvider> _getDefaultProviders() {
    return [
      LlmProvider(
        name: 'OpenAI',
        baseUrl: 'https://api.openai.com/v1',
        alias: 'OpenAI',
        description: 'OpenAI官方API',
        models: [
          Model(id: 'gpt-4', object: 'model', ownedBy: 'openai'),
          Model(id: 'gpt-4-turbo', object: 'model', ownedBy: 'openai'),
          Model(id: 'gpt-3.5-turbo', object: 'model', ownedBy: 'openai'),
        ],
      ),
      LlmProvider(
        name: 'Anthropic',
        baseUrl: 'https://api.anthropic.com',
        alias: 'Claude',
        description: 'Anthropic Claude API',
        models: [
          Model(id: 'claude-3-opus-20240229', object: 'model', ownedBy: 'anthropic'),
          Model(id: 'claude-3-sonnet-20240229', object: 'model', ownedBy: 'anthropic'),
          Model(id: 'claude-3-haiku-20240307', object: 'model', ownedBy: 'anthropic'),
        ],
      ),
      LlmProvider(
        name: 'Google',
        baseUrl: 'https://generativelanguage.googleapis.com',
        alias: 'Gemini',
        description: 'Google Gemini API',
        models: [
          Model(id: 'gemini-pro', object: 'model', ownedBy: 'google'),
          Model(id: 'gemini-pro-vision', object: 'model', ownedBy: 'google'),
        ],
      ),
    ];
  }

  /// 保存模型供应商到本地存储
  Future<void> _saveProviders() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final providersJson = jsonEncode(state.map((p) => p.toJson()).toList());
      await prefs.setString(_providersKey, providersJson);
    } catch (e) {
      debugPrint('Failed to save providers: $e');
    }
  }

  /// 添加新的模型供应商
  Future<void> addProvider(LlmProvider provider) async {
    // 检查是否已存在相同名称的供应商
    if (state.any((p) => p.name == provider.name)) {
      throw Exception('已存在相同名称的模型供应商');
    }
    
    state = [...state, provider];
    await _saveProviders();
  }

  /// 更新模型供应商
  Future<void> updateProvider(String originalName, LlmProvider updatedProvider) async {
    final index = state.indexWhere((p) => p.name == originalName);
    if (index == -1) {
      throw Exception('未找到要更新的模型供应商');
    }
    
    // 检查新名称是否与其他供应商冲突
    if (originalName != updatedProvider.name && 
        state.any((p) => p.name == updatedProvider.name)) {
      throw Exception('已存在相同名称的模型供应商');
    }
    
    final newProviders = List<LlmProvider>.from(state);
    newProviders[index] = updatedProvider;
    state = newProviders;
    await _saveProviders();
  }

  /// 删除模型供应商
  Future<void> deleteProvider(String name) async {
    state = state.where((p) => p.name != name).toList();
    await _saveProviders();
  }

  /// 根据名称获取模型供应商
  LlmProvider? getProviderByName(String name) {
    try {
      return state.firstWhere((p) => p.name == name);
    } catch (e) {
      return null;
    }
  }

  /// 获取所有聊天模型
  List<Model> getAllChatModels() {
    final models = <Model>[];
    for (final provider in state) {
      if (provider.models != null) {
        for (final model in provider.models!) {
          if (model.isChatModel) {
            models.add(model);
          }
        }
      }
    }
    return models;
  }

  /// 获取指定供应商的聊天模型
  List<Model> getChatModelsByProvider(String providerName) {
    final provider = getProviderByName(providerName);
    if (provider?.models == null) return [];
    
    return provider!.models!.where((model) => model.isChatModel).toList();
  }

  /// 更新供应商的API密钥
  Future<void> updateProviderApiKey(String providerName, String apiKey) async {
    final provider = getProviderByName(providerName);
    if (provider == null) {
      throw Exception('未找到指定的模型供应商');
    }
    
    final updatedProvider = provider.copyWith(apiKey: apiKey);
    await updateProvider(providerName, updatedProvider);
  }

  /// 测试供应商连接
  Future<bool> testProviderConnection(LlmProvider provider) async {
    // TODO: 实现实际的连接测试
    // 这里可以发送一个简单的API请求来验证连接
    await Future.delayed(const Duration(seconds: 1)); // 模拟网络请求
    return provider.hasApiKey; // 简单检查是否有API密钥
  }
} 