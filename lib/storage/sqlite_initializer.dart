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
  static const int _databaseVersion = 5; // 更新数据库版本以添加FTS全文搜索支持
  
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
    
    if (oldVersion < 3) {
      // 版本3: 添加deleted字段到conversations和messages表
      try {
        // 检查conversations表是否存在deleted列
        final conversationTableInfo = await db.rawQuery("PRAGMA table_info(${Conversation.tableName()})");
        final conversationHasDeleted = conversationTableInfo.any((column) => column['name'] == 'deleted');
        
        if (!conversationHasDeleted) {
          // 添加deleted字段到conversations表
          await db.execute('ALTER TABLE ${Conversation.tableName()} ADD COLUMN deleted INTEGER NOT NULL DEFAULT 0');
          debugPrint('conversations表添加deleted字段成功');
        }
        
        // 检查messages表是否存在deleted列
        final messageTableInfo = await db.rawQuery("PRAGMA table_info(${Message.tableName()})");
        final messageHasDeleted = messageTableInfo.any((column) => column['name'] == 'deleted');
        
        if (!messageHasDeleted) {
          // 添加deleted字段到messages表
          await db.execute('ALTER TABLE ${Message.tableName()} ADD COLUMN deleted INTEGER NOT NULL DEFAULT 0');
          debugPrint('messages表添加deleted字段成功');
        }
        
        // 检查messages表的updated_at列
        final messageHasUpdatedAt = messageTableInfo.any((column) => column['name'] == 'updated_at');
        if (messageHasUpdatedAt) {
          // SQLite不支持直接删除列，我们保持兼容并在代码中处理
          debugPrint('messages表包含updated_at列，代码已兼容处理');
        }
        
        debugPrint('数据库升级到版本3完成');
      } catch (e) {
        debugPrint('升级到版本3失败: $e');
        rethrow;
      }
    }
    
    if (oldVersion < 4) {
      // 版本4: 添加reasoning_content字段到messages表
      try {
        // 检查messages表是否存在reasoning_content列
        final messageTableInfo = await db.rawQuery("PRAGMA table_info(${Message.tableName()})");
        final messageHasReasoningContent = messageTableInfo.any((column) => column['name'] == 'reasoning_content');
        
        if (!messageHasReasoningContent) {
          // 添加reasoning_content字段到messages表
          await db.execute('ALTER TABLE ${Message.tableName()} ADD COLUMN reasoning_content TEXT');
          debugPrint('messages表添加reasoning_content字段成功');
        }
        
        debugPrint('数据库升级到版本4完成');
      } catch (e) {
        debugPrint('升级到版本4失败: $e');
        rethrow;
      }
    }
    
    if (oldVersion < 5) {
      // 版本5: 添加FTS全文搜索支持
      try {
        // 检查FTS表是否已存在
        final tableList = await db.rawQuery(
          "SELECT name FROM sqlite_master WHERE type='table' AND name='${Message.tableName()}_fts'"
        );
        
        if (tableList.isEmpty) {
          // 创建FTS虚拟表
          await db.execute('''
            CREATE VIRTUAL TABLE ${Message.tableName()}_fts USING fts5(
              id UNINDEXED,
              conversation_id UNINDEXED,
              content,
              reasoning_content,
              content=${Message.tableName()},
              content_rowid=rowid
            )
          ''');
          debugPrint('FTS虚拟表创建成功');
          
          // 创建触发器
          await db.execute('''
            CREATE TRIGGER ${Message.tableName()}_fts_insert AFTER INSERT ON ${Message.tableName()} BEGIN
              INSERT INTO ${Message.tableName()}_fts(id, conversation_id, content, reasoning_content)
              VALUES (new.id, new.conversation_id, new.content, new.reasoning_content);
            END
          ''');
          
          await db.execute('''
            CREATE TRIGGER ${Message.tableName()}_fts_delete AFTER DELETE ON ${Message.tableName()} BEGIN
              DELETE FROM ${Message.tableName()}_fts WHERE id = old.id;
            END
          ''');
          
          await db.execute('''
            CREATE TRIGGER ${Message.tableName()}_fts_update AFTER UPDATE ON ${Message.tableName()} BEGIN
              UPDATE ${Message.tableName()}_fts 
              SET content = new.content, reasoning_content = new.reasoning_content
              WHERE id = new.id;
            END
          ''');
          debugPrint('FTS触发器创建成功');
          
          // 为现有数据重建FTS索引
          await db.execute('INSERT INTO ${Message.tableName()}_fts(${Message.tableName()}_fts) VALUES(\'rebuild\')');
          debugPrint('FTS索引重建成功');
        }
        
        debugPrint('数据库升级到版本5完成');
      } catch (e) {
        debugPrint('升级到版本5失败: $e');
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

  /// 获取当前数据库版本
  Future<int> getCurrentDatabaseVersion() async {
    final db = await database;
    final result = await db.rawQuery('PRAGMA user_version');
    return result.first['user_version'] as int;
  }

  /// 检查FTS表是否存在
  Future<bool> checkFtsTableExists() async {
    try {
      final db = await database;
      final result = await db.rawQuery(
        "SELECT name FROM sqlite_master WHERE type='table' AND name='${Message.tableName()}_fts'"
      );
      return result.isNotEmpty;
    } catch (e) {
      debugPrint('检查FTS表失败: $e');
      return false;
    }
  }

  /// 强制重建FTS表和触发器
  Future<bool> rebuildFtsTable() async {
    try {
      final db = await database;
      
      // 删除现有的FTS表和触发器（如果存在）
      try {
        await db.execute('DROP TRIGGER IF EXISTS ${Message.tableName()}_fts_insert');
        await db.execute('DROP TRIGGER IF EXISTS ${Message.tableName()}_fts_delete');
        await db.execute('DROP TRIGGER IF EXISTS ${Message.tableName()}_fts_update');
        await db.execute('DROP TABLE IF EXISTS ${Message.tableName()}_fts');
        debugPrint('已清理现有FTS结构');
      } catch (e) {
        debugPrint('清理FTS结构时出错（可能不存在）: $e');
      }

      // 创建FTS虚拟表
      await db.execute('''
        CREATE VIRTUAL TABLE ${Message.tableName()}_fts USING fts5(
          id UNINDEXED,
          conversation_id UNINDEXED,
          content,
          reasoning_content,
          content=${Message.tableName()},
          content_rowid=rowid
        )
      ''');
      debugPrint('FTS虚拟表重建成功');
      
      // 创建触发器
      await db.execute('''
        CREATE TRIGGER ${Message.tableName()}_fts_insert AFTER INSERT ON ${Message.tableName()} BEGIN
          INSERT INTO ${Message.tableName()}_fts(id, conversation_id, content, reasoning_content)
          VALUES (new.id, new.conversation_id, new.content, new.reasoning_content);
        END
      ''');
      
      await db.execute('''
        CREATE TRIGGER ${Message.tableName()}_fts_delete AFTER DELETE ON ${Message.tableName()} BEGIN
          DELETE FROM ${Message.tableName()}_fts WHERE id = old.id;
        END
      ''');
      
      await db.execute('''
        CREATE TRIGGER ${Message.tableName()}_fts_update AFTER UPDATE ON ${Message.tableName()} BEGIN
          UPDATE ${Message.tableName()}_fts 
          SET content = new.content, reasoning_content = new.reasoning_content
          WHERE id = new.id;
        END
      ''');
      debugPrint('FTS触发器重建成功');
      
      // 为现有数据重建FTS索引
      await db.execute('INSERT INTO ${Message.tableName()}_fts(${Message.tableName()}_fts) VALUES(\'rebuild\')');
      debugPrint('FTS索引重建成功');
      
      return true;
    } catch (e) {
      debugPrint('重建FTS表失败: $e');
      return false;
    }
  }

  /// 诊断数据库状态
  Future<Map<String, dynamic>> diagnoseDatabaseState() async {
    try {
      final db = await database;
      final version = await getCurrentDatabaseVersion();
      final ftsExists = await checkFtsTableExists();
      
      // 获取所有表列表
      final tables = await db.rawQuery(
        "SELECT name FROM sqlite_master WHERE type='table'"
      );
      
      // 获取消息表的列信息
      final messageColumns = await db.rawQuery(
        "PRAGMA table_info(${Message.tableName()})"
      );
      
      return {
        'version': version,
        'fts_table_exists': ftsExists,
        'tables': tables.map((t) => t['name']).toList(),
        'message_columns': messageColumns.map((c) => c['name']).toList(),
      };
    } catch (e) {
      debugPrint('诊断数据库状态失败: $e');
      return {'error': e.toString()};
    }
  }
} 