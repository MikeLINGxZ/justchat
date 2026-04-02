import React, { useEffect, useMemo, useState } from "react";
import type { ContentSegment } from "../interleave_utils";
import { isRunningStep } from "@/components/chat/execution_trace";
import MarkdownRenderer from "@/components/markdown_renderer";
import InlineToolCard from "../inline_tool_card";

interface InterleavedContentProps {
    segments: ContentSegment[];
    isStreaming?: boolean;
    onApprovalDecision?: (approvalId: string, decision: 'allow' | 'reject') => void;
    onSendApprovalComment?: (approvalId: string, comment: string) => Promise<void> | void;
}

const InterleavedContent: React.FC<InterleavedContentProps> = ({
    segments,
    isStreaming = false,
    onApprovalDecision,
    onSendApprovalComment,
}) => {
    const [nowMs, setNowMs] = useState(() => Date.now());

    // Check if any tool is running for timer
    const hasRunningTool = useMemo(() => {
        return segments.some(seg => {
            if (seg.type !== 'tool') return false;
            if (seg.traceStep) return isRunningStep(seg.traceStep);
            const status = seg.toolUse.status;
            return status === "running" || status === "pending" || status === "awaiting_approval";
        });
    }, [segments]);

    useEffect(() => {
        if (!hasRunningTool) return;
        setNowMs(Date.now());
        const timer = window.setInterval(() => setNowMs(Date.now()), 1000);
        return () => window.clearInterval(timer);
    }, [hasRunningTool]);

    if (segments.length === 0) {
        return null;
    }

    return (
        <>
            {segments.map((seg) => {
                if (seg.type === 'text') {
                    if (!seg.text.trim()) return null;
                    return (
                        <MarkdownRenderer
                            key={seg.key}
                            content={seg.text}
                            variant="assistant"
                        />
                    );
                }
                return (
                    <InlineToolCard
                        key={seg.key}
                        toolUse={seg.toolUse}
                        traceStep={seg.traceStep}
                        childSegments={seg.childSegments}
                        nowMs={nowMs}
                        isStreaming={isStreaming}
                        onApprovalDecision={onApprovalDecision}
                        onSendApprovalComment={onSendApprovalComment}
                    />
                );
            })}
        </>
    );
};

export default InterleavedContent;
