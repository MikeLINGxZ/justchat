
bool isValidDomain(String input) {
  final domainRegex = RegExp(
    r'^((?!-)[A-Za-z0-9-]{1,63}(?<!-)\.)+[A-Za-z]{2,6}$',
    caseSensitive: false,
  );
  return domainRegex.hasMatch(input);
}

bool isValidIPv4(String input) {
  final ipv4Regex = RegExp(
    r'^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$',
  );
  return ipv4Regex.hasMatch(input);
}

bool isPortStr(String port) {
  try {
    final portNum = int.parse(port);
    return portNum >= 1 && portNum <= 65535;
  } catch (e) {
    return false;
  }
}

bool isValidProxyAddress(String input) {
  final regex = RegExp(
      r'^(?:(https?|socks5):\/\/)?'       // optional scheme
      r'(?:([\w\-]+):([\w\-]+)@)?'        // optional username:password@
      r'([a-zA-Z0-9.-]+)'                 // host
      r':'
      r'(\d{1,5})$'                       // port
  );

  final match = regex.firstMatch(input);
  if (match == null) return false;

  final port = int.tryParse(match.group(5)!);
  if (port == null || port < 1 || port > 65535) return false;

  return true;
}

bool isValidEmail(String input) {
  final regex = RegExp(
      r"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$"
  );
  return regex.hasMatch(input);
}

String displayError(Object e) {
  return e.toString().split('Exception: ').last;
}