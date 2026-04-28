import React, {useCallback, useEffect, useRef, useState} from 'react';
import {BackTop, Layout, message} from 'antd';
import {useNavigate, useParams} from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import Index from './sidebar';
import {useViewportHeight} from '@/hooks/useViewportHeight';
import './index.module.scss';
import Chat from '@/pages/home/chat';
import {Service} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/index.ts";
import styles from './index.module.scss';
import { useOPCStore } from '@/stores/opcStore';
import OPCSidebar from './opc/sidebar';
import OPCChatArea from './opc/chat';
import AddPersonDialog from './opc/dialogs/AddPersonDialog';
import CreateGroupDialog from './opc/dialogs/CreateGroupDialog';

const SIDEBAR_WIDTH_KEY = 'home_sidebar_width';
const SIDEBAR_DEFAULT_WIDTH = 280;
const SIDEBAR_MIN_WIDTH = 220;
const SIDEBAR_MAX_WIDTH = 480;

function readStoredSidebarWidth(): number {
    try {
        const raw = localStorage.getItem(SIDEBAR_WIDTH_KEY);
        if (!raw) return SIDEBAR_DEFAULT_WIDTH;
        const parsed = parseInt(raw, 10);
        if (Number.isNaN(parsed)) return SIDEBAR_DEFAULT_WIDTH;
        return Math.min(SIDEBAR_MAX_WIDTH, Math.max(SIDEBAR_MIN_WIDTH, parsed));
    } catch {
        return SIDEBAR_DEFAULT_WIDTH;
    }
}

const {Content, Sider} = Layout;

interface ChatPageProps {
    className?: string;
}

