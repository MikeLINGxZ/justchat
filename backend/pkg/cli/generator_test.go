package cli

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

// TestGenerateParsesManifestAndOverridesExecutable verifies generated JSON is parsed and the executable is caller-controlled.
func TestGenerateParsesManifestAndOverridesExecutable(t *testing.T) {
	t.Helper()

	manifest, err := Generate(context.Background(), GenerateParams{
		HelpText:    "demo --help",
		PackageName: "demo-cli",
		PackageMeta: PackageMeta{Version: "1.0.0", Description: "demo"},
		Executable:  "/tmp/demo-cli",
		Caller: func(ctx context.Context, system, user string) (string, error) {
			return `{
				"name":"demo-cli",
				"version":"1.0.0",
				"description":"demo",
				"executable":"hallucinated",
				"isolation":"isolated",
				"tools":[
					{
						"name":"list_items",
						"description":"list items",
						"input_schema":{"type":"object","properties":{"query":{"type":"string"}}},
						"argv_template":["list","--query","{query}"],
						"output_mode":"json",
						"timeout_seconds":60,
						"requires_confirm":false,
						"enabled":true
					}
				]
			}`, nil
		},
	})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if manifest.Executable != "/tmp/demo-cli" {
		t.Fatalf("expected executable override, got %q", manifest.Executable)
	}
	if len(manifest.Tools) != 1 || manifest.Tools[0].Name != "list_items" {
		t.Fatalf("unexpected tools: %+v", manifest.Tools)
	}
}

// TestGenerateAcceptsFencedJSON verifies markdown code fences are tolerated.
func TestGenerateAcceptsFencedJSON(t *testing.T) {
	t.Helper()

	manifest, err := Generate(context.Background(), GenerateParams{
		HelpText:    "demo --help",
		PackageName: "demo-cli",
		PackageMeta: PackageMeta{Version: "1.0.0"},
		Executable:  "/tmp/demo-cli",
		Caller: func(ctx context.Context, system, user string) (string, error) {
			return "```json\n{\n  \"name\":\"demo-cli\",\n  \"version\":\"1.0.0\",\n  \"description\":\"demo\",\n  \"executable\":\"ignored\",\n  \"isolation\":\"isolated\",\n  \"tools\":[]\n}\n```", nil
		},
	})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if manifest.Name != "demo-cli" || manifest.Executable != "/tmp/demo-cli" {
		t.Fatalf("unexpected manifest: %+v", manifest)
	}
}

// TestGenerateFallsBackToStubManifest verifies malformed model output degrades to an editable empty manifest.
func TestGenerateFallsBackToStubManifest(t *testing.T) {
	t.Helper()

	manifest, err := Generate(context.Background(), GenerateParams{
		HelpText:    "demo --help",
		PackageName: "demo-cli",
		PackageMeta: PackageMeta{Version: "1.0.0", Description: "demo"},
		Executable:  "/tmp/demo-cli",
		Caller: func(ctx context.Context, system, user string) (string, error) {
			return "this is not json", nil
		},
	})
	if err != nil {
		t.Fatalf("expected fallback manifest without error, got %v", err)
	}
	if manifest.Name != "demo-cli" || manifest.Executable != "/tmp/demo-cli" {
		t.Fatalf("unexpected fallback manifest: %+v", manifest)
	}
	if len(manifest.Tools) != 0 {
		t.Fatalf("expected empty tools on fallback, got %+v", manifest.Tools)
	}
}

// TestGenerateBuildsValidStubSchema verifies the fallback manifest remains serializable for later editing.
func TestGenerateBuildsValidStubSchema(t *testing.T) {
	t.Helper()

	manifest, err := Generate(context.Background(), GenerateParams{
		HelpText:    "demo --help",
		PackageName: "demo-cli",
		PackageMeta: PackageMeta{Version: "1.0.0"},
		Executable:  "/tmp/demo-cli",
		Caller: func(ctx context.Context, system, user string) (string, error) {
			return `{"name":"demo-cli","version":"1.0.0","description":"demo","executable":"","isolation":"isolated","tools":[]}`, nil
		},
	})
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if _, marshalErr := json.Marshal(manifest); marshalErr != nil {
		t.Fatalf("expected fallback manifest to remain serializable: %v", marshalErr)
	}
}

