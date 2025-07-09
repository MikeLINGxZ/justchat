import 'dart:convert';
import 'package:shared_preferences/shared_preferences.dart';

/// 单例本地存储类，提供通用的键值对存储功能
class LocalStorage {
  static LocalStorage? _instance;
  static LocalStorage get instance {
    _instance ??= LocalStorage._internal();
    return _instance!;
  }

  LocalStorage._internal();

  late SharedPreferences _prefs;
  bool _initialized = false;

  /// 确保存储已初始化
  Future<void> _ensureInitialized() async {
    if (!_initialized) {
      _prefs = await SharedPreferences.getInstance();
      _initialized = true;
    }
  }

  /// 存储字符串值
  Future<bool> setString(String key, String value) async {
    await _ensureInitialized();
    return await _prefs.setString(key, value);
  }

  /// 获取字符串值
  Future<String?> getString(String key) async {
    await _ensureInitialized();
    return _prefs.getString(key);
  }

  /// 存储布尔值
  Future<bool> setBool(String key, bool value) async {
    await _ensureInitialized();
    return await _prefs.setBool(key, value);
  }

  /// 获取布尔值
  Future<bool?> getBool(String key) async {
    await _ensureInitialized();
    return _prefs.getBool(key);
  }

  /// 存储整数值
  Future<bool> setInt(String key, int value) async {
    await _ensureInitialized();
    return await _prefs.setInt(key, value);
  }

  /// 获取整数值
  Future<int?> getInt(String key) async {
    await _ensureInitialized();
    return _prefs.getInt(key);
  }

  /// 存储双精度浮点数值
  Future<bool> setDouble(String key, double value) async {
    await _ensureInitialized();
    return await _prefs.setDouble(key, value);
  }

  /// 获取双精度浮点数值
  Future<double?> getDouble(String key) async {
    await _ensureInitialized();
    return _prefs.getDouble(key);
  }

  /// 存储字符串列表
  Future<bool> setStringList(String key, List<String> value) async {
    await _ensureInitialized();
    return await _prefs.setStringList(key, value);
  }

  /// 获取字符串列表
  Future<List<String>?> getStringList(String key) async {
    await _ensureInitialized();
    return _prefs.getStringList(key);
  }

  /// 存储JSON对象（自动序列化）
  Future<bool> setJson(String key, Map<String, dynamic> value) async {
    await _ensureInitialized();
    return await _prefs.setString(key, jsonEncode(value));
  }

  /// 获取JSON对象（自动反序列化）
  Future<Map<String, dynamic>?> getJson(String key) async {
    await _ensureInitialized();
    final jsonString = _prefs.getString(key);
    if (jsonString != null) {
      return jsonDecode(jsonString) as Map<String, dynamic>;
    }
    return null;
  }

  /// 存储对象（需要对象支持toJson方法）
  Future<bool> setObject<T>(String key, T object) async {
    await _ensureInitialized();
    if (object is Map<String, dynamic>) {
      return await setJson(key, object);
    } else {
      // 尝试调用toJson方法
      try {
        final json = (object as dynamic).toJson();
        return await setJson(key, json);
      } catch (e) {
        throw Exception('对象必须实现toJson方法或为Map<String, dynamic>类型');
      }
    }
  }

  /// 获取对象（需要提供fromJson构造函数）
  Future<T?> getObject<T>(String key, T Function(Map<String, dynamic>) fromJson) async {
    await _ensureInitialized();
    final json = await getJson(key);
    if (json != null) {
      return fromJson(json);
    }
    return null;
  }

  /// 检查键是否存在
  Future<bool> containsKey(String key) async {
    await _ensureInitialized();
    return _prefs.containsKey(key);
  }

  /// 删除指定键的值
  Future<bool> remove(String key) async {
    await _ensureInitialized();
    return await _prefs.remove(key);
  }

  /// 清空所有数据
  Future<bool> clear() async {
    await _ensureInitialized();
    return await _prefs.clear();
  }

  /// 获取所有键
  Future<Set<String>> getKeys() async {
    await _ensureInitialized();
    return _prefs.getKeys();
  }

  /// 获取存储大小（字节）
  Future<int> getSize() async {
    await _ensureInitialized();
    // 这是一个估算值，实际实现可能需要更复杂的计算
    int size = 0;
    for (String key in _prefs.getKeys()) {
      final value = _prefs.get(key);
      if (value != null) {
        size += key.length + value.toString().length;
      }
    }
    return size;
  }

  /// 批量设置多个键值对
  Future<bool> setMultiple(Map<String, dynamic> values) async {
    await _ensureInitialized();
    bool success = true;
    for (var entry in values.entries) {
      final result = await _setValue(entry.key, entry.value);
      if (!result) {
        success = false;
      }
    }
    return success;
  }

  /// 批量获取多个键的值
  Future<Map<String, dynamic>> getMultiple(List<String> keys) async {
    await _ensureInitialized();
    Map<String, dynamic> result = {};
    for (String key in keys) {
      result[key] = _prefs.get(key);
    }
    return result;
  }

  /// 内部方法：根据值类型设置存储
  Future<bool> _setValue(String key, dynamic value) async {
    if (value is String) {
      return await setString(key, value);
    } else if (value is bool) {
      return await setBool(key, value);
    } else if (value is int) {
      return await setInt(key, value);
    } else if (value is double) {
      return await setDouble(key, value);
    } else if (value is List<String>) {
      return await setStringList(key, value);
    } else if (value is Map<String, dynamic>) {
      return await setJson(key, value);
    } else {
      throw ArgumentError('不支持的数据类型: ${value.runtimeType}');
    }
  }
}