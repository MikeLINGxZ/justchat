import React, { useEffect } from 'react';
import {
    PlusOutlined,
    SearchOutlined,
    TeamOutlined,
    UserAddOutlined,
} from '@ant-design/icons';
import { Dropdown, Input, message } from 'antd';

const { Search } = Input;
import type { MenuProps } from 'antd';
import { useTranslation } from 'react-i18next';
import { useOPCStore } from '@/stores/opcStore';
import type { OPCPersonView, OPCGroupView } from '@/stores/opcStore';
import { Service } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service';
import SidebarHeader from '@/components/sidebar/SidebarHeader';
import SidebarUserMenu from '@/components/sidebar/SidebarUserMenu';
import ContactList from './contact_list';
import './index.scss';

interface OPCSidebarProps {
    className?: string;
    isSidebarCollapsed: boolean;
    onToggleSidebar: () => void;
    onOpenAddPerson: () => void;
    onOpenCreateGroup: () => void;
    onEditPerson?: (uuid: string) => void;
    onEditGroup?: (uuid: string) => void;
}

const OPCSidebar: React.FC<OPCSidebarProps> = ({
    className,
    isSidebarCollapsed,
    onToggleSidebar,
    onOpenAddPerson,
    onOpenCreateGroup,
    onEditPerson,
    onEditGroup,
}) => {
    const { t } = useTranslation();
    const {
        mode, persons, groups, searchQuery, setSearchQuery,
        setPersons, setGroups, selectedUuid, setSelected,
        sidebarTab, setSidebarTab,
    } = useOPCStore();

    const fetchData = async () => {
        try {
            const [personList, groupList] = await Promise.all([
                Service.OPCListPersons(),
                Service.OPCListGroups(),
            ]);
            setPersons(personList || []);
            setGroups(groupList || []);
        } catch (err) {
            console.error('Failed to load OPC data:', err);
        }
    };

    useEffect(() => {
        fetchData();
        const handleRefresh = () => fetchData();
        window.addEventListener('opc-refresh', handleRefresh);
        return () => window.removeEventListener('opc-refresh', handleRefresh);
    }, []);

    const handleSelectContact = (type: 'person' | 'group', uuid: string) => {
        setSelected(type, uuid);
    };

    const handleDeletePerson = async (uuid: string) => {
        try {
            await Service.OPCDeletePerson(uuid);
            message.success(t('opc.person.deleteSuccess'));
            fetchData();
        } catch (err) {
            message.error(t('opc.person.deleteFailed'));
        }
    };

    const handleDeleteGroup = async (uuid: string) => {
        try {
            await Service.OPCDeleteGroup(uuid);
            message.success(t('opc.group.deleteSuccess'));
            fetchData();
        } catch (err) {
            message.error(t('opc.group.deleteFailed'));
        }
    };

    const handleTogglePin = async (type: 'person' | 'group', uuid: string, pinned: boolean) => {
        try {
            if (type === 'person') {
                await Service.OPCTogglePinPerson(uuid, !pinned);
            } else {
                await Service.OPCTogglePinGroup(uuid, !pinned);
            }
            fetchData();
        } catch (err) {
            message.error(t('common.operationFailed'));
        }
    };

    const handleClearConversation = async (chatUuid: string) => {
        try {
            await Service.OPCClearConversation(chatUuid);
            message.success(t('opc.sidebar.clearSuccess'));
            fetchData();
        } catch (err) {
            message.error(t('common.operationFailed'));
        }
    };

    const addMenuItems: MenuProps['items'] = [
        {
            key: 'addPerson',
            icon: <UserAddOutlined />,
            label: t('opc.sidebar.addPerson'),
            onClick: onOpenAddPerson,
        },
        {
            key: 'createGroup',
            icon: <TeamOutlined />,
            label: t('opc.sidebar.createGroup'),
            onClick: onOpenCreateGroup,
        },
    ];

    return (
        <div className={`opc-sidebar ${isSidebarCollapsed ? 'collapsed' : ''} ${className || ''}`}>
            <SidebarHeader
                logoText="OPC"
                isSidebarCollapsed={isSidebarCollapsed}
                onToggleSidebar={onToggleSidebar}
            />

            {/* 功能按钮区域 */}
            {!isSidebarCollapsed && (
                <div className="sidebar-actions">
                    <Dropdown menu={{ items: addMenuItems }} trigger={['click']} placement="bottomRight">
                        <button className="action-btn">
                            <PlusOutlined className="action-icon" />
                            <span className="action-text">{t('opc.sidebar.addNew')}</span>
                        </button>
                    </Dropdown>
                </div>
            )}

            {/* Tab 切换 */}
            {!isSidebarCollapsed && (
                <div className="sidebar-tabs">
                    <div className="tab-switch">
                        <div
                            className={`tab-option ${sidebarTab === 'conversations' ? 'active' : ''}`}
                            onClick={() => setSidebarTab('conversations')}
                        >
                            {t('opc.sidebar.tabConversations')}
                        </div>
                        <div
                            className={`tab-option ${sidebarTab === 'contacts' ? 'active' : ''}`}
                            onClick={() => setSidebarTab('contacts')}
                        >
                            {t('opc.sidebar.tabContacts')}
                        </div>
                        <div className="tab-slider" data-active={sidebarTab}></div>
                    </div>
                </div>
            )}

            {/* 搜索区域 */}
            {!isSidebarCollapsed && (
                <div className="opc-search-area">
                    <Search
                        prefix={<SearchOutlined />}
                        placeholder={
                            sidebarTab === 'conversations'
                                ? t('opc.sidebar.searchConversations')
                                : t('opc.sidebar.searchContacts')
                        }
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        allowClear
                    />
                </div>
            )}

            {/* 联系人列表 */}
            {!isSidebarCollapsed && (
                <div className="opc-contact-area">
                    <ContactList
                        persons={persons}
                        groups={groups}
                        searchQuery={searchQuery}
                        selectedUuid={selectedUuid}
                        activeTab={sidebarTab}
                        onSelect={handleSelectContact}
                        onDeletePerson={handleDeletePerson}
                        onDeleteGroup={handleDeleteGroup}
                        onTogglePin={handleTogglePin}
                        onEditPerson={onEditPerson}
                        onEditGroup={onEditGroup}
                        onClearConversation={handleClearConversation}
                    />
                </div>
            )}

            <div className="sidebar-spacer"></div>

            <SidebarUserMenu
                isSidebarCollapsed={isSidebarCollapsed}
                currentMode={mode}
            />
        </div>
    );
};

export default OPCSidebar;
