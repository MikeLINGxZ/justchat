package plugin

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	pkgcli "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/cli"
	pkgProvider "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/provider"
	pkgRuntimeState "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/runtime_state"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/plugin/plugin_dto"
)

// TestGenerateCliManifestPersistsManifestAndProjectsTools verifies the generated manifest is saved and mirrored back to config/runtime metadata.
func TestGenerateCliManifestPersistsManifestAndProjectsTools(t *testing.T) {
	p := NewPlugin()
	dataDir := withFakeCliManager(t, p)
	if _, err := p.InstallCliFromNpm(context.Background(), plugin_dto.InstallCliFromNpmInput{NpmPackage: "lark-cli", Name: "lark-cli"}); err != nil {
		t.Fatal(err)
	}

	p.probeCliHelp = func(ctx context.Context, executable string, env []string) (string, error) {
		return "lark --help", nil
	}
	p.loadPersistedRuntime = func() (pkgRuntimeState.StateSnapshot, error) {
		return pkgRuntimeState.StateSnapshot{State: "ready", NodePath: "/tmp/node", NpmPath: "/tmp/npm"}, nil
	}
	p.resolveDefaultChatModel = func() (defaultChatModel, error) {
		return defaultChatModel{
			BaseURL:      "http://example.invalid",
			APIKey:       "test-key",
			ModelName:    "qwen-plus",
			ProviderType: pkgProvider.OpenAiCompatibility,
		}, nil
	}
	p.generateManifest = func(ctx context.Context, params pkgcli.GenerateParams) (pkgcli.Manifest, error) {
		return pkgcli.Manifest{
			Name:        params.PackageName,
			Version:     params.PackageMeta.Version,
			Description: params.PackageMeta.Description,
			Executable:  params.Executable,
			Isolation:   pkgcli.IsolationIsolated,
			Tools: []pkgcli.Tool{
				{
					Name:            "send_message",
					Description:     "send a message",
					InputSchema:     []byte(`{"type":"object","properties":{"text":{"type":"string"}}}`),
					ArgvTemplate:    []string{"message", "send", "{text}"},
					OutputMode:      pkgcli.OutputText,
					RequiresConfirm: true,
					Enabled:         true,
				},
			},
		}, nil
	}

	out, err := p.GenerateCliManifest(context.Background(), plugin_dto.GenerateCliManifestInput{ID: "cli:lark-cli"})
	if err != nil {
		t.Fatalf("GenerateCliManifest returned error: %v", err)
	}
	if len(out.Extension.Tools) != 1 || out.Extension.Tools[0].Name != "send_message" {
		t.Fatalf("expected generated tools to be projected, got %+v", out.Extension.Tools)
	}

	manifestPath := filepath.Join(dataDir, "plugins", "cli_data", "lark-cli", "manifest.json")
	manifest, err := pkgcli.LoadManifest(manifestPath)
	if err != nil {
		t.Fatalf("LoadManifest returned error: %v", err)
	}
	if len(manifest.Tools) != 1 || manifest.Tools[0].Name != "send_message" {
		t.Fatalf("expected manifest file to be updated, got %+v", manifest.Tools)
	}
}

// TestGenerateCliManifestRequiresDefaultProvider verifies the service fails cleanly when no default chat model can be resolved.
func TestGenerateCliManifestRequiresDefaultProvider(t *testing.T) {
	p := NewPlugin()
	_ = withFakeCliManager(t, p)
	if _, err := p.InstallCliFromNpm(context.Background(), plugin_dto.InstallCliFromNpmInput{NpmPackage: "lark-cli", Name: "lark-cli"}); err != nil {
		t.Fatal(err)
	}

	p.resolveDefaultChatModel = func() (defaultChatModel, error) {
		return defaultChatModel{}, os.ErrNotExist
	}

	_, err := p.GenerateCliManifest(context.Background(), plugin_dto.GenerateCliManifestInput{ID: "cli:lark-cli"})
	if err == nil {
		t.Fatal("expected GenerateCliManifest to fail without a default provider")
	}
}

