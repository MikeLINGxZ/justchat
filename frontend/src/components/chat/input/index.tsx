import React, {useCallback, useEffect, useRef, useState} from "react";
import styles from "./index.module.scss";
import {useIsMobile} from "@/hooks/useViewportHeight.ts";
import {Service} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service";
import {FileInfo, Model, Tool} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models";
import {notify} from "@/utils/notification.ts";
import RichMarkdownEditor from "@/components/chat/input/rich_markdown_editor";

interface ChatInputProps {
    // 所选模型
    selectedModelId: number;
    // 是否已选择模型
    hasSelectedModel: boolean;
    // 可用模型
    availableModels: Model[];
    // 可用工具
    availableTools: Tool[];
    // 已选中的工具 id 列表
    selectedToolIds: string[];
    // 工具选择变更事件
    onSelectedToolsChange: (toolIds: string[]) => void;
    // 刷新工具列表
    onRefreshTools: () => Promise<Tool[]>;
    // 是否正在生成消息
    isGenerating: boolean;
    // 模型变更事件
    onSelectModelChange: (modelId: number, modelName: string) => void;
    // 输入内容变更事件
    onMessageChange: (message: string) => void;
    // 文件变更事件
    onSelectFileChange: (paths: FileInfo[]) => void;
    // 发送按钮点击事件
    onSendButtonClick: () => void;
    // 点击停止生成按钮事件
    onStopGeneration: () => void;
    // 消息列表滚动到底部按钮点击事件
    onMessageListScrollToBottom?: () => void;
    // 模型选择框点击事件
    onModelSelectorClick?: () => void;
    // 当前审批回复上下文
    approvalInput?: { approvalId: string; title: string; message: string } | null;
    // 发送审批意见
    onSendApprovalComment?: (approvalId: string, comment: string) => Promise<void> | void;
    // 取消审批意见输入
    onCancelApprovalComment?: () => void;
}

const DEFAULT_MODEL_KEY = 'chat_default_model';

interface DefaultModelConfig {
    modelId: number;
    modelName: string;
}

function getDefaultModelConfig(): DefaultModelConfig | null {
    try {
        const raw = localStorage.getItem(DEFAULT_MODEL_KEY);
        if (!raw) return null;
        return JSON.parse(raw) as DefaultModelConfig;
    } catch {
        return null;
    }
}

function setDefaultModelConfig(config: DefaultModelConfig) {
    localStorage.setItem(DEFAULT_MODEL_KEY, JSON.stringify(config));
}