const ChatPage: React.FC<ChatPageProps> = ({className}) => {
    const { t } = useTranslation();
    // 获取路由参数和导航函数
    const {chatUuid: urlChatUuid} = useParams<{ chatUuid?: string }>();
    const navigate = useNavigate();
    // 本地状态管理
    const [currentChatUuid, setCurrentChatUuid] = useState<string>(urlChatUuid ?? '');
    // 移动端侧边栏遮罩显示状态（桌面端永远显示）
    const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false);
    const [sidebarWidth, setSidebarWidth] = useState<number>(() => readStoredSidebarWidth());
    const [refreshChatList, setRefreshChatList] = useState<(() => void) | null>(
        null
    );
    const [generatingChatUuids, setGeneratingChatUuids] = useState<string[]>([]);
    const stopGenerationForChatRef = useRef<(uuid: string) => void>(() => {});
    // 历史聊天首次加载完成（用于立即滚动到底部，无动画）
    const [isFirstHistoricalLoad, setIsFirstHistoricalLoad] = useState(false);
    // 使用视口高度检测 Hook
    const {isMobile} = useViewportHeight();
    // OPC 模式状态
    const { mode: appMode } = useOPCStore();
    const [addPersonOpen, setAddPersonOpen] = useState(false);
    const [createGroupOpen, setCreateGroupOpen] = useState(false);
    const [editPersonData, setEditPersonData] = useState<any>(null);
    const [editGroupData, setEditGroupData] = useState<any>(null);

    const handleEditPerson = async (uuid: string) => {
        try {
            const person = await Service.OPCGetPerson(uuid);
            if (!person) return;
            // 需要获取 agent 详情来拿到 prompt/tools/skills
            const agentDetail = await Service.GetAgent(person.agent_id).catch(() => null);
            setEditPersonData({
                uuid: person.uuid,
                name: person.name,
                role: person.role,
                avatar: person.avatar || '',
                prompt: agentDetail?.prompts?.[0]?.content || '',
                tools: agentDetail?.tools || [],
                skills: agentDetail?.skills || [],
            });
            setAddPersonOpen(true);
        } catch (err) {
            console.error('Failed to load person for edit:', err);
        }
    };

    const handleEditGroup = async (uuid: string) => {
        try {
            const group = await Service.OPCGetGroup(uuid);
            if (!group) return;
            setEditGroupData({
                uuid: group.uuid,
                name: group.name,
                description: group.description,
                member_uuids: group.members?.map((m: any) => m.uuid) || [],
            });
            setCreateGroupOpen(true);
        } catch (err) {
            console.error('Failed to load group for edit:', err);
        }
    };

    // 移动端默认隐藏侧边栏；桌面端强制展开（使用拖拽调整宽度）
    useEffect(() => {
        if (isMobile) {
            setIsSidebarCollapsed(true);
        } else {
            setIsSidebarCollapsed(false);
        }
    }, [isMobile]);

    // 侧边栏拖拽调整宽度（仅在结束时持久化，过程中用 rAF 同步更新）
    const handleSidebarResizeStart = useCallback((event: React.MouseEvent<HTMLDivElement>) => {
        if (isMobile) return;
        event.preventDefault();
        const startX = event.clientX;
        const startWidth = sidebarWidth;
        let latestWidth = startWidth;
        let rafId = 0;
        const onMove = (e: MouseEvent) => {
            latestWidth = Math.min(
                SIDEBAR_MAX_WIDTH,
                Math.max(SIDEBAR_MIN_WIDTH, startWidth + (e.clientX - startX))
            );
            if (rafId) return;
            rafId = requestAnimationFrame(() => {
                rafId = 0;
                setSidebarWidth(latestWidth);
            });
        };
        const onUp = () => {
            if (rafId) cancelAnimationFrame(rafId);
            setSidebarWidth(latestWidth);
            try {
                localStorage.setItem(SIDEBAR_WIDTH_KEY, String(latestWidth));
            } catch {
                // ignore storage errors
            }
            document.removeEventListener('mousemove', onMove);
            document.removeEventListener('mouseup', onUp);
            document.body.style.cursor = '';
            document.body.style.userSelect = '';
        };
        document.addEventListener('mousemove', onMove);
        document.addEventListener('mouseup', onUp);
        document.body.style.cursor = 'ew-resize';
        document.body.style.userSelect = 'none';
    }, [isMobile, sidebarWidth]);

    // 设置页面标题
    useEffect(() => {
        document.title = t('app.chatTitle');
    }, [t]);

    // 同步URL参数与当前聊天UUID
    useEffect(() => {
        const newChatUuid = urlChatUuid || '';
        if (newChatUuid !== currentChatUuid) {
            setCurrentChatUuid(newChatUuid);
        }
    }, [urlChatUuid, currentChatUuid]);

    // 历史聊天首次加载完成后，短暂保留标记供 MessageList 使用，然后重置
    useEffect(() => {
        if (isFirstHistoricalLoad) {
            const timer = setTimeout(() => setIsFirstHistoricalLoad(false), 200);
            return () => clearTimeout(timer);
        }
    }, [isFirstHistoricalLoad]);

    // handleToggleSidebar 展示/隐藏侧边菜单
    const handleToggleSidebar = () => {
        setIsSidebarCollapsed(!isSidebarCollapsed);
    };

    // 处理新建对话
    const handleNewChat = useCallback(() => {
        setCurrentChatUuid(''); // 设置为空字符串表示新对话
        // 更新URL为新对话状态
        navigate('/home', {replace: true});
        // 移动端新建对话后自动隐藏侧边栏
        if (isMobile) {
            setIsSidebarCollapsed(true);
        }
    }, [isMobile, navigate]);

    // 处理对话选择
    const handleChatSelect = useCallback(
        (chatUuid: string) => {
            setCurrentChatUuid(chatUuid);
            // 更新URL但不刷新页面
            navigate(`/home/${chatUuid}`, {replace: true});
            // 移动端选择对话后自动隐藏侧边栏
            if (isMobile) {
                setIsSidebarCollapsed(true);
            }
        },
        [isMobile, navigate]
    );

    // 处理删除聊天
    const handleDeleteChat = useCallback(
        async (chatUuid: string) => {
            try {
                await Service.DeleteChat(chatUuid);
                // 如果删除的是当前聊天，导航到新聊天页面
                if (chatUuid === currentChatUuid) {
                    handleNewChat();
                }
                // 刷新聊天列表
                if (refreshChatList) {
                    refreshChatList();
                }
                message.success(t('home.chat.deleteSuccess'));
            } catch (error) {
                console.error('删除聊天失败:', error);
                message.error(t('home.chat.deleteFailed'));
            }
        },
        [currentChatUuid, handleNewChat, refreshChatList, t]
    );

    // 设置刷新聊天列表的回调
    const handleSetRefreshChatList = useCallback((refreshFn: () => void) => {
        setRefreshChatList(() => refreshFn);
    }, []);

    const handleRegisterStopGenerationForChat = useCallback(
        (fn: (chatUuid: string) => void) => {
            stopGenerationForChatRef.current = fn;
        },
        [],
    );

    // OPC 模式
    if (appMode === 'opc') {
        return (
            <Layout className={`${className || ''} ${styles.chatLayout}`}>
                <Sider
                    className={`${styles.sidebar} ${isSidebarCollapsed ? styles.collapsed : ''}`}
                    width={sidebarWidth}
                    collapsedWidth={0}
                    collapsed={isMobile ? isSidebarCollapsed : false}
                    trigger={null}
                    collapsible
                >
                    <OPCSidebar
                        isSidebarCollapsed={isSidebarCollapsed}
                        onToggleSidebar={handleToggleSidebar}
                        onResizeStart={handleSidebarResizeStart}
                        onOpenAddPerson={() => { setEditPersonData(null); setAddPersonOpen(true); }}
                        onOpenCreateGroup={() => { setEditGroupData(null); setCreateGroupOpen(true); }}
                        onEditPerson={handleEditPerson}
                        onEditGroup={handleEditGroup}
                    />
                </Sider>
                <Layout className={styles.mainLayout}>
                    <Content className={styles.mainContent}>
                        <OPCChatArea
                            isSidebarCollapsed={isSidebarCollapsed}
                            onToggleSidebar={handleToggleSidebar}
                        />
                    </Content>
                </Layout>
                <AddPersonDialog
                    open={addPersonOpen}
                    onClose={() => { setAddPersonOpen(false); setEditPersonData(null); }}
                    onSuccess={() => window.dispatchEvent(new CustomEvent('opc-refresh'))}
                    editData={editPersonData}
                />
                <CreateGroupDialog
                    open={createGroupOpen}
                    onClose={() => { setCreateGroupOpen(false); setEditGroupData(null); }}
                    onSuccess={() => window.dispatchEvent(new CustomEvent('opc-refresh'))}
                    editData={editGroupData}
                />
            </Layout>
        );
    }

    // 聊天模式（原有逻辑）
    return (
        <Layout className={`${className || ''} ${styles.chatLayout}`}>
            <Sider
                className={`${styles.sidebar} ${
                    isSidebarCollapsed ? styles.collapsed : ''
                }`}
                width={sidebarWidth}
                collapsedWidth={0}
                collapsed={isMobile ? isSidebarCollapsed : false}
                trigger={null}
                collapsible
            >
                <Index
                    onNewChat={handleNewChat}
                    onChatSelect={handleChatSelect}
                    onRegisterRefreshCallback={handleSetRefreshChatList}
                    onDeleteChat={handleDeleteChat}
                    currentChatUuid={currentChatUuid}
                    isSidebarCollapsed={isSidebarCollapsed}
                    onToggleSidebar={handleToggleSidebar}
                    onResizeStart={handleSidebarResizeStart}
                    generatingChatUuids={generatingChatUuids}
                    onStopGenerationForChat={(uuid) =>
                        stopGenerationForChatRef.current(uuid)
                    }
                />
            </Sider>
            <Layout className={styles.mainLayout}>
                <Content className={styles.mainContent} hidden={isMobile && !isSidebarCollapsed}>
                    <Chat
                        key={currentChatUuid || 'new-chat'}
                        chatUuid={currentChatUuid}
                        isSidebarCollapsed={isSidebarCollapsed}
                        onToggleSidebar={handleToggleSidebar}
                        refreshChatList={refreshChatList}
                        onChatChange={setCurrentChatUuid}
                        onGeneratingUuidsChange={setGeneratingChatUuids}
                        onRegisterStopGenerationForChat={handleRegisterStopGenerationForChat}
                    />
                </Content>
            </Layout>
            <BackTop/>
        </Layout>
    );
};

export default ChatPage;
