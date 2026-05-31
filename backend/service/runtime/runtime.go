package runtime

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	goRuntime "runtime"
	"strings"
	"sync"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/ierror"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/runtime/runtime_dto"
)

// Runtime manages local Node.js runtime download and lifecycle.
type Runtime struct {
	wailsApp *application.App

	mu         sync.Mutex
	cancelFunc context.CancelFunc
	running    bool
}

// NewRuntime constructs the Runtime service.
func NewRuntime() *Runtime {
	return &Runtime{}
}

// GetStatus returns the current persisted Node runtime state.
func (s *Runtime) GetStatus(ctx context.Context, input runtime_dto.GetStatusInput) (*runtime_dto.GetStatusOutput, error) {
	state, err := LoadPersistedState()
	if err != nil {
		return nil, ierror.Error(ierror.ErrRuntimeReadState, err)
	}
	return &runtime_dto.GetStatusOutput{
		State:      state.State,
		Version:    state.Version,
		InstallDir: state.InstallDir,
		NodePath:   state.NodePath,
		NpmPath:    state.NpmPath,
		ErrorMsg:   state.ErrorMsg,
	}, nil
}

// progressPayload is sent on the "runtime.node.progress" event.
type progressPayload struct {
	Phase    string `json:"phase"`
	Received int64  `json:"received"`
	Total    int64  `json:"total"`
	Percent  int    `json:"percent"`
}

// emitProgress publishes a progress update for the frontend.
func (s *Runtime) emitProgress(phase string, received, total int64) {
	percent := 0
	if total > 0 {
		percent = int(received * 100 / total)
	}
	if s.wailsApp != nil {
		s.wailsApp.Event.Emit("runtime.node.progress", progressPayload{
			Phase:    phase,
			Received: received,
			Total:    total,
			Percent:  percent,
		})
	}
}

// failState records a failure into state.json and emits a terminal progress event.
func (s *Runtime) failState(msg string) {
	_ = saveState(RuntimeState{
		State:    StateFailed,
		Version:  NodeLTSVersion,
		ErrorMsg: msg,
	})
	s.emitProgress("verify", 0, 0)
}

// DownloadNode starts an asynchronous Node.js download, returning immediately.
func (s *Runtime) DownloadNode(ctx context.Context, input runtime_dto.DownloadNodeInput) (*runtime_dto.DownloadNodeOutput, error) {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return &runtime_dto.DownloadNodeOutput{}, nil
	}
	runCtx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel
	s.running = true
	s.mu.Unlock()

	go s.runDownload(runCtx)

	return &runtime_dto.DownloadNodeOutput{}, nil
}

func (s *Runtime) runDownload(ctx context.Context) {
	defer func() {
		s.mu.Lock()
		s.running = false
		s.cancelFunc = nil
		s.mu.Unlock()
	}()

	_ = saveState(RuntimeState{State: StateDownloading, Version: NodeLTSVersion})

	url, err := archiveURL(NodeLTSVersion)
	if err != nil {
		s.failState(err.Error())
		return
	}
	archiveName := filepath.Base(url)

	expectedSum, err := fetchSha256(ctx, sha256SumsURL(NodeLTSVersion), archiveName)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			_ = saveState(RuntimeState{State: StateMissing, Version: NodeLTSVersion})
			return
		}
		s.failState(err.Error())
		return
	}

	base, err := runtimeBaseDir()
	if err != nil {
		s.failState(err.Error())
		return
	}
	archivePath := filepath.Join(base, archiveName)

	actualSum, err := downloadArchive(ctx, url, archivePath, func(received, total int64) {
		s.emitProgress("download", received, total)
	})
	if err != nil {
		if errors.Is(err, context.Canceled) {
			_ = saveState(RuntimeState{State: StateMissing, Version: NodeLTSVersion})
			_ = os.Remove(archivePath)
			return
		}
		s.failState(err.Error())
		return
	}
	if !strings.EqualFold(actualSum, expectedSum) {
		s.failState(fmt.Sprintf("checksum mismatch: expected %s got %s", expectedSum, actualSum))
		return
	}

	installRoot := filepath.Join(base, NodeLTSVersion)
	_ = os.RemoveAll(installRoot)
	if err := os.MkdirAll(installRoot, 0o755); err != nil {
		s.failState(err.Error())
		return
	}

	var topDir string
	ext, _ := archiveExt()
	switch ext {
	case "tar.gz":
		topDir, err = extractTarGz(ctx, archivePath, installRoot, func(received, total int64) {
			s.emitProgress("extract", received, total)
		})
	case "zip":
		topDir, err = extractZip(ctx, archivePath, installRoot, func(received, total int64) {
			s.emitProgress("extract", received, total)
		})
	default:
		err = fmt.Errorf("unsupported ext: %s", ext)
	}
	if err != nil {
		if errors.Is(err, context.Canceled) {
			_ = saveState(RuntimeState{State: StateMissing, Version: NodeLTSVersion})
			return
		}
		s.failState(err.Error())
		return
	}

	s.emitProgress("verify", 1, 1)

	extractedRoot := filepath.Join(installRoot, topDir)
	var nodePath string
	var npmPath string
	if goRuntime.GOOS == "windows" {
		nodePath = filepath.Join(extractedRoot, "node.exe")
		npmPath = filepath.Join(extractedRoot, "npm.cmd")
	} else {
		nodePath = filepath.Join(extractedRoot, "bin", "node")
		npmPath = filepath.Join(extractedRoot, "bin", "npm")
	}

	_ = os.Remove(archivePath)

	if err := saveState(RuntimeState{
		State:       StateReady,
		Version:     NodeLTSVersion,
		InstallDir:  extractedRoot,
		NodePath:    nodePath,
		NpmPath:     npmPath,
		InstalledAt: time.Now(),
	}); err != nil {
		s.failState(err.Error())
		return
	}
}

// CancelDownload aborts an in-progress download (if any).
func (s *Runtime) CancelDownload(ctx context.Context, input runtime_dto.CancelDownloadInput) (*runtime_dto.CancelDownloadOutput, error) {
	s.mu.Lock()
	cancel := s.cancelFunc
	s.mu.Unlock()
	if cancel != nil {
		cancel()
	}
	return &runtime_dto.CancelDownloadOutput{}, nil
}

// MarkDownloadLater records a user choice to download the runtime later.
func (s *Runtime) MarkDownloadLater(ctx context.Context, input runtime_dto.MarkDownloadLaterInput) (*runtime_dto.MarkDownloadLaterOutput, error) {
	current, err := LoadPersistedState()
	if err != nil {
		return nil, ierror.Error(ierror.ErrRuntimeReadState, err)
	}
	current.State = StatePendingLater
	current.Version = NodeLTSVersion
	current.ErrorMsg = ""
	if err := saveState(current); err != nil {
		return nil, ierror.Error(ierror.ErrRuntimeWriteState, err)
	}
	return &runtime_dto.MarkDownloadLaterOutput{}, nil
}
