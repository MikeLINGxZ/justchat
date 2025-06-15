import 'dart:io';

class System {
  static bool get isDesktop {
    return Platform.isWindows || Platform.isMacOS || Platform.isLinux;
  }
  static bool get isMobile {
    return Platform.isAndroid || Platform.isIOS;
  }

  static bool get isAndroid {
    return Platform.isAndroid;
  }

  static bool get isIOS {
    return Platform.isIOS;
  }

  static bool get isWindows {
    return Platform.isWindows;
  }

  static bool get isMacOS {
    return Platform.isMacOS;
  }

  static bool get isLinux {
    return Platform.isLinux;
  }
}