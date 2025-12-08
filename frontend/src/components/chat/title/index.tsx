import React, { useState, useRef, useEffect } from 'react';
import { message } from 'antd';
import styles from '@/components/chat/title/index.module.scss';
import {useIsMobile} from "@/hooks/useViewportHeight.ts";

interface ChatTitleProps {
  // 聊天标题
  title: string;
  // 聊天UUID（用于API调用）
  uuid?: string;
  // 聊天标题变更事件
  onTitleChange: (newTitle: string) => void;
  // 侧边栏是否收起
  isSidebarCollapsed?: boolean;
  // 切换侧边栏事件
  onToggleSidebar?: () => void;
}

const ChatTitle: React.FC<ChatTitleProps> = ({
  title,
  uuid,
  onTitleChange,
  isSidebarCollapsed = false,
  onToggleSidebar,
}) => {
  const [isEditing, setIsEditing] = useState(false);
  const [editValue, setEditValue] = useState(title);
  const [inputWidth, setInputWidth] = useState(120); // 默认宽度
  const inputRef = useRef<HTMLInputElement>(null);
  const measureRef = useRef<HTMLSpanElement>(null);

  // 判断是否允许编辑标题：只有当 chatUuid 不为空且不是空字符串时才允许编辑
  const canEditTitle = Boolean(uuid && uuid.trim() !== '');
  const isMobile =  useIsMobile();

  // 计算输入框宽度
  const calculateInputWidth = (text: string) => {
    if (measureRef.current) {
      measureRef.current.textContent = text || '新建对话';
      const textWidth = measureRef.current.offsetWidth;
      const maxWidth = 50 * 16; // 假设每个字符约16px，20个字符的最大宽度
      const minWidth = 120; // 最小宽度
      return Math.min(Math.max(textWidth + 20, minWidth), maxWidth); // 加20px的padding
    }
    return 120;
  };

  // 开始编辑
  const handleStartEdit = () => {
    // 检查是否允许编辑
    if (!canEditTitle) {
      message.info('请先保存对话后再编辑标题');
      return;
    }

    const initialValue = title || '';
    setEditValue(initialValue);
    setInputWidth(calculateInputWidth(initialValue));
    setIsEditing(true);
  };

  // 确认编辑 (模拟实现)
  const handleConfirm = async () => {
    const trimmedValue = editValue.trim();
    if (trimmedValue && trimmedValue !== title) {
      // 如果有chatUuid，模拟保存标题
      if (uuid) {
        try {
          // 模拟保存延迟
          await new Promise(resolve => setTimeout(resolve, 200));
          onTitleChange(trimmedValue);
          message.success('标题保存成功');
        } catch (error) {
          console.error('保存标题失败:', error);
          message.error('保存标题失败，请重试');
          return; // 保存失败时不关闭编辑状态
        }
      } else {
        // 新对话，只更新本地状态
        onTitleChange(trimmedValue);
      }
    }
    setIsEditing(false);
  };

  // 取消编辑
  const handleCancel = () => {
    setEditValue(title || '');
    setIsEditing(false);
  };

  // 处理输入变化
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value;
    // 限制最多20个字符
    if (newValue.length <= 50) {
      setEditValue(newValue);
      setInputWidth(calculateInputWidth(newValue));
    }
  };

  // 处理键盘事件
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleConfirm();
    } else if (e.key === 'Escape') {
      handleCancel();
    }
  };

  // 编辑状态时自动聚焦
  useEffect(() => {
    if (isEditing && inputRef.current) {
      inputRef.current.focus();
      inputRef.current.select();
    }
  }, [isEditing]);

  return (
    <div className={styles.chatTitlePage}>
      <div className={styles.titleHeader}>
        {/* 移动端菜单按钮 */}
        {isMobile && isSidebarCollapsed && (
          <button
            className={styles.mobileMenuButton}
            onClick={onToggleSidebar}
            title="打开侧边栏"
          >
            <svg
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <line x1="3" y1="6" x2="21" y2="6"></line>
              <line x1="3" y1="12" x2="21" y2="12"></line>
              <line x1="3" y1="18" x2="21" y2="18"></line>
            </svg>
          </button>
        )}

        {isEditing ? (
          <div className={styles.editContainer}>
            <input
              ref={inputRef}
              type="text"
              value={editValue}
              onChange={handleInputChange}
              onKeyDown={handleKeyDown}
              className={styles.titleInput}
              placeholder="新建对话"
              maxLength={50}
              style={{ width: `${inputWidth}px` }}
            />
            <div className={styles.editActions}>
              <button
                className={`${styles.actionButton} ${styles.confirmButton}`}
                onClick={handleConfirm}
                title="确认"
              >
                <svg
                  width="16"
                  height="16"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <polyline points="20,6 9,17 4,12"></polyline>
                </svg>
              </button>
              <button
                className={`${styles.actionButton} ${styles.cancelButton}`}
                onClick={handleCancel}
                title="取消"
              >
                <svg
                  width="16"
                  height="16"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                >
                  <line x1="18" y1="6" x2="6" y2="18"></line>
                  <line x1="6" y1="6" x2="18" y2="18"></line>
                </svg>
              </button>
            </div>
          </div>
        ) : (
          <div className={styles.titleContainer}>
            <div
              className={`${styles.title} ${!title ? styles.defaultTitle : ''}`}
            >
              {title || '新建对话'}
            </div>
            <button
              className={`${styles.editButton} ${!canEditTitle ? styles.editButtonDisabled : ''}`}
              onClick={handleStartEdit}
              title={canEditTitle ? '编辑标题' : '请先保存对话后编辑标题'}
              disabled={!canEditTitle}
            >
              <svg
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="currentColor"
              >
                <path d="M3 17.25V21h3.75L17.81 9.94l-3.75-3.75L3 17.25zM20.71 7.04c.39-.39.39-1.02 0-1.41l-2.34-2.34c-.39-.39-1.02-.39-1.41 0l-1.83 1.83 3.75 3.75 1.83-1.83z" />
              </svg>
            </button>
          </div>
        )}
      </div>

      {/* 用于测量文本宽度的隐藏元素 */}
      <span
        ref={measureRef}
        className={styles.measureSpan}
        aria-hidden="true"
      />
    </div>
  );
};

export default ChatTitle;
