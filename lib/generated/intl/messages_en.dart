// DO NOT EDIT. This is code generated via package:intl/generate_localized.dart
// This is a library that provides messages for a en locale. All the
// messages from the main program should be duplicated here with the same
// function name.

// Ignore issues from commonly used lints in this file.
// ignore_for_file:unnecessary_brace_in_string_interps, unnecessary_new
// ignore_for_file:prefer_single_quotes,comment_references, directives_ordering
// ignore_for_file:annotate_overrides,prefer_generic_function_type_aliases
// ignore_for_file:unused_import, file_names, avoid_escaping_inner_quotes
// ignore_for_file:unnecessary_string_interpolations, unnecessary_string_escapes

import 'package:intl/intl.dart';
import 'package:intl/message_lookup_by_library.dart';

final messages = new MessageLookup();

typedef String MessageIfAbsent(String messageStr, List<dynamic> args);

class MessageLookup extends MessageLookupByLibrary {
  String get localeName => 'en';

  final messages = _notInlinedMessages(_notInlinedMessages);
  static Map<String, Function> _notInlinedMessages(_) => <String, Function>{
    "addAccount": MessageLookupByLibrary.simpleMessage("Add account"),
    "advancedSettings": MessageLookupByLibrary.simpleMessage(
      "Advanced settings",
    ),
    "advancedSettingsMissing": MessageLookupByLibrary.simpleMessage(
      "Advanced settings missing,Please check the Settings",
    ),
    "appName": MessageLookupByLibrary.simpleMessage("Teamail"),
    "authenticationFailedMsg": MessageLookupByLibrary.simpleMessage(
      "Authentication failed, please check the email and password",
    ),
    "back": MessageLookupByLibrary.simpleMessage("Back"),
    "close": MessageLookupByLibrary.simpleMessage("Close"),
    "continuE": MessageLookupByLibrary.simpleMessage("Continue"),
    "editSettings": MessageLookupByLibrary.simpleMessage("Edit settings"),
    "email": MessageLookupByLibrary.simpleMessage("Email"),
    "emailFormatError": MessageLookupByLibrary.simpleMessage(
      "Email format error",
    ),
    "emailInvalid": MessageLookupByLibrary.simpleMessage("Email invalid"),
    "emailSettingsNotMatch": MessageLookupByLibrary.simpleMessage(
      "The login mailbox is different from the selected mailbox_ Continue to log in or modify the configuration?",
    ),
    "enableProxy": MessageLookupByLibrary.simpleMessage("Enable proxy"),
    "login": MessageLookupByLibrary.simpleMessage("Login"),
    "loginFailed": MessageLookupByLibrary.simpleMessage("Login failed"),
    "loginFailedMsg": MessageLookupByLibrary.simpleMessage(
      "Login failed, please check the email and password and make sure the mail server Settings are correct",
    ),
    "loginMessageMissing": MessageLookupByLibrary.simpleMessage(
      "Login message missing",
    ),
    "noSupportedEmailServer": MessageLookupByLibrary.simpleMessage(
      "No supported email server",
    ),
    "otherEmail": MessageLookupByLibrary.simpleMessage("Other Email"),
    "password": MessageLookupByLibrary.simpleMessage("password"),
    "passwordMustNotEmpty": MessageLookupByLibrary.simpleMessage(
      "pass must not empty",
    ),
    "pleaseEnterEmail": MessageLookupByLibrary.simpleMessage(
      "Please enter email",
    ),
    "pleaseEnterPassword": MessageLookupByLibrary.simpleMessage(
      "Please enter password or secret code",
    ),
    "port": MessageLookupByLibrary.simpleMessage("port"),
    "portFormatErr": MessageLookupByLibrary.simpleMessage(
      "port format is incorrect",
    ),
    "proxyAddress": MessageLookupByLibrary.simpleMessage("Proxy address"),
    "secureConnection": MessageLookupByLibrary.simpleMessage(
      "Secure connection",
    ),
    "server": MessageLookupByLibrary.simpleMessage("server"),
    "serverAddressFormatErr": MessageLookupByLibrary.simpleMessage(
      "address or ip format is incorrect",
    ),
    "serverType": MessageLookupByLibrary.simpleMessage("Server type"),
    "smtpErrCode211": MessageLookupByLibrary.simpleMessage(
      "System status, or system help reply",
    ),
    "smtpErrCode214": MessageLookupByLibrary.simpleMessage(
      "Help message (A response to the HELP command)",
    ),
    "smtpErrCode220": MessageLookupByLibrary.simpleMessage("Service ready"),
    "smtpErrCode221": MessageLookupByLibrary.simpleMessage(
      "Service closing transmission channel",
    ),
    "smtpErrCode221_2_0_0": MessageLookupByLibrary.simpleMessage("Goodbye"),
    "smtpErrCode235_2_7_0": MessageLookupByLibrary.simpleMessage(
      "Authentication succeeded",
    ),
    "smtpErrCode240": MessageLookupByLibrary.simpleMessage("QUIT"),
    "smtpErrCode250": MessageLookupByLibrary.simpleMessage(
      "Requested mail action okay, completed",
    ),
    "smtpErrCode251": MessageLookupByLibrary.simpleMessage(
      "User not local; will forward",
    ),
    "smtpErrCode252": MessageLookupByLibrary.simpleMessage(
      "Cannot verify the user, but it will try to deliver the message anyway",
    ),
    "smtpErrCode334": MessageLookupByLibrary.simpleMessage(
      "Server challenge - the text part contains the Base64-encoded challenge",
    ),
    "smtpErrCode354": MessageLookupByLibrary.simpleMessage("Start mail input"),
    "smtpErrCode421": MessageLookupByLibrary.simpleMessage(
      "Service not available, closing transmission channel",
    ),
    "smtpErrCode432_4_7_12": MessageLookupByLibrary.simpleMessage(
      "A password transition is needed",
    ),
    "smtpErrCode450": MessageLookupByLibrary.simpleMessage(
      "Requested mail action not taken: mailbox unavailable",
    ),
    "smtpErrCode451": MessageLookupByLibrary.simpleMessage(
      "Requested action aborted: local error in processing",
    ),
    "smtpErrCode451_4_4_1": MessageLookupByLibrary.simpleMessage(
      "IMAP server unavailable",
    ),
    "smtpErrCode452": MessageLookupByLibrary.simpleMessage(
      "Requested action not taken: insufficient system storage",
    ),
    "smtpErrCode454_4_7_0": MessageLookupByLibrary.simpleMessage(
      "Temporary authentication failure",
    ),
    "smtpErrCode455": MessageLookupByLibrary.simpleMessage(
      "Server unable to accommodate parameters",
    ),
    "smtpErrCode500": MessageLookupByLibrary.simpleMessage(
      "Syntax error, command unrecognized",
    ),
    "smtpErrCode500_5_5_6": MessageLookupByLibrary.simpleMessage(
      "Authentication Exchange line is too long",
    ),
    "smtpErrCode501": MessageLookupByLibrary.simpleMessage(
      "Syntax error in parameters or arguments",
    ),
    "smtpErrCode501_5_5_2": MessageLookupByLibrary.simpleMessage(
      "Cannot Base64-decode Client responses",
    ),
    "smtpErrCode501_5_7_0": MessageLookupByLibrary.simpleMessage(
      "Client initiated Authentication Exchange",
    ),
    "smtpErrCode502": MessageLookupByLibrary.simpleMessage(
      "Command not implemented",
    ),
    "smtpErrCode503": MessageLookupByLibrary.simpleMessage(
      "Bad sequence of commands",
    ),
    "smtpErrCode504": MessageLookupByLibrary.simpleMessage(
      "Command parameter is not implemented",
    ),
    "smtpErrCode504_5_5_4": MessageLookupByLibrary.simpleMessage(
      "Unrecognized authentication type",
    ),
    "smtpErrCode521": MessageLookupByLibrary.simpleMessage(
      "Server does not accept mail",
    ),
    "smtpErrCode523": MessageLookupByLibrary.simpleMessage("Encryption Needed"),
    "smtpErrCode530_5_7_0": MessageLookupByLibrary.simpleMessage(
      "Authentication required",
    ),
    "smtpErrCode534_5_7_9": MessageLookupByLibrary.simpleMessage(
      "Authentication mechanism is too weak",
    ),
    "smtpErrCode535_5_7_8": MessageLookupByLibrary.simpleMessage(
      "Authentication credentials invalid",
    ),
    "smtpErrCode538_5_7_11": MessageLookupByLibrary.simpleMessage(
      "Encryption required for requested authentication mechanism",
    ),
    "smtpErrCode550": MessageLookupByLibrary.simpleMessage(
      "Requested action not taken: mailbox unavailable",
    ),
    "smtpErrCode551": MessageLookupByLibrary.simpleMessage(
      "User not local; please try <forward-path>",
    ),
    "smtpErrCode552": MessageLookupByLibrary.simpleMessage(
      "Requested mail action aborted: exceeded storage allocation",
    ),
    "smtpErrCode553": MessageLookupByLibrary.simpleMessage(
      "Requested action not taken: mailbox name not allowed",
    ),
    "smtpErrCode554": MessageLookupByLibrary.simpleMessage(
      "Transaction has failed (Or, \'No SMTP service here\')",
    ),
    "smtpErrCode554_5_3_4": MessageLookupByLibrary.simpleMessage(
      "Message too big for system",
    ),
    "smtpErrCode556": MessageLookupByLibrary.simpleMessage(
      "Domain does not accept mail",
    ),
    "unknownError": MessageLookupByLibrary.simpleMessage("Unknown error"),
  };
}
