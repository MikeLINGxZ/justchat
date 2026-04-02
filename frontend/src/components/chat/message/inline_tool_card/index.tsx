import React, { useCallback, useEffect, useRef, useState } from "react";
import { useTranslation } from 'react-i18next';
import styles from "./index.module.scss";
import {
    TraceStepStatus,
    type ToolUse,
    type TraceStep,
} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models/models";
import {
    WrenchIcon,
    AgentIcon,
    DetailBlockView,
    formatElapsedMs,
    getTraceStepElapsedMs,
    isRunningStep,
    getApprovalMeta,
} from "@/components/chat/execution_trace";
import type { ToolSegment } from "../interleave_utils";

interface InlineToolCardProps {
    toolUse: ToolUse;
    traceStep?: TraceStep;
    childSegments?: ToolSegment[];
    nowMs: number;
    isStreaming?: boolean;
    onApprovalDecision?: (approvalId: string, decision: 'allow' | 'reject') => void;
    onSendApprovalComment?: (approvalId: string, comment: string) => Promise<void> | void;
}

function getElapsedMs(toolUse: ToolUse, traceStep: TraceStep | undefined, nowMs: number): number {
    if (traceStep) {
        return getTraceStepElapsedMs(traceStep, nowMs);
    }
    // Fallback: compute from toolUse directly
    const startMs = toolUse.started_at ? new Date(toolUse.started_at as unknown as string).getTime() : 0;
    const endMs = toolUse.finished_at ? new Date(toolUse.finished_at as unknown as string).getTime() : 0;
    if (startMs && endMs) {
        return Math.max(toolUse.elapsed_ms ?? 0, endMs - startMs, 0);
    }
    if (startMs && !endMs) {
        return Math.max(toolUse.elapsed_ms ?? 0, nowMs - startMs, 0);
    }
    return Math.max(toolUse.elapsed_ms ?? 0, 0);
}

function getStatusDotClass(status: string | undefined): string {
    switch (status) {
        case TraceStepStatus.TraceStepStatusDone:
        case "done":
            return styles.statusDotDone;
        case TraceStepStatus.TraceStepStatusError:
        case TraceStepStatus.TraceStepStatusRejected:
        case "error":
        case "rejected":
            return styles.statusDotError;
        case TraceStepStatus.TraceStepStatusSkipped:
        case "skipped":
            return styles.statusDotSkipped;
        case TraceStepStatus.TraceStepStatusPending:
        case "pending":
            return styles.statusDotPending;
        case TraceStepStatus.TraceStepStatusAwaitingApproval:
        case "awaiting_approval":
            return styles.statusDotAwaitingApproval;
        default:
            return styles.statusDotRunning;
    }
}