// TestGenerateCliManifestRepairsMissingExecutable verifies regeneration can recover when an older manifest is missing its executable field.
func TestGenerateCliManifestRepairsMissingExecutable(t *testing.T) {
	p := NewPlugin()
	dataDir := withFakeCliManager(t, p)
	if _, err := p.InstallCliFromNpm(context.Background(), plugin_dto.InstallCliFromNpmInput{NpmPackage: "lark-cli", Name: "lark-cli"}); err != nil {
		t.Fatal(err)
	}

	manifestPath := filepath.Join(dataDir, "plugins", "cli_data", "lark-cli", "manifest.json")
	manifest, err := pkgcli.LoadManifest(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	manifest.Executable = ""
	if err := pkgcli.SaveManifest(manifestPath, manifest); err != nil {
		t.Fatal(err)
	}

	p.probeCliHelp = func(ctx context.Context, executable string, env []string) (string, error) {
		if executable == "" {
			t.Fatal("expected executable to be repaired before probing help")
		}
		return "lark --help", nil
	}
	p.loadPersistedRuntime = func() (pkgRuntimeState.StateSnapshot, error) {
		return pkgRuntimeState.StateSnapshot{State: "ready", NodePath: "/tmp/node", NpmPath: "/tmp/npm"}, nil
	}
	p.resolveDefaultChatModel = func() (defaultChatModel, error) {
		return defaultChatModel{
			BaseURL:      "http://example.invalid",
			APIKey:       "test-key",
			ModelName:    "qwen-plus",
			ProviderType: pkgProvider.OpenAiCompatibility,
		}, nil
	}
	p.generateManifest = func(ctx context.Context, params pkgcli.GenerateParams) (pkgcli.Manifest, error) {
		return pkgcli.Manifest{
			Name:        params.PackageName,
			Version:     params.PackageMeta.Version,
			Description: params.PackageMeta.Description,
			Executable:  params.Executable,
			Isolation:   pkgcli.IsolationIsolated,
			Tools:       []pkgcli.Tool{},
		}, nil
	}

	out, err := p.GenerateCliManifest(context.Background(), plugin_dto.GenerateCliManifestInput{ID: "cli:lark-cli"})
	if err != nil {
		t.Fatalf("GenerateCliManifest returned error: %v", err)
	}
	if out.Extension.ConfigFilePath == "" {
		t.Fatalf("expected extension to be returned, got %+v", out.Extension)
	}

	repaired, err := pkgcli.LoadManifest(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if repaired.Executable == "" {
		t.Fatalf("expected executable to be repaired in manifest: %+v", repaired)
	}
}

// TestGenerateCliManifestPreservesExistingLoginSteps verifies regenerate does not clobber
// a multi-step login flow the user (or a prior regen) had carefully assembled. The AI's
// proposed login_command is ignored when the existing manifest already declared login_steps
// or login_command, because the generator can't infer multi-stage auth from --help alone.
func TestGenerateCliManifestPreservesExistingLoginSteps(t *testing.T) {
	p := NewPlugin()
	dataDir := withFakeCliManager(t, p)
	if _, err := p.InstallCliFromNpm(context.Background(), plugin_dto.InstallCliFromNpmInput{NpmPackage: "lark-cli", Name: "lark-cli"}); err != nil {
		t.Fatal(err)
	}

	// Seed the existing manifest with a multi-step login configuration.
	manifestPath := filepath.Join(dataDir, "plugins", "cli_data", "lark-cli", "manifest.json")
	existing, err := pkgcli.LoadManifest(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	existing.LoginSteps = [][]string{
		{"config", "init", "--new"},
		{"auth", "login"},
	}
	if err := pkgcli.SaveManifest(manifestPath, existing); err != nil {
		t.Fatal(err)
	}

	p.probeCliHelp = func(_ context.Context, _ string, _ []string) (string, error) {
		return "lark --help", nil
	}
	p.loadPersistedRuntime = func() (pkgRuntimeState.StateSnapshot, error) {
		return pkgRuntimeState.StateSnapshot{State: "ready", NodePath: "/tmp/node", NpmPath: "/tmp/npm"}, nil
	}
	p.resolveDefaultChatModel = func() (defaultChatModel, error) {
		return defaultChatModel{
			BaseURL: "http://example.invalid", APIKey: "test-key",
			ModelName: "qwen-plus", ProviderType: pkgProvider.OpenAiCompatibility,
		}, nil
	}
	// AI returns a single-step login_command (the common case), trying to overwrite the multi-step setup.
	p.generateManifest = func(_ context.Context, params pkgcli.GenerateParams) (pkgcli.Manifest, error) {
		return pkgcli.Manifest{
			Name:         params.PackageName,
			Version:      params.PackageMeta.Version,
			Description:  params.PackageMeta.Description,
			Executable:   params.Executable,
			Isolation:    pkgcli.IsolationIsolated,
			LoginCommand: []string{"auth", "login"},
			Tools:        []pkgcli.Tool{},
		}, nil
	}

	if _, err := p.GenerateCliManifest(context.Background(), plugin_dto.GenerateCliManifestInput{ID: "cli:lark-cli"}); err != nil {
		t.Fatalf("GenerateCliManifest: %v", err)
	}

	after, err := pkgcli.LoadManifest(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(after.LoginSteps) != 2 {
		t.Fatalf("expected login_steps to be preserved (len 2), got %v", after.LoginSteps)
	}
	if after.LoginSteps[0][0] != "config" || after.LoginSteps[1][0] != "auth" {
		t.Fatalf("expected preserved step order, got %v", after.LoginSteps)
	}
	if len(after.LoginCommand) != 0 {
		t.Fatalf("expected AI-proposed login_command to be discarded when login_steps present, got %v", after.LoginCommand)
	}
}
