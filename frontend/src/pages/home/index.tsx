import React, {useCallback, useEffect, useRef, useState} from 'react';
import {BackTop, Layout, message} from 'antd';
import {useNavigate, useParams} from 'react-router-dom';
import Index from './sidebar';
import {useViewportHeight} from '@/hooks/useViewportHeight';
import './index.module.scss';
import Chat from '@/pages/home/chat';
import {Service} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/index.ts";
import styles from './index.module.scss';

const {Content, Sider} = Layout;

interface ChatPageProps {
    className?: string;
}

const ChatPage: React.FC<ChatPageProps> = ({className}) => {
    // 获取路由参数和导航函数
    const {chatUuid: urlChatUuid} = useParams<{ chatUuid?: string }>();
    const navigate = useNavigate();
    // 本地状态管理
    const [currentChatUuid, setCurrentChatUuid] = useState<string>(urlChatUuid ?? '');
    const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false);
    const [refreshChatList, setRefreshChatList] = useState<(() => void) | null>(
        null
    );
    const [generatingChatUuids, setGeneratingChatUuids] = useState<string[]>([]);
    const stopGenerationForChatRef = useRef<(uuid: string) => void>(() => {});
    // 历史聊天首次加载完成（用于立即滚动到底部，无动画）
    const [isFirstHistoricalLoad, setIsFirstHistoricalLoad] = useState(false);
    // 使用视口高度检测 Hook
    const {isMobile} = useViewportHeight();

    // 移动端默认隐藏侧边栏
    useEffect(() => {
        if (isMobile) {
            setIsSidebarCollapsed(true);
        } else {
            // Safari内核兼容性：从移动端切换回桌面端时，需要强制重置transform属性
            // 添加延迟重新渲染机制，确保Safari正确应用新的CSS规则
            const timer = setTimeout(() => {
                // 强制触发组件重新渲染
                setIsSidebarCollapsed(prev => prev);
            }, 100);

            return () => clearTimeout(timer);
        }
    }, [isMobile]);

    // 设置页面标题
    useEffect(() => {
        document.title = 'AI聊天 - Lemon Tea';
    }, []);

    useEffect(() => {
        console.log('监听页面参数 chatUuid 变化:', urlChatUuid);
    }, [urlChatUuid]); // 👈 关键依赖：仅当 chatUuid 改变时执行

    // 同步URL参数与当前聊天UUID
    useEffect(() => {
        console.log("xxx",urlChatUuid, currentChatUuid)
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
                message.success('聊天删除成功');
            } catch (error) {
                console.error('删除聊天失败:', error);
                message.error('删除聊天失败');
            }
        },
        [currentChatUuid, handleNewChat, refreshChatList]
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

    return (
        <Layout className={`${className || ''} ${styles.chatLayout}`}>
            <Sider
                className={`${styles.sidebar} ${
                    isSidebarCollapsed ? styles.collapsed : ''
                }`}
                width={280}
                collapsedWidth={isMobile ? 0 : 50}
                collapsed={isSidebarCollapsed}
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
                    generatingChatUuids={generatingChatUuids}
                    onStopGenerationForChat={(uuid) =>
                        stopGenerationForChatRef.current(uuid)
                    }
                />
            </Sider>
            <Layout className={styles.mainLayout}>
                <Content className={styles.mainContent} hidden={isMobile && !isSidebarCollapsed}>
                    <Chat
                        chatUuid={currentChatUuid}
                        isSidebarCollapsed={isSidebarCollapsed}
                        onToggleSidebar={handleToggleSidebar}
                        refreshChatList={refreshChatList}
                        onChatChange={(chatUuid)=>{
                            console.log("setCurrentChatUuid",chatUuid)
                            setCurrentChatUuid(chatUuid)
                        }}
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