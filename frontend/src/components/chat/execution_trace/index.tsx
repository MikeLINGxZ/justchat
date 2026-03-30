import React, { useEffect, useMemo, useState } from "react";
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
    "": "处理中",
    "等待执行": "等待执行",
    "准备执行": "准备执行",
    "意图识别": "意图识别",
    "直接回答": "直接回答",
    "任务拆解": "任务拆解",
    "子任务执行": "子任务执行",
    "等待用户确认": "等待用户确认",
    "结果汇总": "结果汇总",
    "结果审核": "结果审核",
    "重新生成": "重新生成",
    "已完成": "已完成",
};

const typeLabelMap: Record<string, string> = {
    [TraceStepType.TraceStepTypeClassify]: "意图识别",
    [TraceStepType.TraceStepTypePlan]: "任务拆解",
    [TraceStepType.TraceStepTypeDispatch]: "分派任务",
    [TraceStepType.TraceStepTypeAgentRun]: "子 Agent",
    [TraceStepType.TraceStepTypeToolCall]: "工具调用",
    [TraceStepType.TraceStepTypeSynthesize]: "结果汇总",
    [TraceStepType.TraceStepTypeReview]: "结果审核",
    [TraceStepType.TraceStepTypeRetry]: "重新生成",
    [TraceStepType.TraceStepTypeFinalize]: "输出答案",
};

function getStatusLabel(status?: string): string {
    switch (status) {
        case TraceStepStatus.TraceStepStatusDone:
            return "已完成";
        case TraceStepStatus.TraceStepStatusError:
            return "失败";
        case TraceStepStatus.TraceStepStatusSkipped:
            return "跳过";
        case TraceStepStatus.TraceStepStatusPending:
            return "准备中";
        case TraceStepStatus.TraceStepStatusAwaitingApproval:
            return "等待确认";
        case TraceStepStatus.TraceStepStatusRejected:
            return "已拒绝";
        default:
            return "执行中";
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

function getDetailBlockDisplayTitle(step: TraceStep, block: TraceDetailBlock): string {
    if (step.type !== TraceStepType.TraceStepTypeAgentRun) {
        return block.title;
    }
    if (block.kind === "input") {
        return "用户输入";
    }
    if (block.kind === "output") {
        return "最终回答";
    }
    return block.title;
}

function renderDetailBlock(step: TraceStep, block: TraceDetailBlock, index: number) {
    const displayTitle = getDetailBlockDisplayTitle(step, block);
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
    const title = typeof metadata.approval_title === "string" ? metadata.approval_title : (step.title || "工具确认");
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
                            <span className={styles.nodeType}>{typeLabelMap[step.type] || step.type || "步骤"}</span>
                            <span className={styles.nodeTitle}>{step.title || "未命名步骤"}</span>
                        </div>
                        {(hasChildren || detailBlocks.length > 0) && (
                            <button className={styles.toggle} type="button" onClick={() => setExpanded(!expanded)}>
                                <span className={styles.toggleIcon}>{expanded ? "▾" : "▸"}</span>
                                <span className={styles.toggleText}>{expanded ? "折叠" : "展开"}</span>
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
                                    {leadingDetailBlocks.map((block, index) => renderDetailBlock(step, block, index))}
                                </div>
                            )}
                            {(inlineChildren.length > 0 || isRunning) && (
                                <div className={styles.inlineChildrenSection}>
                                    <div className={styles.inlineChildrenTitle}>思考过程</div>
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
                                        <div className={styles.agentThinkingPlaceholder}>Agent 正在处理中...</div>
                                    )}
                                </div>
                            )}
                            {trailingDetailBlocks.length > 0 && (
                                <div className={styles.detailSection}>
                                    {trailingDetailBlocks.map((block, index) => renderDetailBlock(step, block, index))}
                                </div>
                            )}
                        </div>
                    )}
                    {expanded && !isAgentStep && detailBlocks.length > 0 && (
                        <div className={styles.detailSection}>
                            {detailBlocks.map((block, index) => renderDetailBlock(step, block, index))}
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
                                    1. 允许
                                </button>
                                <button
                                    type="button"
                                    className={styles.approvalButton}
                                    onClick={() => onApprovalDecision?.(approvalMeta.approvalId, 'reject')}
                                >
                                    2. 拒绝
                                </button>
                                <button
                                    type="button"
                                    className={styles.approvalButton}
                                    onClick={() => onApprovalComment?.(approvalMeta.approvalId, approvalMeta.title, approvalMeta.message)}
                                >
                                    3. 告诉ai应该怎么做
                                </button>
                            </div>
                        </div>
                    )}
                    {expanded && !isAgentStep && inlineChildren.length > 0 && (
                        <div className={styles.inlineChildrenSection}>
                            <div className={styles.inlineChildrenTitle}>工具调用</div>
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
                    <span className={`${styles.status} ${statusClassName}`}>{getStatusLabel(step.status)}</span>
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
    const steps = trace?.steps ?? [];
    const [nowMs, setNowMs] = useState(() => Date.now());
    const [expanded, setExpanded] = useState(false);
    const tree = useMemo(() => buildTree(steps), [steps]);
    const rootSteps = tree.get("__root__") ?? [];
    const hasRunningStep = useMemo(() => steps.some(isRunningStep), [steps]);
    const prevIsStreamingRef = React.useRef(isStreaming);
    const prevHasRunningStepRef = React.useRef(hasRunningStep);

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
                    <span className={styles.headerTitle}>执行轨迹</span>
                    <span className={styles.stageBadge}>{stageLabelMap[currentStage || ""] || currentStage || "处理中"}</span>
                    {retryCount > 0 && <span className={styles.retryBadge}>已重试 {retryCount} 次</span>}
                </div>
                <span className={styles.headerToggle}>{expanded ? "收起" : "展开"}</span>
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
