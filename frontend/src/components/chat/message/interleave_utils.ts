import {
    RouteType,
    TraceStepType,
    type AssistantMessageExtra,
    type ToolUse,
    type TraceStep,
} from "@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models/models";

export type ToolSegment = {
    type: 'tool';
    toolUse: ToolUse;
    traceStep?: TraceStep;
    childSegments: ToolSegment[];
    key: string;
};

export type ContentSegment =
    | { type: 'text'; text: string; key: string }
    | ToolSegment;

/**
 * Detect if the message is in "direct mode" (single-pass, no workflow).
 * Direct mode messages should use interleaved content + tool card display.
 */
export function isDirectMode(extra: AssistantMessageExtra | undefined | null): boolean {
    if (!extra) {
        return true;
    }
    if (extra.route_type === RouteType.RouteTypeWorkflow) {
        return false;
    }
    if (extra.route_type === RouteType.RouteTypeDirectAnswer) {
        return true;
    }
    // Fallback heuristic: if all trace steps are tool_call type, treat as direct mode
    const steps = extra.execution_trace?.steps ?? [];
    if (steps.length === 0) {
        return true;
    }
    return steps.every(s => s.type === TraceStepType.TraceStepTypeToolCall);
}

/**
 * Split message content into interleaved text and tool segments based on
 * each tool use's content_pos (rune position where the tool was invoked).
 * Sub-agent child tools are nested under their parent sub-agent's segment.
 */
export function buildInterleavedSegments(
    content: string,
    toolUses: ToolUse[],
    traceSteps: TraceStep[],
): ContentSegment[] {
    if (toolUses.length === 0) {
        if (content) {
            return [{ type: 'text', text: content, key: 'text-0' }];
        }
        return [];
    }

    // Build lookup from step_id (= call_id) to trace step
    const traceByCallID = new Map<string, TraceStep>();
    for (const step of traceSteps) {
        if (step.type === TraceStepType.TraceStepTypeToolCall && step.step_id) {
            traceByCallID.set(step.step_id, step);
        }
    }

    // Identify sub-agent call_ids (parents that have children)
    const subAgentCallIDs = new Set<string>();
    for (const step of traceSteps) {
        if (step.type === TraceStepType.TraceStepTypeToolCall && step.step_id) {
            const meta = step.metadata as Record<string, unknown> | undefined;
            if (meta?.is_sub_agent === true) {
                subAgentCallIDs.add(step.step_id);
            }
        }
    }

    // Group child tool uses by parent_step_id
    const childrenByParent = new Map<string, ToolUse[]>();
    const rootToolUses: ToolUse[] = [];

    for (const toolUse of toolUses) {
        const traceStep = traceByCallID.get(toolUse.call_id);
        const parentID = traceStep?.parent_step_id ?? "";
        if (parentID && subAgentCallIDs.has(parentID)) {
            const children = childrenByParent.get(parentID) ?? [];
            children.push(toolUse);
            childrenByParent.set(parentID, children);
        } else {
            rootToolUses.push(toolUse);
        }
    }

    // Build child segments for each sub-agent
    function buildChildSegments(parentCallID: string): ToolSegment[] {
        const children = childrenByParent.get(parentCallID) ?? [];
        return children.map((childUse): ToolSegment => ({
            type: 'tool',
            toolUse: childUse,
            traceStep: traceByCallID.get(childUse.call_id),
            childSegments: buildChildSegments(childUse.call_id),
            key: `tool-${childUse.call_id || childUse.index}`,
        }));
    }

    // Sort root tool uses by content_pos ascending
    const sorted = [...rootToolUses].sort((a, b) => a.content_pos - b.content_pos);

    const runes = Array.from(content);
    const segments: ContentSegment[] = [];
    let cursor = 0;

    for (const toolUse of sorted) {
        const pos = Math.max(0, Math.min(toolUse.content_pos, runes.length));

        if (pos > cursor) {
            const text = runes.slice(cursor, pos).join('');
            segments.push({ type: 'text', text, key: `text-${cursor}-${pos}` });
        }

        segments.push({
            type: 'tool',
            toolUse,
            traceStep: traceByCallID.get(toolUse.call_id),
            childSegments: buildChildSegments(toolUse.call_id),
            key: `tool-${toolUse.call_id || toolUse.index}`,
        });

        cursor = pos;
    }

    // Remaining content after the last tool
    if (cursor < runes.length) {
        const text = runes.slice(cursor).join('');
        segments.push({ type: 'text', text, key: `text-${cursor}-end` });
    }

    return segments;
}
