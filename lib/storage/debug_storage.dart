import 'package:lemon_tea/models/debug_config.dart';
import 'package:lemon_tea/storage/sqlite_util.dart';
import 'package:lemon_tea/utils/debug/debug_key.dart';

class DebugStorage {
  static const String _tableName = 'debug_configs';

  /// 保存调试配置
  /// 
  /// [key] 配置键
  /// [value] 配置值
  static Future<bool> saveConfig(String key, String value) async {
    try {
      final config = DebugConfig(key: key, value: value);
      final result = await SqliteUtil.instance.insert(_tableName, config.toMap());
      return result > 0;
    } catch (e) {
      print('保存调试配置失败: $e');
      return false;
    }
  }

  /// 获取调试配置
  /// 
  /// [key] 配置键
  /// 返回配置值，如果不存在则返回null
  static Future<String?> getConfig(String key) async {
    try {
      final result = await SqliteUtil.instance.query(
        _tableName,
        where: 'key = ?',
        whereArgs: [key],
      );
      
      if (result.isNotEmpty) {
        final config = DebugConfig.fromMap(result.first);
        return config.value;
      }
      return null;
    } catch (e) {
      print('获取调试配置失败: $e');
      return null;
    }
  }

  /// 删除调试配置
  /// 
  /// [key] 配置键
  /// 返回是否删除成功
  static Future<bool> deleteConfig(String key) async {
    try {
      final result = await SqliteUtil.instance.delete(
        _tableName,
        where: 'key = ?',
        whereArgs: [key],
      );
      return result > 0;
    } catch (e) {
      print('删除调试配置失败: $e');
      return false;
    }
  }

  /// 更新调试配置
  /// 
  /// [key] 配置键
  /// [value] 新的配置值
  /// 返回是否更新成功
  static Future<bool> updateConfig(String key, String value) async {
    try {
      final config = DebugConfig(key: key, value: value);
      final result = await SqliteUtil.instance.update(
        _tableName,
        config.toMap(),
        where: 'key = ?',
        whereArgs: [key],
      );
      
      // 如果没有更新任何记录，可能是因为记录不存在，尝试插入
      if (result == 0) {
        return await saveConfig(key, value);
      }
      
      return result > 0;
    } catch (e) {
      print('更新调试配置失败: $e');
      return false;
    }
  }

  /// 是否启用debug
  static Future<bool> isEnableDebug() async {
    final result = await getConfig(DebugKey.enableDebugMode);
    return result == 'true';
  }

}