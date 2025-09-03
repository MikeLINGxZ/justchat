import React, {
  useState,
  useEffect,
  useMemo,
  useRef,
  useCallback,
} from 'react';
import {
  List,
  Typography,
  Input,
  Spin,
  Empty,
  Divider,
  message,
  Dropdown,
  Button,
  Modal,
} from 'antd';
import {
  SearchOutlined,
  MoreOutlined,
  StarOutlined,
  DeleteOutlined,
  ExclamationCircleOutlined,
  EditOutlined,
  CheckOutlined,
  CloseOutlined,
} from '@ant-design/icons';
import type { MenuProps } from 'antd';
import styles from './chats_lists.module.scss';

const { Text } = Typography;
const { Search } = Input;

// 模拟聊天信息接口
interface MockChatInfo {
  chatUuid: string;
  title: string;
  updatedAt: string;
  messagesCount?: number;
}

interface GroupedChats {
  today: MockChatInfo[];
  yesterday: MockChatInfo[];
  pastWeek: MockChatInfo[];
  older: MockChatInfo[];
}

interface SidebarChatsProps {
  currentChatId: string | null;
  onChatSelect?: (chatUuid: string, chatTitle?: string) => void;
  onRegisterRefreshCallback?: (callback: () => void) => void;
  onRegisterUpdateTitleCallback?: (
    callback: (chatUuid: string, newTitle: string) => void
  ) => void;
}

