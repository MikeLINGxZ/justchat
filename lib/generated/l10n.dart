// GENERATED CODE - DO NOT MODIFY BY HAND
import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'intl/messages_all.dart';

// **************************************************************************
// Generator: Flutter Intl IDE plugin
// Made by Localizely
// **************************************************************************

// ignore_for_file: non_constant_identifier_names, lines_longer_than_80_chars
// ignore_for_file: join_return_with_assignment, prefer_final_in_for_each
// ignore_for_file: avoid_redundant_argument_values, avoid_escaping_inner_quotes

class S {
  S();

  static S? _current;

  static S get current {
    assert(
      _current != null,
      'No instance of S was loaded. Try to initialize the S delegate before accessing S.current.',
    );
    return _current!;
  }

  static const AppLocalizationDelegate delegate = AppLocalizationDelegate();

  static Future<S> load(Locale locale) {
    final name =
        (locale.countryCode?.isEmpty ?? false)
            ? locale.languageCode
            : locale.toString();
    final localeName = Intl.canonicalizedLocale(name);
    return initializeMessages(localeName).then((_) {
      Intl.defaultLocale = localeName;
      final instance = S();
      S._current = instance;

      return instance;
    });
  }

  static S of(BuildContext context) {
    final instance = S.maybeOf(context);
    assert(
      instance != null,
      'No instance of S present in the widget tree. Did you add S.delegate in localizationsDelegates?',
    );
    return instance!;
  }

  static S? maybeOf(BuildContext context) {
    return Localizations.of<S>(context, S);
  }

  /// `Teamail`
  String get appName {
    return Intl.message('Teamail', name: 'appName', desc: '', args: []);
  }

  /// `Other Email`
  String get otherEmail {
    return Intl.message('Other Email', name: 'otherEmail', desc: '', args: []);
  }

  /// `Back`
  String get back {
    return Intl.message('Back', name: 'back', desc: '', args: []);
  }

  /// `Add account`
  String get addAccount {
    return Intl.message('Add account', name: 'addAccount', desc: '', args: []);
  }

  /// `Login`
  String get login {
    return Intl.message('Login', name: 'login', desc: '', args: []);
  }

  /// `Email format error`
  String get emailFormatError {
    return Intl.message(
      'Email format error',
      name: 'emailFormatError',
      desc: '',
      args: [],
    );
  }

  /// `Advanced settings`
  String get advancedSettings {
    return Intl.message(
      'Advanced settings',
      name: 'advancedSettings',
      desc: '',
      args: [],
    );
  }

  /// `Server type`
  String get serverType {
    return Intl.message('Server type', name: 'serverType', desc: '', args: []);
  }

  /// `server`
  String get server {
    return Intl.message('server', name: 'server', desc: '', args: []);
  }

  /// `port`
  String get port {
    return Intl.message('port', name: 'port', desc: '', args: []);
  }

  /// `Enable proxy`
  String get enableProxy {
    return Intl.message(
      'Enable proxy',
      name: 'enableProxy',
      desc: '',
      args: [],
    );
  }

  /// `Proxy address`
  String get proxyAddress {
    return Intl.message(
      'Proxy address',
      name: 'proxyAddress',
      desc: '',
      args: [],
    );
  }

  /// `Secure connection`
  String get secureConnection {
    return Intl.message(
      'Secure connection',
      name: 'secureConnection',
      desc: '',
      args: [],
    );
  }

  /// `address or ip format is incorrect`
  String get serverAddressFormatErr {
    return Intl.message(
      'address or ip format is incorrect',
      name: 'serverAddressFormatErr',
      desc: '',
      args: [],
    );
  }

  /// `port format is incorrect`
  String get portFormatErr {
    return Intl.message(
      'port format is incorrect',
      name: 'portFormatErr',
      desc: '',
      args: [],
    );
  }

  /// `Please enter email`
  String get pleaseEnterEmail {
    return Intl.message(
      'Please enter email',
      name: 'pleaseEnterEmail',
      desc: '',
      args: [],
    );
  }

  /// `Please enter password or secret code`
  String get pleaseEnterPassword {
    return Intl.message(
      'Please enter password or secret code',
      name: 'pleaseEnterPassword',
      desc: '',
      args: [],
    );
  }

  /// `Email`
  String get email {
    return Intl.message('Email', name: 'email', desc: '', args: []);
  }

  /// `password`
  String get password {
    return Intl.message('password', name: 'password', desc: '', args: []);
  }

