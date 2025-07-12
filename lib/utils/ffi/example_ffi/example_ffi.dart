
import 'dart:ffi';

import 'package:lemon_tea/utils/ffi/example_ffi/example_ffi.g.dart';

class ExampleFfi {

  static String libPath() {
    return "example_ffi_arm64.dylib";
  }

  static final DynamicLibrary _dylib = DynamicLibrary.open(libPath());

  static ExampleFfiGenerate instance() {
    return ExampleFfiGenerate(_dylib);
  }
}