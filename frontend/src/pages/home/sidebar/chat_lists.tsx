import React, {useCallback, useEffect, useMemo, useRef, useState,} from 'react';
import { Events } from '@wailsio/runtime';
import type {MenuProps} from 'antd';
import {Button, Divider, Dropdown, Empty, Input, List, message, Modal, Spin, Typography,} from 'antd';
import {
    CheckOutlined,
    CloseOutlined,
    DeleteOutlined,
    EditOutlined,
    ExclamationCircleOutlined,
    MoreOutlined,
    PauseOutlined,
    SearchOutlined,
    StarFilled,
    StarOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import styles from './chats_lists.module.scss';
import {
    Chat,
    ChatList
} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models/index.ts";
import {Service} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service";

const {Text} = Typography;
const {Search} = Input;

interface GroupedChats {
    today: Chat[];
    yesterday: Chat[];
    pastWeek: Chat[];
    older: Chat[];
}

interface SidebarChatsProps {
    currentChatUuid: string | null;
    onChatSelect?: (chatUuid: string) => void;
    onRegisterRefreshCallback?: (callback: () => void) => void;
    onDeleteChat?: (chatUuid: string) => void;
    activeTab?: 'history' | 'favorites'; // 添加activeTab属性
    generatingChatUuids?: string[];
    onStopGenerationForChat?: (chatUuid: string) => void;
}

const SidebarChats: React.FC<SidebarChatsProps> = ({
                                                       currentChatUuid,
                                                       onChatSelect,
                                                       onRegisterRefreshCallback,
                                                       onDeleteChat,
                                                       activeTab = 'history', // 默认为历史对话
                                                       generatingChatUuids = [],
                                                       onStopGenerationForChat,
                                                   }) => {
    const { t, i18n } = useTranslation();
    const [chats, setChats] = useState<Chat[]>([]);
    const [loading, setLoading] = useState(false);
    const [loadingMore, setLoadingMore] = useState(false);
    const [searchQuery, setSearchQuery] = useState('');
    const [deleteModalVisible, setDeleteModalVisible] = useState(false);
    const [deletingChatUuid, setDeletingChatUuid] = useState<string | null>(null);
    const [deletingChatTitle, setDeletingChatTitle] = useState<string>('');
    // 内联编辑状态
    const [editingChatUuid, setEditingChatUuid] = useState<string | null>(null);
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

    interface ChatTitleEvent {
        chat_uuid: string;
        title: string;
    }

    // 同步hasMore状态到ref
    useEffect(() => {
        hasMoreRef.current = hasMore;
    }, [hasMore]);

    // 通过更新时间获取标记组
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

        if (chatDateOnly.getTime() >= today.getTime()) {
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
                if (chat.updated_at) {
                    const group = getTimeGroup(chat.updated_at);
                    groups[group].push(chat);
                }
                return groups;
            },
            {today: [], yesterday: [], pastWeek: [], older: []}
        );
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
                // 使用ref来获取最新的chats长度
                const currentOffset = isLoadMore ? chatsCountRef.current : 0;
                // 根据activeTab决定是否为收藏列表
                const isFavorites = activeTab === 'favorites';
                const response: ChatList | null = await Service.ChatList(currentOffset, 50, keyword || null, isFavorites);
                console.log("ChatList response:", response)
                
                // 处理响应为空的情况
                if (!response || !response.lists) {
                    if (isLoadMore) {
                        // 加载更多时返回空数据，说明已经加载完所有数据
                        setHasMore(false);
                    } else {
                        // 初始加载时返回空数据
                        setTotalCount(0);
                        setHasMore(false);
                        setChats([]);
                    }
                    return;
                }

                const newChats: Chat[] = response.lists;
                const total: number = response.total || 0;
                const limit = 50; // 每页加载数量
                
                // 如果返回的数据为空，说明已经加载完所有数据
                if (newChats.length === 0) {
                    setHasMore(false);
                    // 如果是初始加载，清空列表
                    if (!isLoadMore) {
                        setTotalCount(0);
                        setChats([]);
                    }
                    return;
                }

                let currentTotal = 0;
                if (isLoadMore) {
                    // 加载更多时追加到现有列表，使用chatUuid去重
                    setChats(prev => {
                        // 创建一个Map来存储已有的聊天记录，以chatUuid为key
                        const existingChatsMap = new Map(
                            prev.map(chat => [chat.uuid, chat])
                        );

                        // 添加新的聊天记录，如果chatUuid已存在则跳过
                        newChats.forEach(newChat => {
                            if (
                                newChat.uuid &&
                                !existingChatsMap.has(newChat.uuid)
                            ) {
                                existingChatsMap.set(newChat.uuid, newChat);
                            }
                        });

                        const mergedChats = Array.from(existingChatsMap.values());
                        currentTotal = mergedChats.length;
                        chatsCountRef.current = currentTotal;

                        return mergedChats;
                    });
                } else {
                    // 初始加载或搜索时替换列表
                    setChats(newChats);
                    currentTotal = newChats.length;
                    // 更新ref中的聊天数量
                    chatsCountRef.current = currentTotal;
                }
                
                setTotalCount(total);
                
                // 判断是否还有更多数据：
                // 1. 当前已加载数量小于总数
                // 2. 返回的数据数量等于limit，说明可能还有更多数据
                const hasMoreData = currentTotal < total && newChats.length >= limit;
                setHasMore(hasMoreData);

            } catch (error) {
                console.error('Failed to load chats:', error);
                message.error(t('home.chatList.loadFailed'));
            } finally {
                setLoading(false);
                setLoadingMore(false);
                loadingRef.current = false;
            }
        }, [activeTab, t]
    );

    // 加载更多聊天
    const loadMoreChats = useCallback(() => {
        if (!hasMoreRef.current || loadingRef.current) {
            return;
        }
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
            const {scrollTop, scrollHeight, clientHeight} = container;

            // 判断是否滚动到底部（留有50px的缓冲区）
            const distanceFromBottom = scrollHeight - scrollTop - clientHeight;

            if (distanceFromBottom <= 50 && (hasMoreRef.current && !loadingRef.current)) {
                loadMoreChats();
            }
        },
        [loadMoreChats, totalCount]
    );

    // 初始加载
    useEffect(() => {
        loadChats();
    }, []); // 使用空依赖数组，只在组件挂载时执行

    // 监听activeTab变化，重新加载数据
    useEffect(() => {
        loadChats();
    }, [activeTab, loadChats]); // 添加activeTab依赖

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
            chatsCountRef.current = 0;
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

        container.addEventListener('scroll', handleScroll, {passive: true});

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
                chat.uuid === chatUuid
                    ? new Chat({...chat, title: newTitle})
                    : chat
            )
        );
    }, []);

    useEffect(() => {
        const eventKey = 'event:chat_title:all';
        const cancel = Events.On(eventKey, (event) => {
            const payload = event.data as ChatTitleEvent;
            if (!payload?.chat_uuid || !payload?.title) {
                return;
            }
            updateChatTitle(payload.chat_uuid, payload.title);
        });

        return () => {
            cancel?.();
            Events.Off(eventKey);
        };
    }, [updateChatTitle]);

    // 处理聊天选择
    const handleChatSelect = (chatUuid: string) => {
        onChatSelect?.(chatUuid);
    };

    // 处理收藏聊天
    const handleFavoriteChat = useCallback(async (chat: Chat, e: React.MouseEvent): Promise<void> => {
        try {
            await Service.CollectionChat(chat.uuid, !chat.is_collection);

            // 更新本地状态
            setChats(prev =>
                prev.map(item =>
                    item.uuid === chat.uuid
                        ? new Chat({...item, is_collection: !chat.is_collection})
                        : item
                )
            );

            // 刷新对话列表以确保UI与后端状态一致
            loadChats();
        } catch (error) {
            console.error('Failed to favorite chat:', error);
            message.error(t('home.chatList.favoriteFailed'));
        }
    }, [loadChats, t]);

    // 开始内联编辑
    const startInlineEdit = (chatUuid: string, chatTitle: string) => {
        // 根据项目规范，只有已保存的对话（有有效的 chatUuid）才允许重命名
        const canRename = chatUuid;
        if (!canRename && chatUuid == "") {
            message.warning(t('home.chatList.renameNeedSaved'));
            return;
        }

        setEditingChatUuid(chatUuid);
        setEditingTitle(chatTitle || t('home.chatList.defaultTitle'));
    };

    // 确认内联编辑 (模拟实现)
    const confirmInlineEdit = async () => {
        if (!editingChatUuid || !editingTitle.trim()) {
            message.error(t('home.chatList.invalidTitle'));
            return;
        }

        try {
            // 调用 RenameChat API 保存标题
            await Service.RenameChat(editingChatUuid, editingTitle.trim());

            const nextTitle = editingTitle.trim();
            updateChatTitle(editingChatUuid, nextTitle);

            message.success(t('home.chatList.renameSuccess'));
            setEditingChatUuid(null);
            setEditingTitle('');
        } catch (error) {
            console.error('Failed to rename chat:', error);
            message.error(t('home.chatList.renameFailed'));
        }
    };

    // 取消内联编辑
    const cancelInlineEdit = () => {
        setEditingChatUuid(null);
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
        setDeletingChatUuid(chatUuid);
        setDeletingChatTitle(chatTitle || t('home.chatList.defaultTitle'));
        setDeleteModalVisible(true);
    };

    // 处理删除聊天 (模拟实现)
    const handleDeleteChat = async () => {
        if (!deletingChatUuid) return;

        try {
            // 调用传入的删除函数
            if (onDeleteChat) {
                await onDeleteChat(deletingChatUuid);
            } else {
                // 如果没有传入删除函数，使用原来的模拟实现
                // 模拟删除延迟
                await new Promise(resolve => setTimeout(resolve, 200));

                setChats(prev => {
                    const newChats = prev.filter(chat => chat.uuid !== deletingChatUuid);
                    // 更新ref中的聊天数量
                    chatsCountRef.current = newChats.length;
                    return newChats;
                });
                message.success(t('home.chatList.deleteSuccess'));
            }
            setDeleteModalVisible(false);
            setDeletingChatUuid(null);
            setDeletingChatTitle('');
        } catch (error) {
            console.error('Failed to delete chat:', error);
            message.error(t('home.chatList.deleteFailed'));
        }
    };

    // 取消删除
    const handleCancelDelete = () => {
        setDeleteModalVisible(false);
        setDeletingChatUuid(null);
        setDeletingChatTitle('');
    };

    // 获取菜单项配置
    const getMenuItems = (chat: Chat): MenuProps['items'] => {
        // 检查是否允许重命名（根据项目规范）
        const canRename = Boolean(chat.uuid && chat.uuid.trim() !== '');

        return [
            {
                key: 'favorite',
                icon: chat.is_collection ? <StarFilled/> : <StarOutlined/>,
                label: chat.is_collection ? t('home.chatList.unfavorite') : t('home.chatList.favorite'),
                onClick: () => handleFavoriteChat(chat, {} as React.MouseEvent),
            },
            {
                key: 'rename',
                icon: <EditOutlined/>,
                label: t('home.chatList.rename'),
                disabled: !canRename,
                onClick: () => startInlineEdit(chat.uuid, chat.title),
            },
            {
                key: 'delete',
                icon: <DeleteOutlined/>,
                label: t('home.chatList.delete'),
                danger: true,
                onClick: () => showDeleteConfirm(chat.uuid, chat.title),
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
            return date.toLocaleTimeString(i18n.language === 'en-US' ? 'en-US' : 'zh-CN', {
                hour: '2-digit',
                minute: '2-digit',
            });
        } else if (chatDateOnly.getTime() === yesterday.getTime()) {
            // 昨天：显示"昨天"
            return t('home.chatList.yesterday');
        } else if (
            chatDateOnly >= new Date(today.getTime() - 7 * 24 * 60 * 60 * 1000)
        ) {
            // 过去7天：显示星期几
            return date.toLocaleDateString(i18n.language === 'en-US' ? 'en-US' : 'zh-CN', {weekday: 'short'});
        } else {
            // 更久以前：显示月日
            return date.toLocaleDateString(i18n.language === 'en-US' ? 'en-US' : 'zh-CN', {
                month: 'short',
                day: 'numeric',
            });
        }
    };

    const generatingSet = useMemo(
        () => new Set(generatingChatUuids.filter(Boolean)),
        [generatingChatUuids],
    );

    // 渲染聊天项
    const renderChatItem = (chat: Chat) => {
        const isEditing = editingChatUuid === chat.uuid;
        const chatUuid = chat.uuid ?? '';
        const isGenerating = Boolean(chatUuid && generatingSet.has(chatUuid));

        return (
            <List.Item
                key={chat.uuid}
                style={{padding: 0}}
                className={`${styles.chatItem}`}
                onClick={() => !isEditing && handleChatSelect(chat.uuid!)}
            >
                <div
                    className={`${styles.chatContent} ${currentChatUuid === chat.uuid ? styles.active : ''} ${isEditing ? styles.editing : ''}`}
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
                                        icon={<CheckOutlined/>}
                                        className={styles.confirmBtn}
                                        onClick={e => {
                                            e.stopPropagation();
                                            confirmInlineEdit();
                                        }}
                                        title={t('home.chatList.editConfirm')}
                                    />
                                    <Button
                                        type="text"
                                        size="small"
                                        icon={<CloseOutlined/>}
                                        className={styles.cancelBtn}
                                        onClick={e => {
                                            e.stopPropagation();
                                            cancelInlineEdit();
                                        }}
                                        title={t('home.chatList.editCancel')}
                                    />
                                </div>
                            </div>
                        ) : (
                            // 正常状态
                            <>
                                <div className={styles.chatTitleRow}>
                                    {isGenerating && (
                                        <button
                                            type="button"
                                            className={styles.generatingStopBtn}
                                            title={t('home.chatList.stopGeneration')}
                                            aria-label={t('home.chatList.stopGenerationAria')}
                                            onClick={e => {
                                                e.stopPropagation();
                                                e.preventDefault();
                                                onStopGenerationForChat?.(chatUuid);
                                            }}
                                        >
                                            <span className={styles.generatingSpinner} aria-hidden />
                                            <PauseOutlined
                                                className={styles.generatingPauseIcon}
                                                aria-hidden
                                            />
                                        </button>
                                    )}
                                    <div
                                        className={styles.chatTitle}
                                        title={chat.title || t('home.chatList.defaultTitle')}
                                    >
                                        {chat.title || t('home.chatList.defaultTitle')}
                                    </div>
                                </div>
                                <div className={styles.chatActions}>
                                    <Text className={styles.chatTime} hidden={true}>
                                        {chat.updated_at && formatTime(chat.updated_at)}
                                    </Text>
                                    <Dropdown
                                        menu={{
                                            items: getMenuItems(chat),
                                        }}
                                        trigger={['click']}
                                        placement="bottomRight"
                                    >
                                        <Button
                                            type="text"
                                            size="small"
                                            icon={<MoreOutlined/>}
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
    const renderGroup = (title: string, chats: Chat[]) => {
        if (chats.length === 0) return null;


        return (
            <div key={title} className={styles.chatGroup}>
                <div className={styles.groupTitleGroup}>
                    <Divider orientation="horizontal" className={styles.groupTitle}>
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

    // 修改渲染部分以根据activeTab显示不同内容
    const renderContent = () => {
        // 历史对话和收藏tab都使用相同的聊天列表渲染逻辑
        if (loading && chats.length === 0) {
            return (
                <div className={styles.loadingContainer}>
                    <Spin size="small"/>
                </div>
            );
        }

        return (
            <div className={styles.chatsContainer} ref={containerRef}>
                {Object.keys(groupedChats).some(
                    key => groupedChats[key as keyof GroupedChats].length > 0
                ) ? (
                    <>
                        {renderGroup(t('home.chatList.today'), groupedChats.today)}
                        {renderGroup(t('home.chatList.yesterday'), groupedChats.yesterday)}
                        {renderGroup(t('home.chatList.pastWeek'), groupedChats.pastWeek)}
                        {renderGroup(t('home.chatList.older'), groupedChats.older)}

                        {/* 加载更多按钮或已加载全部提示 */}
                        {loadingMore && (
                            <div className={styles.loadingContainer}>
                                <Spin size="small"/>
                                <span
                                    style={{
                                        marginLeft: '8px',
                                        fontSize: '14px',
                                        color: 'var(--text-color-secondary)',
                                    }}
                                >
                  {t('common.loading')}
                </span>
                            </div>
                        )}

                        {!hasMore && chats.length > 0 && (
                            <div className={styles.endContainer}>
                                <Text type="secondary" className={styles.endText}>
                                    {activeTab === 'favorites'
                                        ? t('home.chatList.allFavoritesLoaded', { count: totalCount })
                                        : t('home.chatList.allHistoryLoaded', { count: totalCount })}
                                </Text>
                            </div>
                        )}
                    </>
                ) : (
                    <Empty
                        image={Empty.PRESENTED_IMAGE_SIMPLE}
                        description={searchQuery
                            ? t('home.chatList.noResults')
                            : (activeTab === 'favorites' ? t('home.chatList.emptyFavorites') : t('home.chatList.emptyHistory'))}
                    >
                        {activeTab === 'favorites' && (
                            <p style={{color: 'var(--text-color-secondary)', fontSize: '14px'}}>
                                {t('home.chatList.favoritesHint')}
                            </p>
                        )}
                    </Empty>
                )}
            </div>
        );
    };

    return (
        <div className={styles.sidebarChats}>
            <div className={styles.searchContainer}>
                <Search
                    ref={searchInputRef}
                    placeholder={t('home.chatList.searchPlaceholder')}
                    value={searchQuery}
                    onChange={handleSearchChange}
                    prefix={<SearchOutlined/>}
                    allowClear
                />
            </div>

            {renderContent()}

            {/* 删除确认对话框 */}
            <Modal
                title={
                    <div style={{display: 'flex', alignItems: 'center', gap: '8px'}}>
                        <ExclamationCircleOutlined
                            style={{color: '#faad14', fontSize: '16px'}}
                        />
                        <span>{t('home.chatList.deleteModalTitle')}</span>
                    </div>
                }
                open={deleteModalVisible}
                onOk={handleDeleteChat}
                onCancel={handleCancelDelete}
                okText={t('home.chatList.deleteModalOk')}
                cancelText={t('home.chatList.deleteModalCancel')}
                okButtonProps={{danger: true}}
                confirmLoading={loading}
            >
                <p>{t('home.chatList.deleteModalContent', { title: deletingChatTitle })}</p>
                <p style={{color: '#666', fontSize: '14px', marginTop: '8px'}}>
                    {t('home.chatList.deleteWarning')}
                </p>
            </Modal>
        </div>
    );
};

export default SidebarChats;