  /// `Email invalid`
  String get emailInvalid {
    return Intl.message(
      'Email invalid',
      name: 'emailInvalid',
      desc: '',
      args: [],
    );
  }

  /// `pass must not empty`
  String get passwordMustNotEmpty {
    return Intl.message(
      'pass must not empty',
      name: 'passwordMustNotEmpty',
      desc: '',
      args: [],
    );
  }

  /// `Login message missing`
  String get loginMessageMissing {
    return Intl.message(
      'Login message missing',
      name: 'loginMessageMissing',
      desc: '',
      args: [],
    );
  }

  /// `Advanced settings missing,Please check the Settings`
  String get advancedSettingsMissing {
    return Intl.message(
      'Advanced settings missing,Please check the Settings',
      name: 'advancedSettingsMissing',
      desc: '',
      args: [],
    );
  }

  /// `Close`
  String get close {
    return Intl.message('Close', name: 'close', desc: '', args: []);
  }

  /// `The login mailbox is different from the selected mailbox_ Continue to log in or modify the configuration?`
  String get emailSettingsNotMatch {
    return Intl.message(
      'The login mailbox is different from the selected mailbox_ Continue to log in or modify the configuration?',
      name: 'emailSettingsNotMatch',
      desc: '',
      args: [],
    );
  }

  /// `Continue`
  String get continuE {
    return Intl.message('Continue', name: 'continuE', desc: '', args: []);
  }

  /// `Edit settings`
  String get editSettings {
    return Intl.message(
      'Edit settings',
      name: 'editSettings',
      desc: '',
      args: [],
    );
  }

  /// `Login failed, please check the email and password and make sure the mail server Settings are correct`
  String get loginFailedMsg {
    return Intl.message(
      'Login failed, please check the email and password and make sure the mail server Settings are correct',
      name: 'loginFailedMsg',
      desc: '',
      args: [],
    );
  }

  /// `Authentication failed, please check the email and password`
  String get authenticationFailedMsg {
    return Intl.message(
      'Authentication failed, please check the email and password',
      name: 'authenticationFailedMsg',
      desc: '',
      args: [],
    );
  }

  /// `Unknown error`
  String get unknownError {
    return Intl.message(
      'Unknown error',
      name: 'unknownError',
      desc: '',
      args: [],
    );
  }

  /// `Login failed`
  String get loginFailed {
    return Intl.message(
      'Login failed',
      name: 'loginFailed',
      desc: '',
      args: [],
    );
  }

  /// `No supported email server`
  String get noSupportedEmailServer {
    return Intl.message(
      'No supported email server',
      name: 'noSupportedEmailServer',
      desc: '',
      args: [],
    );
  }

  /// `System status, or system help reply`
  String get smtpErrCode211 {
    return Intl.message(
      'System status, or system help reply',
      name: 'smtpErrCode211',
      desc: '',
      args: [],
    );
  }

  /// `Help message (A response to the HELP command)`
  String get smtpErrCode214 {
    return Intl.message(
      'Help message (A response to the HELP command)',
      name: 'smtpErrCode214',
      desc: '',
      args: [],
    );
  }

  /// `Service ready`
  String get smtpErrCode220 {
    return Intl.message(
      'Service ready',
      name: 'smtpErrCode220',
      desc: '',
      args: [],
    );
  }

  /// `Service closing transmission channel`
  String get smtpErrCode221 {
    return Intl.message(
      'Service closing transmission channel',
      name: 'smtpErrCode221',
      desc: '',
      args: [],
    );
  }

  /// `Goodbye`
  String get smtpErrCode221_2_0_0 {
    return Intl.message(
      'Goodbye',
      name: 'smtpErrCode221_2_0_0',
      desc: '',
      args: [],
    );
  }

  /// `Authentication succeeded`
  String get smtpErrCode235_2_7_0 {
    return Intl.message(
      'Authentication succeeded',
      name: 'smtpErrCode235_2_7_0',
      desc: '',
      args: [],
    );
  }

  /// `QUIT`
  String get smtpErrCode240 {
    return Intl.message('QUIT', name: 'smtpErrCode240', desc: '', args: []);
  }

  /// `Requested mail action okay, completed`
  String get smtpErrCode250 {
    return Intl.message(
      'Requested mail action okay, completed',
      name: 'smtpErrCode250',
      desc: '',
      args: [],
    );
  }

