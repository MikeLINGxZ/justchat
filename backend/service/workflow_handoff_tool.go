package service

import (
	"context"
	"strings"

	einotool "github.com/cloudwego/eino/components/tool"
	toolutils "github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

const workflowHandoffToolName = "create_workflow_task"

type workflowHandoff struct {
	Reason   string
	Summary  string
	Source   routeSource
	RuleName string
}

type workflowHandoffToolParams struct {
	Reason  string `json:"reason"`
	Summary string `json:"summary"`
}

func newWorkflowHandoffTool(onHandoff func(workflowHandoff)) einotool.BaseTool {
	return toolutils.NewTool(
		&schema.ToolInfo{
			Name: workflowHandoffToolName,
			Desc: "当用户请求需要任务拆解、多步骤处理、文件分析、工具调用或工作流执行时，调用此工具交付任务编排。调用前不要先输出面向用户的最终答案。参数 reason 描述为何不能直接回答，summary 描述后续 workflow 需要完成什么。",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"reason": {
					Type:     "string",
					Desc:     "为什么当前请求不能直接回答",
					Required: true,
				},
				"summary": {
					Type:     "string",
					Desc:     "后续 workflow 需要解决什么",
					Required: true,
				},
			}),
		},
		func(ctx context.Context, params workflowHandoffToolParams) (string, error) {
			handoff := workflowHandoff{
				Reason:  strings.TrimSpace(params.Reason),
				Summary: strings.TrimSpace(params.Summary),
				Source:  routeSourceMainModel,
			}
			if handoff.Reason == "" {
				handoff.Reason = "主模型判断该请求需要进入任务编排"
			}
			if handoff.Summary == "" {
				handoff.Summary = handoff.Reason
			}
			if onHandoff != nil {
				onHandoff(handoff)
			}
			return "workflow task created", nil
		},
	)
}
