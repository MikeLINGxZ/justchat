package agent

import (
	"context"
	"fmt"
	"os"
	"strings"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	pkgAgent "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/agent"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/agent/tools"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/id/event_id"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/id/prompt_id"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/ierror"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/prompt"
	pkgterminal "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/terminal"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/agent/agent_dto"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

// Agent exposes backend agent operations to the frontend.
type Agent struct {
	manager *pkgAgent.Manager
}

// Dependencies groups optional agent integrations that are only needed at app startup.
type Dependencies struct {
	SkillProvider        tools.SkillProvider
	AttentionRequester   tools.AttentionRequester
	SkillCreator         tools.SkillCreator
	CliInstaller         tools.CliInstaller
	CliInstallProgress   tools.CliInstallProgressReporter
	CliManifestGenerator tools.CliManifestGenerator
	CliCommandRunner     tools.CliCommandRunner
	TerminalRunner       *pkgterminal.Manager
}

// NewAgent creates an agent service backed by application storage.
func NewAgent(istorage *storage.Storage, deps ...Dependencies) *Agent {
	svc := &Agent{manager: pkgAgent.NewManager(istorage)}
	if len(deps) > 0 {
		if deps[0].SkillProvider != nil {
			svc.manager.SetSkillProvider(deps[0].SkillProvider)
		}
		if deps[0].AttentionRequester != nil {
			svc.manager.SetAttentionRequester(deps[0].AttentionRequester)
		}
		if deps[0].SkillCreator != nil {
			svc.manager.SetSkillCreator(deps[0].SkillCreator)
		}
		if deps[0].CliInstaller != nil {
			svc.manager.SetCliInstaller(deps[0].CliInstaller)
		}
		if deps[0].CliInstallProgress != nil {
			svc.manager.SetCliInstallProgressReporter(deps[0].CliInstallProgress)
		}
		if deps[0].CliManifestGenerator != nil {
			svc.manager.SetCliManifestGenerator(deps[0].CliManifestGenerator)
		}
		if deps[0].CliCommandRunner != nil {
			svc.manager.SetCliCommandRunner(deps[0].CliCommandRunner)
		}
		if deps[0].TerminalRunner != nil {
			svc.manager.SetTerminalRunner(deps[0].TerminalRunner)
		}
	}
	return svc
}

func (a *Agent) SendMessage(ctx context.Context, input agent_dto.SendMessageInput) (*agent_dto.SendMessageOutput, error) {
	atts, err := convertAttachments(input.Attachments)
	if err != nil {
		return nil, ierror.Error(ierror.ErrAgentSendMessage, err)
	}

	handler := a.manager.NewChatHandler()
	err = handler.SendMessage(ctx, pkgAgent.SendMessageParams{
		SessionID:        input.SessionID,
		Content:          input.Content,
		SystemPrompt:     input.SystemPrompt,
		PrimedSkill:      input.SkillName,
		BaseURL:          input.BaseURL,
		ApiKey:           input.ApiKey,
		ModelName:        input.ModelName,
		ProviderType:     input.ProviderType,
		EnabledUserTools: input.EnabledUserTools,
		Attachments:      atts,
	})
	if err != nil {
		return nil, ierror.Error(ierror.ErrAgentSendMessage, err)
	}
	return &agent_dto.SendMessageOutput{}, nil
}

func (a *Agent) StopGeneration(ctx context.Context, input agent_dto.StopGenerationInput) (*agent_dto.StopGenerationOutput, error) {
	a.manager.Streams().Stop(input.SessionID)
	_ = a.manager.Storage().UpdateSessionStatus(input.SessionID, "idle")
	return &agent_dto.StopGenerationOutput{}, nil
}

func (a *Agent) RespondToConfirm(ctx context.Context, input agent_dto.RespondToConfirmInput) (*agent_dto.RespondToConfirmOutput, error) {
	ok := a.manager.Streams().SendConfirmResponse(input.SessionID, pkgAgent.ConfirmResponse{
		Approved: input.Approved,
		Message:  input.Message,
		Action:   input.Action,
	})
	if !ok {
		return nil, ierror.Error(ierror.ErrAgentNoConfirm, fmt.Errorf("no active confirmation for session %d", input.SessionID))
	}
	return &agent_dto.RespondToConfirmOutput{}, nil
}

func (a *Agent) CreateSession(ctx context.Context, input agent_dto.CreateSessionInput) (*agent_dto.CreateSessionOutput, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		title = "New Chat"
	}
	tags := normalizeSessionTags(input.Tags)
	kind := "user"
	if tagsContain(tags, "task") {
		kind = "task"
	}

	session, err := a.manager.Storage().CreateSession(data_models.Session{
		Title:  title,
		Status: "idle",
		Kind:   kind,
		Tags:   marshalSessionTags(tags),
	})
	if err != nil {
		return nil, ierror.Error(ierror.ErrAgentCreateSession, err)
	}

	return &agent_dto.CreateSessionOutput{
		SessionID: session.ID,
		Title:     session.Title,
		Tags:      tags,
	}, nil
}

