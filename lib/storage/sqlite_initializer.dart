import 'dart:async';
import 'package:lemon_tea/models/conversation.dart';
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
  static const int _databaseVersion = 1;
  
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
      
      // 打开数据库
      _database = await openDatabase(
        path,
        version: _databaseVersion,
        onCreate: _onCreate,
        onUpgrade: _onUpgrade,
      );
      
      // 完成初始化
      if (!_initDBCompleter.isCompleted) {
        _initDBCompleter.complete(_database);
      }
    } catch (e) {
      if (!_initDBCompleter.isCompleted) {
        _initDBCompleter.completeError(e);
      }
      rethrow;
    }
  }
  
  /// 创建数据库表
  Future<void> _onCreate(Database db, int version) async {
    // 创建会话表
    await db.execute(Conversation.createTableSql());
    
    // 创建消息表
    await db.execute(Message.createTableSql());
    
    // 创建模型表
    await db.execute(Model.createTableSql());
    
    // 创建提供商表
    await db.execute(LlmProvider.createTableSql());
  }
  
  /// 数据库升级
  Future<void> _onUpgrade(Database db, int oldVersion, int newVersion) async {
    // 处理数据库版本升级
    if (oldVersion < 2) {
      // 未来版本2的升级代码
    }
  }
  
  /// 关闭数据库连接
  Future<void> close() async {
    final db = _database;
    if (db != null) {
      await db.close();
      _database = null;
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
  }
} 