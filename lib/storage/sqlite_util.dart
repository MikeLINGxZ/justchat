import 'dart:async';
import 'package:sqflite/sqflite.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'sqlite_initializer.dart';

/// SQLite数据库存储提供者
final sqliteStorageProvider = Provider<SqliteStorage>((ref) {
  return SqliteStorage.instance;
});

/// SQLite数据库存储类
/// 提供数据库操作功能
class SqliteStorage {
  static SqliteStorage? _instance;
  final SqliteDatabaseInitializer _initializer = SqliteDatabaseInitializer();
  
  // 单例模式
  static SqliteStorage get instance {
    _instance ??= SqliteStorage._internal();
    return _instance!;
  }
  
  // 私有构造函数
  SqliteStorage._internal();
  
  /// 获取数据库实例
  Future<Database> get database => _initializer.database;
  
  /// 关闭数据库连接
  Future<void> close() => _initializer.close();
  
  /// 删除数据库
  Future<void> deleteDatabase() => _initializer.deleteDatabase();
  
  /// 执行原始SQL查询
  Future<List<Map<String, dynamic>>> rawQuery(String sql, [List<dynamic>? arguments]) async {
    final db = await database;
    return await db.rawQuery(sql, arguments);
  }
  
  /// 执行原始SQL插入
  Future<int> rawInsert(String sql, [List<dynamic>? arguments]) async {
    final db = await database;
    return await db.rawInsert(sql, arguments);
  }
  
  /// 执行原始SQL更新
  Future<int> rawUpdate(String sql, [List<dynamic>? arguments]) async {
    final db = await database;
    return await db.rawUpdate(sql, arguments);
  }
  
  /// 执行原始SQL删除
  Future<int> rawDelete(String sql, [List<dynamic>? arguments]) async {
    final db = await database;
    return await db.rawDelete(sql, arguments);
  }
  
  /// 插入数据到表
  Future<int> insert(String table, Map<String, dynamic> values) async {
    final db = await database;
    return await db.insert(table, values, conflictAlgorithm: ConflictAlgorithm.replace);
  }
  
  /// 查询表数据
  Future<List<Map<String, dynamic>>> query(
    String table, {
    bool? distinct,
    List<String>? columns,
    String? where,
    List<dynamic>? whereArgs,
    String? groupBy,
    String? having,
    String? orderBy,
    int? limit,
    int? offset,
  }) async {
    final db = await database;
    return await db.query(
      table,
      distinct: distinct,
      columns: columns,
      where: where,
      whereArgs: whereArgs,
      groupBy: groupBy,
      having: having,
      orderBy: orderBy,
      limit: limit,
      offset: offset,
    );
  }
  
  /// 更新表数据
  Future<int> update(
    String table,
    Map<String, dynamic> values, {
    String? where,
    List<dynamic>? whereArgs,
  }) async {
    final db = await database;
    return await db.update(
      table,
      values,
      where: where,
      whereArgs: whereArgs,
    );
  }
  
  /// 删除表数据
  Future<int> delete(
    String table, {
    String? where,
    List<dynamic>? whereArgs,
  }) async {
    final db = await database;
    return await db.delete(
      table,
      where: where,
      whereArgs: whereArgs,
    );
  }
  
  /// 获取表中记录数量
  Future<int> count(String table, {String? where, List<dynamic>? whereArgs}) async {
    final db = await database;
    final result = await db.rawQuery(
      'SELECT COUNT(*) as count FROM $table ${where != null ? 'WHERE $where' : ''}',
      whereArgs,
    );
    return Sqflite.firstIntValue(result) ?? 0;
  }
  
  /// 检查表是否存在
  Future<bool> tableExists(String tableName) async {
    final db = await database;
    final result = await db.rawQuery(
      "SELECT name FROM sqlite_master WHERE type='table' AND name=?",
      [tableName],
    );
    return result.isNotEmpty;
  }
  
  /// 执行事务
  Future<T> transaction<T>(Future<T> Function(Transaction txn) action) async {
    final db = await database;
    return await db.transaction(action);
  }
}