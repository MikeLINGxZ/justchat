package agent

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	pkgProvider "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/provider"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/model/openai"
)

const defaultOneshotTimeout = 60 * time.Second

// OneshotRequest describes one non-streaming LLM call.
type OneshotRequest struct {
	BaseURL      string
	APIKey       string
	ModelName    string
	ProviderType pkgProvider.Type
	System       string
	User         string
	MaxTokens    int
	Timeout      time.Duration
}

// OneshotResponse returns one model response plus basic usage.
type OneshotResponse struct {
	Text             string
	PromptTokens     int
	CompletionTokens int
}

// OneshotComplete runs one non-streaming model request without binding any tools.
func OneshotComplete(ctx context.Context, req OneshotRequest) (OneshotResponse, error) {
	timeout := req.Timeout
	if timeout <= 0 {
		timeout = defaultOneshotTimeout
	}
	callCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	mdl := openai.New(
		req.ModelName,
		openai.WithBaseURL(req.BaseURL),
		openai.WithAPIKey(req.APIKey),
		openai.WithVariant(inferOpenAIVariant(req.ProviderType)),
	)

	messages := make([]model.Message, 0, 2)
	if strings.TrimSpace(req.System) != "" {
		messages = append(messages, model.NewSystemMessage(req.System))
	}
	messages = append(messages, model.NewUserMessage(req.User))

	request := &model.Request{
		Messages: messages,
		GenerationConfig: model.GenerationConfig{
			Stream: false,
		},
	}
	if req.MaxTokens > 0 {
		request.MaxTokens = &req.MaxTokens
	}

	responseChan, err := mdl.GenerateContent(callCtx, request)
	if err != nil {
		if errors.Is(callCtx.Err(), context.DeadlineExceeded) {
			return OneshotResponse{}, callCtx.Err()
		}
		return OneshotResponse{}, err
	}

	var builder strings.Builder
	var usage *model.Usage
	for response := range responseChan {
		if response == nil {
			continue
		}
		if response.Error != nil {
			return OneshotResponse{}, fmt.Errorf("oneshot llm error: %s", response.Error.Message)
		}
		if response.Usage != nil {
			usage = response.Usage
		}
		for _, choice := range response.Choices {
			if choice.Message.Content != "" {
				builder.WriteString(choice.Message.Content)
				continue
			}
			if choice.Delta.Content != "" {
				builder.WriteString(choice.Delta.Content)
			}
		}
	}
	if errors.Is(callCtx.Err(), context.DeadlineExceeded) {
		return OneshotResponse{}, callCtx.Err()
	}

	out := OneshotResponse{
		Text: strings.TrimSpace(builder.String()),
	}
	if usage != nil {
		out.PromptTokens = usage.PromptTokens
		out.CompletionTokens = usage.CompletionTokens
	}
	return out, nil
}