  /// `User not local; will forward`
  String get smtpErrCode251 {
    return Intl.message(
      'User not local; will forward',
      name: 'smtpErrCode251',
      desc: '',
      args: [],
    );
  }

  /// `Cannot verify the user, but it will try to deliver the message anyway`
  String get smtpErrCode252 {
    return Intl.message(
      'Cannot verify the user, but it will try to deliver the message anyway',
      name: 'smtpErrCode252',
      desc: '',
      args: [],
    );
  }

  /// `Server challenge - the text part contains the Base64-encoded challenge`
  String get smtpErrCode334 {
    return Intl.message(
      'Server challenge - the text part contains the Base64-encoded challenge',
      name: 'smtpErrCode334',
      desc: '',
      args: [],
    );
  }

  /// `Start mail input`
  String get smtpErrCode354 {
    return Intl.message(
      'Start mail input',
      name: 'smtpErrCode354',
      desc: '',
      args: [],
    );
  }

  /// `Service not available, closing transmission channel`
  String get smtpErrCode421 {
    return Intl.message(
      'Service not available, closing transmission channel',
      name: 'smtpErrCode421',
      desc: '',
      args: [],
    );
  }

  /// `A password transition is needed`
  String get smtpErrCode432_4_7_12 {
    return Intl.message(
      'A password transition is needed',
      name: 'smtpErrCode432_4_7_12',
      desc: '',
      args: [],
    );
  }

  /// `Requested mail action not taken: mailbox unavailable`
  String get smtpErrCode450 {
    return Intl.message(
      'Requested mail action not taken: mailbox unavailable',
      name: 'smtpErrCode450',
      desc: '',
      args: [],
    );
  }

  /// `Requested action aborted: local error in processing`
  String get smtpErrCode451 {
    return Intl.message(
      'Requested action aborted: local error in processing',
      name: 'smtpErrCode451',
      desc: '',
      args: [],
    );
  }

  /// `IMAP server unavailable`
  String get smtpErrCode451_4_4_1 {
    return Intl.message(
      'IMAP server unavailable',
      name: 'smtpErrCode451_4_4_1',
      desc: '',
      args: [],
    );
  }

  /// `Requested action not taken: insufficient system storage`
  String get smtpErrCode452 {
    return Intl.message(
      'Requested action not taken: insufficient system storage',
      name: 'smtpErrCode452',
      desc: '',
      args: [],
    );
  }

  /// `Temporary authentication failure`
  String get smtpErrCode454_4_7_0 {
    return Intl.message(
      'Temporary authentication failure',
      name: 'smtpErrCode454_4_7_0',
      desc: '',
      args: [],
    );
  }

  /// `Server unable to accommodate parameters`
  String get smtpErrCode455 {
    return Intl.message(
      'Server unable to accommodate parameters',
      name: 'smtpErrCode455',
      desc: '',
      args: [],
    );
  }

  /// `Syntax error, command unrecognized`
  String get smtpErrCode500 {
    return Intl.message(
      'Syntax error, command unrecognized',
      name: 'smtpErrCode500',
      desc: '',
      args: [],
    );
  }

  /// `Authentication Exchange line is too long`
  String get smtpErrCode500_5_5_6 {
    return Intl.message(
      'Authentication Exchange line is too long',
      name: 'smtpErrCode500_5_5_6',
      desc: '',
      args: [],
    );
  }

  /// `Syntax error in parameters or arguments`
  String get smtpErrCode501 {
    return Intl.message(
      'Syntax error in parameters or arguments',
      name: 'smtpErrCode501',
      desc: '',
      args: [],
    );
  }

  /// `Cannot Base64-decode Client responses`
  String get smtpErrCode501_5_5_2 {
    return Intl.message(
      'Cannot Base64-decode Client responses',
      name: 'smtpErrCode501_5_5_2',
      desc: '',
      args: [],
    );
  }

  /// `Client initiated Authentication Exchange`
  String get smtpErrCode501_5_7_0 {
    return Intl.message(
      'Client initiated Authentication Exchange',
      name: 'smtpErrCode501_5_7_0',
      desc: '',
      args: [],
    );
  }

  /// `Command not implemented`
  String get smtpErrCode502 {
    return Intl.message(
      'Command not implemented',
      name: 'smtpErrCode502',
      desc: '',
      args: [],
    );
  }

  /// `Bad sequence of commands`
  String get smtpErrCode503 {
    return Intl.message(
      'Bad sequence of commands',
      name: 'smtpErrCode503',
      desc: '',
      args: [],
    );
  }

