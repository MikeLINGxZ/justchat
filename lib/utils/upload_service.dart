import 'package:flutter/material.dart';
import 'package:lemon_tea/utils/llm/models/message.dart';
import 'package:lemon_tea/utils/cli/client/client.dart';
import 'package:lemon_tea/rpc/service.pb.dart' as grpc_service;

/// 文件上传服务
class UploadService {
  static final UploadService _instance = UploadService._internal();
  factory UploadService() => _instance;
  UploadService._internal();

  final Client _grpcClient = Client();

  /// 上传单个文件
  /// 
  /// 参数：
  /// - filePath: 文件路径
  /// 
  /// 返回：
  /// - 成功时返回resource_key，失败时返回null
  Future<String?> uploadFile(String filePath) async {
    try {
      // 检查gRPC客户端
      if (_grpcClient.stub == null) {
        await _grpcClient.init();
        if (_grpcClient.stub == null) {
          throw Exception("gRPC客户端初始化失败，请检查服务是否启动");
        }
      }

      // 创建上传请求
      final request = grpc_service.UploadFileRequest(
        filePath: filePath,
      );

      // 调用gRPC上传接口
      final response = await _grpcClient.stub!.uploadFile(request);
      
      debugPrint('文件上传成功: $filePath -> ${response.resourceKey}');
      return response.resourceKey;
    } catch (e) {
      debugPrint('文件上传失败: $e');
      return null;
    }
  }

  /// 批量上传文件
  /// 
  /// 参数：
  /// - filePaths: 文件路径列表
  /// 
  /// 返回：
  /// - 成功上传的文件resource_key列表
  Future<List<String>> uploadFiles(List<String> filePaths) async {
    List<String> resourceKeys = [];
    
    for (String filePath in filePaths) {
      String? resourceKey = await uploadFile(filePath);
      if (resourceKey != null) {
        resourceKeys.add(resourceKey);
      }
    }
    
    return resourceKeys;
  }

  /// 上传FileContent中的文件
  /// 
  /// 这个方法会将FileContent中的文件数据写入临时文件，
  /// 然后上传到服务器，返回resource_key
  /// 
  /// 参数：
  /// - fileContent: 文件内容对象
  /// 
  /// 返回：
  /// - 成功时返回resource_key，失败时返回null
  Future<String?> uploadFileContent(FileContent fileContent) async {
    // TODO: 实现将FileContent中的base64数据写入临时文件并上传
    // 这里需要根据实际需求来实现，可能需要：
    // 1. 将base64数据解码为二进制数据
    // 2. 写入临时文件
    // 3. 调用uploadFile方法
    // 4. 删除临时文件
    
    debugPrint('FileContent上传功能待实现: ${fileContent.name}');
    return null;
  }
}