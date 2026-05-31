package tools

import (
	"context"
	"encoding/json"
	"testing"
)

type fakeTaskStateStore struct {
	values map[string]string
}

func (f *fakeTaskStateStore) SaveTaskState(_ uint, key, value string) error {
	if f.values == nil {
		f.values = map[string]string{}
	}
	f.values[key] = value
	return nil
}

func (f *fakeTaskStateStore) LoadTaskState(_ uint, key string) (string, bool, error) {
	value, ok := f.values[key]
	return value, ok, nil
}

func TestInvokeSaveTaskStatePersistsValue(t *testing.T) {
	store := &fakeTaskStateStore{}
	out, err := InvokeSaveTaskState(context.Background(), store, 7, json.RawMessage(`{"key":"device_code","value":"abc123"}`))
	if err != nil {
		t.Fatal(err)
	}
	if store.values["device_code"] != "abc123" {
		t.Fatalf("unexpected saved value: %#v", store.values)
	}
	if out == "" {
		t.Fatal("expected non-empty result")
	}
}

func TestInvokeLoadTaskStateReturnsStructuredPayload(t *testing.T) {
	store := &fakeTaskStateStore{values: map[string]string{"verification_url": "https://example.com"}}
	out, err := InvokeLoadTaskState(context.Background(), store, 7, json.RawMessage(`{"key":"verification_url"}`))
	if err != nil {
		t.Fatal(err)
	}
	var parsed map[string]any
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}
	if parsed["found"] != true || parsed["value"] != "https://example.com" {
		t.Fatalf("unexpected output: %#v", parsed)
	}
}
