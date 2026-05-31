package tools

import (
	"context"
	"encoding/json"
	"testing"
)

type fakeAttentionRequester struct {
	id     uint
	waited bool
}

// NotifyAttention records one generated notification id for tests.
func (f *fakeAttentionRequester) NotifyAttention(_ context.Context, _ uint, _, _ string) (uint, error) {
	f.id++
	return f.id, nil
}

// WaitForResolution marks the requester as having blocked for user input.
func (f *fakeAttentionRequester) WaitForResolution(_ context.Context, _ uint) error {
	f.waited = true
	return nil
}

// TestInvokeRequestAttentionReturnsAfterResolution verifies the tool waits before returning success.
func TestInvokeRequestAttentionReturnsAfterResolution(t *testing.T) {
	requester := &fakeAttentionRequester{}

	out, err := InvokeRequestAttention(context.Background(), requester, 42, json.RawMessage(`{"title":"X","message":"Y"}`))
	if err != nil {
		t.Fatal(err)
	}
	if !requester.waited {
		t.Fatal("expected tool to wait for resolution")
	}
	if out == "" {
		t.Fatal("expected non-empty tool result")
	}
}
