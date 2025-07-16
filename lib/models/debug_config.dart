import 'package:json_annotation/json_annotation.dart';

part 'debug_config.g.dart';

@JsonSerializable()
class DebugConfig {
  final String key;
  final String value;

  DebugConfig({
    required this.key,
    required this.value,
  });

  /// 表创建sql
  static String createTableSql() {
    return '''
      CREATE TABLE debug_configs (
        key TEXT PRIMARY KEY,
        value TEXT
      )
    ''';
  }

  factory DebugConfig.fromJson(Map<String, dynamic> json) => _$DebugConfigFromJson(json);
  Map<String, dynamic> toJson() => _$DebugConfigToJson(this);

  // 转换为数据库 Map
  Map<String, dynamic> toMap() {
    return {
      'key': key,
      'value': value,
    };
  }

  // 从数据库 Map 构造对象
  factory DebugConfig.fromMap(Map<String, dynamic> map) {
    return DebugConfig(
      key: map['key'],
      value: map['value'],
    );
  }
}