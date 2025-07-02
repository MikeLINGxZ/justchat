import 'package:json_annotation/json_annotation.dart';
import 'package:lemon_tea/models/model.dart';

part 'llm_provider.g.dart';

@JsonSerializable()
class LlmProvider {
  final String name;
  final String baseUrl;
  final String? apiKey;
  final String? alias;
  final String? description;
  final List<Model>? models;

  const LlmProvider({
    required this.name, 
    required this.baseUrl, 
    this.apiKey, 
    this.alias, 
    this.description, 
    this.models,
  });

  /// 从 JSON 创建 LlmProvider 实例
  factory LlmProvider.fromJson(Map<String, dynamic> json) => 
      _$LlmProviderFromJson(json);

  /// 转换为 JSON
  Map<String, dynamic> toJson() => _$LlmProviderToJson(this);

  /// 创建 LlmProvider 的副本，可选择性地更新某些字段
  LlmProvider copyWith({
    String? name,
    String? baseUrl,
    String? apiKey,
    String? alias,
    String? description,
    List<Model>? models,
  }) {
    return LlmProvider(
      name: name ?? this.name,
      baseUrl: baseUrl ?? this.baseUrl,
      apiKey: apiKey ?? this.apiKey,
      alias: alias ?? this.alias,
      description: description ?? this.description,
      models: models ?? this.models,
    );
  }

  /// 检查是否有 API 密钥
  bool get hasApiKey => apiKey != null && apiKey!.isNotEmpty;

  /// 获取显示名称（优先使用别名，否则使用名称）
  String get displayName => alias ?? name;

  /// 检查是否有效（必须有名称和基础URL）
  bool get isValid => name.isNotEmpty && baseUrl.isNotEmpty;

  @override
  String toString() {
    return 'LlmProvider(name: $name, baseUrl: $baseUrl, alias: $alias, description: $description, modelsCount: ${models?.length ?? 0})';
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is LlmProvider &&
        other.name == name &&
        other.baseUrl == baseUrl &&
        other.apiKey == apiKey &&
        other.alias == alias &&
        other.description == description &&
        _listEquals(other.models, models);
  }

  @override
  int get hashCode {
    return Object.hash(
      name,
      baseUrl,
      apiKey,
      alias,
      description,
      Object.hashAll(models ?? []),
    );
  }

  /// 比较两个列表是否相等
  bool _listEquals<T>(List<T>? a, List<T>? b) {
    if (a == null && b == null) return true;
    if (a == null || b == null) return false;
    if (a.length != b.length) return false;
    for (int i = 0; i < a.length; i++) {
      if (a[i] != b[i]) return false;
    }
    return true;
  }
}