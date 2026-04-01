import React, { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useTranslation } from 'react-i18next';
import styles from "./index.module.scss";
import {
    TraceDetailFormat,
    TraceStepStatus,
    TraceStepType,
    type ExecutionTrace,
    type TraceDetailBlock,
    type TraceStep,
} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models/models";
import MarkdownRenderer from "@/components/markdown_renderer";

interface ExecutionTraceProps {
    trace?: ExecutionTrace | null;
    currentStage?: string;
    retryCount?: number;
    isStreaming?: boolean;
    onApprovalDecision?: (approvalId: string, decision: 'allow' | 'reject') => void;
    onSendApprovalComment?: (approvalId: string, comment: string) => Promise<void> | void;
}

const stageLabelMap: Record<string, string> = {
    "": "chat.executionTrace.processing",
    "chat.stage.pending": "chat.executionTrace.waiting",
    "chat.stage.preparing": "chat.executionTrace.preparing",
    "chat.stage.classify": "chat.executionTrace.classify",
    "chat.stage.direct_answer": "chat.executionTrace.directAnswer",
    "chat.stage.plan": "chat.executionTrace.plan",
    "chat.stage.running_tasks": "chat.executionTrace.runningTasks",
    "chat.stage.awaiting_approval": "chat.executionTrace.awaitingApproval",
    "chat.stage.synthesize": "chat.executionTrace.synthesize",
    "chat.stage.review": "chat.executionTrace.review",
    "chat.stage.retry": "chat.executionTrace.retry",
    "chat.stage.finished": "chat.executionTrace.finished",
    "等待执行": "chat.executionTrace.waiting",
    "准备执行": "chat.executionTrace.preparing",
    "意图识别": "chat.executionTrace.classify",
    "直接回答": "chat.executionTrace.directAnswer",
    "任务拆解": "chat.executionTrace.plan",
    "子任务执行": "chat.executionTrace.runningTasks",
    "等待用户确认": "chat.executionTrace.awaitingApproval",
    "结果汇总": "chat.executionTrace.synthesize",
    "结果审核": "chat.executionTrace.review",
    "重新生成": "chat.executionTrace.retry",
    "已完成": "chat.executionTrace.finished",
};

const typeLabelMap: Record<string, string> = {
    [TraceStepType.TraceStepTypeClassify]: "chat.executionTrace.classify",
    [TraceStepType.TraceStepTypePlan]: "chat.executionTrace.plan",
    [TraceStepType.TraceStepTypeDispatch]: "chat.executionTrace.dispatch",
    [TraceStepType.TraceStepTypeAgentRun]: "chat.executionTrace.agentRun",
    [TraceStepType.TraceStepTypeToolCall]: "chat.executionTrace.toolCall",
    [TraceStepType.TraceStepTypeSynthesize]: "chat.executionTrace.synthesize",
    [TraceStepType.TraceStepTypeReview]: "chat.executionTrace.review",
    [TraceStepType.TraceStepTypeRetry]: "chat.executionTrace.retry",
    [TraceStepType.TraceStepTypeFinalize]: "chat.executionTrace.finalize",
};

function getStatusDotClass(status: string | undefined): string {
    switch (status) {
        case TraceStepStatus.TraceStepStatusDone:
            return styles.statusDotDone;
        case TraceStepStatus.TraceStepStatusError:
        case TraceStepStatus.TraceStepStatusRejected:
            return styles.statusDotError;
        case TraceStepStatus.TraceStepStatusSkipped:
            return styles.statusDotSkipped;
        case TraceStepStatus.TraceStepStatusPending:
            return styles.statusDotPending;
        case TraceStepStatus.TraceStepStatusAwaitingApproval:
            return styles.statusDotAwaitingApproval;
        default:
            return styles.statusDotRunning;
    }
}

