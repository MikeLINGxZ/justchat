
import 'dart:ffi';

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
}