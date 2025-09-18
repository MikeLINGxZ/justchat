package internal

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/flow/agent/multiagent/host"
)

func NewHost(ctx context.Context, baseURL, apiKey, modelName string) (*host.Host, error) {
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: baseURL,
		Model:   modelName,
		APIKey:  apiKey,
	})
	if err != nil {
		return nil, err
	}

	return &host.Host{
		ToolCallingModel: chatModel,
		SystemPrompt:     "You can help the user store and retrieve memories by creating, reading, and organizing journal entries on their behalf. When the user asks a question or shares a thought, always respond using relevant journal content, and save meaningful experiences into memory with context and emotion preserved.",
	}, nil
}
