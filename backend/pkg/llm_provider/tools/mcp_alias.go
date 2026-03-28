package tools

import (
	"context"
	"fmt"

	einotool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type AliasedTool struct {
	AliasID     string
	DisplayName string
	Desc        string
	Base        einotool.BaseTool
}

func (a *AliasedTool) Id() string {
	return a.AliasID
}

func (a *AliasedTool) Name() string {
	return a.DisplayName
}

func (a *AliasedTool) Description() string {
	return a.Desc
}

func (a *AliasedTool) Tool() einotool.BaseTool {
	return a
}

func (a *AliasedTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	if a.Base == nil {
		return nil, fmt.Errorf("base tool is nil")
	}
	info, err := a.Base.Info(ctx)
	if err != nil {
		return nil, err
	}
	cloned := *info
	cloned.Name = a.AliasID
	if a.Desc != "" {
		cloned.Desc = a.Desc
	}
	return &cloned, nil
}

func (a *AliasedTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...einotool.Option) (string, error) {
	invokable, ok := a.Base.(einotool.InvokableTool)
	if !ok {
		return "", fmt.Errorf("base tool %s is not invokable", a.AliasID)
	}
	return invokable.InvokableRun(ctx, argumentsInJSON, opts...)
}