function getStatusLabel(status: string | undefined, t: (key: string) => string): string {
    switch (status) {
        case TraceStepStatus.TraceStepStatusDone:
            return t('chat.executionTrace.finished');
        case TraceStepStatus.TraceStepStatusError:
            return t('chat.message.failed');
        case TraceStepStatus.TraceStepStatusSkipped:
            return t('chat.executionTrace.skipped');
        case TraceStepStatus.TraceStepStatusPending:
            return t('chat.message.pending');
        case TraceStepStatus.TraceStepStatusAwaitingApproval:
            return t('chat.message.awaitingApproval');
        case TraceStepStatus.TraceStepStatusRejected:
            return t('chat.message.rejected');
        default:
            return t('chat.message.running');
    }
}

function formatElapsedMs(value?: number): string {
    const ms = value ?? 0;
    if (ms < 1000) {
        return `${ms}ms`;
    }
    const sec = Math.floor(ms / 1000);
    if (sec < 60) {
        const tenths = Math.floor((ms % 1000) / 100);
        return tenths > 0 ? `${sec}.${tenths}s` : `${sec}s`;
    }
    return `${Math.floor(sec / 60)}m ${sec % 60}s`;
}

function parseTime(value: unknown): number | null {
    if (!value) {
        return null;
    }
    const date = new Date(value as string);
    const timestamp = date.getTime();
    return Number.isNaN(timestamp) ? null : timestamp;
}

function isRunningStep(step: TraceStep): boolean {
    return step.status === TraceStepStatus.TraceStepStatusRunning
        || step.status === TraceStepStatus.TraceStepStatusPending
        || step.status === TraceStepStatus.TraceStepStatusAwaitingApproval
        || step.status === "";
}

function getTraceStepElapsedMs(step: TraceStep, nowMs: number): number {
    const startedAtMs = parseTime(step.started_at);
    const finishedAtMs = parseTime(step.finished_at);
    const baseElapsedMs = step.elapsed_ms ?? 0;
    if (startedAtMs !== null && finishedAtMs !== null) {
        return Math.max(baseElapsedMs, finishedAtMs - startedAtMs, 0);
    }
    if (isRunningStep(step) && startedAtMs !== null) {
        return Math.max(baseElapsedMs, nowMs - startedAtMs, 0);
    }
    return Math.max(baseElapsedMs, 0);
}

function buildTree(steps: TraceStep[]): Map<string, TraceStep[]> {
    const tree = new Map<string, TraceStep[]>();
    const sortedSteps = [...steps].sort((a, b) => {
        const aTime = a.started_at ? new Date(a.started_at as unknown as string).getTime() : 0;
        const bTime = b.started_at ? new Date(b.started_at as unknown as string).getTime() : 0;
        return aTime - bTime;
    });
    const latestAgentStepByName = new Map<string, string>();

    sortedSteps.forEach((step) => {
        let parent = step.parent_step_id || "__root__";
        if (
            parent === "__root__" &&
            step.type === TraceStepType.TraceStepTypeToolCall &&
            step.agent_name &&
            latestAgentStepByName.has(step.agent_name)
        ) {
            parent = latestAgentStepByName.get(step.agent_name) || "__root__";
        }
        const group = tree.get(parent) ?? [];
        group.push(step);
        tree.set(parent, group);

        if (step.type === TraceStepType.TraceStepTypeAgentRun && step.step_id) {
            latestAgentStepByName.set(step.agent_name || step.step_id, step.step_id);
        }
    });
    tree.forEach((items) => {
        items.sort((a, b) => {
            const aTime = a.started_at ? new Date(a.started_at as unknown as string).getTime() : 0;
            const bTime = b.started_at ? new Date(b.started_at as unknown as string).getTime() : 0;
            return aTime - bTime;
        });
    });
    return tree;
}

function getDetailBlockDisplayTitle(step: TraceStep, block: TraceDetailBlock, t: (key: string) => string): string {
    if (step.type !== TraceStepType.TraceStepTypeAgentRun) {
        return block.title;
    }
    if (block.kind === "input") {
        return t('chat.executionTrace.userInput');
    }
    if (block.kind === "output") {
        return t('chat.executionTrace.finalAnswer');
    }
    return block.title;
}

