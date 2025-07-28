import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:lemon_tea/utils/llm/models/message.dart';
import 'package:lemon_tea/utils/file_service.dart';
import 'package:lemon_tea/utils/font_size_utils.dart';
import 'package:lemon_tea/utils/style.dart';

class FilePreview extends ConsumerWidget {
  final List<FileContent> files;
  final bool isUserMessage;

  const FilePreview({
    super.key,
    required this.files,
    this.isUserMessage = false,
  });

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    if (files.isEmpty) return const SizedBox.shrink();

    return Container(
      margin: const EdgeInsets.only(bottom: 8.0),
      child: Wrap(
        spacing: 8.0,
        runSpacing: 8.0,
        children: files.map((file) => _buildFileItem(context, ref, file)).toList(),
      ),
    );
  }

  Widget _buildFileItem(BuildContext context, WidgetRef ref, FileContent file) {
    return Container(
      constraints: const BoxConstraints(maxWidth: 300),
      decoration: BoxDecoration(
        color: isUserMessage 
            ? Theme.of(context).colorScheme.surfaceVariant.withOpacity(0.3)
            : Theme.of(context).colorScheme.surfaceVariant.withOpacity(0.1),
        borderRadius: BorderRadius.circular(8.0),
        border: Border.all(
          color: Style.primaryBorder(context),
          width: 0.5,
        ),
      ),
      child: FileService.isImage(file) 
          ? _buildImagePreview(context, ref, file)
          : _buildDocumentPreview(context, ref, file),
    );
  }

  Widget _buildImagePreview(BuildContext context, WidgetRef ref, FileContent file) {
    final imageData = FileService.getImagePreviewData(file);
    
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        // 图片预览区域
        Container(
          height: 200,
          width: double.infinity,
          decoration: BoxDecoration(
            borderRadius: const BorderRadius.vertical(top: Radius.circular(8.0)),
            color: Theme.of(context).colorScheme.surface,
          ),
          child: ClipRRect(
            borderRadius: const BorderRadius.vertical(top: Radius.circular(8.0)),
            child: imageData != null
                ? Image.memory(
                    imageData,
                    fit: BoxFit.cover,
                    errorBuilder: (context, error, stackTrace) {
                      return _buildErrorWidget(context);
                    },
                  )
                : _buildErrorWidget(context),
          ),
        ),
        
        // 文件信息
        Padding(
          padding: const EdgeInsets.all(8.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                file.name,
                style: TextStyle(
                  fontSize: FontSizeUtils.getSmallSize(ref),
                  fontWeight: FontWeight.w500,
                ),
                maxLines: 2,
                overflow: TextOverflow.ellipsis,
              ),
              const SizedBox(height: 4),
              Text(
                '${FileService.formatFileSize(file.size)} • ${file.mimeType}',
                style: TextStyle(
                  fontSize: FontSizeUtils.getSmallSize(ref) - 2,
                  color: Theme.of(context).colorScheme.onSurfaceVariant,
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildDocumentPreview(BuildContext context, WidgetRef ref, FileContent file) {
    return Padding(
      padding: const EdgeInsets.all(12.0),
      child: Row(
        children: [
          // 文件图标
          Container(
            width: 48,
            height: 48,
            decoration: BoxDecoration(
              color: Theme.of(context).colorScheme.primary.withOpacity(0.1),
              borderRadius: BorderRadius.circular(8.0),
            ),
            child: Icon(
              FileService.getFileTypeIcon(file.type),
              size: 24,
              color: Theme.of(context).colorScheme.primary,
            ),
          ),
          
          const SizedBox(width: 12),
          
          // 文件信息
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  file.name,
                  style: TextStyle(
                    fontSize: FontSizeUtils.getSmallSize(ref),
                    fontWeight: FontWeight.w500,
                  ),
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
                const SizedBox(height: 4),
                Text(
                  '${FileService.formatFileSize(file.size)} • ${file.mimeType}',
                  style: TextStyle(
                    fontSize: FontSizeUtils.getSmallSize(ref) - 2,
                    color: Theme.of(context).colorScheme.onSurfaceVariant,
                  ),
                ),
                if (file.description != null && file.description!.isNotEmpty) ...[
                  const SizedBox(height: 4),
                  Text(
                    file.description!,
                    style: TextStyle(
                      fontSize: FontSizeUtils.getSmallSize(ref) - 2,
                      color: Theme.of(context).colorScheme.onSurfaceVariant,
                      fontStyle: FontStyle.italic,
                    ),
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                  ),
                ],
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildErrorWidget(BuildContext context) {
    return Container(
      width: double.infinity,
      height: double.infinity,
      color: Theme.of(context).colorScheme.errorContainer.withOpacity(0.1),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            Icons.broken_image,
            size: 48,
            color: Theme.of(context).colorScheme.error,
          ),
          const SizedBox(height: 8),
          Text(
            '图片加载失败',
            style: TextStyle(
              color: Theme.of(context).colorScheme.error,
              fontSize: 12,
            ),
          ),
        ],
      ),
    );
  }
} 