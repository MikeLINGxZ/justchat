import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';

/// 应用设置数据模型
class AppSettings {
  final bool autoSave;
  final String language;
  final int selectedMenuIndex;
  final Map<String, dynamic> modelSettings;
  final Map<String, dynamic> dataSettings;

  AppSettings({
    this.autoSave = true,
    this.language = '中文',
    this.selectedMenuIndex = 0,
    this.modelSettings = const {},
    this.dataSettings = const {},
  });

  AppSettings copyWith({
    bool? autoSave,
    String? language,
    int? selectedMenuIndex,
    Map<String, dynamic>? modelSettings,
    Map<String, dynamic>? dataSettings,
  }) {
    return AppSettings(
      autoSave: autoSave ?? this.autoSave,
      language: language ?? this.language,
      selectedMenuIndex: selectedMenuIndex ?? this.selectedMenuIndex,
      modelSettings: modelSettings ?? this.modelSettings,
      dataSettings: dataSettings ?? this.dataSettings,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'autoSave': autoSave,
      'language': language,
      'selectedMenuIndex': selectedMenuIndex,
      'modelSettings': modelSettings,
      'dataSettings': dataSettings,
    };
  }

  factory AppSettings.fromJson(Map<String, dynamic> json) {
    return AppSettings(
      autoSave: json['autoSave'] ?? true,
      language: json['language'] ?? '中文',
      selectedMenuIndex: json['selectedMenuIndex'] ?? 0,
      modelSettings: Map<String, dynamic>.from(json['modelSettings'] ?? {}),
      dataSettings: Map<String, dynamic>.from(json['dataSettings'] ?? {}),
    );
  }
}

/// 设置管理器Provider
final settingsManagerProvider = StateNotifierProvider<SettingsManager, AppSettings>((ref) {
  return SettingsManager();
});

/// 设置管理器
class SettingsManager extends StateNotifier<AppSettings> {
  static const String _settingsKey = 'app_settings';
  
  SettingsManager() : super(AppSettings()) {
    loadSettings();
  }

  /// 加载保存的设置
  Future<void> loadSettings() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final settingsJson = prefs.getString(_settingsKey);
      
      if (settingsJson != null) {
        final Map<String, dynamic> json = Map<String, dynamic>.from(
          jsonDecode(settingsJson)
        );
        state = AppSettings.fromJson(json);
      }
    } catch (e) {
      debugPrint('Failed to load settings: $e');
      // 使用默认设置
      state = AppSettings();
    }
  }

  /// 保存设置
  Future<void> _saveSettings() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final settingsJson = jsonEncode(state.toJson());
      await prefs.setString(_settingsKey, settingsJson);
    } catch (e) {
      debugPrint('Failed to save settings: $e');
    }
  }

  /// 设置自动保存
  Future<void> setAutoSave(bool autoSave) async {
    state = state.copyWith(autoSave: autoSave);
    await _saveSettings();
  }

  /// 设置语言
  Future<void> setLanguage(String language) async {
    state = state.copyWith(language: language);
    await _saveSettings();
  }

  /// 设置选中的菜单索引
  Future<void> setSelectedMenuIndex(int index) async {
    state = state.copyWith(selectedMenuIndex: index);
    await _saveSettings();
  }

  /// 更新模型设置
  Future<void> updateModelSettings(Map<String, dynamic> modelSettings) async {
    state = state.copyWith(modelSettings: modelSettings);
    await _saveSettings();
  }

  /// 更新数据设置
  Future<void> updateDataSettings(Map<String, dynamic> dataSettings) async {
    state = state.copyWith(dataSettings: dataSettings);
    await _saveSettings();
  }

  /// 重置所有设置
  Future<void> resetAllSettings() async {
    state = AppSettings();
    await _saveSettings();
  }

  /// 清除所有设置数据
  Future<void> clearAllSettings() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      await prefs.remove(_settingsKey);
      state = AppSettings();
    } catch (e) {
      debugPrint('Failed to clear settings: $e');
    }
  }
} 