const InlineToolCard: React.FC<InlineToolCardProps> = ({
    toolUse,
    traceStep,
    childSegments = [],
    nowMs,
    isStreaming = false,
    onApprovalDecision,
    onSendApprovalComment,
}) => {
    const { t } = useTranslation();
    const [expanded, setExpanded] = useState(false);
    const [showInlineInput, setShowInlineInput] = useState(false);
    const [inlineComment, setInlineComment] = useState('');
    const [isSending, setIsSending] = useState(false);
    const [resultExpanded, setResultExpanded] = useState(false);
    const inlineInputRef = useRef<HTMLTextAreaElement>(null);

    const status = traceStep?.status ?? toolUse.status ?? "";
    const isSubAgent = traceStep ? (traceStep.metadata as Record<string, unknown>)?.is_sub_agent === true : false;
    const toolName = traceStep?.tool_name || toolUse.tool_name || t('chat.executionTrace.unnamedStep');
    const elapsedMs = getElapsedMs(toolUse, traceStep, nowMs);
    const detailBlocks = (traceStep?.detail_blocks ?? []).filter(b => !b.collapsed && b.content?.trim());
    const approvalMeta = traceStep ? getApprovalMeta(traceStep) : null;
    const isAwaitingApproval = (status === TraceStepStatus.TraceStepStatusAwaitingApproval || status === "awaiting_approval") && approvalMeta != null;
    const hasChildren = childSegments.length > 0;
    const hasExpandableContent = detailBlocks.length > 0 || (toolUse.tool_result?.trim().length ?? 0) > 0 || isAwaitingApproval || hasChildren;

    const prevIsStreamingRef = useRef(isStreaming);

    // Auto-expand when awaiting approval
    useEffect(() => {
        if (isAwaitingApproval) {
            setExpanded(true);
        }
    }, [isAwaitingApproval]);

    // Auto-expand when running (to show progress)
    useEffect(() => {
        if (traceStep && isRunningStep(traceStep)) {
            setExpanded(true);
        }
    }, [traceStep]);

    // Auto-collapse when streaming finishes
    useEffect(() => {
        const wasStreaming = prevIsStreamingRef.current;
        prevIsStreamingRef.current = isStreaming;
        if (wasStreaming && !isStreaming) {
            setExpanded(false);
        }
    }, [isStreaming]);

    useEffect(() => {
        if (showInlineInput && inlineInputRef.current) {
            inlineInputRef.current.focus();
        }
    }, [showInlineInput]);

    const handleSendComment = useCallback(async (approvalId: string) => {
        const trimmed = inlineComment.trim();
        if (!trimmed || isSending) return;
        setIsSending(true);
        try {
            await onSendApprovalComment?.(approvalId, trimmed);
            setShowInlineInput(false);
            setInlineComment('');
        } finally {
            setIsSending(false);
        }
    }, [inlineComment, isSending, onSendApprovalComment]);

    const toolResult = toolUse.tool_result?.trim() || '';
    const isResultLong = toolResult.length > 500;

    return (
        <div className={styles.inlineToolCard}>
            {/* Header row */}
            <div
                className={styles.cardHeader}
                onClick={() => hasExpandableContent && setExpanded(!expanded)}
                role={hasExpandableContent ? "button" : undefined}
                tabIndex={hasExpandableContent ? 0 : undefined}
                onKeyDown={hasExpandableContent ? (e) => {
                    if (e.key === 'Enter' || e.key === ' ') setExpanded(!expanded);
                } : undefined}
            >
                <span className={isSubAgent ? styles.agentIcon : styles.toolIcon}>
                    {isSubAgent ? <AgentIcon /> : <WrenchIcon />}
                </span>
                {isSubAgent && (
                    <span className={styles.subAgentBadge}>{t('chat.executionTrace.agent')}</span>
                )}
                <span className={styles.toolName}>{toolName}</span>
                <span className={styles.sep}>&middot;</span>
                <span className={styles.elapsed}>{formatElapsedMs(elapsedMs)}</span>
                <span className={`${styles.statusDot} ${getStatusDotClass(status)}`} />
                {hasExpandableContent && (
                    <span className={`${styles.toggleHint} ${expanded ? styles.toggleHintExpanded : ''}`}>
                        ▾
                    </span>
                )}
            </div>

            {/* Expandable body */}
            <div className={`${styles.cardBody} ${expanded ? styles.cardBodyExpanded : styles.cardBodyCollapsed}`}>
                {expanded && (
                    <div className={styles.cardBodyContent}>
                        {/* Detail blocks from trace step */}
                        {detailBlocks.length > 0 && traceStep && (
                            <div className={styles.detailSection}>
                                {detailBlocks.map((block, index) => (
                                    <DetailBlockView
                                        key={`${block.kind}-${index}`}
                                        step={traceStep}
                                        block={block}
                                        index={index}
                                        t={t}
                                    />
                                ))}
                            </div>
                        )}

                        {/* Summary */}
                        {traceStep && (traceStep.summary || traceStep.output_preview) && (
                            <div className={styles.summary}>
                                {traceStep.summary || traceStep.output_preview}
                            </div>
                        )}

                        {/* Fallback: show tool_result when no trace detail blocks */}
                        {detailBlocks.length === 0 && toolResult && (
                            <div className={styles.fallbackResult}>
                                <div className={styles.fallbackResultTitle}>
                                    {t('chat.executionTrace.toolCall')}
                                </div>
                                <pre className={`${styles.fallbackResultContent} ${isResultLong && !resultExpanded ? styles.fallbackResultClamped : ''}`}>
                                    {toolResult}
                                </pre>
                                {isResultLong && (
                                    <button
                                        type="button"
                                        className={styles.showMoreBtn}
                                        onClick={(e) => {
                                            e.stopPropagation();
                                            setResultExpanded(!resultExpanded);
                                        }}
                                    >
                                        {resultExpanded ? t('chat.executionTrace.collapse') : t('chat.executionTrace.expand')}
                                    </button>
                                )}
                            </div>
                        )}

                        {/* Child tool cards (sub-agent's tools) */}
                        {hasChildren && (
                            <div className={styles.childToolsSection}>
                                {childSegments.map((child) => (
                                    <InlineToolCard
                                        key={child.key}
                                        toolUse={child.toolUse}
                                        traceStep={child.traceStep}
                                        childSegments={child.childSegments}
                                        nowMs={nowMs}
                                        isStreaming={isStreaming}
                                        onApprovalDecision={onApprovalDecision}
                                        onSendApprovalComment={onSendApprovalComment}
                                    />
                                ))}
                            </div>
                        )}

                        {/* Approval UI */}
                        {isAwaitingApproval && approvalMeta && (
                            <div className={styles.approvalCard}>
                                <div className={styles.approvalTitle}>{approvalMeta.title}</div>
                                <div className={styles.approvalBody}>{approvalMeta.message}</div>
                                <div className={styles.approvalActions}>
                                    <button
                                        type="button"
                                        className={`${styles.approvalButton} ${styles.approvalButtonAllow}`}
                                        onClick={(e) => {
                                            e.stopPropagation();
                                            onApprovalDecision?.(approvalMeta.approvalId, 'allow');
                                        }}
                                    >
                                        {t('chat.executionTrace.allow')}
                                    </button>
                                    <button
                                        type="button"
                                        className={`${styles.approvalButton} ${styles.approvalButtonReject}`}
                                        onClick={(e) => {
                                            e.stopPropagation();
                                            onApprovalDecision?.(approvalMeta.approvalId, 'reject');
                                        }}
                                    >
                                        {t('chat.executionTrace.reject')}
                                    </button>
                                    <button
                                        type="button"
                                        className={`${styles.approvalButton} ${styles.approvalButtonGuide}`}
                                        onClick={(e) => {
                                            e.stopPropagation();
                                            setShowInlineInput(true);
                                        }}
                                    >
                                        {t('chat.executionTrace.guideAi')}
                                    </button>
                                </div>
                                {showInlineInput && (
                                    <div className={styles.approvalInlineInput}>
                                        <textarea
                                            ref={inlineInputRef}
                                            className={styles.approvalTextarea}
                                            value={inlineComment}
                                            onChange={(e) => setInlineComment(e.target.value)}
                                            placeholder={t('chat.executionTrace.inlinePlaceholder')}
                                            rows={2}
                                            onKeyDown={(e) => {
                                                if (e.key === 'Enter' && !e.shiftKey) {
                                                    e.preventDefault();
                                                    handleSendComment(approvalMeta.approvalId);
                                                }
                                                if (e.key === 'Escape') {
                                                    setShowInlineInput(false);
                                                    setInlineComment('');
                                                }
                                            }}
                                            disabled={isSending}
                                        />
                                        <div className={styles.approvalInlineActions}>
                                            <button
                                                type="button"
                                                className={styles.approvalInlineSend}
                                                onClick={() => handleSendComment(approvalMeta.approvalId)}
                                                disabled={!inlineComment.trim() || isSending}
                                            >
                                                {t('chat.executionTrace.sendComment')}
                                            </button>
                                            <button
                                                type="button"
                                                className={styles.approvalInlineCancel}
                                                onClick={() => {
                                                    setShowInlineInput(false);
                                                    setInlineComment('');
                                                }}
                                                disabled={isSending}
                                            >
                                                {t('chat.executionTrace.cancelComment')}
                                            </button>
                                        </div>
                                    </div>
                                )}
                            </div>
                        )}
                    </div>
                )}
            </div>
        </div>
    );
};

export default InlineToolCard;
