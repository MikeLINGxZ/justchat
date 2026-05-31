package agent

import (
	"context"

	pkgProvider "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/provider"
)

type runContextKey struct{}

type runContext struct {
	SessionID    uint
	ModelName    string
	BaseURL      string
	APIKey       string
	ProviderType pkgProvider.Type
}

type toolConfirmation struct {
	Approved bool
	Message  string
	Action   string
}

type toolConfirmationContextKey struct{}

func withRunContext(ctx context.Context, info runContext) context.Context {
	return context.WithValue(ctx, runContextKey{}, info)
}

func getRunContext(ctx context.Context) (runContext, bool) {
	info, ok := ctx.Value(runContextKey{}).(runContext)
	return info, ok
}

func withToolConfirmation(ctx context.Context, confirmation toolConfirmation) context.Context {
	return context.WithValue(ctx, toolConfirmationContextKey{}, confirmation)
}

func getToolConfirmation(ctx context.Context) (toolConfirmation, bool) {
	info, ok := ctx.Value(toolConfirmationContextKey{}).(toolConfirmation)
	return info, ok
}
