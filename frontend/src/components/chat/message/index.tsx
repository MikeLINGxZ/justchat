import React, {useEffect, useMemo, useState} from "react";
import styles from "./index.module.scss";
import ReasoningContent from "@/components/chat/reasoning_message";
import ExecutionTracePanel from "@/components/chat/execution_trace";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import {Prism as SyntaxHighlighter} from "react-syntax-highlighter";
import {tomorrow} from "react-syntax-highlighter/dist/esm/styles/prism";
import {Service} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service";
import type {Message, Tool as ViewTool} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models";
import {ToolUseStatus, type ToolUse} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models/models";

interface ChatMessageProps {
    message: Message
    isLoading?: boolean
}

let cachedToolsPromise: Promise<ViewTool[]> | null = null;

function loadAvailableTools(): Promise<ViewTool[]> {
    if (!cachedToolsPromise) {
        cachedToolsPromise = Service.GetTools()
            .then((tools) => tools ?? [])
            .catch(() => []);
    }
    return cachedToolsPromise;
}

function resolveToolMeta(toolUse: ToolUse, toolDefinitions: Map<string, ViewTool>): { id: string; name: string; description: string } {
    let matchedTool: ViewTool | undefined;

    if (toolUse.tool_id) {
        matchedTool = toolDefinitions.get(toolUse.tool_id);
    }
    if (!matchedTool && toolUse.tool_name) {
        matchedTool = [...toolDefinitions.values()].find((tool) =>
            tool.id === toolUse.tool_name || tool.name === toolUse.tool_name
        );
    }

    return {
        id: matchedTool?.id || toolUse.tool_id || toolUse.tool_name || 'unknown',
        name: matchedTool?.name || toolUse.tool_name || '未命名工具',
        description: matchedTool?.description || toolUse.tool_description || '暂无描述',
    };
}

function getToolUseDisplayIndex(toolUse: ToolUse, fallbackIndex: number): number {
    return toolUse.index > 0 ? toolUse.index : fallbackIndex + 1;
}

function parseTime(value: unknown): number | null {
    if (!value) {
        return null;
    }
    const date = new Date(value as string);
    const timestamp = date.getTime();
    return Number.isNaN(timestamp) ? null : timestamp;
}

function isToolUseRunning(toolUse: ToolUse): boolean {
    return toolUse.status === ToolUseStatus.ToolUseStatusRunning || toolUse.status === "running";
}

function getToolUseElapsedMs(toolUse: ToolUse, nowMs: number): number {
    const startedAtMs = parseTime(toolUse.started_at);
    const finishedAtMs = parseTime(toolUse.finished_at);
    const baseElapsedMs = toolUse.elapsed_ms ?? 0;

    if (startedAtMs !== null && finishedAtMs !== null) {
        return Math.max(baseElapsedMs, finishedAtMs - startedAtMs, 0);
    }
    if (isToolUseRunning(toolUse) && startedAtMs !== null) {
        return Math.max(baseElapsedMs, nowMs - startedAtMs, 0);
    }
    return Math.max(baseElapsedMs, 0);
}

function formatDuration(elapsedMs: number): string {
    const seconds = Math.max(0, Math.floor(elapsedMs / 1000));
    if (seconds < 60) {
        return `${seconds}s`;
    }
    const minutes = Math.floor(seconds / 60);
    const remainSeconds = seconds % 60;
    return `${minutes}m ${remainSeconds}s`;
}

function getStatusLabel(toolUse: ToolUse): string {
    if (toolUse.status === ToolUseStatus.ToolUseStatusDone || toolUse.status === "done") {
        return "已完成";
    }
    if (toolUse.status === ToolUseStatus.ToolUseStatusError || toolUse.status === "error") {
        return "失败";
    }
    if (toolUse.status === ToolUseStatus.ToolUseStatusPending || toolUse.status === "pending") {
        return "准备中";
    }
    return "执行中";
}

