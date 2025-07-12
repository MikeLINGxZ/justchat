
import 'dart:ffi';
import 'package:ffi/ffi.dart';
import 'package:lemon_tea/utils/ffi/example_ffi/example_ffi.g.dart';

class ExampleFfi {
  static ExampleFfi? _instance;
  static ExampleFfiGenerate? _generateInstance;

  static String libPath() {
    return "example_ffi_arm64.dylib";
  }

  static final DynamicLibrary _dylib = DynamicLibrary.open(libPath());

  factory ExampleFfi() {
    _instance ??= ExampleFfi._internal();
    return _instance!;
  }

  ExampleFfi._internal();

  static ExampleFfiGenerate instance() {
    _generateInstance ??= ExampleFfiGenerate(_dylib);
    return _generateInstance!;
  }
  
  /// 封装ProcessString方法，处理字符串转换
  /// 
  /// 参数:
  /// - [input]: 输入字符串
  /// 
  /// 返回:
  /// - 处理后的字符串结果
  static String processString(String input) {
    // 获取FFI实例
    final ffiInstance = instance();
    
    // 将Dart字符串转换为C字符串
    final inputPtr = input.toNativeUtf8().cast<Char>();
    
    // 调用本地函数
    final resultPtr = ffiInstance.ProcessString(inputPtr);
    
    // 将结果转换回Dart字符串
    String result = "";
    if (resultPtr != nullptr) {
      result = resultPtr.cast<Utf8>().toDartString();
    }
    
    // 释放内存
    calloc.free(inputPtr);
    
    return result;
  }
}