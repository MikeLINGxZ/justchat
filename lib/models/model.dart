import 'package:json_annotation/json_annotation.dart';

part 'model.g.dart';

@JsonSerializable()
class Model {
  // 供应商id
  final String llm_provider_id;

  // 模型id
  final String id;

  // 模型类型（默认值 "model"）
  @JsonKey(defaultValue: "model")
  final String object;

  // 模型提供者
  @JsonKey(name: 'owned_by')
  final String ownedBy;

  // 是否启用（默认值 true）
  @JsonKey(defaultValue: true)
  final bool enabled;

  Model({
    required this.llm_provider_id,
    required this.id,
    this.object = "model",
    required this.ownedBy,
    this.enabled = true,
  });

  // 生成 fromJson 方法
  factory Model.fromJson(Map<String, dynamic> json) => _$ModelFromJson(json);

  // 生成 toJson 方法
  Map<String, dynamic> toJson() => _$ModelToJson(this);

  // 转换为数据库 Map
  Map<String, dynamic> toMap() {
    return {
      'llm_provider_id': llm_provider_id,
      'id': id,
      'object': object,
      'owned_by': ownedBy,
      'enabled': enabled ? 1 : 0, // SQLite 中布尔用 0/1 表示
    };
  }

  // 从数据库 Map 构造对象
  factory Model.fromMap(Map<String, dynamic> map) {
    return Model(
      llm_provider_id: map['llm_provider_id'],
      id: map['id'],
      object: map['object'] ?? "model",
      ownedBy: map['owned_by'],
      enabled: map['enabled'] == 1,
    );
  }
}