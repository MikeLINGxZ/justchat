
import 'nativeflplug_platform_interface.dart';

class Nativeflplug {
  Future<String?> getPlatformVersion() {
    return NativeflplugPlatform.instance.getPlatformVersion();
  }
}
