import 'package:json_annotation/json_annotation.dart';

part 'llm_provider.g.dart';

@JsonSerializable()
class LlmProvider {

  // 模型供应商id
  final String id;
  // 模型供应商名称
  final String name;
  // 模型供应商接口url
  final String baseUrl;
  // 模型供应商api密钥
  final String? apiKey;
  // 模型供应商别名
  final String? alias;
  // 模型供应商描述
  final String? description;
  // 是否启用
  @JsonKey(defaultValue: true)
  final bool enable;
  // 是否验证
  @JsonKey(defaultValue: false)
  final bool checked;

  LlmProvider({
    required this.id,
    required this.name,
    required this.baseUrl,
    this.apiKey,
    this.alias,
    this.description,
    this.enable = true,
    this.checked = false,
  });


  /// 从 JSON 创建 LlmProvider 实例
  factory LlmProvider.fromJson(Map<String, dynamic> json) => _$LlmProviderFromJson(json);

  /// 转换为 JSON
  Map<String, dynamic> toJson() => _$LlmProviderToJson(this);

  Map<String, dynamic> toMap() {
    return {
      'id': id,
      'name': name,
      'base_url': baseUrl,
      'api_key': apiKey,
      'alias': alias,
      'description': description,
      'enable': enable ? 1 : 0,
      'checked': checked ? 1 : 0,
    };
  }

  factory LlmProvider.fromMap(Map<String, dynamic> map) {
    return LlmProvider(
      id: map['id'],
      name: map['name'],
      baseUrl: map['base_url'],
      apiKey: map['api_key'],
      alias: map['alias'],
      description: map['description'],
      enable: map['enable'] == 1,
      checked: map['checked'] == 1,
    );
  }
}