package cli

import (
	"context"
	"encoding/json"
	"log"
	"strings"
)

const maxGeneratorHelpChars = 12000

const manifestGeneratorSystemPrompt = `You convert a CLI tool's --help output into a JSON manifest for the Lemontea CLI plugin system.

The manifest schema is:
{
  "name": "<package name>",
  "version": "<package version>",
  "description": "<short>",
  "executable": "<provided>",
  "login_command": ["<subcommand path that starts an interactive auth flow; see rules below>"],
  "isolation": "isolated",
  "tools": [
    {
      "name": "<snake_case_tool_name>",
      "description": "<what it does>",
      "input_schema": { "type": "object", "properties": { "<key>": {"type":"string"} }, "required": [] },
      "argv_template": ["<subcmd>", "--flag", "{key}"],
      "output_mode": "json" | "text" | "lines",
      "timeout_seconds": 60,
      "requires_confirm": true|false,
      "enabled": true
    }
  ]
}

Rules:
- Map each useful subcommand to a tool. Skip "help" / "version" / "completion".
- login_command: scan top-level subcommands for anything related to authentication — names like "auth", "login", "signin", "oauth", "config" (when the description mentions credentials or app setup). If you find one with sub-subcommands, dig in: e.g. an "auth" group containing "login" should yield ["auth", "login"]. A standalone "login" command yields ["login"]. A config-style setup like "config init" yields ["config", "init"]. Only leave login_command as an empty array if the CLI clearly has no authentication step at all. Prefer the OAuth / device-flow command over credential-stdin setup when both exist.
- requires_confirm: true for verbs that mutate state (send/delete/create/update/upload). false for read verbs (list/get/show).
- argv_template must use {placeholder} tokens that match input_schema.properties keys.
- Any {placeholder} the CLI cannot accept as missing must be listed in input_schema.required. Placeholders left out of required are treated as optional: at runtime, if the caller omits the value, the engine drops that argv segment AND the immediately preceding "--flag" segment, so optional pairs vanish cleanly. Pair optional values with a "--flag" segment (e.g. ["--data", "{data}"]); never put an optional placeholder in a positional slot.
- output_mode: prefer "json" if the CLI has a --json flag for that subcommand, else "text".
- Return ONLY the JSON object, no prose.`

// OneshotCaller is the injected LLM dependency for manifest generation.
type OneshotCaller func(ctx context.Context, system, user string) (string, error)

// GenerateParams contains the inputs required to draft a CLI manifest.
type GenerateParams struct {
	HelpText    string
	PackageName string
	PackageMeta PackageMeta
	Executable  string
	Caller      OneshotCaller
}

// Generate asks an LLM to draft a manifest and falls back to an empty-tools stub when parsing/validation fails.
func Generate(ctx context.Context, p GenerateParams) (Manifest, error) {
	userPrompt := buildGeneratorUserPrompt(p)
	raw, err := p.Caller(ctx, manifestGeneratorSystemPrompt, userPrompt)
	if err != nil {
		log.Printf("cli manifest generator caller failed: %v", err)
		return stubManifest(p), nil
	}

	manifest, parseErr := parseGeneratedManifest(raw)
	if parseErr != nil {
		log.Printf("cli manifest generator parse failed: %v", parseErr)
		return stubManifest(p), nil
	}
	manifest.Executable = p.Executable
	if manifest.Isolation == "" {
		manifest.Isolation = IsolationIsolated
	}
	if len(manifest.LoginCommand) == 0 {
		manifest.LoginCommand = guessLoginCommand(p.HelpText)
	}
	if validateErr := Validate(manifest); validateErr != nil {
		log.Printf("cli manifest generator validation failed: %v", validateErr)
		return stubManifest(p), nil
	}
	return manifest, nil
}

// guessLoginCommand scans CLI --help output for a top-level auth-like subcommand
// and returns the canonical login invocation. auth/oauth are expanded to
// ["<name>", "login"] (gh, gcloud, lark-cli convention); login/signin stay as one element.
// Returns nil when nothing matches.
func guessLoginCommand(helpText string) []string {
	for _, line := range strings.Split(helpText, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 1 {
			continue
		}
		switch strings.ToLower(fields[0]) {
		case "auth", "oauth":
			return []string{fields[0], "login"}
		case "login", "signin":
			return []string{fields[0]}
		}
	}
	return nil
}

// buildGeneratorUserPrompt renders the task-specific prompt payload for one CLI package.
func buildGeneratorUserPrompt(p GenerateParams) string {
	return "Package: " + p.PackageName + " v" + p.PackageMeta.Version + "\n" +
		"Description: " + p.PackageMeta.Description + "\n" +
		"Executable: " + p.Executable + "\n\n" +
		"--help output:\n```\n" + clampGeneratorHelpText(p.HelpText) + "\n```"
}

// clampGeneratorHelpText trims oversized help output so the model prompt stays responsive on verbose CLIs.
func clampGeneratorHelpText(helpText string) string {
	trimmed := strings.TrimSpace(helpText)
	if len(trimmed) <= maxGeneratorHelpChars {
		return trimmed
	}

	const omittedMarker = "\n...\n[help output truncated]\n...\n"
	headLen := 8000
	tailLen := maxGeneratorHelpChars - headLen - len(omittedMarker)
	if tailLen < 0 {
		tailLen = 0
	}
	if headLen > len(trimmed) {
		headLen = len(trimmed)
	}
	start := trimmed[:headLen]
	end := ""
	if tailLen > 0 && tailLen < len(trimmed)-headLen {
		end = trimmed[len(trimmed)-tailLen:]
	}
	return start + omittedMarker + end
}

// parseGeneratedManifest decodes raw or fenced JSON from the model output.
func parseGeneratedManifest(raw string) (Manifest, error) {
	trimmed := strings.TrimSpace(raw)
	var manifest Manifest
	if err := json.Unmarshal([]byte(trimmed), &manifest); err == nil {
		return manifest, nil
	}

	start := strings.Index(trimmed, "{")
	end := strings.LastIndex(trimmed, "}")
	if start < 0 || end < 0 || end < start {
		return Manifest{}, json.Unmarshal([]byte(trimmed), &manifest)
	}
	if err := json.Unmarshal([]byte(trimmed[start:end+1]), &manifest); err != nil {
		return Manifest{}, err
	}
	return manifest, nil
}

// stubManifest creates the smallest editable manifest we can safely persist.
func stubManifest(p GenerateParams) Manifest {
	return Manifest{
		Name:        p.PackageName,
		Version:     p.PackageMeta.Version,
		Description: p.PackageMeta.Description,
		Executable:  p.Executable,
		Isolation:   IsolationIsolated,
		Tools:       []Tool{},
	}
}
