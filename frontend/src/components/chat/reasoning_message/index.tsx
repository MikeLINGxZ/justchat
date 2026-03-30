import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import styles from './index.module.scss';
import MarkdownRenderer from "@/components/markdown_renderer";

interface ReasoningContentProps {
  content: string;
  className?: string;
  isStreaming?: boolean; // 是否正在流式输入
  title?: string;
}

const ReasoningContent: React.FC<ReasoningContentProps> = ({
  content,
  className,
  isStreaming = false,
  title,
}) => {
  const { t } = useTranslation();
  // 初始状态：如果正在生成且有内容，则展开
  const [isExpanded, setIsExpanded] = useState(() => isStreaming && content.trim().length > 0);
  const prevIsStreamingRef = React.useRef(isStreaming);

  // 当开始流式输入思考过程时自动展开，生成完成后自动折叠
  useEffect(() => {
    const wasStreaming = prevIsStreamingRef.current;
    const isNowStreaming = isStreaming;
    
    let timeoutId: NodeJS.Timeout | null = null;
    
    // 如果正在生成中且有内容，自动展开
    if (isNowStreaming && content.trim()) {
      if (!isExpanded) {
        setIsExpanded(true);
      }
    }
    // 如果从生成中变为完成，自动折叠（无条件）
    else if (wasStreaming && !isNowStreaming && content.trim()) {
      // 延迟一下折叠，确保用户能看到最后的更新
      timeoutId = setTimeout(() => {
        setIsExpanded(false);
      }, 300);
    }
    
    // 更新引用
    prevIsStreamingRef.current = isStreaming;
    
    // 清理函数
    return () => {
      if (timeoutId) {
        clearTimeout(timeoutId);
      }
    };
  }, [isStreaming, content, isExpanded]);

  if (!content.trim()) {
    return null;
  }

  const toggleExpanded = () => {
    setIsExpanded(!isExpanded);
  };

  return (
    <div className={`${styles.reasoningContainer} ${className || ''}`}>
      <div className={styles.reasoningHeader} onClick={toggleExpanded}>
        <div className={styles.reasoningTitle}>
          <span className={styles.reasoningIcon}>🧠</span>
          <span>{title || t('chat.messageAction.reasoningTitle')}</span>
          {isStreaming && (
            <span className={styles.streamingIndicator}>
              <span className={styles.streamingDots}>
                <span></span>
                <span></span>
                <span></span>
              </span>
            </span>
          )}
        </div>
        <div className={`${styles.toggleIcon} ${isExpanded ? styles.expanded : ''}`}>
          <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
            <path d="M4 6l4 4 4-4" stroke="currentColor" strokeWidth="2" fill="none" strokeLinecap="round" strokeLinejoin="round"/>
          </svg>
        </div>
      </div>
      
      {isExpanded && (
        <div className={styles.reasoningContent}>
          <MarkdownRenderer content={content} variant="reasoning" />
        </div>
      )}
    </div>
  );
};

export default ReasoningContent;
