import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/storage/llm_storage.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:lemon_tea/generated/l10n.dart';
import 'package:lemon_tea/models/llm_provider.dart';

class ModelSettings extends ConsumerStatefulWidget {
  const ModelSettings({super.key});

  @override
  ConsumerState<ModelSettings> createState() => _ModelSettingsState();
}

class _ModelSettingsState extends ConsumerState<ModelSettings>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.fromLTRB(24, 24, 24, 0),
          child: Text(
            S.of(context).modelSettings,
            style: TextStyle(
              fontSize: FontSizeUtils.getHeadingSize(ref),
              fontWeight: FontWeight.bold,
            ),
          ),
        ),
        const SizedBox(height: 24),
        Container(
          margin: const EdgeInsets.symmetric(horizontal: 24),
          decoration: BoxDecoration(
            color: theme.colorScheme.surfaceContainerHighest,
            borderRadius: BorderRadius.circular(8),
          ),
          child: TabBar(
            controller: _tabController,
            labelColor: theme.colorScheme.primary,
            unselectedLabelColor: theme.colorScheme.onSurface,
            indicatorSize: TabBarIndicatorSize.tab,
            dividerColor: Colors.transparent,
            indicator: BoxDecoration(
              color: theme.colorScheme.surface,
              borderRadius: BorderRadius.circular(8),
              boxShadow: [
                BoxShadow(
                  color: Colors.black.withAlpha(13),
                  blurRadius: 4,
                  offset: const Offset(0, 2),
                ),
              ],
            ),
            splashBorderRadius: BorderRadius.circular(8),
            padding: const EdgeInsets.all(4),
            labelStyle: TextStyle(
              fontSize: FontSizeUtils.getBodySize(ref),
              fontWeight: FontWeight.w600,
            ),
            unselectedLabelStyle: TextStyle(
              fontSize: FontSizeUtils.getBodySize(ref),
              fontWeight: FontWeight.normal,
            ),
            tabs: [
              Tab(
                icon: const Icon(Icons.cloud),
                text: '模型供应商',
                iconMargin: const EdgeInsets.only(bottom: 4),
                height: 64,
              ),
              Tab(
                icon: const Icon(Icons.text_fields),
                text: '提示词',
                iconMargin: const EdgeInsets.only(bottom: 4),
                height: 64,
              ),
            ],
          ),
        ),
        const SizedBox(height: 16),
        Expanded(
          child: TabBarView(
            controller: _tabController,
            children: [
              FutureBuilder<List<LlmProvider>>(
                future: LlmStorage.getAllProviders(),
                builder: (context, snapshot) {
                  if (snapshot.connectionState == ConnectionState.waiting) {
                    return const Center(child: CircularProgressIndicator());
                  }
                  
                  if (snapshot.hasError) {
                    return Center(child: Text('加载失败: ${snapshot.error}'));
                  }
                  
                  final providers = snapshot.data ?? [];
                  
                  return SingleChildScrollView(
                    padding: const EdgeInsets.all(24),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start, 
                      children: [
                        for (final provider in providers)
                          ListTile(
                            title: Text(provider.name),
                            subtitle: Text(provider.baseUrl),
                            trailing: Switch(
                              value: provider.enable,
                              onChanged: (value) {
                                // TODO: 更新提供商启用状态
                              },
                            ),
                          ),
                      ],
                    ),
                  );
                },
              ),
              _buildPromptsTab()
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildPromptsTab() {
    return Center(
      child: Text(
        '提示词',
        style: TextStyle(
          fontSize: FontSizeUtils.getHeadingSize(ref),
          fontWeight: FontWeight.bold,
        ),
      ),
    );
  }
}
