import React, { useState, useRef, useEffect, useCallback } from "react";
import styles from "./index.module.scss";
import type { ModelOption } from '@/hooks/useModels';
import {useIsMobile} from "@/hooks/useViewportHeight.ts";

interface ChatInputProps {
    // 类名
    className?: string;
    // 所选模型
    selectedModel: string;
    // 可用模型
    availableModels: ModelOption[];
    // 是否显示滚动到底部按钮
    showScrollToBottom?: boolean;
    // 是否正在生成消息
    isGenerating?: boolean;
    // 点击发送按钮事件（现在会传递输入的文本）
    onSendMessage: (message: string) => void;
    // 点击停止生成按钮事件
    onStopGeneration?: () => void;
    // 模型变更事件
    onModelChange: (model: string) => void;
    // 文件上传事件
    onFileUpload?: (files: File[]) => void;
    // 图片上传事件
    onImageUpload?: (files: File[]) => void;
    // 消息列表滚动到底部按钮点击事件
    onMessageListScrollToBottom?: () => void;
    // 清空输入框的回调
    onClearInput?: () => void;
    // 模型选择框点击事件（用于刷新模型数据）
    onModelSelectorClick?: () => void;
}

const ChatInput: React.FC<ChatInputProps> = ({
    selectedModel,
    availableModels,
    isGenerating = false,
    onSendMessage,
    onStopGeneration,
    onModelChange,
    onFileUpload,
    onImageUpload,
    onMessageListScrollToBottom,
    onClearInput,
    onModelSelectorClick,
    showScrollToBottom = true,
    className
}) => {

    const [showAddMenu, setShowAddMenu] = useState(false);
    const [showModelMenu, setShowModelMenu] = useState(false);
    const [inputValue, setInputValue] = useState('');
    const [isComposing, setIsComposing] = useState(false);
    const [modelSearchValue, setModelSearchValue] = useState(''); // 新增：模型搜索值
    const fileInputRef = useRef<HTMLInputElement>(null);
    const imageInputRef = useRef<HTMLInputElement>(null);
    const addMenuRef = useRef<HTMLDivElement>(null);
    const modelMenuRef = useRef<HTMLDivElement>(null);
    const textareaRef = useRef<HTMLTextAreaElement>(null);
    const lastHeightRef = useRef<number>(0);
    const isMobile =  useIsMobile();

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
        if (textareaRef.current) {
            textareaRef.current.style.height = 'auto';
            textareaRef.current.style.height = '40px';
        }
        onClearInput?.();
    }, [onClearInput]);

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

    // 图像上传事件
    const handleImageUpload = useCallback(() => {
        imageInputRef.current?.click();
        setShowAddMenu(false);
    }, []);

    // 文件上传事件
    const handleFileUpload = useCallback(() => {
        fileInputRef.current?.click();
        setShowAddMenu(false);
    }, []);

    const handleImageChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
        const files = Array.from(e.target.files || []);
        if (files.length > 0 && onImageUpload) {
            onImageUpload(files);
        }
    }, [onImageUpload]);

    //
    const handleFileChange = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
        const files = Array.from(e.target.files || []);
        if (files.length > 0 && onFileUpload) {
            onFileUpload(files);
        }
    }, [onFileUpload]);

    // 模型选择事件
    const handleModelSelect = useCallback((model: string) => {
        onModelChange(model);
        setShowModelMenu(false);
        setModelSearchValue(''); // 选择后清空搜索值
    }, [onModelChange]);

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
            model.name.toLowerCase().includes(modelSearchValue.toLowerCase()) ||
            model.id.toLowerCase().includes(modelSearchValue.toLowerCase())
        );
    }, [availableModels, modelSearchValue]);

    // 消息发送事件
    const handleSend = useCallback(() => {
        if (onMessageListScrollToBottom != null) {
            onMessageListScrollToBottom();
        }
        const trimmedValue = inputValue.trim();
        if (trimmedValue) {
            onSendMessage(trimmedValue);
            clearInput();
        }
    }, [inputValue, onSendMessage, onMessageListScrollToBottom, clearInput]);

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
        
        // 使用 requestAnimationFrame 优化性能
        requestAnimationFrame(() => {
            adjustTextareaHeight();
        });
    }, [adjustTextareaHeight]);

    // 暴露清空输入框的方法供外部调用
    useEffect(() => {
        if (onClearInput) {
            // 将clearInput方法暴露给父组件，但由于React的限制，这里我们通过回调的方式处理
            // 实际的清空逻辑在handleSend中已经处理
        }
    }, [onClearInput]);

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
        <div className={`${styles.chatInput} ${className || ''}`}>
            <div className={styles.inputContainer}>
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
                                        onClick={handleImageUpload}
                                        type="button"
                                    >
                                        <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                                            <path d="M21 19V5c0-1.1-.9-2-2-2H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2zM8.5 13.5l2.5 3.01L14.5 12l4.5 6H5l3.5-4.5z"/>
                                        </svg>
                                        添加图片
                                    </button>
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
                                {availableModels.find(model => model.id === selectedModel)?.name || selectedModel || '选择模型'}
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
                                                className={`${styles.menuItem} ${model.id === selectedModel ? styles.selected : ''}`}
                                                onClick={() => handleModelSelect(model.id)}
                                                type="button"
                                            >
                                                <span className={styles.modelName}>{model.name}</span>
                                                {model.name !== model.id && (
                                                    <span className={styles.modelId}>{model.id}</span>
                                                )}
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
            
            {/* 隐藏的文件输入 */}
            <input
                ref={imageInputRef}
                type="file"
                accept="image/*"
                multiple
                onChange={handleImageChange}
                style={{ display: 'none' }}
            />
            <input
                ref={fileInputRef}
                type="file"
                multiple
                onChange={handleFileChange}
                style={{ display: 'none' }}
            />

            {/* 滚动到底部按钮 */}
            {showScrollToBottom && (
                <div className={`${styles.bottomAction}`}>
                    <div className={`${styles.bottomArrow}`} onClick={onMessageListScrollToBottom}>
                        <svg
                            width="24"
                            height="24"
                            viewBox="0 0 24 24"
                            fill="none"
                            xmlns="http://www.w3.org/2000/svg"
                        >
                            <path
                                d="M7 10L12 15L17 10"
                                stroke="currentColor"
                                strokeWidth="2"
                                strokeLinecap="round"
                                strokeLinejoin="round"
                            />
                        </svg>
                    </div>
                </div>
            )}

        </div>
    );
};

export default ChatInput;