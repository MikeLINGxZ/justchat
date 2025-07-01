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

  /// `Task Management`
  String get taskManagement {
    return Intl.message(
      'Task Management',
      name: 'taskManagement',
      desc: '',
      args: [],
    );
  }

  /// `No Tasks`
  String get noTasks {
    return Intl.message('No Tasks', name: 'noTasks', desc: '', args: []);
  }

  /// `Manage your tasks here`
  String get manageTasksHere {
    return Intl.message(
      'Manage your tasks here',
      name: 'manageTasksHere',
      desc: '',
      args: [],
    );
  }

  /// `Confirm Delete`
  String get confirmDelete {
    return Intl.message(
      'Confirm Delete',
      name: 'confirmDelete',
      desc: '',
      args: [],
    );
  }

  /// `Are you sure you want to delete this conversation? This cannot be undone.`
  String get confirmDeleteConversation {
    return Intl.message(
      'Are you sure you want to delete this conversation? This cannot be undone.',
      name: 'confirmDeleteConversation',
      desc: '',
      args: [],
    );
  }

  /// `Cancel`
  String get cancel {
    return Intl.message('Cancel', name: 'cancel', desc: '', args: []);
  }

  /// `Delete`
  String get delete {
    return Intl.message('Delete', name: 'delete', desc: '', args: []);
  }

  /// `Conversation History`
  String get conversationHistory {
    return Intl.message(
      'Conversation History',
      name: 'conversationHistory',
      desc: '',
      args: [],
    );
  }

  /// `New Conversation`
  String get newConversation {
    return Intl.message(
      'New Conversation',
      name: 'newConversation',
      desc: '',
      args: [],
    );
  }

  /// `Search conversations...`
  String get searchConversations {
    return Intl.message(
      'Search conversations...',
      name: 'searchConversations',
      desc: '',
      args: [],
    );
  }

  /// `AI Assistant`
  String get aiAssistant {
    return Intl.message(
      'AI Assistant',
      name: 'aiAssistant',
      desc: '',
      args: [],
    );
  }

  /// `Welcome to Markdown`
  String get welcomeToMarkdown {
    return Intl.message(
      'Welcome to Markdown',
      name: 'welcomeToMarkdown',
      desc: '',
      args: [],
    );
  }

  /// `This is a simple Markdown example document showing common syntax:`
  String get markdownExampleDoc {
    return Intl.message(
      'This is a simple Markdown example document showing common syntax:',
      name: 'markdownExampleDoc',
      desc: '',
      args: [],
    );
  }

  /// `Heading Levels`
  String get headingLevels {
    return Intl.message(
      'Heading Levels',
      name: 'headingLevels',
      desc: '',
      args: [],
    );
  }

  /// `Level 2 heading ('##') to level 6 heading ('######')`
  String get headingLevelsDesc {
    return Intl.message(
      'Level 2 heading (`##`) to level 6 heading (`######`)',
      name: 'headingLevelsDesc',
      desc: '',
      args: [],
    );
  }

  /// `Text Styles`
  String get textStyles {
    return Intl.message('Text Styles', name: 'textStyles', desc: '', args: []);
  }

  /// `bold text`
  String get boldTextExample {
    return Intl.message(
      'bold text',
      name: 'boldTextExample',
      desc: '',
      args: [],
    );
  }

  /// `Chinese`
  String get chinese {
    return Intl.message('Chinese', name: 'chinese', desc: '', args: []);
  }

  /// `General`
  String get general {
    return Intl.message('General', name: 'general', desc: '', args: []);
  }

  /// `Model`
  String get model {
    return Intl.message('Model', name: 'model', desc: '', args: []);
  }

  /// `Data`
  String get data {
    return Intl.message('Data', name: 'data', desc: '', args: []);
  }

  /// `About`
  String get about {
    return Intl.message('About', name: 'about', desc: '', args: []);
  }

  /// `Settings`
  String get settings {
    return Intl.message('Settings', name: 'settings', desc: '', args: []);
  }

  /// `General Settings`
  String get generalSettings {
    return Intl.message(
      'General Settings',
      name: 'generalSettings',
      desc: '',
      args: [],
    );
  }

  /// `Theme`
  String get theme {
    return Intl.message('Theme', name: 'theme', desc: '', args: []);
  }

  /// `Theme Mode`
  String get themeMode {
    return Intl.message('Theme Mode', name: 'themeMode', desc: '', args: []);
  }

  /// `Light Mode`
  String get lightMode {
    return Intl.message('Light Mode', name: 'lightMode', desc: '', args: []);
  }

  /// `Dark Mode`
  String get darkMode {
    return Intl.message('Dark Mode', name: 'darkMode', desc: '', args: []);
  }

  /// `System Mode`
  String get systemMode {
    return Intl.message('System Mode', name: 'systemMode', desc: '', args: []);
  }

  /// `Font Size`
  String get fontSize {
    return Intl.message('Font Size', name: 'fontSize', desc: '', args: []);
  }

  /// `Interface Font`
  String get interfaceFont {
    return Intl.message(
      'Interface Font',
      name: 'interfaceFont',
      desc: '',
      args: [],
    );
  }

  /// `Language`
  String get language {
    return Intl.message('Language', name: 'language', desc: '', args: []);
  }

  /// `Interface Language`
  String get interfaceLanguage {
    return Intl.message(
      'Interface Language',
      name: 'interfaceLanguage',
      desc: '',
      args: [],
    );
  }

  /// `Model Settings`
  String get modelSettings {
    return Intl.message(
      'Model Settings',
      name: 'modelSettings',
      desc: '',
      args: [],
    );
  }

  /// `Add Model`
  String get addModel {
    return Intl.message('Add Model', name: 'addModel', desc: '', args: []);
  }

  /// `Add New Model`
  String get addNewModel {
    return Intl.message(
      'Add New Model',
      name: 'addNewModel',
      desc: '',
      args: [],
    );
  }

  /// `Model List`
  String get modelList {
    return Intl.message('Model List', name: 'modelList', desc: '', args: []);
  }

  /// `Data Settings`
  String get dataSettings {
    return Intl.message(
      'Data Settings',
      name: 'dataSettings',
      desc: '',
      args: [],
    );
  }

  /// `Data Storage`
  String get dataStorage {
    return Intl.message(
      'Data Storage',
      name: 'dataStorage',
      desc: '',
      args: [],
    );
  }

  /// `Auto Save Data`
  String get autoSaveData {
    return Intl.message(
      'Auto Save Data',
      name: 'autoSaveData',
      desc: '',
      args: [],
    );
  }

  /// `Conversation`
  String get conversation {
    return Intl.message(
      'Conversation',
      name: 'conversation',
      desc: '',
      args: [],
    );
  }

  /// `No Conversation History`
  String get noConversationHistory {
    return Intl.message(
      'No Conversation History',
      name: 'noConversationHistory',
      desc: '',
      args: [],
    );
  }

  /// `{count} messages`
  String messagesCount(Object count) {
    return Intl.message(
      '$count messages',
      name: 'messagesCount',
      desc: '',
      args: [count],
    );
  }

  /// `Delete Conversation`
  String get deleteConversation {
    return Intl.message(
      'Delete Conversation',
      name: 'deleteConversation',
      desc: '',
      args: [],
    );
  }

  /// `Type a message...`
  String get inputMessage {
    return Intl.message(
      'Type a message...',
      name: 'inputMessage',
      desc: '',
      args: [],
    );
  }

  /// `Upload Image`
  String get uploadImage {
    return Intl.message(
      'Upload Image',
      name: 'uploadImage',
      desc: '',
      args: [],
    );
  }

  /// `Upload File`
  String get uploadFile {
    return Intl.message('Upload File', name: 'uploadFile', desc: '', args: []);
  }

  /// `Extra Small`
  String get fontSizeExtraSmall {
    return Intl.message(
      'Extra Small',
      name: 'fontSizeExtraSmall',
      desc: '',
      args: [],
    );
  }

  /// `Small`
  String get fontSizeSmall {
    return Intl.message('Small', name: 'fontSizeSmall', desc: '', args: []);
  }

  /// `Medium`
  String get fontSizeMedium {
    return Intl.message('Medium', name: 'fontSizeMedium', desc: '', args: []);
  }

  /// `Large`
  String get fontSizeLarge {
    return Intl.message('Large', name: 'fontSizeLarge', desc: '', args: []);
  }

  /// `Extra Large`
  String get fontSizeExtraLarge {
    return Intl.message(
      'Extra Large',
      name: 'fontSizeExtraLarge',
      desc: '',
      args: [],
    );
  }

  /// `Plugin Management`
  String get pluginManagement {
    return Intl.message(
      'Plugin Management',
      name: 'pluginManagement',
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
