package ierror

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestError_NilReturnsNil(t *testing.T) {
	result := Error(ErrAgentSendMessage, nil)
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestError_CreatesIErrorWithDetail(t *testing.T) {
	underlying := errors.New("connection refused")
	result := Error(ErrAgentSendMessage, underlying)

	var iErr *IError
	if !errors.As(result, &iErr) {
		t.Fatalf("expected *IError, got %T", result)
	}
	if iErr.Detail != "connection refused" {
		t.Errorf("Detail = %q, want %q", iErr.Detail, "connection refused")
	}
	if iErr.Msg == "" {
		t.Error("Msg should not be empty")
	}
}

func TestError_PassthroughExistingIError(t *testing.T) {
	original := &IError{Detail: "original detail", Msg: "original msg"}
	result := Error(ErrAgentSendMessage, original)

	var iErr *IError
	if !errors.As(result, &iErr) {
		t.Fatalf("expected *IError")
	}
	if iErr.Msg != "original msg" {
		t.Errorf("expected passthrough of original IError, got Msg=%q", iErr.Msg)
	}
}

func TestIError_JSONSerializesDetailAndMsg(t *testing.T) {
	e := &IError{Detail: "raw error text", Msg: "user message"}
	data, err := json.Marshal(e)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out map[string]string
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out["detail"] != "raw error text" {
		t.Errorf("detail = %q, want %q", out["detail"], "raw error text")
	}
	if out["msg"] != "user message" {
		t.Errorf("msg = %q, want %q", out["msg"], "user message")
	}
}

func TestIError_IsMatchesByMsg(t *testing.T) {
	a := &IError{Detail: "d1", Msg: "same msg"}
	b := &IError{Detail: "d2", Msg: "same msg"}
	c := &IError{Detail: "d3", Msg: "different msg"}
	if !errors.Is(a, b) {
		t.Error("expected a.Is(b) == true for same Msg")
	}
	if errors.Is(a, c) {
		t.Error("expected a.Is(c) == false for different Msg")
	}
}
