package plugin

import (
	"io"
	"os"
	"os/exec"
	"sync"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
)

type ExtensionHost struct {
	hostScriptPath string
	pluginsDir     string
	cmd            *exec.Cmd
	stdin          io.WriteCloser
	stdout         io.ReadCloser
	rpc            *JsonRpcConn
	mu             sync.Mutex
	running        bool
	onCrash        func() // callback when process exits unexpectedly
}

func NewExtensionHost(hostScriptPath, pluginsDir string) *ExtensionHost {
	return &ExtensionHost{
		hostScriptPath: hostScriptPath,
		pluginsDir:     pluginsDir,
	}
}

func (h *ExtensionHost) Start() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.running {
		return nil
	}

	cmd := exec.Command("node", h.hostScriptPath)
	cmd.Env = append(os.Environ(), "LEMONTEA_PLUGINS_DIR="+h.pluginsDir)
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		stdin.Close()
		return err
	}

	if err := cmd.Start(); err != nil {
		stdin.Close()
		stdout.Close()
		return err
	}

	h.cmd = cmd
	h.stdin = stdin
	h.stdout = stdout
	h.rpc = NewJsonRpcConn(stdout, stdin)
	h.rpc.Start()
	h.running = true

	go func() {
		err := cmd.Wait()
		h.mu.Lock()
		wasRunning := h.running
		if wasRunning {
			h.running = false
		}
		crashFn := h.onCrash
		h.mu.Unlock()

		if wasRunning {
			logger.Error("extension host exited unexpectedly", err)
			if crashFn != nil {
				crashFn()
			}
		}
	}()

	return nil
}

func (h *ExtensionHost) Stop() error {
	h.mu.Lock()
	if !h.running {
		h.mu.Unlock()
		return nil
	}
	h.running = false
	h.mu.Unlock()

	_ = h.rpc.Notify("shutdown", nil)
	h.rpc.Close()
	h.stdin.Close()

	done := make(chan struct{})
	go func() {
		_ = h.cmd.Wait()
		close(done)
	}()

	select {
	case <-done:
		// process exited cleanly
	case <-time.After(5 * time.Second):
		logger.Warm("extension host did not exit in time, killing")
		_ = h.cmd.Process.Kill()
		<-done
	}

	return nil
}

func (h *ExtensionHost) Restart() error {
	_ = h.Stop()
	return h.Start()
}

func (h *ExtensionHost) IsRunning() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.running
}

func (h *ExtensionHost) RPC() *JsonRpcConn {
	return h.rpc
}

func (h *ExtensionHost) SetOnCrash(fn func()) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.onCrash = fn
}
