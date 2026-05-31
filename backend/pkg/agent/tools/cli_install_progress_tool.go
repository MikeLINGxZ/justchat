package tools

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

// CliInstallProgressReporter reports install/init progress for smart CLI setup tasks.
type CliInstallProgressReporter interface {
	ReportCliInstallProgress(ctx context.Context, sessionID uint, item CliInstallProgressItem) error
}

// CliInstallProgressItem is the generic progress payload surfaced to the UI.
type CliInstallProgressItem struct {
	NpmPackage  string `json:"npm_package,omitempty"`
	Name        string `json:"name,omitempty"`
	ExtensionID string `json:"extension_id,omitempty"`
	Phase       string `json:"phase"`
	Detail      string `json:"detail,omitempty"`
	ActionURL   string `json:"action_url,omitempty"`
	ActionLabel string `json:"action_label,omitempty"`
	ExpiresAt   string `json:"expires_at,omitempty"`
}

const ReportCliInstallProgressToolName = "ReportCliInstallProgress"

type reportCliInstallProgressInput struct {
	NpmPackage  string `json:"npm_package" jsonschema:"description=installed npm package name when known"`
	Name        string `json:"name" jsonschema:"description=CLI display name when known"`
	ExtensionID string `json:"extension_id" jsonschema:"description=installed extension id when known"`
	Phase       string `json:"phase" jsonschema:"description=current phase such as downloading|installed|generating|initializing|waiting_auth|verifying|done|failed,required"`
	Detail      string `json:"detail" jsonschema:"description=short user-facing detail for the current phase"`
	ActionURL   string `json:"action_url" jsonschema:"description=optional URL the user should open or click"`
	ActionLabel string `json:"action_label" jsonschema:"description=optional short label for the URL action"`
	ExpiresAt   string `json:"expires_at" jsonschema:"description=optional expiration timestamp or user-facing expiry note"`
}

func BuildReportCliInstallProgressTool() ToolMeta {
	return ToolMeta{
		Name:        ReportCliInstallProgressToolName,
		Description: "Report smart CLI install or initialization progress to the UI. Use this at key checkpoints such as install started, initializing, waiting for auth, verifying, done, or failed.",
		Category:    CategoryBuiltin,
		FormatPurpose: func(args json.RawMessage) string {
			var input reportCliInstallProgressInput
			_ = json.Unmarshal(args, &input)
			phase := strings.TrimSpace(input.Phase)
			if phase == "" {
				phase = "updating progress"
			}
			if input.Name != "" {
				return "Report CLI progress: " + input.Name + " (" + phase + ")"
			}
			if input.NpmPackage != "" {
				return "Report CLI progress: " + input.NpmPackage + " (" + phase + ")"
			}
			return "Report CLI progress: " + phase
		},
	}
}

func InvokeReportCliInstallProgress(ctx context.Context, reporter CliInstallProgressReporter, sessionID uint, args json.RawMessage) (string, error) {
	var input reportCliInstallProgressInput
	if err := json.Unmarshal(args, &input); err != nil {
		return "", err
	}
	if strings.TrimSpace(input.Phase) == "" {
		return "", errors.New("phase is required")
	}
	item := CliInstallProgressItem{
		NpmPackage:  strings.TrimSpace(input.NpmPackage),
		Name:        strings.TrimSpace(input.Name),
		ExtensionID: strings.TrimSpace(input.ExtensionID),
		Phase:       strings.TrimSpace(input.Phase),
		Detail:      strings.TrimSpace(input.Detail),
		ActionURL:   strings.TrimSpace(input.ActionURL),
		ActionLabel: strings.TrimSpace(input.ActionLabel),
		ExpiresAt:   strings.TrimSpace(input.ExpiresAt),
	}
	if err := reporter.ReportCliInstallProgress(ctx, sessionID, item); err != nil {
		return "", err
	}
	return "cli install progress updated", nil
}

func NewReportCliInstallProgressTool(reporter CliInstallProgressReporter, sessionID uint) *function.FunctionTool[reportCliInstallProgressInput, string] {
	meta := BuildReportCliInstallProgressTool()
	return function.NewFunctionTool(
		func(ctx context.Context, input reportCliInstallProgressInput) (string, error) {
			payload, err := json.Marshal(input)
			if err != nil {
				return "", err
			}
			return InvokeReportCliInstallProgress(ctx, reporter, sessionID, payload)
		},
		function.WithName(ReportCliInstallProgressToolName),
		function.WithDescription(meta.Description),
	)
}