function buildToolMetaTooltip(toolUse: ToolUse, fallbackIndex: number, toolDefinitions: Map<string, ViewTool>): string {
    const displayIndex = getToolUseDisplayIndex(toolUse, fallbackIndex);
    const toolMeta = resolveToolMeta(toolUse, toolDefinitions);
    const lines = [
        `工具 #${displayIndex}`,
        `ID: ${toolMeta.id}`,
        `名称: ${toolMeta.name}`,
        `描述: ${toolMeta.description}`,
    ];

    return lines.join("\n");
}

function buildContentWithToolMarkers(content: string, toolUses: ToolUse[]): string {
    if (!content || toolUses.length === 0) {
        return content;
    }

    const runes = Array.from(content);
    const markersByPos = new Map<number, string[]>();

    toolUses.forEach((toolUse, idx) => {
        const displayIndex = getToolUseDisplayIndex(toolUse, idx);
        const rawPos = typeof toolUse.content_pos === "number" ? toolUse.content_pos : runes.length;
        const pos = Math.max(0, Math.min(rawPos, runes.length));
        const currentMarkers = markersByPos.get(pos) ?? [];
        currentMarkers.push(`[${displayIndex}]`);
        markersByPos.set(pos, currentMarkers);
    });

    const chunks: string[] = [];
    for (let i = 0; i <= runes.length; i++) {
        const markers = markersByPos.get(i);
        if (markers?.length) {
            chunks.push(markers.join(""));
        }
        if (i < runes.length) {
            chunks.push(runes[i]);
        }
    }

    return chunks.join("");
}

function renderTextWithToolMarkers(
    value: string,
    toolUsesByIndex: Map<number, { toolUse: ToolUse; fallbackIndex: number }>,
    toolDefinitions: Map<string, ViewTool>
): React.ReactNode[] {
    const parts = value.split(/(\[\d+\])/g);

    return parts.filter(Boolean).map((part, idx) => {
        const match = /^\[(\d+)\]$/.exec(part);
        if (!match) {
            return <React.Fragment key={`text-${idx}`}>{part}</React.Fragment>;
        }

        const displayIndex = Number(match[1]);
        const toolUseInfo = toolUsesByIndex.get(displayIndex);
        if (!toolUseInfo) {
            return <React.Fragment key={`text-${idx}`}>{part}</React.Fragment>;
        }

        return (
            <span
                key={`marker-${displayIndex}-${idx}`}
                className={styles.inlineToolMarkerWrap}
            >
                <sup className={styles.inlineToolMarker}>
                    {part}
                </sup>
                <span className={styles.inlineToolTooltip} role="tooltip">
                    {buildToolMetaTooltip(toolUseInfo.toolUse, toolUseInfo.fallbackIndex, toolDefinitions)}
                </span>
            </span>
        );
    });
}

function withInlineToolMarkers(
    children: React.ReactNode,
    toolUsesByIndex: Map<number, { toolUse: ToolUse; fallbackIndex: number }>,
    toolDefinitions: Map<string, ViewTool>
): React.ReactNode {
    return React.Children.map(children, (child) => {
        if (typeof child === 'string') {
            return renderTextWithToolMarkers(child, toolUsesByIndex, toolDefinitions);
        }
        if (React.isValidElement(child) && child.props?.children) {
            return React.cloneElement(child as React.ReactElement<any>, {
                ...child.props,
                children: withInlineToolMarkers(child.props.children, toolUsesByIndex, toolDefinitions),
            });
        }
        return child;
    });
}

