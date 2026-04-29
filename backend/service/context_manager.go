package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/llm_provider"
)

type ContextLayer int

const (
	LayerTask ContextLayer = iota
	LayerStage
	LayerTool
)

type ContextConfig struct {
	TaskTimeout  time.Duration
	StageTimeout time.Duration
	ToolTimeout  time.Duration
}

func DefaultContextConfig() ContextConfig {
	return ContextConfig{
		TaskTimeout:  30 * time.Minute,
		StageTimeout: 5 * time.Minute,
		ToolTimeout:  5 * time.Minute,
	}
}

type ContextManager struct {
	config ContextConfig
}

func NewContextManager(config ContextConfig) *ContextManager {
	return &ContextManager{config: config}
}

func (m *ContextManager) NewTaskContext(parent context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, m.config.TaskTimeout)
}

func (m *ContextManager) NewStageContext(parent context.Context, stage string) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, m.config.StageTimeout)
}

func (m *ContextManager) NewToolContext(parent context.Context, toolName string) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parent, m.config.ToolTimeout)
}

type TokenEstimator struct {
	charsPerToken float64
}

func NewTokenEstimator() *TokenEstimator {
	return &TokenEstimator{charsPerToken: 3.5}
}

func (e *TokenEstimator) EstimateTokens(messages []schema.Message) int {
	total := 0
	for _, msg := range messages {
		total += e.EstimateMessageTokens(msg)
	}
	return total
}

func (e *TokenEstimator) EstimateMessageTokens(msg schema.Message) int {
	chars := len([]rune(msg.Content))
	chars += len([]rune(msg.ReasoningContent))
	if msg.UserInputMultiContent != nil {
		for _, part := range msg.UserInputMultiContent {
			chars += len([]rune(part.Text))
		}
	}
	for _, tc := range msg.ToolCalls {
		chars += len([]rune(tc.Function.Name))
		chars += len([]rune(tc.Function.Arguments))
	}
	return int(float64(chars) / e.charsPerToken)
}

type ContextWindowConfig struct {
	MaxTokens     int
	ReserveTokens int
}

func DefaultContextWindowConfig() ContextWindowConfig {
	return ContextWindowConfig{
		MaxTokens:     128000,
		ReserveTokens: 4096,
	}
}

type ContextCompressor struct {
	config         ContextWindowConfig
	estimator      *TokenEstimator
	compressPrompt string
}

func NewContextCompressor(config ContextWindowConfig) *ContextCompressor {
	return &ContextCompressor{
		config:    config,
		estimator: NewTokenEstimator(),
		compressPrompt: `请将以下对话历史压缩为简洁摘要，保留关键信息和上下文：

%s

请输出一个紧凑的摘要（不超过500字）：`,
	}
}

func (c *ContextCompressor) NeedsCompression(messages []schema.Message) bool {
	tokens := c.estimator.EstimateTokens(messages)
	return tokens > c.config.MaxTokens-c.config.ReserveTokens
}

func (c *ContextCompressor) CompressIfNeeded(
	ctx context.Context,
	messages []schema.Message,
	provider *llm_provider.Provider,
) ([]schema.Message, error) {
	if !c.NeedsCompression(messages) {
		return messages, nil
	}
	availableTokens := c.config.MaxTokens - c.config.ReserveTokens
	recentCount := c.countMessagesForTokens(messages, availableTokens/2)
	if recentCount <= 2 {
		recentCount = 2
	}
	olderMessages := messages[:len(messages)-recentCount]
	recentMessages := messages[len(messages)-recentCount:]
	if len(olderMessages) == 0 {
		return messages, nil
	}
	summary, err := c.compressMessages(ctx, olderMessages, provider)
	if err != nil {
		return messages, nil
	}
	result := append([]schema.Message{{
		Role:    schema.System,
		Content: fmt.Sprintf("对话历史摘要：\n%s", summary),
	}}, recentMessages...)
	return result, nil
}

func (c *ContextCompressor) compressMessages(
	ctx context.Context,
	messages []schema.Message,
	provider *llm_provider.Provider,
) (string, error) {
	history := c.messagesToText(messages)
	prompt := fmt.Sprintf(c.compressPrompt, history)
	return "", fmt.Errorf("compress not fully implemented: %s", prompt)
}

func (c *ContextCompressor) messagesToText(messages []schema.Message) string {
	var sb strings.Builder
	for _, msg := range messages {
		role := string(msg.Role)
		content := msg.Content
		if content == "" {
			continue
		}
		sb.WriteString(fmt.Sprintf("[%s]: %s\n", role, content))
	}
	return sb.String()
}

func (c *ContextCompressor) countMessagesForTokens(messages []schema.Message, maxTokens int) int {
	tokens := 0
	for i := len(messages) - 1; i >= 0; i-- {
		tokens += c.estimator.EstimateMessageTokens(messages[i])
		if tokens > maxTokens {
			return len(messages) - i - 1
		}
	}
	return len(messages)
}
