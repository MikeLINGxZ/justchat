package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/agent/tools"
	pkgProvider "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/provider"
)

// ToolConfirmDecisionInput describes the tool call under review for runtime confirmation.
type ToolConfirmDecisionInput struct {
	ToolName     string
	Arguments    string
	Purpose      string
	Meta         tools.ToolMeta
	ModelName    string
	BaseURL      string
	APIKey       string
	ProviderType pkgProvider.Type
	SessionID    uint
}

// ToolConfirmDecision is the structured result returned by a runtime confirmation decider.
type ToolConfirmDecision struct {
	RequireConfirm bool
	Deny           bool
	Reason         string
	Source         string
}

// ToolConfirmDecider decides whether one tool call should require user confirmation.
type ToolConfirmDecider func(ctx context.Context, input ToolConfirmDecisionInput) (ToolConfirmDecision, error)

const toolConfirmDecisionSystemPrompt = `You decide whether a tool call should require user confirmation before execution.

Return strict JSON:
{
  "require_confirm": true|false,
  "reason": "<short reason under 12 words>"
}

Guidelines:
- require_confirm=true for destructive, externally visible, irreversible, costly, privileged, or environment-modifying actions.
- require_confirm=false for read-only inspection, low-risk analysis, and harmless metadata queries.
- If the risk is uncertain, choose true.
- Judge this specific invocation, not just the tool name.`

func defaultToolConfirmDecider(ctx context.Context, input ToolConfirmDecisionInput) (ToolConfirmDecision, error) {
	if decision, ok := deterministicToolConfirmDecision(input); ok {
		return decision, nil
	}
	if strings.TrimSpace(input.ModelName) == "" || strings.TrimSpace(input.BaseURL) == "" {
		return heuristicToolConfirmDecision(input), nil
	}

	userPrompt := buildToolConfirmDecisionPrompt(input)
	resp, err := OneshotComplete(ctx, OneshotRequest{
		BaseURL:      input.BaseURL,
		APIKey:       input.APIKey,
		ModelName:    input.ModelName,
		ProviderType: input.ProviderType,
		System:       toolConfirmDecisionSystemPrompt,
		User:         userPrompt,
		MaxTokens:    120,
		Timeout:      20 * time.Second,
	})
	if err != nil {
		fallback := heuristicToolConfirmDecision(input)
		fallback.Reason = "fallback after decider error"
		return fallback, nil
	}

	var parsed struct {
		RequireConfirm bool   `json:"require_confirm"`
		Reason         string `json:"reason"`
	}
	if err := json.Unmarshal([]byte(extractJSONObject(resp.Text)), &parsed); err != nil {
		fallback := heuristicToolConfirmDecision(input)
		fallback.Reason = "fallback after parse error"
		return fallback, nil
	}

	return ToolConfirmDecision{
		RequireConfirm: parsed.RequireConfirm,
		Reason:         strings.TrimSpace(parsed.Reason),
		Source:         "model",
	}, nil
}

func buildToolConfirmDecisionPrompt(input ToolConfirmDecisionInput) string {
	var builder strings.Builder
	builder.WriteString("Tool name: ")
	builder.WriteString(input.ToolName)
	builder.WriteString("\n")
	if strings.TrimSpace(input.Purpose) != "" {
		builder.WriteString("Purpose: ")
		builder.WriteString(strings.TrimSpace(input.Purpose))
		builder.WriteString("\n")
	}
	builder.WriteString("Arguments JSON: ")
	builder.WriteString(input.Arguments)
	builder.WriteString("\n")
	builder.WriteString("Static requires_confirm fallback: ")
	builder.WriteString(fmt.Sprintf("%t", input.Meta.RequiresConfirm))
	return builder.String()
}

func heuristicToolConfirmDecision(input ToolConfirmDecisionInput) ToolConfirmDecision {
	if input.Meta.RequiresConfirm {
		return ToolConfirmDecision{
			RequireConfirm: true,
			Reason:         "static risky action",
			Source:         "static",
		}
	}

	lower := strings.ToLower(strings.Join([]string{input.ToolName, input.Purpose, input.Arguments}, " "))
	riskyTerms := []string{
		"delete", "remove", "destroy", "drop", "send", "post ", "publish", "create",
		"update", "write", "upload", "execute", "shell", "command", "prod", "production",
		"deploy", "restart", "kill", "token", "secret", "config", "permission",
	}
	safeTerms := []string{
		"list", "get", "show", "read", "fetch", "query", "status", "inspect", "describe",
	}

	for _, term := range riskyTerms {
		if strings.Contains(lower, term) {
			return ToolConfirmDecision{
				RequireConfirm: true,
				Reason:         "heuristic risky action",
				Source:         "heuristic",
			}
		}
	}
	for _, term := range safeTerms {
		if strings.Contains(lower, term) {
			return ToolConfirmDecision{
				RequireConfirm: false,
				Reason:         "heuristic read only",
				Source:         "heuristic",
			}
		}
	}
	return ToolConfirmDecision{
		RequireConfirm: true,
		Reason:         "heuristic uncertain risk",
		Source:         "heuristic",
	}
}