const ChatInput: React.FC<ChatInputProps> = ({
    selectedModelId,
    hasSelectedModel,
    availableModels,
    availableTools,
    selectedToolIds,
    onSelectedToolsChange,
    onRefreshTools,
    isGenerating = false,
    onMessageChange,
    onSendButtonClick,
    onSelectFileChange,
    onStopGeneration,
    onSelectModelChange,
    onMessageListScrollToBottom,
    onModelSelectorClick,
    approvalInput,
    onSendApprovalComment,
    onCancelApprovalComment,
}) => {

    const [showAddMenu, setShowAddMenu] = useState(false);
    const [showModelMenu, setShowModelMenu] = useState(false);
    const [showToolMenu, setShowToolMenu] = useState(false);
    const [inputValue, setInputValue] = useState('');
    const [modelSearchValue, setModelSearchValue] = useState('');
    const [isAddingMCPTool, setIsAddingMCPTool] = useState(false);
    const [defaultModelId, setDefaultModelId] = useState<number | null>(() => {
        return getDefaultModelConfig()?.modelId ?? null;
    });
    const addMenuRef = useRef<HTMLDivElement>(null);
    const modelMenuRef = useRef<HTMLDivElement>(null);
    const toolMenuRef = useRef<HTMLDivElement>(null);
    const isMobile =  useIsMobile();
    const [selectFiles, setSelectFiles] = useState<FileInfo[]>([]);

    // 清空输入框
    const clearInput = useCallback(() => {
        setInputValue('');
        setSelectFiles([]);
        onSelectFileChange([]);
        onMessageChange('');
    }, [onMessageChange, onSelectFileChange]);

    // 添加按钮点击事件
    const handleAddClick = useCallback(() => {
        setShowAddMenu(!showAddMenu);
    }, [showAddMenu]);

    // 模型选择框点击事件
    const handleModelClick = useCallback(() => {
        const willOpen = !showModelMenu;
        setShowModelMenu(willOpen);
        // 打开菜单时清空搜索值并刷新模型数据
        if (willOpen) {
            setModelSearchValue('');
            // 调用回调刷新模型数据
            onModelSelectorClick?.();
        }
    }, [showModelMenu, onModelSelectorClick]);

    // Tool 选择框点击事件
    const handleToolClick = useCallback(() => {
        setShowToolMenu(!showToolMenu);
    }, [showToolMenu]);

    // Tool 开关切换
    const handleToolToggle = useCallback((toolId: string, enabled: boolean) => {
        if (enabled) {
            onSelectedToolsChange([...selectedToolIds, toolId]);
        } else {
            onSelectedToolsChange(selectedToolIds.filter(id => id !== toolId));
        }
    }, [selectedToolIds, onSelectedToolsChange]);

    const handleCustomToolToggle = useCallback(async (tool: Tool, enabled: boolean) => {
        try {
            await Service.UpdateMCPToolEnabled(tool.id, enabled);
            await onRefreshTools();
            if (enabled) {
                onSelectedToolsChange([...new Set([...selectedToolIds, tool.id])]);
            } else {
                onSelectedToolsChange(selectedToolIds.filter(id => id !== tool.id));
            }
        } catch (error) {
            notify.error("更新失败", `无法${enabled ? '启用' : '禁用'} ${tool.name}`);
        }
    }, [onRefreshTools, onSelectedToolsChange, selectedToolIds]);

    const handleAddMCPTool = useCallback(async () => {
        if (isAddingMCPTool) {
            return;
        }
        setIsAddingMCPTool(true);
        try {
            const folderPath = await Service.SelectMCPFolder();
            if (!folderPath) {
                return;
            }
            const createdTool = await Service.AddMCPToolFromFolder(folderPath);
            if (!createdTool) {
                return;
            }
            await onRefreshTools();
            onSelectedToolsChange([...new Set([...selectedToolIds, createdTool.id])]);
            notify.success("添加成功", `${createdTool.name} 已加入工具列表`);
        } catch (error: any) {
            notify.error("添加失败", error?.message || "所选目录不是有效的 MCP 服务目录");
        } finally {
            setIsAddingMCPTool(false);
        }
    }, [isAddingMCPTool, onRefreshTools, onSelectedToolsChange, selectedToolIds]);

    const handleDeleteMCPTool = useCallback(async (tool: Tool) => {
        try {
            await Service.DeleteMCPTool(tool.id);
            await onRefreshTools();
            onSelectedToolsChange(selectedToolIds.filter(id => id !== tool.id));
            notify.success("删除成功", `${tool.name} 已从工具列表移除`);
        } catch (error) {
            notify.error("删除失败", `无法删除 ${tool.name}`);
        }
    }, [onRefreshTools, onSelectedToolsChange, selectedToolIds]);

    // 文件上传事件
    const handleFileUpload = useCallback(() => {
        setShowAddMenu(false);
        Service.SelectFiles().then(async (files: FileInfo[]) => {
            if (files.length === 0) return;
            setSelectFiles(prevFiles => {
                const mergedFiles = [...prevFiles];
                files.forEach(file => {
                    if (!mergedFiles.some(item => item.path === file.path)) {
                        mergedFiles.push(file);
                    }
                });
                onSelectFileChange(mergedFiles);
                return mergedFiles;
            });
        }).catch(() => {
        })
    }, [onSelectFileChange]);

    // 删除文件
    const handleRemoveFile = useCallback((filePath: string) => {
        setSelectFiles(prevFiles => {
            const nextFiles = prevFiles.filter(f => f.path !== filePath);
            onSelectFileChange(nextFiles);
            return nextFiles;
        });
    }, [onSelectFileChange]);

    // 初始化时自动选中默认模型
    useEffect(() => {
        if (selectedModelId > 0 || availableModels.length === 0) return;
        const config = getDefaultModelConfig();
        if (!config) return;
        const exists = availableModels.some(m => m.id === config.modelId);
        if (exists) {
            onSelectModelChange(config.modelId, config.modelName);
        }
    }, [availableModels]);

    // 模型选择事件
    const handleModelSelect = useCallback((modelId: number, modelName: string) => {
        onSelectModelChange(modelId,modelName);
        setShowModelMenu(false);
        setModelSearchValue('');
    }, [onSelectModelChange]);

    // 设为/取消默认模型
    const handleSetDefaultModel = useCallback((e: React.MouseEvent, modelId: number, modelName: string) => {
        e.stopPropagation();
        if (defaultModelId === modelId) {
            localStorage.removeItem(DEFAULT_MODEL_KEY);
            setDefaultModelId(null);
        } else {
            setDefaultModelConfig({ modelId, modelName });
            setDefaultModelId(modelId);
        }
    }, [defaultModelId]);

    // 处理模型搜索
    const handleModelSearch = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
        setModelSearchValue(e.target.value);
    }, []);

    // 过滤模型列表
    const filteredModels = useCallback(() => {
        if (!modelSearchValue.trim()) {
            return availableModels;
        }
        return availableModels.filter(model => 
            model.model.toLowerCase().includes(modelSearchValue.toLowerCase()) ||
            model.alias?.toLowerCase().includes(modelSearchValue.toLowerCase())
        );
    }, [availableModels, modelSearchValue]);

    // 消息发送事件
    const handleSend = useCallback(async () => {
        if (!approvalInput && !hasSelectedModel) {
            return;
        }
        if (onMessageListScrollToBottom != null) {
            onMessageListScrollToBottom();
        }
        const trimmedValue = inputValue.trim();
        if (approvalInput) {
            if (!trimmedValue) {
                return;
            }
            await onSendApprovalComment?.(approvalInput.approvalId, trimmedValue);
            clearInput();
            return;
        }
        if (trimmedValue || selectFiles.length > 0) {
            onSendButtonClick();
            clearInput(); // 清空输入框和文件列表
        }
    }, [approvalInput, hasSelectedModel, inputValue, selectFiles, onSendApprovalComment, onSendButtonClick, onMessageListScrollToBottom, clearInput]);

    const handleInputChange = useCallback((value: string) => {
        setInputValue(value);
        onMessageChange(value);
    }, [onMessageChange]);

    // 处理点击外部区域关闭菜单
    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (addMenuRef.current && !addMenuRef.current.contains(event.target as Node)) {
                setShowAddMenu(false);
            }
            if (modelMenuRef.current && !modelMenuRef.current.contains(event.target as Node)) {
                setShowModelMenu(false);
                setModelSearchValue(''); // 关闭菜单时清空搜索值
            }
            if (toolMenuRef.current && !toolMenuRef.current.contains(event.target as Node)) {
                setShowToolMenu(false);
            }
        };

        if (showAddMenu || showModelMenu || showToolMenu) {
            document.addEventListener('mousedown', handleClickOutside);
        }

        return () => {
            document.removeEventListener('mousedown', handleClickOutside);
        };
    }, [showAddMenu, showModelMenu, showToolMenu]);

    const isSendDisabled = approvalInput
        ? !inputValue.trim()
        : !hasSelectedModel || (!inputValue.trim() && selectFiles.length === 0);

    return (
        <div className={`${styles.chatInput}`}>
            {!approvalInput && !hasSelectedModel && (
                <div className={styles.modelWarning}>请先选择模型</div>
            )}
            {approvalInput && (
                <div className={styles.approvalNotice}>
                    <div className={styles.approvalNoticeText}>
                        <div className={styles.approvalNoticeTitle}>正在回复审批意见：{approvalInput.title}</div>
                        <div className={styles.approvalNoticeBody}>{approvalInput.message}</div>
                    </div>
                    <button
                        type="button"
                        className={styles.approvalNoticeClose}
                        onClick={onCancelApprovalComment}
                    >
                        取消
                    </button>
                </div>
            )}
            <div className={styles.inputContainer}>
                 {/* 文件列表显示区域 */}
                {selectFiles.length > 0 && (
                    <div className={styles.filesContainer}>
                        {selectFiles.map((file) => (
                            <div key={file.path} className={styles.fileItem}>
                                <div className={styles.filePreview}>
                                    {file.preview ? (
                                        <img
                                            src={file.preview}
                                            alt={file.name}
                                            className={styles.fileImage}
                                        />
                                    ) : (
                                        <div className={styles.fileIcon}>
                                            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                                <path d="M14 2H6C4.89 2 4 2.89 4 4V20C4 21.11 4.89 22 6 22H18C19.11 22 20 21.11 20 20V8L14 2ZM18 20H6V4H13V9H18V20Z" fill="currentColor"/>
                                            </svg>
                                        </div>
                                    )}
                                    <button
                                        className={styles.fileRemove}
                                        onClick={() => handleRemoveFile(file.path)}
                                        type="button"
                                        title="删除文件"
                                    >
                                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                            <path d="M19 6.41L17.59 5L12 10.59L6.41 5L5 6.41L10.59 12L5 17.59L6.41 19L12 13.41L17.59 19L19 17.59L13.41 12L19 6.41Z" fill="currentColor"/>
                                        </svg>
                                    </button>
                                </div>
                                <div className={styles.fileName} title={file.name}>
                                    {file.name}
                                </div>
                            </div>
                        ))}
                    </div>
                )}
                <div className={styles.richEditorWrapper}>
                    <RichMarkdownEditor
                        value={inputValue}
                        onChange={handleInputChange}
                        onSend={handleSend}
                        placeholder={approvalInput ? "输入你对这次工具请求的意见..." : "输入消息... (支持 Markdown 格式)"}
                    />
                </div>
                
                <div className={styles.bottomBar}>
                    <div className={styles.leftActions}>
                        <div className={styles.addButtonContainer} ref={addMenuRef}>
                            <button 
                                className={styles.addButton}
                                onClick={handleAddClick}
                                type="button"
                            >
                                <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                                    <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
                                </svg>
                            </button>
                            
                            {showAddMenu && (
                                <div className={`${styles.addMenu} ${isMobile ? styles.mobileMenu : ''}`}>
                                    <button 
                                        className={styles.menuItem}
                                        onClick={handleFileUpload}
                                        type="button"
                                    >
                                        <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                                            <path d="M14,2H6A2,2 0 0,0 4,4V20A2,2 0 0,0 6,22H18A2,2 0 0,0 20,20V8L14,2M18,20H6V4H13V9H18V20Z"/>
                                        </svg>
                                        上传文件
                                    </button>
                                </div>
                            )}
                        </div>

                        <div className={styles.toolSelector} ref={toolMenuRef}>
                            <button
                                className={styles.toolButton}
                                onClick={handleToolClick}
                                type="button"
                                title="选择工具"
                            >
                                <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                                    <path d="M22.7 19l-9.1-9.1c.9-2.3.4-5-1.5-6.9-2-2-5-2.4-7.4-1.3L9 6 6 9 1.6 4.7C.4 7.1.9 10.1 2.9 12.1c1.9 1.9 4.6 2.4 6.9 1.5l9.1 9.1c.4.4 1 .4 1.4 0l2.3-2.3c.5-.4.5-1.1.1-1.4z"/>
                                </svg>
                                <span className={styles.toolButtonText}>Tool</span>
                                {selectedToolIds.length > 0 && (
                                    <span className={styles.toolBadge}>{selectedToolIds.length}</span>
                                )}
                                <svg width="12" height="12" viewBox="0 0 24 24" fill="currentColor">
                                    <path d="M7 10l5 5 5-5z"/>
                                </svg>
                            </button>

                            {showToolMenu && (
                                <div className={`${styles.toolMenu} ${isMobile ? styles.mobileMenu : ''}`}>
                                    <div className={styles.toolList}>
                                        {availableTools.length === 0 ? (
                                            <div className={styles.noResults}>暂无可用工具</div>
                                        ) : (
                                            availableTools.map((tool) => {
                                                const isCustomMCP = tool.source_type === 'mcp_custom';
                                                const isEnabled = isCustomMCP ? tool.enabled : selectedToolIds.includes(tool.id);
                                                return (
                                                    <div
                                                        key={tool.id}
                                                        className={styles.toolItem}
                                                    >
                                                        <div className={styles.toolItemInfo}>
                                                            <div className={styles.toolItemHeader}>
                                                                <span className={styles.toolItemName}>{tool.name}</span>
                                                                {isCustomMCP && (
                                                                    <span className={styles.toolSourceTag}>MCP</span>
                                                                )}
                                                            </div>
                                                            {tool.description && (
                                                                <span
                                                                    className={styles.toolItemDesc}
                                                                    title={tool.description}
                                                                >
                                                                    {tool.description}
                                                                </span>
                                                            )}
                                                        </div>
                                                        <div className={styles.toolItemActions}>
                                                            {tool.is_deletable && (
                                                                <button
                                                                    className={styles.deleteToolButton}
                                                                    type="button"
                                                                    title={`删除 ${tool.name}`}
                                                                    onClick={() => {
                                                                        void handleDeleteMCPTool(tool);
                                                                    }}
                                                                >
                                                                    <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
                                                                        <path d="M6 7h12l-1 14H7L6 7zm3-3h6l1 2h4v2H4V6h4l1-2z"/>
                                                                    </svg>
                                                                </button>
                                                            )}
                                                            <label className={styles.toggleSwitch}>
                                                                <input
                                                                    type="checkbox"
                                                                    checked={isEnabled}
                                                                    onChange={(e) => {
                                                                        if (isCustomMCP) {
                                                                            void handleCustomToolToggle(tool, e.target.checked);
                                                                            return;
                                                                        }
                                                                        handleToolToggle(tool.id, e.target.checked);
                                                                    }}
                                                                />
                                                                <span className={styles.toggleSlider} />
                                                            </label>
                                                        </div>
                                                    </div>
                                                );
                                            })
                                        )}
                                    </div>
                                    <div className={styles.toolMenuFooter}>
                                        <button
                                            className={styles.addMcpButton}
                                            onClick={() => {
                                                void handleAddMCPTool();
                                            }}
                                            type="button"
                                            disabled={isAddingMCPTool}
                                        >
                                            {isAddingMCPTool ? (
                                                <>
                                                    <span className={styles.buttonSpinner} />
                                                    添加中...
                                                </>
                                            ) : (
                                                <>
                                                    <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
                                                        <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
                                                    </svg>
                                                    添加 MCP 服务
                                                </>
                                            )}
                                        </button>
                                    </div>
                                </div>
                            )}
                        </div>
                    </div>
                    
                    <div className={styles.rightActions}>
                        <div className={styles.modelSelector} ref={modelMenuRef}>
                            <button 
                                className={styles.modelButton}
                                onClick={handleModelClick}
                                type="button"
                            >
                                {availableModels.find(model => model.id === selectedModelId)?.model  || '选择模型'}
                                <svg width="12" height="12" viewBox="0 0 24 24" fill="currentColor">
                                    <path d="M7 10l5 5 5-5z"/>
                                </svg>
                            </button>
                            
                            {showModelMenu && (
                                <div className={`${styles.modelMenu} ${isMobile ? styles.mobileMenu : ''}`}>
                                    {/* 模型列表 */}
                                    <div className={styles.modelList}>
                                        {filteredModels().map((model) => (
                                            <button
                                                key={model.id}
                                                className={`${styles.menuItem} ${model.id === selectedModelId ? styles.selected : ''}`}
                                                onClick={() => handleModelSelect(model.id,model.model)}
                                                type="button"
                                            >
                                                <span className={styles.modelName}>
                                                    {model.model}
                                                    {model.id === defaultModelId && (
                                                        <span className={styles.defaultBadge}>默认</span>
                                                    )}
                                                </span>
                                                <span className={styles.modelItemRight}>
                                                    {model.alias != null && (
                                                        <span className={styles.modelId}>{model.alias}</span>
                                                    )}
                                                    <span
                                                        className={styles.setDefaultBtn}
                                                        onClick={(e) => handleSetDefaultModel(e, model.id, model.model)}
                                                    >
                                                        {model.id === defaultModelId ? '取消默认' : '设为默认'}
                                                    </span>
                                                </span>
                                            </button>
                                        ))}
                                        {/* 无结果提示 */}
                                        {filteredModels().length === 0 && (
                                            <div className={styles.noResults}>
                                                没有找到匹配的模型
                                            </div>
                                        )}
                                    </div>
                                    {/* 搜索输入框 */}
                                    <div className={styles.searchContainer}>
                                        <input
                                            type="text"
                                            placeholder="搜索模型..."
                                            value={modelSearchValue}
                                            onChange={handleModelSearch}
                                            className={styles.searchInput}
                                            autoFocus
                                        />
                                    </div>
                                </div>
                            )}
                        </div>
                        
                        {isGenerating ? (
                            <button 
                                className={`${styles.sendButton} ${styles.stopButton}`}
                                onClick={onStopGeneration}
                                type="button"
                                title="停止生成"
                            >
                                <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                                    <rect x="6" y="6" width="12" height="12" rx="2"/>
                                </svg>
                            </button>
                        ) : (
                            <button 
                                className={`${styles.sendButton} ${isSendDisabled ? styles.disabled : ''}`}
                                onClick={handleSend}
                                disabled={isSendDisabled}
                                type="button"
                                title={approvalInput ? "发送审批意见" : (hasSelectedModel ? "发送消息" : "请先选择模型")}
                            >
                                <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                                    <path d="M2.01 21L23 12 2.01 3 2 10l15 2-15 2z"/>
                                </svg>
                            </button>
                        )}
                    </div>
                </div>
            </div>

        </div>
    );
};

export default ChatInput;
