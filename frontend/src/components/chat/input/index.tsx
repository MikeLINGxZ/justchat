import React, {useCallback, useEffect, useRef, useState} from "react";
import styles from "./index.module.scss";
import {useIsMobile} from "@/hooks/useViewportHeight.ts";
import {Service} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service";
import {FileInfo, Model} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models";

interface ChatInputProps {
    // 所选模型
    selectedModelId: number;
    // 可用模型
    availableModels: Model[];
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
    availableModels,
    isGenerating = false,
    onMessageChange,
    onSendButtonClick,
    onSelectFileChange,
    onStopGeneration,
    onSelectModelChange,
    onMessageListScrollToBottom,
    onModelSelectorClick,
}) => {

    const [showAddMenu, setShowAddMenu] = useState(false);
    const [showModelMenu, setShowModelMenu] = useState(false);
    const [inputValue, setInputValue] = useState('');
    const [isComposing, setIsComposing] = useState(false);
    const [modelSearchValue, setModelSearchValue] = useState('');
    const [defaultModelId, setDefaultModelId] = useState<number | null>(() => {
        return getDefaultModelConfig()?.modelId ?? null;
    });
    const fileInputRef = useRef<HTMLInputElement>(null);
    const imageInputRef = useRef<HTMLInputElement>(null);
    const addMenuRef = useRef<HTMLDivElement>(null);
    const modelMenuRef = useRef<HTMLDivElement>(null);
    const textareaRef = useRef<HTMLTextAreaElement>(null);
    const lastHeightRef = useRef<number>(0);
    const isMobile =  useIsMobile();
    const [selectFiles, setSelectFiles] = useState<FileInfo[]>([]);

    // 优化的高度调整函数
    const adjustTextareaHeight = useCallback(() => {
        const textarea = textareaRef.current;
        if (!textarea) return;

        // 避免不必要的DOM操作
        const currentScrollHeight = textarea.scrollHeight;
        if (currentScrollHeight === lastHeightRef.current) {
            return;
        }

        // 使用requestAnimationFrame优化DOM操作
        requestAnimationFrame(() => {
            textarea.style.height = 'auto';
            const scrollHeight = textarea.scrollHeight;
            const maxHeight = parseFloat(getComputedStyle(textarea).maxHeight);
            
            if (scrollHeight <= maxHeight) {
                textarea.style.height = `${scrollHeight}px`;
                textarea.style.overflowY = 'hidden';
            } else {
                textarea.style.height = `${maxHeight}px`;
                textarea.style.overflowY = 'auto';
            }
            
            lastHeightRef.current = scrollHeight;
        });
    }, []);

    // 清空输入框
    const clearInput = useCallback(() => {
        setInputValue('');
        setSelectFiles([]); // 清空文件列表
        onSelectFileChange(selectFiles);
        onMessageChange(inputValue);
        if (textareaRef.current) {
            textareaRef.current.style.height = 'auto';
            textareaRef.current.style.height = '40px';
        }
    }, []);

    // 处理中文输入开始
    const handleCompositionStart = useCallback(() => {
        setIsComposing(true);
    }, []);

    // 处理中文输入结束
    const handleCompositionEnd = useCallback((e: React.CompositionEvent<HTMLTextAreaElement>) => {
        setIsComposing(false);
        adjustTextareaHeight();
    }, [adjustTextareaHeight]);

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

    // 文件上传事件
    const handleFileUpload = useCallback(() => {
        fileInputRef.current?.click();
        Service.SelectFiles().then(async (files: FileInfo[]) => {
            if (files.length === 0) return;
            const newPaths = [...selectFiles, ...files];
            setSelectFiles(newPaths);
            onSelectFileChange(newPaths);
        }).catch(() => {
        }).finally(() => {
            setShowAddMenu(false);
        })
    }, [selectFiles, onSelectFileChange]);

    // 删除文件
    const handleRemoveFile = useCallback((filePath: string) => {
        setSelectFiles(prevFiles => prevFiles.filter(f => f.path !== filePath));
        setSelectFiles(prevState => prevState.filter(f => f.path !== filePath))
    }, []);

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
    const handleSend = useCallback(() => {
        if (onMessageListScrollToBottom != null) {
            onMessageListScrollToBottom();
        }
        const trimmedValue = inputValue.trim();
        if (trimmedValue) {
            onSendButtonClick();
            clearInput(); // 清空输入框和文件列表
        }
    }, [inputValue, selectFiles, onSendButtonClick, onMessageListScrollToBottom, clearInput]);

    //
    const handleKeyDown = useCallback((e: React.KeyboardEvent) => {
        // 如果正在使用中文输入法（composition状态），不处理回车键
        if (isComposing) {
            return;
        }
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            handleSend();
        }
    }, [handleSend, isComposing]);

    //
    const handleInputChange = useCallback((e: React.ChangeEvent<HTMLTextAreaElement>) => {
        const value = e.target.value;
        setInputValue(value);
        onMessageChange(value);
        // 使用 requestAnimationFrame 优化性能
        requestAnimationFrame(() => {
            adjustTextareaHeight();
        });
    }, [adjustTextareaHeight]);

    // 初始化时调整高度
    useEffect(() => {
        adjustTextareaHeight();
    }, [adjustTextareaHeight]);

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
        };

        if (showAddMenu || showModelMenu) {
            document.addEventListener('mousedown', handleClickOutside);
        }

        return () => {
            document.removeEventListener('mousedown', handleClickOutside);
        };
    }, [showAddMenu, showModelMenu]);

    return (
        <div className={`${styles.chatInput}`}>
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
                <textarea
                    ref={textareaRef}
                    className={styles.textInput}
                    value={inputValue}
                    onChange={handleInputChange}
                    onCompositionStart={handleCompositionStart}
                    onCompositionEnd={handleCompositionEnd}
                    onKeyDown={handleKeyDown}
                    placeholder="输入消息... (支持 Markdown 格式)"
                    rows={1}
                    data-markdown="true"
                    style={{
                        height: 'auto',
                        minHeight: '40px',
                        maxHeight: 'calc(1.5em * 8 + 16px)'
                    }}
                />
                
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
                                className={`${styles.sendButton} ${!inputValue.trim() ? styles.disabled : ''}`}
                                onClick={handleSend}
                                disabled={!inputValue.trim()}
                                type="button"
                                title="发送消息"
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