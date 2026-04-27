package agents

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
)

func uniqueToolsByInfoName(ctx context.Context, tools []tool.BaseTool) ([]tool.BaseTool, error) {
	if len(tools) < 2 {
		return tools, nil
	}

	seen := make(map[string]struct{}, len(tools))
	unique := make([]tool.BaseTool, 0, len(tools))
	for _, item := range tools {
		if item == nil {
			continue
		}

		info, err := item.Info(ctx)
		if err != nil {
			return nil, fmt.Errorf("read tool info for duplicate check: %w", err)
		}
		if info == nil {
			return nil, fmt.Errorf("tool info is nil")
		}

		if _, ok := seen[info.Name]; ok {
			logger.Warmf("skip duplicate tool name: %s", info.Name)
			continue
		}
		seen[info.Name] = struct{}{}
		unique = append(unique, item)
	}

	return unique, nil
}