const ToolUseItem: React.FC<{ toolUse: ToolUse; fallbackIndex: number; nowMs: number; toolDefinitions: Map<string, ViewTool> }> = ({ toolUse, fallbackIndex, nowMs, toolDefinitions }) => {
    const [expanded, setExpanded] = useState(false);
    const result = toolUse.tool_result?.trim() || '';
    const isLong = result.length > 120;
    const displayResult = isLong && !expanded ? result.slice(0, 120) + '…' : result;
    const elapsedLabel = formatDuration(getToolUseElapsedMs(toolUse, nowMs));
    const displayIndex = getToolUseDisplayIndex(toolUse, fallbackIndex);
    const toolMeta = resolveToolMeta(toolUse, toolDefinitions);
    const statusLabel = getStatusLabel(toolUse);
    const statusClassName = isToolUseRunning(toolUse)
        ? styles.toolUseStatusRunning
        : (toolUse.status === ToolUseStatus.ToolUseStatusError || toolUse.status === "error")
            ? styles.toolUseStatusError
            : styles.toolUseStatusDone;
    const tooltip = buildToolMetaTooltip(toolUse, fallbackIndex, toolDefinitions);

    return (
        <div className={styles.toolUseItem}>
            <div
                className={`${styles.toolUseHeader} ${isLong ? styles.toolUseHeaderClickable : ''}`}
                onClick={() => isLong && setExpanded(!expanded)}
                role={isLong ? 'button' : undefined}
            >
                <div className={styles.toolUseMain}>
                    <span className={styles.toolUseBadge}>#{displayIndex}</span>
                    <span className={styles.toolUseName}>{toolMeta.name}</span>
                </div>
                <div className={styles.toolUseMeta}>
                    <span className={`${styles.toolUseStatus} ${statusClassName}`}>
                        {elapsedLabel} · {statusLabel}
                    </span>
                    {isLong && (
                        <span className={styles.toolUseToggle}>
                            {expanded ? '收起' : '展开'}
                        </span>
                    )}
                </div>
            </div>
            {result && (
                <pre className={styles.toolUseResult}>{displayResult}</pre>
            )}
            <div className={styles.toolUseTooltip} role="tooltip">
                {tooltip}
            </div>
        </div>
    );
};

const ToolUsesSection: React.FC<{ toolUses: ToolUse[]; toolDefinitions: Map<string, ViewTool> }> = ({ toolUses, toolDefinitions }) => {
    const [nowMs, setNowMs] = useState(() => Date.now());
    const hasRunningTool = toolUses.some(isToolUseRunning);

    useEffect(() => {
        if (!hasRunningTool) {
            return;
        }
        setNowMs(Date.now());
        const timer = window.setInterval(() => {
            setNowMs(Date.now());
        }, 1000);
        return () => window.clearInterval(timer);
    }, [hasRunningTool]);

    return (
        <div className={styles.toolUsesSection}>
            <div className={styles.toolUsesHeader}>
                <svg className={styles.toolUsesIcon} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                    <path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/>
                </svg>
                <span>工具调用</span>
                <span className={styles.toolUsesCount}>({toolUses.length})</span>
            </div>
            <div className={styles.toolUsesList}>
                {toolUses.map((toolUse, idx) => (
                    <ToolUseItem
                        key={toolUse.call_id || `${toolUse.tool_name}-${idx}`}
                        toolUse={toolUse}
                        fallbackIndex={idx}
                        nowMs={nowMs}
                        toolDefinitions={toolDefinitions}
                    />
                ))}
            </div>
        </div>
    );
};