// TestGenerateFallsBackWhenCallerFails verifies transport/model failures still return an editable stub manifest.
func TestGenerateFallsBackWhenCallerFails(t *testing.T) {
	t.Helper()

	manifest, err := Generate(context.Background(), GenerateParams{
		HelpText:    "demo --help",
		PackageName: "demo-cli",
		PackageMeta: PackageMeta{Version: "1.0.0", Description: "demo"},
		Executable:  "/tmp/demo-cli",
		Caller: func(ctx context.Context, system, user string) (string, error) {
			return "", context.DeadlineExceeded
		},
	})
	if err != nil {
		t.Fatalf("expected fallback manifest without error, got %v", err)
	}
	if manifest.Name != "demo-cli" || manifest.Executable != "/tmp/demo-cli" {
		t.Fatalf("unexpected fallback manifest: %+v", manifest)
	}
	if len(manifest.Tools) != 0 {
		t.Fatalf("expected empty tools on caller failure, got %+v", manifest.Tools)
	}
}

// TestBuildGeneratorUserPromptTruncatesVerboseHelp verifies extremely large help text is clamped before prompting the model.
func TestBuildGeneratorUserPromptTruncatesVerboseHelp(t *testing.T) {
	t.Helper()

	helpText := strings.Repeat("subcommand --flag description\n", 800)
	prompt := buildGeneratorUserPrompt(GenerateParams{
		HelpText:    helpText,
		PackageName: "demo-cli",
		PackageMeta: PackageMeta{Version: "1.0.0", Description: "demo"},
		Executable:  "/tmp/demo-cli",
	})
	if len(prompt) >= len(helpText)+200 {
		t.Fatalf("expected prompt to truncate help text, prompt length=%d raw length=%d", len(prompt), len(helpText))
	}
	if !strings.Contains(prompt, "[help output truncated]") {
		t.Fatalf("expected prompt to contain truncation marker: %q", prompt)
	}
}

// TestGuessLoginCommand verifies the heuristic expands auth/oauth into the canonical login pair
// and leaves login/signin one-element. Unrelated CLIs return nil.
func TestGuessLoginCommand(t *testing.T) {
	cases := []struct {
		name string
		help string
		want []string
	}{
		{"lark auth", "  auth        OAuth credentials and authorization management\n", []string{"auth", "login"}},
		{"gh oauth", "  oauth       OAuth helpers\n", []string{"oauth", "login"}},
		{"plain login", "  login   Sign in to the service\n", []string{"login"}},
		{"signin", "  signin   Authenticate\n", []string{"signin"}},
		{"none", "  list   List items\n  delete   Delete item\n", nil},
	}
	for _, c := range cases {
		got := guessLoginCommand(c.help)
		if len(got) != len(c.want) {
			t.Fatalf("%s: got %v, want %v", c.name, got, c.want)
		}
		for i := range got {
			if got[i] != c.want[i] {
				t.Fatalf("%s: got %v, want %v", c.name, got, c.want)
			}
		}
	}
}

// TestGenerateAppliesLoginCommandHeuristicWhenAIMisses verifies that when the LLM omits
// login_command, Generate fills it from the --help heuristic.
func TestGenerateAppliesLoginCommandHeuristicWhenAIMisses(t *testing.T) {
	manifest, err := Generate(context.Background(), GenerateParams{
		HelpText:    "USAGE: lark-cli <cmd>\n\nCOMMANDS:\n  auth      OAuth credentials\n  list      list things\n",
		PackageName: "lark-cli",
		PackageMeta: PackageMeta{Version: "1.0.0", Description: "feishu"},
		Executable:  "/tmp/lark-cli",
		Caller: func(ctx context.Context, system, user string) (string, error) {
			return `{"name":"lark-cli","version":"1.0.0","description":"feishu","executable":"/tmp/lark-cli","tools":[]}`, nil
		},
	})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if len(manifest.LoginCommand) != 2 || manifest.LoginCommand[0] != "auth" || manifest.LoginCommand[1] != "login" {
		t.Fatalf("expected [auth login] from heuristic, got %v", manifest.LoginCommand)
	}
}
