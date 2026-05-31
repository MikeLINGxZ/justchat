package tools

import (
	"context"
	"encoding/json"
	"testing"
)

type fakeCliCommandRunner struct {
	gotID         string
	gotArgv       []string
	gotOutputMode string
	gotTimeout    int
	progress      []string
}

func (f *fakeCliCommandRunner) RunCliCommandSync(_ context.Context, id string, argv []string, outputMode string, timeoutSeconds int, onProgress func(result string)) (string, error) {
	f.gotID = id
	f.gotArgv = append([]string(nil), argv...)
	f.gotOutputMode = outputMode
	f.gotTimeout = timeoutSeconds
	if onProgress != nil {
		onProgress("streaming output")
		f.progress = append(f.progress, "streaming output")
	}
	return `{"exit_code":0,"stdout":"ok"}`, nil
}

func TestInvokeRunCliCommand_PassesInputs(t *testing.T) {
	runner := &fakeCliCommandRunner{}
	args := json.RawMessage(`{"id":"cli:lark-cli","argv":["auth","login","--no-wait","--json"],"output_mode":"json","timeout_seconds":90}`)

	out, err := InvokeRunCliCommand(context.Background(), runner, args, nil)
	if err != nil {
		t.Fatal(err)
	}
	if runner.gotID != "cli:lark-cli" {
		t.Fatalf("got id %q", runner.gotID)
	}
	if len(runner.gotArgv) != 4 || runner.gotArgv[0] != "auth" || runner.gotArgv[3] != "--json" {
		t.Fatalf("got argv %v", runner.gotArgv)
	}
	if runner.gotOutputMode != "json" {
		t.Fatalf("got output mode %q", runner.gotOutputMode)
	}
	if runner.gotTimeout != 90 {
		t.Fatalf("got timeout %d", runner.gotTimeout)
	}
	if len(runner.progress) != 0 {
		t.Fatalf("expected no stored progress without callback, got %v", runner.progress)
	}
	if out == "" {
		t.Fatal("empty result")
	}
}

func TestInvokeRunCliCommand_ForwardsProgressCallback(t *testing.T) {
	runner := &fakeCliCommandRunner{}
	var got []string
	_, err := InvokeRunCliCommand(context.Background(), runner, json.RawMessage(`{"id":"cli:lark-cli","argv":["auth","login"]}`), func(result string) {
		got = append(got, result)
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != "streaming output" {
		t.Fatalf("unexpected progress callback payload: %v", got)
	}
}

func TestInvokeRunCliCommand_DefaultsTimeoutForLongInitFlows(t *testing.T) {
	runner := &fakeCliCommandRunner{}
	_, err := InvokeRunCliCommand(context.Background(), runner, json.RawMessage(`{"id":"cli:lark-cli","argv":["auth","login"]}`), nil)
	if err != nil {
		t.Fatal(err)
	}
	if runner.gotTimeout != defaultRunCliCommandTimeoutSeconds {
		t.Fatalf("expected default timeout %d, got %d", defaultRunCliCommandTimeoutSeconds, runner.gotTimeout)
	}
}
