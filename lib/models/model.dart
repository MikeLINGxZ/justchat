import 'package:json_annotation/json_annotation.dart';

part 'model.g.dart';

@JsonSerializable()
class Model {
  final String id;
  final String object;
  @JsonKey(name: 'owned_by')
  final String ownedBy;

  const Model({
    required this.id, 
    required this.object, 
    required this.ownedBy,
  });

  /// 从 JSON 创建 Model 实例
  factory Model.fromJson(Map<String, dynamic> json) => _$ModelFromJson(json);

  /// 转换为 JSON
  Map<String, dynamic> toJson() => _$ModelToJson(this);

  /// 创建 Model 的副本，可选择性地更新某些字段
  Model copyWith({
    String? id,
    String? object,
    String? ownedBy,
  }) {
    return Model(
      id: id ?? this.id,
      object: object ?? this.object,
      ownedBy: ownedBy ?? this.ownedBy,
    );
  }

  /// 检查模型是否有效（必须有ID）
  bool get isValid => id.isNotEmpty;

  /// 获取模型的显示名称（使用ID作为显示名称）
  String get displayName => id;

  /// 检查是否为特定类型的模型
  bool isType(String type) => object.toLowerCase() == type.toLowerCase();

  /// 检查是否为聊天模型
  bool get isChatModel => isType('chat.completions') || isType('model');

  /// 检查是否为嵌入模型
  bool get isEmbeddingModel => isType('embeddings') || isType('text-embedding');

  /// 检查是否为图像模型
  bool get isImageModel => isType('images') || isType('image');

  @override
  String toString() {
    return 'Model(id: $id, object: $object, ownedBy: $ownedBy)';
  }

  @override
  bool operator ==(Object other) {
    if (identical(this, other)) return true;
    return other is Model &&
        other.id == id &&
        other.object == object &&
        other.ownedBy == ownedBy;
  }

  @override
  int get hashCode {
    return Object.hash(id, object, ownedBy);
  }
}