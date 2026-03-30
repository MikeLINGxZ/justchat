import React, { useEffect, useMemo, useState } from "react";
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
    onApprovalComment?: (approvalId: string, title: string, message: string) => void;
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
        return `${sec}s`;
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

function renderDetailBlock(step: TraceStep, block: TraceDetailBlock, index: number, t: (key: string) => string) {
    const displayTitle = getDetailBlockDisplayTitle(step, block, t);
    if (!block.content?.trim()) {
        return null;
    }
    if (block.format === TraceDetailFormat.TraceDetailFormatJSON) {
        return (
            <div key={`${block.kind}-${index}`} className={styles.detailBlock}>
                <div className={styles.detailTitle}>{displayTitle}</div>
                <pre className={`${styles.detailContent} ${styles.jsonBlock}`}>{block.content}</pre>
            </div>
        );
    }
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
            <pre className={styles.detailContent}>{block.content}</pre>
        </div>
    );
}

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

const TraceNode: React.FC<{
    step: TraceStep;
    tree: Map<string, TraceStep[]>;
    nowMs: number;
    depth?: number;
    autoExpand?: boolean;
    onApprovalDecision?: (approvalId: string, decision: 'allow' | 'reject') => void;
    onApprovalComment?: (approvalId: string, title: string, message: string) => void;
}> = ({ step, tree, nowMs, depth = 0, autoExpand = false, onApprovalDecision, onApprovalComment }) => {
    const { t } = useTranslation();
    const children = tree.get(step.step_id) ?? [];
    const inlineChildren = step.type === TraceStepType.TraceStepTypeAgentRun ? children : [];
    const nestedChildren = step.type === TraceStepType.TraceStepTypeAgentRun ? [] : children;
    const [expanded, setExpanded] = useState(autoExpand);
    const hasChildren = children.length > 0;
    const isRunning = isRunningStep(step);
    const statusClassName = isRunning
        ? styles.running
        : step.status === TraceStepStatus.TraceStepStatusError || step.status === TraceStepStatus.TraceStepStatusRejected
            ? styles.error
            : styles.done;
    const detailBlocks = step.detail_blocks ?? [];
    const approvalMeta = getApprovalMeta(step);
    const isAwaitingApproval = step.status === TraceStepStatus.TraceStepStatusAwaitingApproval && approvalMeta != null;
    const isAgentStep = step.type === TraceStepType.TraceStepTypeAgentRun;
    const leadingDetailBlocks = isAgentStep
        ? detailBlocks.filter((block) => block.kind !== "output" && block.kind !== "tool_result")
        : detailBlocks;
    const trailingDetailBlocks = isAgentStep
        ? detailBlocks.filter((block) => block.kind === "output" || block.kind === "tool_result")
        : [];

    useEffect(() => {
        if (isRunning) {
            setExpanded(true);
        }
    }, [isRunning]);

    return (
        <div className={styles.node} style={{ marginLeft: depth * 14 }}>
            <div className={styles.nodeHeader}>
                <span className={styles.toggleSpacer} />
                <div className={styles.nodeMain}>
                    <div className={styles.nodeTopRow}>
                        <div className={styles.nodeTitleRow}>
                            <span className={styles.nodeType}>{t(typeLabelMap[step.type] || 'chat.executionTrace.step')}</span>
                            <span className={styles.nodeTitle}>{step.title || t('chat.executionTrace.unnamedStep')}</span>
                        </div>
                        {(hasChildren || detailBlocks.length > 0) && (
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
                    {expanded && isAgentStep && (
                        <div className={styles.agentFlowSection}>
                            {leadingDetailBlocks.length > 0 && (
                                <div className={styles.detailSection}>
                                    {leadingDetailBlocks.map((block, index) => renderDetailBlock(step, block, index, t))}
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
                                                    onApprovalComment={onApprovalComment}
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
                                    {trailingDetailBlocks.map((block, index) => renderDetailBlock(step, block, index, t))}
                                </div>
                            )}
                        </div>
                    )}
                    {expanded && !isAgentStep && detailBlocks.length > 0 && (
                        <div className={styles.detailSection}>
                            {detailBlocks.map((block, index) => renderDetailBlock(step, block, index, t))}
                        </div>
                    )}
                    {expanded && isAwaitingApproval && approvalMeta && (
                        <div className={styles.approvalCard}>
                            <div className={styles.approvalTitle}>{approvalMeta.title}</div>
                            <div className={styles.approvalBody}>{approvalMeta.message}</div>
                            <div className={styles.approvalActions}>
                                <button
                                    type="button"
                                    className={styles.approvalButton}
                                    onClick={() => onApprovalDecision?.(approvalMeta.approvalId, 'allow')}
                                >
                                    {t('chat.executionTrace.allow')}
                                </button>
                                <button
                                    type="button"
                                    className={styles.approvalButton}
                                    onClick={() => onApprovalDecision?.(approvalMeta.approvalId, 'reject')}
                                >
                                    {t('chat.executionTrace.reject')}
                                </button>
                                <button
                                    type="button"
                                    className={styles.approvalButton}
                                    onClick={() => onApprovalComment?.(approvalMeta.approvalId, approvalMeta.title, approvalMeta.message)}
                                >
                                    {t('chat.executionTrace.guideAi')}
                                </button>
                            </div>
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
                                        onApprovalComment={onApprovalComment}
                                    />
                                ))}
                            </div>
                        </div>
                    )}
                </div>
                <div className={styles.nodeMeta}>
                    <span className={`${styles.status} ${statusClassName}`}>{getStatusLabel(step.status, t)}</span>
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
                            onApprovalComment={onApprovalComment}
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
    onApprovalComment,
}) => {
    const { t } = useTranslation();
    const steps = trace?.steps ?? [];
    const [nowMs, setNowMs] = useState(() => Date.now());
    const [expanded, setExpanded] = useState(false);
    const tree = useMemo(() => buildTree(steps), [steps]);
    const rootSteps = tree.get("__root__") ?? [];
    const hasRunningStep = useMemo(() => steps.some(isRunningStep), [steps]);
    const prevIsStreamingRef = React.useRef(isStreaming);
    const prevHasRunningStepRef = React.useRef(hasRunningStep);
    const currentStageLabel = (() => {
        const key = stageLabelMap[currentStage || ""];
        if (key) {
            return t(key);
        }
        if (currentStage?.includes('.')) {
            return t(currentStage);
        }
        if (currentStage) {
            return currentStage;
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

    if (steps.length === 0 && !currentStage) {
        return null;
    }

    return (
        <div className={styles.tracePanel}>
            <button type="button" className={styles.header} onClick={() => setExpanded(!expanded)}>
                <div className={styles.headerMain}>
                    <span className={styles.headerTitle}>{t('chat.executionTrace.title')}</span>
                    <span className={styles.stageBadge}>{currentStageLabel}</span>
                    {retryCount > 0 && <span className={styles.retryBadge}>{t('chat.executionTrace.retryCount', { count: retryCount })}</span>}
                </div>
                <span className={styles.headerToggle}>{expanded ? t('chat.executionTrace.collapse') : t('chat.executionTrace.expand')}</span>
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
                            onApprovalComment={onApprovalComment}
                        />
                    ))}
                </div>
            )}
        </div>
    );
};

export default ExecutionTracePanel;
