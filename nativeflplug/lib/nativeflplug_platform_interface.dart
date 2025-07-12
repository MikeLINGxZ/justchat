import 'package:plugin_platform_interface/plugin_platform_interface.dart';

import 'nativeflplug_method_channel.dart';

abstract class NativeflplugPlatform extends PlatformInterface {
  /// Constructs a NativeflplugPlatform.
  NativeflplugPlatform() : super(token: _token);

  static final Object _token = Object();

  static NativeflplugPlatform _instance = MethodChannelNativeflplug();

  /// The default instance of [NativeflplugPlatform] to use.
  ///
  /// Defaults to [MethodChannelNativeflplug].
  static NativeflplugPlatform get instance => _instance;

  /// Platform-specific implementations should set this with their own
  /// platform-specific class that extends [NativeflplugPlatform] when
  /// they register themselves.
  static set instance(NativeflplugPlatform instance) {
    PlatformInterface.verifyToken(instance, _token);
    _instance = instance;
  }

  Future<String?> getPlatformVersion() {
    throw UnimplementedError('platformVersion() has not been implemented.');
  }
}
