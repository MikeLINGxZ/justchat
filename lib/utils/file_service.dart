import 'dart:convert';
import 'dart:typed_data';
import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:lemon_tea/utils/llm/models/message.dart';

class FileService {
  /// 选择图片文件
  static Future<List<FileContent>?> pickImages() async {
    try {
      FilePickerResult? result = await FilePicker.platform.pickFiles(
        type: FileType.image,
        allowMultiple: true,
        withData: true,
      );

      if (result != null) {
        List<FileContent> files = [];
        for (PlatformFile file in result.files) {
          if (file.bytes != null) {
            final fileContent = await _createFileContent(
              file.name,
              file.bytes!,
              'image',
              file.size,
            );
            files.add(fileContent);
          }
        }
        return files;
      }
    } catch (e) {
      debugPrint('选择图片时出错: $e');
    }
    return null;
  }

  /// 选择文档文件
  static Future<List<FileContent>?> pickDocuments() async {
    try {
      FilePickerResult? result = await FilePicker.platform.pickFiles(
        type: FileType.custom,
        allowedExtensions: ['pdf', 'doc', 'docx', 'txt', 'rtf', 'md', 'csv', 'xls', 'xlsx'],
        allowMultiple: true,
        withData: true,
      );

      if (result != null) {
        List<FileContent> files = [];
        for (PlatformFile file in result.files) {
          if (file.bytes != null) {
            final fileContent = await _createFileContent(
              file.name,
              file.bytes!,
              'document',
              file.size,
            );
            files.add(fileContent);
          }
        }
        return files;
      }
    } catch (e) {
      debugPrint('选择文档时出错: $e');
    }
    return null;
  }

  /// 选择任意类型文件
  static Future<List<FileContent>?> pickAnyFiles() async {
    try {
      FilePickerResult? result = await FilePicker.platform.pickFiles(
        type: FileType.any,
        allowMultiple: true,
        withData: true,
      );

      if (result != null) {
        List<FileContent> files = [];
        for (PlatformFile file in result.files) {
          if (file.bytes != null) {
            final fileType = _getFileType(file.extension ?? '');
            final fileContent = await _createFileContent(
              file.name,
              file.bytes!,
              fileType,
              file.size,
            );
            files.add(fileContent);
          }
        }
        return files;
      }
    } catch (e) {
      debugPrint('选择文件时出错: $e');
    }
    return null;
  }

  /// 根据文件扩展名判断文件类型
  static String _getFileType(String extension) {
    extension = extension.toLowerCase();
    
    // 图片类型
    if (['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp', 'svg'].contains(extension)) {
      return 'image';
    }
    
    // 文档类型
    if (['pdf', 'doc', 'docx', 'txt', 'rtf', 'md', 'csv', 'xls', 'xlsx', 'ppt', 'pptx'].contains(extension)) {
      return 'document';
    }
    
    // 音频类型
    if (['mp3', 'wav', 'aac', 'flac', 'ogg', 'm4a'].contains(extension)) {
      return 'audio';
    }
    
    // 视频类型
    if (['mp4', 'avi', 'mov', 'wmv', 'flv', 'webm', 'mkv'].contains(extension)) {
      return 'video';
    }
    
    return 'other';
  }

  /// 创建FileContent对象
  static Future<FileContent> _createFileContent(
    String name,
    Uint8List bytes,
    String type,
    int size,
  ) async {
    // 将文件数据编码为Base64
    final base64Data = base64Encode(bytes);
    
    // 根据文件扩展名推断MIME类型
    final mimeType = _getMimeType(name);
    
    return FileContent(
      name: name,
      mimeType: mimeType,
      type: type,
      data: base64Data,
      size: size,
    );
  }

  /// 根据文件名推断MIME类型
  static String _getMimeType(String fileName) {
    final extension = fileName.split('.').last.toLowerCase();
    
    switch (extension) {
      // 图片
      case 'jpg':
      case 'jpeg':
        return 'image/jpeg';
      case 'png':
        return 'image/png';
      case 'gif':
        return 'image/gif';
      case 'bmp':
        return 'image/bmp';
      case 'webp':
        return 'image/webp';
      case 'svg':
        return 'image/svg+xml';
      
      // 文档
      case 'pdf':
        return 'application/pdf';
      case 'doc':
        return 'application/msword';
      case 'docx':
        return 'application/vnd.openxmlformats-officedocument.wordprocessingml.document';
      case 'txt':
        return 'text/plain';
      case 'rtf':
        return 'application/rtf';
      case 'md':
        return 'text/markdown';
      case 'csv':
        return 'text/csv';
      case 'xls':
        return 'application/vnd.ms-excel';
      case 'xlsx':
        return 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet';
      case 'ppt':
        return 'application/vnd.ms-powerpoint';
      case 'pptx':
        return 'application/vnd.openxmlformats-officedocument.presentationml.presentation';
      
      // 音频
      case 'mp3':
        return 'audio/mpeg';
      case 'wav':
        return 'audio/wav';
      case 'aac':
        return 'audio/aac';
      case 'flac':
        return 'audio/flac';
      case 'ogg':
        return 'audio/ogg';
      case 'm4a':
        return 'audio/mp4';
      
      // 视频
      case 'mp4':
        return 'video/mp4';
      case 'avi':
        return 'video/x-msvideo';
      case 'mov':
        return 'video/quicktime';
      case 'wmv':
        return 'video/x-ms-wmv';
      case 'flv':
        return 'video/x-flv';
      case 'webm':
        return 'video/webm';
      case 'mkv':
        return 'video/x-matroska';
      
      default:
        return 'application/octet-stream';
    }
  }

  /// 格式化文件大小显示
  static String formatFileSize(int bytes) {
    if (bytes < 1024) {
      return '${bytes}B';
    } else if (bytes < 1024 * 1024) {
      return '${(bytes / 1024).toStringAsFixed(1)}KB';
    } else if (bytes < 1024 * 1024 * 1024) {
      return '${(bytes / (1024 * 1024)).toStringAsFixed(1)}MB';
    } else {
      return '${(bytes / (1024 * 1024 * 1024)).toStringAsFixed(1)}GB';
    }
  }

  /// 获取文件类型的图标
  static IconData getFileTypeIcon(String type) {
    switch (type) {
      case 'image':
        return Icons.image;
      case 'document':
        return Icons.description;
      case 'audio':
        return Icons.audiotrack;
      case 'video':
        return Icons.videocam;
      default:
        return Icons.insert_drive_file;
    }
  }

  /// 检查文件是否为图片
  static bool isImage(FileContent file) {
    return file.type == 'image';
  }

  /// 获取图片的预览数据
  static Uint8List? getImagePreviewData(FileContent file) {
    if (isImage(file) && file.data != null) {
      try {
        return base64Decode(file.data!);
      } catch (e) {
        debugPrint('解码图片数据时出错: $e');
      }
    }
    return null;
  }
} 