const SidebarChats: React.FC<SidebarChatsProps> = ({
  currentChatId,
  onChatSelect,
  onRegisterRefreshCallback,
  onRegisterUpdateTitleCallback,
}) => {
  const [chats, setChats] = useState<MockChatInfo[]>([]);
  const [loading, setLoading] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [deleteModalVisible, setDeleteModalVisible] = useState(false);
  const [deletingChatId, setDeletingChatId] = useState<string | null>(null);
  const [deletingChatTitle, setDeletingChatTitle] = useState<string>('');
  // 内联编辑状态
  const [editingChatId, setEditingChatId] = useState<string | null>(null);
  const [editingTitle, setEditingTitle] = useState<string>('');
  const [totalCount, setTotalCount] = useState<number>(0);
  const [hasMore, setHasMore] = useState(true);
  const loadingRef = useRef(false);
  const searchTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const searchQueryRef = useRef<string>('');
  const chatsCountRef = useRef<number>(0);
  const hasMoreRef = useRef<boolean>(true);
  const searchInputRef = useRef<any>(null);

  // 模拟聊天数据
  const getMockChats = useCallback((): MockChatInfo[] => {
    const now = Date.now();
    return [
      {
        chatUuid: 'mock-chat-1',
        title: 'AI 助手介绍',
        updatedAt: new Date(now - 1000 * 60 * 30).toISOString(), // 30分钟前
        messagesCount: 3,
      },
      {
        chatUuid: 'mock-chat-2', 
        title: '编程问题咨询',
        updatedAt: new Date(now - 1000 * 60 * 60 * 2).toISOString(), // 2小时前
        messagesCount: 8,
      },
      {
        chatUuid: 'mock-chat-3',
        title: '日常闲聊',
        updatedAt: new Date(now - 1000 * 60 * 60 * 24).toISOString(), // 昨天
        messagesCount: 12,
      },
      {
        chatUuid: 'mock-chat-4',
        title: '学习计划制定',
        updatedAt: new Date(now - 1000 * 60 * 60 * 24 * 3).toISOString(), // 3天前
        messagesCount: 15,
      },
      {
        chatUuid: 'mock-chat-5',
        title: '工作项目讨论',
        updatedAt: new Date(now - 1000 * 60 * 60 * 24 * 10).toISOString(), // 10天前
        messagesCount: 25,
      },
    ];
  }, []);

  // 同步hasMore状态到ref
  useEffect(() => {
    hasMoreRef.current = hasMore;
  }, [hasMore]);
  const getTimeGroup = (updatedAt: string): keyof GroupedChats => {
    const now = new Date();
    let chatDate: Date;

    // 处理不同的时间格式
    if (updatedAt.includes('-')) {
      // ISO格式: "2024-01-15T10:30:00Z"
      chatDate = new Date(updatedAt);
    } else {
      // Unix时间戳格式: "1705312200"
      chatDate = new Date(parseInt(updatedAt) * 1000);
    }

    // 重置时间为当天的00:00:00
    const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
    const yesterday = new Date(today.getTime() - 24 * 60 * 60 * 1000);
    const pastWeek = new Date(today.getTime() - 7 * 24 * 60 * 60 * 1000);

    // 重置聊天日期为当天的00:00:00进行比较
    const chatDateOnly = new Date(
      chatDate.getFullYear(),
      chatDate.getMonth(),
      chatDate.getDate()
    );

    if (chatDateOnly.getTime() === today.getTime()) {
      return 'today';
    } else if (chatDateOnly.getTime() === yesterday.getTime()) {
      return 'yesterday';
    } else if (chatDateOnly >= pastWeek) {
      return 'pastWeek';
    } else {
      return 'older';
    }
  };

  // 对聊天记录进行分组
  const groupedChats = useMemo((): GroupedChats => {
    const result = chats.reduce<GroupedChats>(
      (groups, chat) => {
        if (chat.updatedAt) {
          const group = getTimeGroup(chat.updatedAt);
          groups[group].push(chat);
          console.log(
            `聊天 "${chat.title}" 分组到: ${group}, 时间: ${chat.updatedAt}`
          );
        }
        return groups;
      },
      { today: [], yesterday: [], pastWeek: [], older: [] }
    );

    console.log('分组结果:', {
      today: result.today.length,
      yesterday: result.yesterday.length,
      pastWeek: result.pastWeek.length,
      older: result.older.length,
    });

    return result;
  }, [chats]);

  // 加载聊天列表 (模拟实现)
  const loadChats = useCallback(
    async (keyword?: string, isLoadMore = false) => {
      // 防止重复请求
      if (loadingRef.current) return;

      loadingRef.current = true;
      if (isLoadMore) {
        setLoadingMore(true);
      } else {
        setLoading(true);
      }

      try {
        // 模拟加载延迟
        await new Promise(resolve => setTimeout(resolve, 300));
        
        // 获取模拟数据
        const allMockChats = getMockChats();
        
        // 根据关键词过滤
        const filteredChats = keyword 
          ? allMockChats.filter(chat => 
              chat.title.toLowerCase().includes(keyword.toLowerCase())
            )
          : allMockChats;
        
        const currentOffset = isLoadMore ? chatsCountRef.current : 0;
        const pageSize = 50;
        const paginatedChats = filteredChats.slice(currentOffset, currentOffset + pageSize);
        
        if (isLoadMore) {
          // 加载更多时追加到现有列表
          setChats(prev => {
            const existingChatsMap = new Map(
              prev.map(chat => [chat.chatUuid, chat])
            );
            
            paginatedChats.forEach(newChat => {
              if (!existingChatsMap.has(newChat.chatUuid)) {
                existingChatsMap.set(newChat.chatUuid, newChat);
              }
            });
            
            const mergedChats = Array.from(existingChatsMap.values());
            chatsCountRef.current = mergedChats.length;
            return mergedChats;
          });
        } else {
          // 初始加载或搜索时替换列表
          setChats(paginatedChats);
          chatsCountRef.current = paginatedChats.length;
        }
        
        setTotalCount(filteredChats.length);
        setHasMore(currentOffset + pageSize < filteredChats.length);
        
        console.log('模拟聊天列表加载完成', {
          新加载: paginatedChats.length,
          当前总数: isLoadMore ? chatsCountRef.current : paginatedChats.length,
          总数量: filteredChats.length,
          还有更多: currentOffset + pageSize < filteredChats.length,
          isLoadMore,
        });
        
      } catch (error) {
        console.error('Failed to load chats:', error);
        message.error('加载聊天列表失败');
      } finally {
        setLoading(false);
        setLoadingMore(false);
        loadingRef.current = false;
      }
    },
    [getMockChats]
  );

  // 加载更多聊天
  const loadMoreChats = useCallback(() => {
    if (!hasMoreRef.current || loadingRef.current) {
      console.log('loadMoreChats: 条件不满足', {
        hasMore: hasMoreRef.current,
        loading: loadingRef.current,
      });
      return;
    }
    console.log('loadMoreChats: 开始加载更多');
    loadChats(searchQueryRef.current || undefined, true);
  }, []); // 移除loadChats依赖，直接调用

  // 处理搜索输入变化
  const handleSearchChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const newValue = e.target.value;
      setSearchQuery(newValue);
    },
    []
  );

  // 滚动事件处理
  const handleScroll = useCallback(
    (e: Event) => {
      const container = e.target as HTMLDivElement;
      const { scrollTop, scrollHeight, clientHeight } = container;

      // 判断是否滚动到底部（留有50px的缓冲区）
      const distanceFromBottom = scrollHeight - scrollTop - clientHeight;

      // 添加更详细的调试日志
      console.log('滚动事件触发:', {
        scrollTop,
        scrollHeight,
        clientHeight,
        distanceFromBottom,
        hasMore: hasMoreRef.current,
        loading: loadingRef.current,
        当前聊天数量: chatsCountRef.current,
        总数量: totalCount,
        searchQuery: searchQueryRef.current,
      });

      if (distanceFromBottom <= 50) {
        console.log('到达底部，准备加载更多:', {
          hasMore: hasMoreRef.current,
          loading: loadingRef.current,
          当前聊天数量: chatsCountRef.current,
          distanceFromBottom,
        });
        if (hasMoreRef.current && !loadingRef.current) {
          console.log('开始加载更多数据');
          loadMoreChats();
        } else {
          console.log('跳过加载更多，原因:', {
            hasMore: hasMoreRef.current,
            loading: loadingRef.current,
          });
        }
      }
    },
    [loadMoreChats, totalCount]
  );

  // 初始加载
  useEffect(() => {
    loadChats();
  }, []); // 使用空依赖数组，只在组件挂载时执行

  // 搜索防抖处理
  useEffect(() => {
    // 更新搜索查询的ref
    searchQueryRef.current = searchQuery;

    // 清除之前的定时器
    if (searchTimeoutRef.current) {
      clearTimeout(searchTimeoutRef.current);
    }

    searchTimeoutRef.current = setTimeout(async () => {
      // 搜索时重置分页状态
      setHasMore(true);
      chatsCountRef.current = 0; // 重置聊天数量
      console.log('搜索开始，重置分页状态:', {
        searchQuery,
        hasMore: true,
        chatsCount: 0,
      });
      // 直接调用loadChats，传入搜索词
      loadChats(searchQuery || undefined);

      // 搜索完成后恢复输入框焦点
      if (searchInputRef.current) {
        const inputElement =
          searchInputRef.current.input || searchInputRef.current;
        if (inputElement && inputElement.focus) {
          // 使用 setTimeout 确保在渲染完成后恢复焦点
          setTimeout(() => {
            inputElement.focus();
            // 将光标移动到文本末尾
            const length = inputElement.value?.length || 0;
            inputElement.setSelectionRange(length, length);
          }, 0);
        }
      }
    }, 300); // 300ms 防抖延迟

    return () => {
      if (searchTimeoutRef.current) {
        clearTimeout(searchTimeoutRef.current);
      }
    };
  }, [searchQuery]);

  // 滚动事件监听
  useEffect(() => {
    const container = containerRef.current;
    if (!container) {
      console.error('滚动容器未找到');
      return;
    }

    console.log('滚动事件监听器已绑定', {
      containerHeight: container.clientHeight,
      containerScrollHeight: container.scrollHeight,
      hasScroll: container.scrollHeight > container.clientHeight,
    });

    container.addEventListener('scroll', handleScroll, { passive: true });

    return () => {
      container.removeEventListener('scroll', handleScroll);
    };
  }, [handleScroll]);

  // 注册刷新回调函数
  useEffect(() => {
    if (onRegisterRefreshCallback) {
      onRegisterRefreshCallback(() => loadChats());
    }
  }, [onRegisterRefreshCallback]); // 移除loadChats依赖，直接调用

  // 更新聊天标题
  const updateChatTitle = useCallback((chatUuid: string, newTitle: string) => {
    setChats(prev =>
      prev.map(chat =>
        chat.chatUuid === chatUuid ? { ...chat, title: newTitle } : chat
      )
    );
  }, []);

  // 注册更新标题回调函数
  useEffect(() => {
    if (onRegisterUpdateTitleCallback) {
      onRegisterUpdateTitleCallback(updateChatTitle);
    }
  }, [onRegisterUpdateTitleCallback, updateChatTitle]);

  // 处理聊天选择
  const handleChatSelect = (chatUuid: string, chatTitle?: string) => {
    onChatSelect?.(chatUuid, chatTitle);
  };

  // 处理收藏聊天
  const handleFavoriteChat = async (_chatUuid: string, e: React.MouseEvent) => {
    e.stopPropagation();
    try {
      // TODO: 实现收藏功能，目前先显示提示
      message.info('收藏功能开发中...');
    } catch (error) {
      console.error('Failed to favorite chat:', error);
      message.error('收藏失败');
    }
  };

  // 开始内联编辑
  const startInlineEdit = (chatUuid: string, chatTitle: string) => {
    // 根据项目规范，只有已保存的对话（有有效的 chatUuid）才允许重命名
    const canRename = Boolean(chatUuid && chatUuid.trim() !== '');
    if (!canRename) {
      message.warning('请先保存对话后再重命名');
      return;
    }
    
    setEditingChatId(chatUuid);
    setEditingTitle(chatTitle || '新对话');
  };

  // 确认内联编辑 (模拟实现)
  const confirmInlineEdit = async () => {
    if (!editingChatId || !editingTitle.trim()) {
      message.error('请输入有效的对话标题');
      return;
    }

    try {
      // 模拟保存延迟
      await new Promise(resolve => setTimeout(resolve, 200));
      
      // 更新本地状态
      setChats(prev =>
        prev.map(chat =>
          chat.chatUuid === editingChatId
            ? { ...chat, title: editingTitle.trim() }
            : chat
        )
      );
      
      // 调用外部更新回调
      updateChatTitle(editingChatId, editingTitle.trim());
      
      message.success('重命名成功');
      setEditingChatId(null);
      setEditingTitle('');
    } catch (error) {
      console.error('Failed to rename chat:', error);
      message.error('重命名失败');
    }
  };

  // 取消内联编辑
  const cancelInlineEdit = () => {
    setEditingChatId(null);
    setEditingTitle('');
  };

  // 处理输入框按键事件
  const handleEditKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    e.stopPropagation(); // 阻止事件冒泡，避免触发聊天选中
    if (e.key === 'Enter') {
      confirmInlineEdit();
    } else if (e.key === 'Escape') {
      cancelInlineEdit();
    }
  };

  // 显示删除确认对话框
  const showDeleteConfirm = (chatUuid: string, chatTitle: string) => {
    setDeletingChatId(chatUuid);
    setDeletingChatTitle(chatTitle || '新对话');
    setDeleteModalVisible(true);
  };

  // 处理删除聊天 (模拟实现)
  const handleDeleteChat = async () => {
    if (!deletingChatId) return;

    try {
      // 模拟删除延迟
      await new Promise(resolve => setTimeout(resolve, 200));
      
      setChats(prev => {
        const newChats = prev.filter(chat => chat.chatUuid !== deletingChatId);
        // 更新ref中的聊天数量
        chatsCountRef.current = newChats.length;
        return newChats;
      });
      message.success('删除成功');
      setDeleteModalVisible(false);
      setDeletingChatId(null);
      setDeletingChatTitle('');
    } catch (error) {
      console.error('Failed to delete chat:', error);
      message.error('删除失败');
    }
  };

  // 取消删除
  const handleCancelDelete = () => {
    setDeleteModalVisible(false);
    setDeletingChatId(null);
    setDeletingChatTitle('');
  };

  // 获取菜单项配置
  const getMenuItems = (
    chatUuid: string,
    chatTitle: string
  ): MenuProps['items'] => {
    // 检查是否允许重命名（根据项目规范）
    const canRename = Boolean(chatUuid && chatUuid.trim() !== '');
    
    return [
      {
        key: 'favorite',
        icon: <StarOutlined />,
        label: '收藏',
        onClick: () => handleFavoriteChat(chatUuid, {} as React.MouseEvent),
      },
      {
        key: 'rename',
        icon: <EditOutlined />,
        label: '重命名',
        disabled: !canRename,
        onClick: () => startInlineEdit(chatUuid, chatTitle),
      },
      {
        key: 'delete',
        icon: <DeleteOutlined />,
        label: '删除',
        danger: true,
        onClick: () => showDeleteConfirm(chatUuid, chatTitle),
      },
    ];
  };

  // 格式化时间显示
  const formatTime = (updatedAt: string) => {
    let date: Date;

    // 处理不同的时间格式
    if (updatedAt.includes('-')) {
      // ISO格式: "2024-01-15T10:30:00Z"
      date = new Date(updatedAt);
    } else {
      // Unix时间戳格式: "1705312200"
      date = new Date(parseInt(updatedAt) * 1000);
    }

    const now = new Date();
    const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
    const yesterday = new Date(today.getTime() - 24 * 60 * 60 * 1000);

    // 重置聊天日期为当天的00:00:00进行比较
    const chatDateOnly = new Date(
      date.getFullYear(),
      date.getMonth(),
      date.getDate()
    );

    if (chatDateOnly.getTime() === today.getTime()) {
      // 今天：显示时间
      return date.toLocaleTimeString('zh-CN', {
        hour: '2-digit',
        minute: '2-digit',
      });
    } else if (chatDateOnly.getTime() === yesterday.getTime()) {
      // 昨天：显示"昨天"
      return '昨天';
    } else if (
      chatDateOnly >= new Date(today.getTime() - 7 * 24 * 60 * 60 * 1000)
    ) {
      // 过去7天：显示星期几
      return date.toLocaleDateString('zh-CN', { weekday: 'short' });
    } else {
      // 更久以前：显示月日
      return date.toLocaleDateString('zh-CN', {
        month: 'short',
        day: 'numeric',
      });
    }
  };

  // 渲染聊天项
  const renderChatItem = (chat: MockChatInfo) => {
    const isEditing = editingChatId === chat.chatUuid;
    
    return (
      <List.Item
        key={chat.chatUuid}
        style={{ padding: 0 }}
        className={`${styles.chatItem}`}
        onClick={() => !isEditing && handleChatSelect(chat.chatUuid!, chat.title)}
      >
        <div
          className={`${styles.chatContent} ${currentChatId === chat.chatUuid ? styles.active : ''} ${isEditing ? styles.editing : ''}`}
        >
          <div className={styles.chatHeader}>
            {isEditing ? (
              // 编辑状态
              <div className={styles.editContainer}>
                <Input
                  value={editingTitle}
                  onChange={e => setEditingTitle(e.target.value)}
                  onKeyDown={handleEditKeyDown}
                  className={styles.editInput}
                  maxLength={100}
                  autoFocus
                  onClick={e => e.stopPropagation()}
                />
                <div className={styles.editActions}>
                  <Button
                    type="text"
                    size="small"
                    icon={<CheckOutlined />}
                    className={styles.confirmBtn}
                    onClick={e => {
                      e.stopPropagation();
                      confirmInlineEdit();
                    }}
                    title="确认"
                  />
                  <Button
                    type="text"
                    size="small"
                    icon={<CloseOutlined />}
                    className={styles.cancelBtn}
                    onClick={e => {
                      e.stopPropagation();
                      cancelInlineEdit();
                    }}
                    title="取消"
                  />
                </div>
              </div>
            ) : (
              // 正常状态
              <>
                <div className={styles.chatTitle} title={chat.title || '新对话'}>
                  {chat.title || '新对话'}
                </div>
                <div className={styles.chatActions}>
                  <Text className={styles.chatTime} hidden={true}>
                    {chat.updatedAt && formatTime(chat.updatedAt)}
                  </Text>
                  <Dropdown
                    menu={{
                      items: getMenuItems(chat.chatUuid!, chat.title || '新对话'),
                    }}
                    trigger={['click']}
                    placement="bottomRight"
                  >
                    <Button
                      type="text"
                      size="small"
                      icon={<MoreOutlined />}
                      className={styles.moreButton}
                      onClick={(e: React.MouseEvent) => e.stopPropagation()}
                    />
                  </Dropdown>
                </div>
              </>
            )}
          </div>
        </div>
      </List.Item>
    );
  };

  // 渲染分组
  const renderGroup = (title: string, chats: CommonChatInfo[]) => {
    if (chats.length === 0) return null;

    return (
      <div key={title} className={styles.chatGroup}>
        <div className={styles.groupTitleGroup}>
          <Divider orientation="left" className={styles.groupTitle}>
            <Text type="secondary">{title}</Text>
          </Divider>
        </div>
        <List
          dataSource={chats}
          renderItem={renderChatItem}
          className={styles.chatList}
        />
      </div>
    );
  };

  if (loading && chats.length === 0) {
    return (
      <div className={styles.loadingContainer}>
        <Spin size="small" />
      </div>
    );
  }

  return (
    <div className={styles.sidebarChats}>
      <div className={styles.searchContainer}>
        <Search
          ref={searchInputRef}
          placeholder="搜索对话..."
          value={searchQuery}
          onChange={handleSearchChange}
          prefix={<SearchOutlined />}
          allowClear
        />
      </div>

      <div className={styles.chatsContainer} ref={containerRef}>
        {Object.keys(groupedChats).some(
          key => groupedChats[key as keyof GroupedChats].length > 0
        ) ? (
          <>
            {renderGroup('今天', groupedChats.today)}
            {renderGroup('昨天', groupedChats.yesterday)}
            {renderGroup('过去7天', groupedChats.pastWeek)}
            {renderGroup('更久以前', groupedChats.older)}

            {/* 加载更多按钮或已加载全部提示 */}
            {loadingMore && (
              <div className={styles.loadingContainer}>
                <Spin size="small" />
                <span
                  style={{
                    marginLeft: '8px',
                    fontSize: '14px',
                    color: 'var(--text-color-secondary)',
                  }}
                >
                  加载中...
                </span>
              </div>
            )}

            {!hasMore && chats.length > 0 && (
              <div className={styles.endContainer}>
                <Text type="secondary" className={styles.endText}>
                  已加载全部聊天记录 ({totalCount} 条)
                </Text>
              </div>
            )}
          </>
        ) : (
          <Empty
            image={Empty.PRESENTED_IMAGE_SIMPLE}
            description={searchQuery ? '未找到匹配的对话' : '暂无对话记录'}
          />
        )}
      </div>

      {/* 删除确认对话框 */}
      <Modal
        title={
          <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
            <ExclamationCircleOutlined
              style={{ color: '#faad14', fontSize: '16px' }}
            />
            <span>确认删除</span>
          </div>
        }
        open={deleteModalVisible}
        onOk={handleDeleteChat}
        onCancel={handleCancelDelete}
        okText="删除"
        cancelText="取消"
        okButtonProps={{ danger: true }}
        confirmLoading={loading}
      >
        <p>确定要删除对话 "{deletingChatTitle}" 吗？</p>
        <p style={{ color: '#666', fontSize: '14px', marginTop: '8px' }}>
          删除后无法恢复，请谨慎操作。
        </p>
      </Modal>
    </div>
  );
};

export default SidebarChats;