/** Wrench icon for tool call rows */
const WrenchIcon: React.FC = () => (
    <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
        <path d="M10.5 2.5a3.5 3.5 0 0 0-3.27 4.73L2.5 12l1.5 1.5 4.77-4.73A3.5 3.5 0 0 0 14 5.5l-2 2-1.5-.5-.5-1.5 2-2a3.5 3.5 0 0 0-1.5-.5z" />
    </svg>
);

/** Agent icon for sub-agent call rows */
const AgentIcon: React.FC = () => (
    <svg viewBox="0 0 16 16" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
        <circle cx="8" cy="5" r="2.5" />
        <path d="M3 14c0-2.76 2.24-5 5-5s5 2.24 5 5" />
    </svg>
);

/** Detail block with optional "show more" for long content */
const DetailBlockView: React.FC<{
    step: TraceStep;
    block: TraceDetailBlock;
    index: number;
    t: (key: string) => string;
}> = ({ step, block, index, t }) => {
    const displayTitle = getDetailBlockDisplayTitle(step, block, t);
    const [contentExpanded, setContentExpanded] = useState(false);

    if (block.collapsed) {
        return null;
    }

    if (!block.content?.trim()) {
        return null;
    }

    const isLong = (block.content?.length ?? 0) > 500;
    const contentClasses = [
        styles.detailContent,
        block.format === TraceDetailFormat.TraceDetailFormatJSON ? styles.jsonBlock : '',
        isLong && !contentExpanded ? styles.detailContentClamped : '',
        isLong && contentExpanded ? styles.detailContentExpanded : '',
    ].filter(Boolean).join(' ');

    if (block.format === TraceDetailFormat.TraceDetailFormatMarkdown) {
        return (
            <div key={`${block.kind}-${index}`} className={styles.detailBlock}>
                <div className={styles.detailTitle}>{displayTitle}</div>
                <div className={`${styles.detailContent} ${styles.detailMarkdown}`}>
                    <MarkdownRenderer content={block.content} variant="trace" />
                </div>
            </div>
        );
    }

    return (
        <div key={`${block.kind}-${index}`} className={styles.detailBlock}>
            <div className={styles.detailTitle}>{displayTitle}</div>
            <pre className={contentClasses}>{block.content}</pre>
            {isLong && (
                <button
                    type="button"
                    className={styles.showMoreBtn}
                    onClick={() => setContentExpanded(!contentExpanded)}
                >
                    {contentExpanded ? t('chat.executionTrace.collapse') : t('chat.executionTrace.expand')}
                </button>
            )}
        </div>
    );
};

function getApprovalMeta(step: TraceStep): { approvalId: string; title: string; message: string } | null {
    const metadata = step.metadata ?? {};
    const approvalId = typeof metadata.approval_id === "string" ? metadata.approval_id : "";
    if (!approvalId) {
        return null;
    }
    const title = typeof metadata.approval_title === "string" ? metadata.approval_title : (step.title || "");
    const message = typeof metadata.approval_message === "string" ? metadata.approval_message : (step.summary || "");
    return { approvalId, title, message };
}

/** Get agent step CSS classes for the left border accent */
function getAgentMainClasses(step: TraceStep): string {
    const classes = [styles.nodeMain, styles.nodeMainAgent];
    if (isRunningStep(step)) {
        classes.push(styles.nodeMainAgentRunning);
    } else if (step.status === TraceStepStatus.TraceStepStatusError || step.status === TraceStepStatus.TraceStepStatusRejected) {
        classes.push(styles.nodeMainAgentError);
    } else {
        classes.push(styles.nodeMainAgentDone);
    }
    return classes.join(' ');
}