func (a *Agent) ListSessions(ctx context.Context, input agent_dto.ListSessionsInput) (*agent_dto.ListSessionsOutput, error) {
	limit := input.Limit
	if limit <= 0 {
		limit = 20
	}

	sessions, err := a.manager.Storage().ListSessions(input.Cursor, limit+1, input.StarredOnly, input.IncludeHidden)
	if err != nil {
		return nil, ierror.Error(ierror.ErrAgentListSessions, err)
	}

	hasMore := len(sessions) > limit
	if hasMore {
		sessions = sessions[:limit]
	}

	items := make([]agent_dto.SessionItem, 0, len(sessions))
	for _, session := range sessions {
		items = append(items, toSessionItem(session))
	}

	var nextCursor uint
	if hasMore && len(sessions) > 0 {
		nextCursor = sessions[len(sessions)-1].ID
	}

	return &agent_dto.ListSessionsOutput{
		Sessions:   items,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

// SpawnTaskSession creates a task session and asynchronously submits its first turn.
func (a *Agent) SpawnTaskSession(ctx context.Context, input agent_dto.SpawnTaskSessionInput) (*agent_dto.SpawnTaskSessionOutput, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		title = "Task"
	}

	session, err := a.manager.Storage().CreateSession(data_models.Session{
		Title:  title,
		Status: "idle",
		Kind:   "task",
		Tags:   marshalSessionTags([]string{"task"}),
	})
	if err != nil {
		return nil, ierror.Error(ierror.ErrAgentCreateSession, err)
	}

	if app := a.manager.App(); app != nil {
		app.Event.Emit(event_id.AgentSessionSpawned, map[string]any{
			"sessionId":   session.ID,
			"title":       session.Title,
			"kind":        session.Kind,
			"tags":        []string{"task"},
			"userMessage": input.UserMessage,
		})
	}

	go func() {
		handler := a.manager.NewChatHandler()
		_ = handler.SendMessage(context.Background(), pkgAgent.SendMessageParams{
			SessionID:    session.ID,
			Content:      input.UserMessage,
			SystemPrompt: input.SystemPrompt,
			PrimedSkill:  input.SkillName,
			BaseURL:      input.BaseURL,
			ApiKey:       input.ApiKey,
			ModelName:    input.ModelName,
			ProviderType: input.ProviderType,
		})
	}()

	return &agent_dto.SpawnTaskSessionOutput{
		SessionID: session.ID,
		Title:     session.Title,
	}, nil
}

func (a *Agent) LoadSessionMessages(ctx context.Context, input agent_dto.LoadSessionMessagesInput) (*agent_dto.LoadSessionMessagesOutput, error) {
	limit := input.Limit
	if limit <= 0 {
		limit = 50
	}

	msgs, err := a.manager.Storage().ListMessagesForSession(input.SessionID, input.Offset, limit)
	if err != nil {
		return nil, ierror.Error(ierror.ErrAgentListMessages, err)
	}

	total, err := a.manager.Storage().CountMessagesForSession(input.SessionID)
	if err != nil {
		return nil, ierror.Error(ierror.ErrAgentCountMessages, err)
	}

	items := make([]agent_dto.MessageItem, 0, len(msgs))
	for _, msg := range msgs {
		items = append(items, toMessageItem(msg))
	}

	return &agent_dto.LoadSessionMessagesOutput{
		Messages: items,
		Total:    total,
		HasMore:  int64(input.Offset+len(msgs)) < total,
	}, nil
}

func (a *Agent) MarkSessionRead(ctx context.Context, input agent_dto.MarkSessionReadInput) (*agent_dto.MarkSessionReadOutput, error) {
	if err := a.manager.Storage().UpdateSessionStatus(input.SessionID, "idle"); err != nil {
		return nil, ierror.Error(ierror.ErrAgentMarkRead, err)
	}
	return &agent_dto.MarkSessionReadOutput{}, nil
}

func (a *Agent) RenameSession(ctx context.Context, input agent_dto.RenameSessionInput) (*agent_dto.RenameSessionOutput, error) {
	if err := a.manager.Storage().UpdateSessionTitle(input.SessionID, strings.TrimSpace(input.Title)); err != nil {
		return nil, ierror.Error(ierror.ErrAgentRename, err)
	}
	return &agent_dto.RenameSessionOutput{}, nil
}

func (a *Agent) DeleteSession(ctx context.Context, input agent_dto.DeleteSessionInput) (*agent_dto.DeleteSessionOutput, error) {
	a.manager.Streams().Stop(input.SessionID)
	if err := a.manager.Storage().DeleteSession(input.SessionID); err != nil {
		return nil, ierror.Error(ierror.ErrAgentDelete, err)
	}
	return &agent_dto.DeleteSessionOutput{}, nil
}

func (a *Agent) ToggleStarSession(ctx context.Context, input agent_dto.ToggleStarSessionInput) (*agent_dto.ToggleStarSessionOutput, error) {
	if err := a.manager.Storage().UpdateSessionStarred(input.SessionID, input.Starred); err != nil {
		return nil, ierror.Error(ierror.ErrAgentToggleStar, err)
	}
	return &agent_dto.ToggleStarSessionOutput{}, nil
}

func (a *Agent) GenerateTitle(ctx context.Context, input agent_dto.GenerateTitleInput) (*agent_dto.GenerateTitleOutput, error) {
	fmt.Println("Generating title for session:", input.SessionID)

	msgs, err := a.manager.Storage().ListMessagesForSession(input.SessionID, 0, 4)
	if err != nil || len(msgs) == 0 {
		return &agent_dto.GenerateTitleOutput{Title: "New Chat"}, nil
	}

	var preview strings.Builder
	for _, msg := range msgs {
		if msg.ContentType != "text" || msg.Content == "" {
			continue
		}
		content := msg.Content
		if len(content) > 200 {
			content = content[:200]
		}
		preview.WriteString(fmt.Sprintf("%s: %s\n", msg.Role, content))
	}

	if preview.Len() == 0 {
		return &agent_dto.GenerateTitleOutput{Title: "New Chat"}, nil
	}

	r, err := a.manager.GetOrCreateRunner(input.BaseURL, input.ApiKey, input.ModelName, input.ProviderType, nil, "user", input.SessionID)
	if err != nil {
		return &agent_dto.GenerateTitleOutput{Title: "New Chat"}, nil
	}

	titlePrompt := buildTitlePrompt(preview.String())
	events, err := r.Run(ctx, "local", fmt.Sprintf("title-%d", input.SessionID), model.NewUserMessage(titlePrompt))
	if err != nil {
		return &agent_dto.GenerateTitleOutput{Title: "New Chat"}, nil
	}

	result := collectTitleFromEvents(events)
	if result == "" {
		result = "New Chat"
	}

	_ = a.manager.Storage().UpdateSessionTitle(input.SessionID, result)
	return &agent_dto.GenerateTitleOutput{Title: result}, nil
}

func collectTitleFromEvents(events <-chan *event.Event) string {
	var streamedTitle strings.Builder
	var finalTitle string

	for evt := range events {
		if evt == nil || evt.Error != nil {
			continue
		}
		for _, choice := range evt.Choices {
			if choice.Delta.Content != "" {
				streamedTitle.WriteString(choice.Delta.Content)
			}
			if finalTitle == "" && choice.Message.Content != "" {
				finalTitle = choice.Message.Content
			}
		}
	}

	result := strings.TrimSpace(streamedTitle.String())
	if result != "" {
		return result
	}

	return strings.TrimSpace(finalTitle)
}

func buildTitlePrompt(preview string) string {
	titleInstruction, err := prompt.Load(prompt_id.GenChatTitle)
	if err != nil || strings.TrimSpace(titleInstruction) == "" {
		titleInstruction = "Generate a concise title under 20 characters for this conversation. Return only the title."
	}

	return fmt.Sprintf("%s\n\n%s", titleInstruction, preview)
}

const (
	maxAttachmentCount = 10
	maxAttachmentBytes = 20 * 1024 * 1024
)

// convertAttachments validates attachment count, file size, and existence,
// then converts DTOs into normalized pkgAgent.Attachment values.
func convertAttachments(in []agent_dto.AttachmentInput) ([]pkgAgent.Attachment, error) {
	if len(in) == 0 {
		return nil, nil
	}
	if len(in) > maxAttachmentCount {
		return nil, ierror.Error(ierror.ErrAgentTooManyAttachments, fmt.Errorf("too many attachments: %d (max %d)", len(in), maxAttachmentCount))
	}
	out := make([]pkgAgent.Attachment, 0, len(in))
	for _, item := range in {
		path := strings.TrimSpace(item.Path)
		if path == "" {
			return nil, ierror.Error(ierror.ErrAgentAttachmentPath, fmt.Errorf("attachment path is empty"))
		}
		info, err := os.Stat(path)
		if err != nil {
			return nil, ierror.Error(ierror.ErrAgentAttachmentNotFound, fmt.Errorf("attachment not found: %s", path))
		}
		if info.Size() > maxAttachmentBytes {
			return nil, ierror.Error(ierror.ErrAgentAttachmentSize, fmt.Errorf("attachment %q exceeds %d bytes", info.Name(), maxAttachmentBytes))
		}
		out = append(out, pkgAgent.NormalizeAttachment(pkgAgent.Attachment{
			Path: path,
			Name: item.Name,
			Mime: item.Mime,
		}))
	}
	return out, nil
}
