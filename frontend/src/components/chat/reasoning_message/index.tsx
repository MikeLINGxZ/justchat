import React, { useState, useEffect } from 'react';
import ReactMarkdown from 'react-markdown';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { tomorrow } from 'react-syntax-highlighter/dist/esm/styles/prism';
import remarkGfm from 'remark-gfm';
import styles from './index.module.scss';

interface ReasoningContentProps {
  content: string;
  className?: string;
  isStreaming?: boolean; // 是否正在流式输入
}

const ReasoningContent: React.FC<ReasoningContentProps> = ({ content, className, isStreaming = false }) => {
  // 初始状态：如果正在生成且有内容，则展开
  const [isExpanded, setIsExpanded] = useState(() => isStreaming && content.trim().length > 0);
  const prevIsStreamingRef = React.useRef(isStreaming);

  console.log('ReasoningContent渲染:', { content, isStreaming, contentLength: content?.length });

  // 当开始流式输入思考过程时自动展开，生成完成后自动折叠
  useEffect(() => {
    const wasStreaming = prevIsStreamingRef.current;
    const isNowStreaming = isStreaming;
    
    let timeoutId: NodeJS.Timeout | null = null;
    
    // 如果正在生成中且有内容，自动展开
    if (isNowStreaming && content.trim()) {
      if (!isExpanded) {
        console.log('生成中，自动展开思考过程');
        setIsExpanded(true);
      }
    }
    // 如果从生成中变为完成，自动折叠（无条件）
    else if (wasStreaming && !isNowStreaming && content.trim()) {
      console.log('生成完成，自动折叠思考过程');
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
    console.log('内容为空，不渲染ReasoningContent');
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
          <span>思考过程</span>
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
          <ReactMarkdown
            remarkPlugins={[remarkGfm]}
            components={{
              code(props: any) {
                const { node, inline, className, children, ...rest } = props;
                const match = /language-(\w+)/.exec(className || '');
                const language = match ? match[1] : '';
                
                return !inline && language ? (
                  <SyntaxHighlighter
                    style={tomorrow}
                    language={language}
                    PreTag="div"
                    customStyle={{
                      margin: '8px 0',
                      borderRadius: '6px',
                      fontSize: '13px',
                      background: 'var(--code-background)',
                    } as any}
                    {...rest}
                  >
                    {String(children).replace(/\n$/, '')}
                  </SyntaxHighlighter>
                ) : (
                  <code className={`${className} ${styles.inlineCode}`} {...rest}>
                    {children}
                  </code>
                );
              },
              // 自定义表格样式
              table: ({children}) => (
                <div className={styles.tableWrapper}>
                  <table className={styles.markdownTable}>{children}</table>
                </div>
              ),
              // 自定义链接样式
              a: ({children, href}) => (
                <a 
                  href={href} 
                  target="_blank" 
                  rel="noopener noreferrer" 
                  className={styles.markdownLink}
                >
                  {children}
                </a>
              ),
              // 自定义引用块样式
              blockquote: ({children}) => (
                <blockquote className={styles.markdownBlockquote}>
                  {children}
                </blockquote>
              ),
              // 自定义列表样式
              ul: ({children}) => (
                <ul className={styles.markdownList}>{children}</ul>
              ),
              ol: ({children}) => (
                <ol className={styles.markdownList}>{children}</ol>
              ),
            }}
          >
            {content}
          </ReactMarkdown>
        </div>
      )}
    </div>
  );
};

export default ReasoningContent;
