import 'package:flutter_test/flutter_test.dart';
import 'package:nativeflplug/nativeflplug.dart';
import 'package:nativeflplug/nativeflplug_platform_interface.dart';
import 'package:nativeflplug/nativeflplug_method_channel.dart';
import 'package:plugin_platform_interface/plugin_platform_interface.dart';

class MockNativeflplugPlatform
    with MockPlatformInterfaceMixin
    implements NativeflplugPlatform {

  @override
  Future<String?> getPlatformVersion() => Future.value('42');
}

void main() {
  final NativeflplugPlatform initialPlatform = NativeflplugPlatform.instance;

  test('$MethodChannelNativeflplug is the default instance', () {
    expect(initialPlatform, isInstanceOf<MethodChannelNativeflplug>());
  });

  test('getPlatformVersion', () async {
    Nativeflplug nativeflplugPlugin = Nativeflplug();
    MockNativeflplugPlatform fakePlatform = MockNativeflplugPlatform();
    NativeflplugPlatform.instance = fakePlatform;

    expect(await nativeflplugPlugin.getPlatformVersion(), '42');
  });
}
