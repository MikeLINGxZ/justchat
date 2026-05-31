package tools

import (
	"context"
	"encoding/json"
	"testing"
)

type fakeInstaller struct{ gotPkg string }

func (f *fakeInstaller) InstallCliSync(_ context.Context, npmPackage string, _ string, _ func(string, string)) (string, error) {
	f.gotPkg = npmPackage
	return `{"id":"cli:test","name":"test","kind":"cli"}`, nil
}

func TestInvokeInstallCli_PassesPackage(t *testing.T) {
	fi := &fakeInstaller{}
	args := json.RawMessage(`{"npm_package":"@scope/test-cli","name":"test"}`)
	out, err := InvokeInstallCli(context.Background(), fi, args)
	if err != nil {
		t.Fatal(err)
	}
	if fi.gotPkg != "@scope/test-cli" {
		t.Fatalf("got %q", fi.gotPkg)
	}
	if out == "" {
		t.Fatal("empty result")
	}
}
