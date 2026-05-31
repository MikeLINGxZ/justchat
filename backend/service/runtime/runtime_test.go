package runtime

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/runtime/runtime_dto"
)

func withTempDataDir(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("LEMONTEA_DATA_DIR", tmp)
	return tmp
}

func TestArchiveURL_PlatformShape(t *testing.T) {
	url, err := archiveURL(NodeLTSVersion)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if !strings.HasPrefix(url, NodeDistBaseURL+"/"+NodeLTSVersion+"/node-"+NodeLTSVersion+"-") {
		t.Fatalf("unexpected url prefix: %s", url)
	}
}

func TestLoadState_MissingFile(t *testing.T) {
	withTempDataDir(t)

	s, err := LoadPersistedState()
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if s.State != StateMissing {
		t.Fatalf("expected missing, got %s", s.State)
	}
}

func TestSaveAndLoadState(t *testing.T) {
	withTempDataDir(t)

	if err := saveState(RuntimeState{State: StateReady, Version: NodeLTSVersion}); err != nil {
		t.Fatalf("save: %v", err)
	}
	s, err := LoadPersistedState()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if s.State != StateReady || s.Version != NodeLTSVersion {
		t.Fatalf("unexpected state: %+v", s)
	}
}

func TestFetchSha256_FindsArchiveChecksum(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("abc123  node-v22.11.0-darwin-arm64.tar.gz\nzzz999  other-file.tar.gz\n"))
	}))
	defer server.Close()

	sum, err := fetchSha256(context.Background(), server.URL, "node-v22.11.0-darwin-arm64.tar.gz")
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if sum != "abc123" {
		t.Fatalf("expected abc123, got %s", sum)
	}
}

func TestMarkDownloadLater_PersistsPendingState(t *testing.T) {
	withTempDataDir(t)

	svc := NewRuntime()
	if _, err := svc.MarkDownloadLater(context.Background(), runtime_dto.MarkDownloadLaterInput{}); err != nil {
		t.Fatalf("unexpected: %v", err)
	}

	state, err := LoadPersistedState()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if state.State != StatePendingLater {
		t.Fatalf("expected %s, got %s", StatePendingLater, state.State)
	}
}