const TraceNode: React.FC<{
    step: TraceStep;
    tree: Map<string, TraceStep[]>;
    nowMs: number;
    depth?: number;
    autoExpand?: boolean;
    onApprovalDecision?: (approvalId: string, decision: 'allow' | 'reject') => void;
    onSendApprovalComment?: (approvalId: string, comment: string) => Promise<void> | void;
}> = ({ step, tree, nowMs, depth = 0, autoExpand = false, onApprovalDecision, onSendApprovalComment }) => {
    const { t } = useTranslation();
    const children = tree.get(step.step_id) ?? [];
    const inlineChildren = step.type === TraceStepType.TraceStepTypeAgentRun ? children : [];
    const nestedChildren = step.type === TraceStepType.TraceStepTypeAgentRun ? [] : children;
    const [expanded, setExpanded] = useState(autoExpand);
    const [showInlineInput, setShowInlineInput] = useState(false);
    const [inlineComment, setInlineComment] = useState('');
    const [isSending, setIsSending] = useState(false);
    const inlineInputRef = useRef<HTMLTextAreaElement>(null);

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

    useEffect(() => {
        if (showInlineInput && inlineInputRef.current) {
            inlineInputRef.current.focus();
        }
    }, [showInlineInput]);
    const hasChildren = children.length > 0;
    const isRunning = isRunningStep(step);
    const detailBlocks = step.detail_blocks ?? [];
    const approvalMeta = getApprovalMeta(step);
    const isAwaitingApproval = step.status === TraceStepStatus.TraceStepStatusAwaitingApproval && approvalMeta != null;
    const isAgentStep = step.type === TraceStepType.TraceStepTypeAgentRun;
    const isToolCallStep = step.type === TraceStepType.TraceStepTypeToolCall;
    const isSubAgent = isToolCallStep && (step.metadata as Record<string, unknown>)?.is_sub_agent === true;
    const leadingDetailBlocks = isAgentStep
        ? detailBlocks.filter((block) => block.kind !== "output" && block.kind !== "tool_result")
        : detailBlocks;
    const trailingDetailBlocks = isAgentStep
        ? detailBlocks.filter((block) => block.kind === "output" || block.kind === "tool_result")
        : [];

    const hasExpandableContent = hasChildren || detailBlocks.length > 0;

    useEffect(() => {
        if (isRunning) {
            setExpanded(true);
        }
    }, [isRunning]);

    // ── Tool Call: compact single-line row ───────────────────────────────
    if (isToolCallStep) {
        return (
            <div className={styles.node} style={{ marginLeft: depth * 14 }}>
                <div className={styles.nodeHeader}>
                    <div
                        className={styles.toolCallRow}
                        onClick={() => hasExpandableContent && setExpanded(!expanded)}
                        role={hasExpandableContent ? "button" : undefined}
                        tabIndex={hasExpandableContent ? 0 : undefined}
                        onKeyDown={hasExpandableContent ? (e) => { if (e.key === 'Enter' || e.key === ' ') setExpanded(!expanded); } : undefined}
                    >
                        <span className={isSubAgent ? styles.agentCallIcon : styles.toolCallIcon}>{isSubAgent ? <AgentIcon /> : <WrenchIcon />}</span>
                        {isSubAgent && <span className={styles.subAgentBadge}>{t('chat.executionTrace.agent')}</span>}
                        {isSubAgent && step.agent_name && <span className={styles.callerInfo}>{t('chat.executionTrace.caller', { name: step.agent_name })}</span>}
                        <span className={styles.toolCallName}>{step.tool_name || step.title || t('chat.executionTrace.unnamedStep')}</span>
                        <span className={styles.toolCallSep}>&middot;</span>
                        <span className={styles.toolCallElapsed}>{formatElapsedMs(getTraceStepElapsedMs(step, nowMs))}</span>
                        <span className={`${styles.statusDot} ${getStatusDotClass(step.status)}`} />
                    </div>
                </div>
                {/* Expandable detail content */}
                <div className={`${styles.expandableContent} ${expanded ? styles.expandableContentExpanded : styles.expandableContentCollapsed}`}>
                    {expanded && detailBlocks.length > 0 && (
                        <div className={styles.detailSection} style={{ marginLeft: 8 }}>
                            {detailBlocks.map((block, index) => (
                                <DetailBlockView key={`${block.kind}-${index}`} step={step} block={block} index={index} t={t} />
                            ))}
                        </div>
                    )}
                    {expanded && (step.summary || step.output_preview) && (
                        <div className={styles.nodeSummary} style={{ marginLeft: 8, marginTop: 4 }}>
                            {step.summary || step.output_preview}
                        </div>
                    )}
                    {expanded && isAwaitingApproval && approvalMeta && (
                        <div className={styles.approvalCard} style={{ marginLeft: 8 }}>
                            <div className={styles.approvalTitle}>{approvalMeta.title}</div>
                            <div className={styles.approvalBody}>{approvalMeta.message}</div>
                            <div className={styles.approvalActions}>
                                <button type="button" className={`${styles.approvalButton} ${styles.approvalButtonAllow}`} onClick={() => onApprovalDecision?.(approvalMeta.approvalId, 'allow')}>
                                    {t('chat.executionTrace.allow')}
                                </button>
                                <button type="button" className={`${styles.approvalButton} ${styles.approvalButtonReject}`} onClick={() => onApprovalDecision?.(approvalMeta.approvalId, 'reject')}>
                                    {t('chat.executionTrace.reject')}
                                </button>
                                <button type="button" className={`${styles.approvalButton} ${styles.approvalButtonGuide}`} onClick={() => setShowInlineInput(true)}>
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
                                            onClick={() => { setShowInlineInput(false); setInlineComment(''); }}
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
                {nestedChildren.length > 0 && expanded && (
                    <div className={styles.children}>
                        {nestedChildren.map((child) => (
                            <TraceNode
                                key={child.step_id}
                                step={child}
                                tree={tree}
                                nowMs={nowMs}
                                depth={depth + 1}
                                autoExpand={isRunning}
                                onApprovalDecision={onApprovalDecision}
                                onSendApprovalComment={onSendApprovalComment}
                            />
                        ))}
                    </div>
                )}
            </div>
        );
    }

    // ── Agent Run & other step types: card layout ────────────────────────
    const mainClasses = isAgentStep ? getAgentMainClasses(step) : styles.nodeMain;

    return (
        <div className={styles.node} style={{ marginLeft: depth * 14 }}>
            <div className={styles.nodeHeader}>
                <div className={mainClasses}>
                    <div className={styles.nodeTopRow}>
                        <div className={styles.nodeTitleRow}>
                            <span className={styles.nodeType}>{t(typeLabelMap[step.type] || 'chat.executionTrace.step')}</span>
                            <span className={styles.nodeTitle}>{step.title || t('chat.executionTrace.unnamedStep')}</span>
                        </div>
                        {hasExpandableContent && (
                            <button className={styles.toggle} type="button" onClick={() => setExpanded(!expanded)}>
                                <span className={styles.toggleIcon}>{expanded ? "▾" : "▸"}</span>
                                <span className={styles.toggleText}>{expanded ? t('chat.executionTrace.collapse') : t('chat.executionTrace.expand')}</span>
                            </button>
                        )}
                    </div>
                    {(step.summary || step.output_preview || step.input_preview) && (
                        <div className={styles.nodeSummary}>
                            {step.summary || step.output_preview || step.input_preview}
                        </div>
                    )}
                    {(step.agent_name || step.tool_name) && (
                        <div className={styles.nodeMetaText}>
                            {step.agent_name ? `Agent: ${step.agent_name}` : ""}
                            {step.agent_name && step.tool_name ? " · " : ""}
                            {step.tool_name ? `Tool: ${step.tool_name}` : ""}
                        </div>
                    )}
                    {/* Expandable content with animation */}
                    <div className={`${styles.expandableContent} ${expanded ? styles.expandableContentExpanded : styles.expandableContentCollapsed}`}>
                        {expanded && isAgentStep && (
                            <div className={styles.agentFlowSection}>
                                {leadingDetailBlocks.length > 0 && (
                                    <div className={styles.detailSection}>
                                        {leadingDetailBlocks.map((block, index) => (
                                            <DetailBlockView key={`${block.kind}-${index}`} step={step} block={block} index={index} t={t} />
                                        ))}
                                    </div>
                                )}
                                {(inlineChildren.length > 0 || isRunning) && (
                                    <div className={styles.inlineChildrenSection}>
                                        <div className={styles.inlineChildrenTitle}>{t('chat.executionTrace.reasoning')}</div>
                                        {inlineChildren.length > 0 ? (
                                            <div className={styles.inlineChildrenList}>
                                                {inlineChildren.map((child) => (
                                                    <TraceNode
                                                        key={child.step_id}
                                                        step={child}
                                                        tree={tree}
                                                        nowMs={nowMs}
                                                        depth={0}
                                                        autoExpand={isRunning}
                                                        onApprovalDecision={onApprovalDecision}
                                                        onSendApprovalComment={onSendApprovalComment}
                                                    />
                                                ))}
                                            </div>
                                        ) : (
                                            <div className={styles.agentThinkingPlaceholder}>{t('chat.executionTrace.agentThinking')}</div>
                                        )}
                                    </div>
                                )}
                                {trailingDetailBlocks.length > 0 && (
                                    <div className={styles.detailSection}>
                                        {trailingDetailBlocks.map((block, index) => (
                                            <DetailBlockView key={`${block.kind}-${index}`} step={step} block={block} index={index} t={t} />
                                        ))}
                                    </div>
                                )}
                            </div>
                        )}
                        {expanded && !isAgentStep && detailBlocks.length > 0 && (
                            <div className={styles.detailSection}>
                                {detailBlocks.map((block, index) => (
                                    <DetailBlockView key={`${block.kind}-${index}`} step={step} block={block} index={index} t={t} />
                                ))}
                            </div>
                        )}
                        {expanded && isAwaitingApproval && approvalMeta && (
                            <div className={styles.approvalCard}>
                                <div className={styles.approvalTitle}>{approvalMeta.title}</div>
                                <div className={styles.approvalBody}>{approvalMeta.message}</div>
                                <div className={styles.approvalActions}>
                                    <button
                                        type="button"
                                        className={`${styles.approvalButton} ${styles.approvalButtonAllow}`}
                                        onClick={() => onApprovalDecision?.(approvalMeta.approvalId, 'allow')}
                                    >
                                        {t('chat.executionTrace.allow')}
                                    </button>
                                    <button
                                        type="button"
                                        className={`${styles.approvalButton} ${styles.approvalButtonReject}`}
                                        onClick={() => onApprovalDecision?.(approvalMeta.approvalId, 'reject')}
                                    >
                                        {t('chat.executionTrace.reject')}
                                    </button>
                                    <button
                                        type="button"
                                        className={`${styles.approvalButton} ${styles.approvalButtonGuide}`}
                                        onClick={() => setShowInlineInput(true)}
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
                                                onClick={() => { setShowInlineInput(false); setInlineComment(''); }}
                                                disabled={isSending}
                                            >
                                                {t('chat.executionTrace.cancelComment')}
                                            </button>
                                        </div>
                                    </div>
                                )}
                            </div>
                        )}
                        {expanded && !isAgentStep && inlineChildren.length > 0 && (
                            <div className={styles.inlineChildrenSection}>
                                <div className={styles.inlineChildrenTitle}>{t('chat.executionTrace.tools')}</div>
                                <div className={styles.inlineChildrenList}>
                                    {inlineChildren.map((child) => (
                                        <TraceNode
                                            key={child.step_id}
                                            step={child}
                                            tree={tree}
                                            nowMs={nowMs}
                                            depth={0}
                                            autoExpand={isRunning}
                                            onApprovalDecision={onApprovalDecision}
                                            onSendApprovalComment={onSendApprovalComment}
                                        />
                                    ))}
                                </div>
                            </div>
                        )}
                    </div>
                </div>
                <div className={styles.nodeMeta}>
                    <span className={`${styles.statusDot} ${getStatusDotClass(step.status)}`} />
                    <span className={styles.elapsed}>{formatElapsedMs(getTraceStepElapsedMs(step, nowMs))}</span>
                </div>
            </div>
            {nestedChildren.length > 0 && expanded && (
                <div className={styles.children}>
                    {nestedChildren.map((child) => (
                        <TraceNode
                            key={child.step_id}
                            step={child}
                            tree={tree}
                            nowMs={nowMs}
                            depth={depth + 1}
                            autoExpand={isRunning}
                            onApprovalDecision={onApprovalDecision}
                            onSendApprovalComment={onSendApprovalComment}
                        />
                    ))}
                </div>
            )}
        </div>
    );
};