  /// `Command parameter is not implemented`
  String get smtpErrCode504 {
    return Intl.message(
      'Command parameter is not implemented',
      name: 'smtpErrCode504',
      desc: '',
      args: [],
    );
  }

  /// `Unrecognized authentication type`
  String get smtpErrCode504_5_5_4 {
    return Intl.message(
      'Unrecognized authentication type',
      name: 'smtpErrCode504_5_5_4',
      desc: '',
      args: [],
    );
  }

  /// `Server does not accept mail`
  String get smtpErrCode521 {
    return Intl.message(
      'Server does not accept mail',
      name: 'smtpErrCode521',
      desc: '',
      args: [],
    );
  }

  /// `Encryption Needed`
  String get smtpErrCode523 {
    return Intl.message(
      'Encryption Needed',
      name: 'smtpErrCode523',
      desc: '',
      args: [],
    );
  }

  /// `Authentication required`
  String get smtpErrCode530_5_7_0 {
    return Intl.message(
      'Authentication required',
      name: 'smtpErrCode530_5_7_0',
      desc: '',
      args: [],
    );
  }

  /// `Authentication mechanism is too weak`
  String get smtpErrCode534_5_7_9 {
    return Intl.message(
      'Authentication mechanism is too weak',
      name: 'smtpErrCode534_5_7_9',
      desc: '',
      args: [],
    );
  }

  /// `Authentication credentials invalid`
  String get smtpErrCode535_5_7_8 {
    return Intl.message(
      'Authentication credentials invalid',
      name: 'smtpErrCode535_5_7_8',
      desc: '',
      args: [],
    );
  }

  /// `Encryption required for requested authentication mechanism`
  String get smtpErrCode538_5_7_11 {
    return Intl.message(
      'Encryption required for requested authentication mechanism',
      name: 'smtpErrCode538_5_7_11',
      desc: '',
      args: [],
    );
  }

  /// `Requested action not taken: mailbox unavailable`
  String get smtpErrCode550 {
    return Intl.message(
      'Requested action not taken: mailbox unavailable',
      name: 'smtpErrCode550',
      desc: '',
      args: [],
    );
  }

  /// `User not local; please try <forward-path>`
  String get smtpErrCode551 {
    return Intl.message(
      'User not local; please try <forward-path>',
      name: 'smtpErrCode551',
      desc: '',
      args: [],
    );
  }

  /// `Requested mail action aborted: exceeded storage allocation`
  String get smtpErrCode552 {
    return Intl.message(
      'Requested mail action aborted: exceeded storage allocation',
      name: 'smtpErrCode552',
      desc: '',
      args: [],
    );
  }

  /// `Requested action not taken: mailbox name not allowed`
  String get smtpErrCode553 {
    return Intl.message(
      'Requested action not taken: mailbox name not allowed',
      name: 'smtpErrCode553',
      desc: '',
      args: [],
    );
  }

  /// `Transaction has failed (Or, 'No SMTP service here')`
  String get smtpErrCode554 {
    return Intl.message(
      'Transaction has failed (Or, \'No SMTP service here\')',
      name: 'smtpErrCode554',
      desc: '',
      args: [],
    );
  }

  /// `Message too big for system`
  String get smtpErrCode554_5_3_4 {
    return Intl.message(
      'Message too big for system',
      name: 'smtpErrCode554_5_3_4',
      desc: '',
      args: [],
    );
  }

  /// `Domain does not accept mail`
  String get smtpErrCode556 {
    return Intl.message(
      'Domain does not accept mail',
      name: 'smtpErrCode556',
      desc: '',
      args: [],
    );
  }
}

class AppLocalizationDelegate extends LocalizationsDelegate<S> {
  const AppLocalizationDelegate();

  List<Locale> get supportedLocales {
    return const <Locale>[
      Locale.fromSubtags(languageCode: 'en'),
      Locale.fromSubtags(languageCode: 'zh', countryCode: 'CN'),
    ];
  }

  @override
  bool isSupported(Locale locale) => _isSupported(locale);
  @override
  Future<S> load(Locale locale) => S.load(locale);
  @override
  bool shouldReload(AppLocalizationDelegate old) => false;

  bool _isSupported(Locale locale) {
    for (var supportedLocale in supportedLocales) {
      if (supportedLocale.languageCode == locale.languageCode) {
        return true;
      }
    }
    return false;
  }
}
