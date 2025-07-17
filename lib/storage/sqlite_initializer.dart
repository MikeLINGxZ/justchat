import 'dart:async';
import 'package:flutter/foundation.dart';
import 'package:lemon_tea/models/conversation.dart';
import 'package:lemon_tea/models/debug_config.dart';
import 'package:lemon_tea/models/message.dart';
import 'package:lemon_tea/models/model.dart';
import 'package:lemon_tea/models/llm_provider.dart';
import 'package:path/path.dart';
import 'package:sqflite/sqflite.dart';
import 'package:path_provider/path_provider.dart';
import 'dart:io';

/// SQLite数据库初始化器
/// 负责数据库的初始化、表创建和版本升级
class SqliteDatabaseInitializer {
  static const String _databaseName = "lemon_tea.db";
  static const int _databaseVersion = 2; // 更新数据库版本
  
  Database? _database;
  final _initDBCompleter = Completer<Database>();
  
  /// 获取数据库实例
  Future<Database> get database async {
    if (_database != null) return _database!;
    
    // 如果数据库为空，初始化数据库
    await _initDatabase();
    return _database!;
  }
  
  /// 初始化数据库
  Future<void> _initDatabase() async {
    try {
      // 获取数据库路径
      Directory documentsDirectory = await getApplicationDocumentsDirectory();
      String path = join(documentsDirectory.path, _databaseName);
      
      debugPrint('初始化数据库: $path');
      
      // 打开数据库
      _database = await openDatabase(
        path,
        version: _databaseVersion,
        onCreate: _onCreate,
        onUpgrade: _onUpgrade,
        onOpen: (db) {
          debugPrint('数据库已打开: ${db.path}');
        },
      );
      
      // 完成初始化
      if (!_initDBCompleter.isCompleted) {
        _initDBCompleter.complete(_database);
        debugPrint('数据库初始化完成');
      }
    } catch (e) {
      debugPrint('数据库初始化错误: $e');
      if (!_initDBCompleter.isCompleted) {
        _initDBCompleter.completeError(e);
      }
      rethrow;
    }
  }
  
  /// 创建数据库表
  Future<void> _onCreate(Database db, int version) async {
    debugPrint('创建数据库表，版本: $version');
    try {
      // 创建会话表
      await db.execute(Conversation.createTableSql());
      debugPrint('会话表创建成功');
      
      // 创建消息表
      await db.execute(Message.createTableSql());
      debugPrint('消息表创建成功');
      
      // 创建模型表
      await db.execute(Model.createTableSql());
      debugPrint('模型表创建成功');
      
      // 创建提供商表
      await db.execute(LlmProvider.createTableSql());
      debugPrint('提供商表创建成功');

      // debug配置表
      await db.execute(DebugConfig.createTableSql());
      debugPrint('debug配置表创建成功');
    } catch (e) {
      debugPrint('创建数据库表失败: $e');
      rethrow;
    }
  }
  
  /// 数据库升级
  Future<void> _onUpgrade(Database db, int oldVersion, int newVersion) async {
    // 处理数据库版本升级
    debugPrint('升级数据库: $oldVersion -> $newVersion');
    
    if (oldVersion < 2) {
      // 版本2: 添加seq_id字段到llm_providers和models表
      try {
        // 添加seq_id字段到llm_providers表
        await db.execute('ALTER TABLE ${LlmProvider.tableName()} ADD COLUMN seq_id INTEGER NOT NULL DEFAULT 0');
        debugPrint('llm_providers表添加seq_id字段成功');
        
        // 添加seq_id字段到models表
        await db.execute('ALTER TABLE ${Model.tableName()} ADD COLUMN seq_id INTEGER NOT NULL DEFAULT 0');
        debugPrint('models表添加seq_id字段成功');
        
        // 更新现有记录的seq_id
        int providerSeqId = 1;
        final providers = await db.query(LlmProvider.tableName());
        for (var provider in providers) {
          await db.update(
            LlmProvider.tableName(),
            {'seq_id': providerSeqId++},
            where: 'id = ?',
            whereArgs: [provider['id']],
          );
          
          // 更新该提供商下的所有模型
          int modelSeqId = 1;
          final models = await db.query(
            Model.tableName(),
            where: 'llm_provider_id = ?',
            whereArgs: [provider['id']],
          );
          
          for (var model in models) {
            await db.update(
              Model.tableName(),
              {'seq_id': modelSeqId++},
              where: 'id = ?',
              whereArgs: [model['id']],
            );
          }
        }
        
        debugPrint('更新现有记录的seq_id成功');
      } catch (e) {
        debugPrint('添加seq_id字段失败: $e');
        rethrow;
      }
    }
  }
  
  /// 关闭数据库连接
  Future<void> close() async {
    final db = _database;
    if (db != null) {
      await db.close();
      _database = null;
      debugPrint('数据库连接已关闭');
    }
  }
  
  /// 获取数据库路径
  Future<String> getDatabasePath() async {
    Directory documentsDirectory = await getApplicationDocumentsDirectory();
    return join(documentsDirectory.path, _databaseName);
  }
  
  /// 删除数据库
  Future<void> deleteDatabase() async {
    await close();
    final path = await getDatabasePath();
    await databaseFactory.deleteDatabase(path);
    _database = null;
    debugPrint('数据库已删除: $path');
  }
} 