// deterministicToolConfirmDecision applies product safety rules that do not need model judgment.
func deterministicToolConfirmDecision(input ToolConfirmDecisionInput) (ToolConfirmDecision, bool) {
	if input.ToolName != "shell" {
		return ToolConfirmDecision{}, false
	}
	command := shellCommandFromArgs(input.Arguments)
	if command == "" {
		return ToolConfirmDecision{RequireConfirm: true, Reason: "missing shell command", Source: "policy"}, true
	}
	normalized := normalizeShellCommand(command)
	if isDeniedShellCommand(normalized) {
		return ToolConfirmDecision{
			Deny:           true,
			RequireConfirm: false,
			Reason:         "destructive shell command",
			Source:         "policy",
		}, true
	}
	if isSafeShellCommand(normalized) {
		return ToolConfirmDecision{
			RequireConfirm: false,
			Reason:         "safe shell inspection",
			Source:         "policy",
		}, true
	}
	if isRiskyShellCommand(normalized) {
		return ToolConfirmDecision{
			RequireConfirm: true,
			Reason:         "file system change",
			Source:         "policy",
		}, true
	}
	return ToolConfirmDecision{}, false
}

// shellCommandFromArgs extracts the command field from a shell tool argument payload.
func shellCommandFromArgs(args string) string {
	var parsed struct {
		Command string `json:"command"`
	}
	if err := json.Unmarshal([]byte(args), &parsed); err != nil {
		return ""
	}
	return strings.TrimSpace(parsed.Command)
}

// normalizeShellCommand trims shell wrappers so simple policy matching remains stable.
func normalizeShellCommand(command string) string {
	lower := strings.ToLower(strings.TrimSpace(command))
	lower = strings.TrimPrefix(lower, "command ")
	return strings.Join(strings.Fields(lower), " ")
}

// isSafeShellCommand recognizes read-only commands useful for ordinary user tasks.
func isSafeShellCommand(command string) bool {
	safeExact := []string{
		"pwd",
	}
	for _, exact := range safeExact {
		if command == exact {
			return true
		}
	}
	safePrefixes := []string{
		"ls ",
		"find ",
		"du ",
		"wc ",
		"stat ",
		"file ",
		"mdls ",
	}
	for _, prefix := range safePrefixes {
		if command == strings.TrimSpace(prefix) || strings.HasPrefix(command, prefix) {
			return true
		}
	}
	return false
}

// isRiskyShellCommand recognizes commands that can change files and should ask first.
func isRiskyShellCommand(command string) bool {
	riskyPrefixes := []string{
		"mv ",
		"cp ",
		"mkdir ",
		"touch ",
		"chmod ",
		"chown ",
		"open ",
		"osascript ",
		"zip ",
		"unzip ",
		"tar ",
	}
	for _, prefix := range riskyPrefixes {
		if strings.HasPrefix(command, prefix) {
			return true
		}
	}
	return false
}

// isDeniedShellCommand blocks clearly destructive commands by default.
func isDeniedShellCommand(command string) bool {
	deniedPrefixes := []string{
		"rm ",
		"rm -",
		"rmdir ",
		"shred ",
		"mkfs",
		"diskutil erase",
	}
	for _, prefix := range deniedPrefixes {
		if strings.HasPrefix(command, prefix) {
			return true
		}
	}
	deniedFragments := []string{
		"; rm ",
		"&& rm ",
		"|| rm ",
		"| rm ",
		"; rmdir ",
		"&& rmdir ",
		"|| rmdir ",
		"| rmdir ",
	}
	for _, fragment := range deniedFragments {
		if strings.Contains(command, fragment) {
			return true
		}
	}
	return false
}

func extractJSONObject(raw string) string {
	trimmed := strings.TrimSpace(raw)
	start := strings.Index(trimmed, "{")
	end := strings.LastIndex(trimmed, "}")
	if start >= 0 && end >= start {
		return trimmed[start : end+1]
	}
	return trimmed
}