const ExecutionTracePanel: React.FC<ExecutionTraceProps> = ({
    trace,
    currentStage,
    retryCount = 0,
    isStreaming = false,
    onApprovalDecision,
    onSendApprovalComment,
}) => {
    const { t } = useTranslation();
    const steps = trace?.steps ?? [];
    const normalizedStage = currentStage?.trim() ?? "";
    const [nowMs, setNowMs] = useState(() => Date.now());
    const [expanded, setExpanded] = useState(false);
    const tree = useMemo(() => buildTree(steps), [steps]);
    const rootSteps = tree.get("__root__") ?? [];
    const hasRunningStep = useMemo(() => steps.some(isRunningStep), [steps]);
    const shouldShowStageOnly = steps.length === 0 && isStreaming && normalizedStage !== "" && normalizedStage !== "chat.stage.finished";
    const hasVisibleTraceContent = steps.length > 0 || shouldShowStageOnly;
    const prevIsStreamingRef = React.useRef(isStreaming);
    const prevHasRunningStepRef = React.useRef(hasRunningStep);
    const currentStageLabel = (() => {
        const key = stageLabelMap[normalizedStage];
        if (key) {
            return t(key);
        }
        if (normalizedStage.includes('.')) {
            return t(normalizedStage);
        }
        if (normalizedStage) {
            return normalizedStage;
        }
        return t('chat.executionTrace.processing');
    })();

    useEffect(() => {
        if (!hasRunningStep) {
            return;
        }
        setNowMs(Date.now());
        const timer = window.setInterval(() => setNowMs(Date.now()), 1000);
        return () => window.clearInterval(timer);
    }, [hasRunningStep]);

    useEffect(() => {
        const wasStreaming = prevIsStreamingRef.current;
        const wasRunning = prevHasRunningStepRef.current;

        if (isStreaming || hasRunningStep) {
            setExpanded(true);
        } else if ((wasStreaming || wasRunning) && steps.length > 0) {
            setExpanded(false);
        }

        prevIsStreamingRef.current = isStreaming;
        prevHasRunningStepRef.current = hasRunningStep;
    }, [hasRunningStep, isStreaming, steps.length]);

    if (!hasVisibleTraceContent) {
        return null;
    }

    return (
        <div className={styles.tracePanel}>
            <button type="button" className={styles.header} onClick={() => setExpanded(!expanded)}>
                <div className={styles.headerMain}>
                    <span className={styles.headerToggle}>{expanded ? "▾" : "▸"}</span>
                    <span className={styles.headerTitle}>{t('chat.executionTrace.title')}</span>
                    <span className={styles.stageBadge}>{currentStageLabel}</span>
                    {retryCount > 0 && <span className={styles.retryBadge}>{t('chat.executionTrace.retryCount', { count: retryCount })}</span>}
                </div>
            </button>
            {expanded && (
                <div className={styles.body}>
                    {rootSteps.map((step) => (
                        <TraceNode
                            key={step.step_id}
                            step={step}
                            tree={tree}
                            nowMs={nowMs}
                            autoExpand={step.status !== TraceStepStatus.TraceStepStatusDone}
                            onApprovalDecision={onApprovalDecision}
                            onSendApprovalComment={onSendApprovalComment}
                        />
                    ))}
                </div>
            )}
        </div>
    );
};

export default ExecutionTracePanel;
