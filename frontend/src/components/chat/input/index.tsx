import React, {useCallback, useEffect, useRef, useState} from "react";
import { useTranslation } from 'react-i18next';
import { Events } from '@wailsio/runtime';
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
    // 初始输入内容（用于欢迎页建议卡片预填）
    initialValue?: string;
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
    initialValue = '',
    onMessageChange,
    onSendButtonClick,
    onSelectFileChange,
    onStopGeneration,
    onSelectModelChange,
    onMessageListScrollToBottom,
    onModelSelectorClick,
}) => {
    const { t } = useTranslation();

    const [showAddMenu, setShowAddMenu] = useState(false);
    const [showModelMenu, setShowModelMenu] = useState(false);
    const [showToolMenu, setShowToolMenu] = useState(false);
    const [inputValue, setInputValue] = useState(initialValue);
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
    const [isDragOver, setIsDragOver] = useState(false);

    // 监听文件拖拽事件
    useEffect(() => {
        const cancel = Events.On("files-dropped", (event: any) => {
            const raw = event.data;
            const files: FileInfo[] = Array.isArray(raw) ? raw : [];
            if (files.length === 0) return;
            setSelectFiles(prevFiles => {
                const mergedFiles = [...prevFiles];
                files.forEach((file: FileInfo) => {
                    if (!mergedFiles.some(item => item.path === file.path)) {
                        mergedFiles.push(file);
                    }
                });
                onSelectFileChange(mergedFiles);
                return mergedFiles;
            });
        });
        return () => {
            if (cancel) cancel();
        };
    }, [onSelectFileChange]);

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
            notify.error(
                t('chat.input.updateFailed'),
                t('chat.input.updateFailedDesc', {
                    action: enabled ? t('chat.input.toolActions.enable') : t('chat.input.toolActions.disable'),
                    name: tool.name,
                }),
            );
        }
    }, [onRefreshTools, onSelectedToolsChange, selectedToolIds, t]);

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
            notify.success(t('chat.input.addSuccess'), t('chat.input.addSuccessDesc', { name: createdTool.name }));
        } catch (error: any) {
            notify.error(t('chat.input.addFailed'), error?.message || t('chat.input.addFailedDesc'));
        } finally {
            setIsAddingMCPTool(false);
        }
    }, [isAddingMCPTool, onRefreshTools, onSelectedToolsChange, selectedToolIds, t]);

    const handleDeleteMCPTool = useCallback(async (tool: Tool) => {
        try {
            await Service.DeleteMCPTool(tool.id);
            await onRefreshTools();
            onSelectedToolsChange(selectedToolIds.filter(id => id !== tool.id));
            notify.success(t('chat.input.deleteSuccess'), t('chat.input.deleteSuccessDesc', { name: tool.name }));
        } catch (error) {
            notify.error(t('chat.input.deleteFailed'), t('chat.input.deleteFailedDesc', { name: tool.name }));
        }
    }, [onRefreshTools, onSelectedToolsChange, selectedToolIds, t]);

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
        if (!hasSelectedModel) {
            return;
        }
        if (onMessageListScrollToBottom != null) {
            onMessageListScrollToBottom();
        }
        const trimmedValue = inputValue.trim();
        if (trimmedValue || selectFiles.length > 0) {
            onSendButtonClick();
            clearInput();
        }
    }, [hasSelectedModel, inputValue, selectFiles, onSendButtonClick, onMessageListScrollToBottom, clearInput]);

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

    const isSendDisabled = !hasSelectedModel || (!inputValue.trim() && selectFiles.length === 0);

    return (
        <div className={`${styles.chatInput}`}>
            {!hasSelectedModel && (
                <div className={styles.modelWarning}>{t('chat.input.selectModelFirst')}</div>
            )}
            <div
                className={`${styles.inputContainer} ${isDragOver ? styles.dragOver : ''}`}
                data-file-drop-target="true"
                onDragEnter={(e) => {
                    e.preventDefault();
                    setIsDragOver(true);
                }}
                onDragOver={(e) => {
                    e.preventDefault();
                    setIsDragOver(true);
                }}
                onDragLeave={(e) => {
                    // Only set false when leaving the container itself
                    if (!e.currentTarget.contains(e.relatedTarget as Node)) {
                        setIsDragOver(false);
                    }
                }}
                onDrop={() => {
                    setIsDragOver(false);
                }}
            >
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
                                        title={t('chat.input.removeFile')}
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
                        key={`default-${t('chat.input.editorPlaceholder')}`}
                        value={inputValue}
                        onChange={handleInputChange}
                        onSend={handleSend}
                        placeholder={t('chat.input.editorPlaceholder')}
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
                                        {t('chat.input.uploadFile')}
                                    </button>
                                </div>
                            )}
                        </div>

                        <div className={styles.toolSelector} ref={toolMenuRef}>
                            <button
                                className={styles.toolButton}
                                onClick={handleToolClick}
                                type="button"
                                title={t('chat.input.selectTool')}
                            >
                                <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                                    <path d="M22.7 19l-9.1-9.1c.9-2.3.4-5-1.5-6.9-2-2-5-2.4-7.4-1.3L9 6 6 9 1.6 4.7C.4 7.1.9 10.1 2.9 12.1c1.9 1.9 4.6 2.4 6.9 1.5l9.1 9.1c.4.4 1 .4 1.4 0l2.3-2.3c.5-.4.5-1.1.1-1.4z"/>
                                </svg>
                                <span className={styles.toolButtonText}>
                                    {selectedToolIds.length === 0
                                        ? t('chat.input.selectTool')
                                        : t('chat.input.toolsSelected', { count: selectedToolIds.length })}
                                </span>
                                <svg width="12" height="12" viewBox="0 0 24 24" fill="currentColor">
                                    <path d="M7 10l5 5 5-5z"/>
                                </svg>
                            </button>

                            {showToolMenu && (
                                <div className={`${styles.toolMenu} ${isMobile ? styles.mobileMenu : ''}`}>
                                    <div className={styles.toolList}>
                                        {availableTools.length === 0 ? (
                                            <div className={styles.noResults}>{t('chat.input.noTools')}</div>
                                        ) : (
                                            availableTools.map((tool) => {
                                                const isCustomMCP = tool.source_type === 'mcp_custom';
                                                const isPlugin = tool.source_type === 'plugin';
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
                                                                {isPlugin && (
                                                                    <span className={styles.toolSourceTag} title={tool.plugin_name || tool.id}>
                                                                        {t('chat.input.pluginTag', '插件')}
                                                                    </span>
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
                                                                        title={t('chat.input.deleteTool', { name: tool.name })}
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
                                                    {t('chat.input.addingMcp')}
                                                </>
                                            ) : (
                                                <>
                                                    <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
                                                        <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
                                                    </svg>
                                                    {t('chat.input.addMcp')}
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
                                {availableModels.find(model => model.id === selectedModelId)?.model  || t('chat.input.selectModel')}
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
                                                        <span className={styles.defaultBadge}>{t('chat.input.default')}</span>
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
                                                        {model.id === defaultModelId ? t('chat.input.unsetDefault') : t('chat.input.setDefault')}
                                                    </span>
                                                </span>
                                            </button>
                                        ))}
                                        {/* 无结果提示 */}
                                        {filteredModels().length === 0 && (
                                            <div className={styles.noResults}>
                                                {t('chat.input.noMatchingModels')}
                                            </div>
                                        )}
                                    </div>
                                    {/* 搜索输入框 */}
                                    <div className={styles.searchContainer}>
                                        <input
                                            type="text"
                                            placeholder={t('chat.input.searchModels')}
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
                                title={t('chat.input.stopGeneration')}
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
                                title={hasSelectedModel ? t('chat.input.sendMessage') : t('chat.input.sendDisabled')}
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
