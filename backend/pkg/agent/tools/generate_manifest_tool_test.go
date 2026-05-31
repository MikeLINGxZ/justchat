package tools

import (
	"context"
	"encoding/json"
	"testing"
)

type fakeGenerator struct{ gotID string }

func (f *fakeGenerator) GenerateCliManifestSync(_ context.Context, id string) (string, error) {
	f.gotID = id
	return `{"id":"cli:test","name":"test","kind":"cli","runtime_status":"ready"}`, nil
}

func TestInvokeGenerateCliManifest_PassesID(t *testing.T) {
	fg := &fakeGenerator{}
	args := json.RawMessage(`{"id":"cli:test"}`)
	out, err := InvokeGenerateCliManifest(context.Background(), fg, args)
	if err != nil {
		t.Fatal(err)
	}
	if fg.gotID != "cli:test" {
		t.Fatalf("got %q", fg.gotID)
	}
	if out == "" {
		t.Fatal("empty result")
	}
}