const ChatMessage: React.FC<ChatMessageProps> = ({
    message,
    isLoading = false,
}: ChatMessageProps) => {
    const [toolDefinitions, setToolDefinitions] = useState<Map<string, ViewTool>>(new Map());
    const isUser = message.role === 'user';
    const wrapperClass = isUser ? styles.userMessageWrapper : styles.assistantMessageWrapper;
    const toolUses = useMemo(() => {
        const currentToolUses = message.assistant_message_extra?.tool_uses ?? [];
        return [...currentToolUses].sort((a, b) => {
            const aIndex = getToolUseDisplayIndex(a, 0);
            const bIndex = getToolUseDisplayIndex(b, 0);
            return aIndex - bIndex;
        });
    }, [message.assistant_message_extra?.tool_uses]);
    const traceSteps = message.assistant_message_extra?.execution_trace?.steps ?? [];
    const toolUsesByIndex = useMemo(() => {
        const map = new Map<number, { toolUse: ToolUse; fallbackIndex: number }>();
        toolUses.forEach((toolUse, idx) => {
            map.set(getToolUseDisplayIndex(toolUse, idx), { toolUse, fallbackIndex: idx });
        });
        return map;
    }, [toolUses]);

    useEffect(() => {
        let active = true;
        loadAvailableTools().then((tools) => {
            if (!active) {
                return;
            }
            setToolDefinitions(new Map(tools.map((tool) => [tool.id, tool])));
        });
        return () => {
            active = false;
        };
    }, []);

    const messageContent = message.content?.trim() ?? "";
    const reasoningContent = message.reasoning_content?.trim() ?? "";
    const finishReason = message.assistant_message_extra?.finish_reason?.trim() ?? "";
    const currentStage = message.assistant_message_extra?.current_stage?.trim() ?? "";
    const isReasoningStreaming = !isUser && isLoading && !finishReason && messageContent.length === 0;
    const isEmptyAssistant = !isUser &&
        !messageContent &&
        !reasoningContent &&
        traceSteps.length === 0 &&
        toolUses.length === 0 &&
        (message.assistant_message_extra?.finish_error == "");
    const hasTrace = traceSteps.length > 0;
    const hasVisibleProgress = hasTrace || currentStage.length > 0 || reasoningContent.length > 0;
    const shouldShowHeadLoading = isLoading && isEmptyAssistant && !hasVisibleProgress;
    const shouldShowTailLoading = !isUser && isLoading && !finishReason && hasVisibleProgress;

    const getDisplayContent = () => {
        if (messageContent) {
            return message.content;
        }
        if (!messageContent && message.assistant_message_extra?.finish_error != "") {
            return message.assistant_message_extra?.finish_error;
        }
        return '';
    };

    const displayMarkdownContent = useMemo(
        () => buildContentWithToolMarkers(message.content, toolUses),
        [message.content, toolUses]
    );

    if (isEmptyAssistant && !isLoading) {
        return null;
    }

    const handleFileClick = (filePath: string) => {
        if (filePath) {
            Service.OpenFile(filePath).catch((err) => {
                console.error('打开文件失败:', err);
            });
        }
    };

    return (
        <div className={styles.ChatMessage}>
            <div className={`${styles.message} ${wrapperClass}`}>
                <div className={styles.messageContainer}>
                    {isUser ? (
                        <>
                            <div className={styles.messageContent}>
                                {getDisplayContent()}
                            </div>
                            {(message.user_message_extra?.files?.length ?? 0) > 0 && (
                                <div className={styles.fileList}>
                                    {message.user_message_extra!.files!.map((file, index) => (
                                        <div
                                            key={index}
                                            className={styles.fileItem}
                                            onClick={() => handleFileClick(file.path)}
                                            title={`点击打开: ${file.name}`}
                                        >
                                            <span className={styles.fileType}>{file.mine_type}</span>
                                            <span className={styles.fileName}>{file.name}</span>
                                            {file.mine_type && (
                                                <span className={styles.fileMimeType}>{file.mine_type}</span>
                                            )}
                                        </div>
                                    ))}
                                </div>
                            )}
                        </>
                    ) : (
                        <div>
                            {shouldShowHeadLoading && (
                                <div className={styles.loadingIndicator}>
                                    <span className={styles.loadingDot} />
                                    <span className={styles.loadingDot} />
                                    <span className={styles.loadingDot} />
                                </div>
                            )}

                            {message.reasoning_content && (
                                <ReasoningContent
                                    content={message.reasoning_content}
                                    isStreaming={isReasoningStreaming}
                                />
                            )}

                            <ExecutionTracePanel
                                trace={message.assistant_message_extra?.execution_trace}
                                currentStage={message.assistant_message_extra?.current_stage}
                                retryCount={message.assistant_message_extra?.retry_count}
                                isStreaming={isLoading}
                            />

                            <div className={`${styles.messageContent} ${styles.markdownContent}`}>
                                <ReactMarkdown
                                    remarkPlugins={[remarkGfm]}
                                    components={{
                                        code(props: any) {
                                            const { inline, className, children, ...rest } = props;
                                            const match = /language-(\w+)/.exec(className || '');
                                            const language = match ? match[1] : '';

                                            return !inline && language ? (
                                                <SyntaxHighlighter
                                                    style={tomorrow}
                                                    language={language}
                                                    PreTag="div"
                                                    customStyle={{
                                                        margin: '8px 0',
                                                        borderRadius: '8px',
                                                        fontSize: '14px'
                                                    } as any}
                                                    {...rest}
                                                >
                                                    {String(children).replace(/\n$/, '')}
                                                </SyntaxHighlighter>
                                            ) : (
                                                <code className={`${className} ${styles.inlineCode}`} {...rest}>
                                                    {children}
                                                </code>
                                            );
                                        },
                                        table: ({children}) => (
                                            <div className={styles.tableWrapper}>
                                                <table className={styles.markdownTable}>{children}</table>
                                            </div>
                                        ),
                                        a: ({children, href}) => (
                                            <a
                                                href={href}
                                                target="_blank"
                                                rel="noopener noreferrer"
                                                className={styles.markdownLink}
                                            >
                                                {children}
                                            </a>
                                        ),
                                        blockquote: ({children}) => (
                                            <blockquote className={styles.markdownBlockquote}>
                                                {withInlineToolMarkers(children, toolUsesByIndex, toolDefinitions)}
                                            </blockquote>
                                        ),
                                        ul: ({children}) => (
                                            <ul className={styles.markdownList}>{children}</ul>
                                        ),
                                        ol: ({children}) => (
                                            <ol className={styles.markdownList}>{children}</ol>
                                        ),
                                        p: ({children}) => (
                                            <p>{withInlineToolMarkers(children, toolUsesByIndex, toolDefinitions)}</p>
                                        ),
                                        li: ({children}) => (
                                            <li>{withInlineToolMarkers(children, toolUsesByIndex, toolDefinitions)}</li>
                                        ),
                                        h1: ({children}) => (
                                            <h1>{withInlineToolMarkers(children, toolUsesByIndex, toolDefinitions)}</h1>
                                        ),
                                        h2: ({children}) => (
                                            <h2>{withInlineToolMarkers(children, toolUsesByIndex, toolDefinitions)}</h2>
                                        ),
                                        h3: ({children}) => (
                                            <h3>{withInlineToolMarkers(children, toolUsesByIndex, toolDefinitions)}</h3>
                                        ),
                                        h4: ({children}) => (
                                            <h4>{withInlineToolMarkers(children, toolUsesByIndex, toolDefinitions)}</h4>
                                        ),
                                        h5: ({children}) => (
                                            <h5>{withInlineToolMarkers(children, toolUsesByIndex, toolDefinitions)}</h5>
                                        ),
                                        h6: ({children}) => (
                                            <h6>{withInlineToolMarkers(children, toolUsesByIndex, toolDefinitions)}</h6>
                                        ),
                                    }}
                                >
                                    {displayMarkdownContent}
                                </ReactMarkdown>
                            </div>

                            {toolUses.length > 0 && traceSteps.length === 0 && (
                                <ToolUsesSection toolUses={toolUses} toolDefinitions={toolDefinitions} />
                            )}

                            {shouldShowTailLoading && (
                                <div className={`${styles.loadingIndicator} ${styles.tailLoadingIndicator}`}>
                                    <span className={styles.loadingDot} />
                                    <span className={styles.loadingDot} />
                                    <span className={styles.loadingDot} />
                                </div>
                            )}
                        </div>
                    )}
                    {message.assistant_message_extra?.finish_reason === 'error' && (
                        <div className={styles.finishReasonError}>⚠ 因错误终止</div>
                    )}
                    {message.assistant_message_extra?.finish_reason === 'user stop' && (
                        <div className={styles.finishReasonUserStop}>⚠ 用户终止生成</div>
                    )}
                </div>
            </div>
        </div>
    );
};

ChatMessage.displayName = 'ChatMessage';
export default ChatMessage;
