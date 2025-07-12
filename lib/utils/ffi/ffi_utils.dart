import 'dart:ffi';
import 'package:ffi/ffi.dart';
import 'package:lemon_tea/utils/ffi/ffi.dart';

/// FFI工具类，用于加载动态库并创建Nativefl实例
class FfiUtils {
  /// 加载动态库并返回Nativefl实例
  /// 
  /// [libPath] 动态库路径，如果为空则尝试使用默认路径
  /// 返回创建的Nativefl实例，加载失败则抛出异常
  static Nativefl loadNativefl({String? libPath}) {
    DynamicLibrary? dylib;
    
    try {
      if (libPath != null) {
        // 使用指定路径加载
        dylib = DynamicLibrary.open(libPath);
      } else {
        // 尝试使用默认路径或进程加载
        dylib = DynamicLibrary.process();
      }
    } catch (e) {
      throw Exception('加载动态库失败: $e');
    }
    
    if (dylib == null) {
      throw Exception('无法加载动态库');
    }
    
    // 创建并返回Nativefl实例
    return Nativefl(dylib);
  }
  
  /// 尝试从多个可能的路径加载动态库并返回Nativefl实例
  /// 
  /// [possiblePaths] 可能的动态库路径列表
  /// 返回创建的Nativefl实例，如果所有路径都加载失败则抛出异常
  static Nativefl tryLoadNativefl(List<String> possiblePaths) {
    DynamicLibrary? dylib;
    String? loadedPath;
    
    for (final path in possiblePaths) {
      try {
        dylib = DynamicLibrary.open(path);
        loadedPath = path;
        break;
      } catch (e) {
        // 继续尝试下一个路径
      }
    }
    
    if (dylib == null) {
      throw Exception('无法加载动态库，所有路径尝试均失败');
    }
    
    // 创建并返回Nativefl实例
    return Nativefl(dylib);
  }
  
  /// 加载多个动态库路径并返回所有成功加载的Nativefl实例
  /// 
  /// [libPaths] 动态库路径列表
  /// 返回一个Map，键为成功加载的路径，值为对应的Nativefl实例
  static Map<String, Nativefl> loadMultipleLibs(List<String> libPaths) {
    final Map<String, Nativefl> result = {};
    final List<String> failedPaths = [];
    
    for (final path in libPaths) {
      try {
        final dylib = DynamicLibrary.open(path);
        result[path] = Nativefl(dylib);
      } catch (e) {
        failedPaths.add(path);
      }
    }
    
    if (result.isEmpty) {
      throw Exception('所有动态库路径加载失败: $failedPaths');
    }
    
    return result;
  }
  
  /// 尝试加载多个动态库，直到成功加载一个为止，并返回加载信息
  /// 
  /// [libPaths] 动态库路径列表
  /// 返回一个包含加载结果的对象，包括成功加载的路径和Nativefl实例
  static ({String path, Nativefl nativefl}) loadFirstSuccessful(List<String> libPaths) {
    for (final path in libPaths) {
      try {
        final dylib = DynamicLibrary.open(path);
        return (path: path, nativefl: Nativefl(dylib));
      } catch (e) {
        // 继续尝试下一个路径
      }
    }
    
    throw Exception('所有动态库路径加载失败');
  }
} 