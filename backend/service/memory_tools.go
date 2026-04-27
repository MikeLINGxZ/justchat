package service

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	memorytools "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/tools"
	appstorage "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
)

type SessionSearchToolRequest struct {
	Query string `json:"query" jsonschema:"title=查询,description=要在历史聊天中搜索的关键词或自然语言问题。required=true"`
	Limit int    `json:"limit,omitempty" jsonschema:"title=数量,description=返回条数，默认 5，最大 20"`
}

type SessionSearchToolResponse struct {
	Success bool                             `json:"success"`
	Message string                           `json:"message,omitempty"`
	Results []appstorage.SessionSearchResult `json:"results"`
}

func newSessionSearchTool(storage *appstorage.Storage) (tool.InvokableTool, error) {
	return utils.InferTool(
		"session_search",
		"Search past chat sessions for details that are not part of bounded core memory, such as previous file discussions or one-off task context.",
		func(ctx context.Context, in SessionSearchToolRequest) (SessionSearchToolResponse, error) {
			resp := SessionSearchToolResponse{Success: false, Results: []appstorage.SessionSearchResult{}}
			query := strings.TrimSpace(in.Query)
			if query == "" {
				resp.Message = "query is required"
				return resp, nil
			}
			results, err := storage.SearchSessions(ctx, query, in.Limit)
			if err != nil {
				resp.Message = err.Error()
				return resp, nil
			}
			for i := range results {
				results[i].Content = compactSessionSearchContent(results[i].Content, 360)
			}
			resp.Success = true
			resp.Results = results
			if len(results) == 0 {
				resp.Message = "no matching past session messages found"
			} else {
				resp.Message = "found matching past session messages"
			}
			return resp, nil
		})
}

func compactSessionSearchContent(content string, maxRunes int) string {
	content = strings.TrimSpace(content)
	if maxRunes <= 0 || utf8.RuneCountInString(content) <= maxRunes {
		return content
	}
	runes := []rune(content)
	return string(runes[:maxRunes]) + "..."
}

func (s *Service) buildMemoryRuntimeTools() []tool.BaseTool {
	if s.memoryStorage == nil || s.storage == nil {
		return nil
	}
	coreMemoryTool, err := memorytools.NewCoreMemoryTool(s.memoryStorage)
	if err != nil {
		return nil
	}
	sessionSearchTool, err := newSessionSearchTool(s.storage)
	if err != nil {
		return []tool.BaseTool{coreMemoryTool}
	}
	return []tool.BaseTool{coreMemoryTool, sessionSearchTool}
}
