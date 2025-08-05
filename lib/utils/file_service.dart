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
            final fileType = _getFileType(file.extension ?? '', file.bytes!);
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

  /// 根据文件扩展名和文件内容魔术字节判断文件类型
  static String _getFileType(String extension, Uint8List bytes) {
    extension = extension.toLowerCase();
    
    // 首先通过文件魔术检查是否为文本文件
    if (_isTextFile(bytes)) {
      return 'text';
    }
    
    // 图片类型（通过文件魔术检测）
    if (_isImageFile(bytes)) {
      return 'image';
    }
    
    // 音频类型（通过文件魔术检测）
    if (_isAudioFile(bytes)) {
      return 'audio';
    }
    
    // 视频类型（通过文件魔术检测）
    if (_isVideoFile(bytes)) {
      return 'video';
    }
    
    // 通过扩展名判断文档类型
    if (['pdf', 'doc', 'docx', 'rtf', 'xls', 'xlsx', 'ppt', 'pptx'].contains(extension)) {
      return 'document';
    }
    
    return 'other';
  }
  
  /// 通过文件魔术检测是否为文本文件
  static bool _isTextFile(Uint8List bytes) {
    if (bytes.isEmpty) return false;
    
    // 检查常见的二进制文件魔术字节
    final magicBytes = bytes.take(16).toList();
    
    // 检查是否包含常见的二进制文件头
    final binarySignatures = [
      [0xFF, 0xD8, 0xFF], // JPEG
      [0x89, 0x50, 0x4E, 0x47], // PNG
      [0x47, 0x49, 0x46, 0x38], // GIF
      [0x42, 0x4D], // BMP
      [0x52, 0x49, 0x46, 0x46], // RIFF (AVI, WAV等)
      [0x50, 0x4B, 0x03, 0x04], // ZIP
      [0x50, 0x4B, 0x05, 0x06], // ZIP (empty)
      [0x50, 0x4B, 0x07, 0x08], // ZIP (spanned)
      [0x25, 0x50, 0x44, 0x46], // PDF
      [0xD0, 0xCF, 0x11, 0xE0], // Microsoft Office
      [0x4F, 0x67, 0x67, 0x53], // OGG
      [0x66, 0x74, 0x79, 0x70], // MP4 (offset 4)
      [0x00, 0x00, 0x00, 0x18, 0x66, 0x74, 0x79, 0x70], // MP4
      [0x00, 0x00, 0x00, 0x1C, 0x66, 0x74, 0x79, 0x70], // MP4
      [0x49, 0x44, 0x33], // MP3
      [0xFF, 0xFB], // MP3
      [0xFF, 0xF3], // MP3
      [0xFF, 0xF2], // MP3
    ];
    
    // 如果匹配任何二进制文件头，则不是文本文件
    for (final signature in binarySignatures) {
      if (_matchesSignature(magicBytes, signature)) {
        return false;
      }
    }
    
    // 检查文件是否包含有效的UTF-8字符
    try {
      final text = utf8.decode(bytes, allowMalformed: false);
      
      // 计算非打印字符的比例
      int nonPrintableCount = 0;
      int totalChars = text.length;
      
      for (int i = 0; i < text.length; i++) {
        final charCode = text.codeUnitAt(i);
        // 检查是否为非打印字符（排除常见的空白字符）
        if (charCode < 32 && charCode != 9 && charCode != 10 && charCode != 13) {
          nonPrintableCount++;
        } else if (charCode == 0 || charCode > 126 && charCode < 160) {
          nonPrintableCount++;
        }
      }
      
      // 如果非打印字符比例小于5%，认为是文本文件
      return totalChars > 0 && (nonPrintableCount / totalChars) < 0.05;
    } catch (e) {
      // UTF-8解码失败，可能是二进制文件
      return false;
    }
  }
  
  /// 通过文件魔术检测是否为图片文件
  static bool _isImageFile(Uint8List bytes) {
    if (bytes.length < 4) return false;
    
    final magicBytes = bytes.take(16).toList();
    
    return _matchesSignature(magicBytes, [0xFF, 0xD8, 0xFF]) || // JPEG
           _matchesSignature(magicBytes, [0x89, 0x50, 0x4E, 0x47]) || // PNG
           _matchesSignature(magicBytes, [0x47, 0x49, 0x46, 0x38]) || // GIF
           _matchesSignature(magicBytes, [0x42, 0x4D]) || // BMP
           _matchesSignature(magicBytes, [0x52, 0x49, 0x46, 0x46]); // WebP (RIFF)
  }
  
  /// 通过文件魔术检测是否为音频文件
  static bool _isAudioFile(Uint8List bytes) {
    if (bytes.length < 4) return false;
    
    final magicBytes = bytes.take(16).toList();
    
    return _matchesSignature(magicBytes, [0x49, 0x44, 0x33]) || // MP3 (ID3)
           _matchesSignature(magicBytes, [0xFF, 0xFB]) || // MP3
           _matchesSignature(magicBytes, [0xFF, 0xF3]) || // MP3
           _matchesSignature(magicBytes, [0xFF, 0xF2]) || // MP3
           _matchesSignature(magicBytes, [0x52, 0x49, 0x46, 0x46]) || // WAV (RIFF)
           _matchesSignature(magicBytes, [0x4F, 0x67, 0x67, 0x53]) || // OGG
           _matchesSignature(magicBytes, [0x66, 0x4C, 0x61, 0x43]); // FLAC
  }
  
  /// 通过文件魔术检测是否为视频文件
  static bool _isVideoFile(Uint8List bytes) {
    if (bytes.length < 8) return false;
    
    final magicBytes = bytes.take(16).toList();
    
    // MP4系列格式检测
    if (bytes.length >= 8) {
      final fourCC = bytes.sublist(4, 8);
      if (_matchesSignature(fourCC, [0x66, 0x74, 0x79, 0x70])) { // ftyp
        return true;
      }
    }
    
    return _matchesSignature(magicBytes, [0x52, 0x49, 0x46, 0x46]) || // AVI (RIFF)
           _matchesSignature(magicBytes, [0x00, 0x00, 0x01, 0xBA]) || // MPEG
           _matchesSignature(magicBytes, [0x00, 0x00, 0x01, 0xB3]) || // MPEG
           _matchesSignature(magicBytes, [0x1A, 0x45, 0xDF, 0xA3]); // MKV
  }
  
  /// 检查字节序列是否匹配签名
  static bool _matchesSignature(List<int> bytes, List<int> signature) {
    if (bytes.length < signature.length) return false;
    
    for (int i = 0; i < signature.length; i++) {
      if (bytes[i] != signature[i]) return false;
    }
    return true;
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
    
    // 根据文件扩展名和类型推断MIME类型
    final mimeType = _getMimeType(name, type);
    
    return FileContent(
      name: name,
      mimeType: mimeType,
      type: type,
      data: base64Data,
      size: size,
    );
  }

  /// 根据文件名和类型推断MIME类型
  static String _getMimeType(String fileName, [String? detectedType]) {
    // 如果检测到是文本文件但没有明确扩展名，返回通用文本类型
    if (detectedType == 'text') {
      final extension = fileName.split('.').last.toLowerCase();
      // 根据扩展名提供更具体的文本MIME类型
      switch (extension) {
        case 'txt':
          return 'text/plain';
        case 'md':
          return 'text/markdown';
        case 'csv':
          return 'text/csv';
        case 'json':
          return 'application/json';
        case 'xml':
          return 'application/xml';
        case 'html':
        case 'htm':
          return 'text/html';
        case 'css':
          return 'text/css';
        case 'js':
          return 'application/javascript';
        case 'ts':
          return 'application/typescript';
        case 'py':
          return 'text/x-python';
        case 'java':
          return 'text/x-java-source';
        case 'cpp':
        case 'cc':
        case 'cxx':
          return 'text/x-c++src';
        case 'c':
          return 'text/x-csrc';
        case 'h':
          return 'text/x-chdr';
        case 'dart':
          return 'application/dart';
        default:
          return 'text/plain';
      }
    }
    
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
      case 'text':
        return Icons.text_snippet;
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