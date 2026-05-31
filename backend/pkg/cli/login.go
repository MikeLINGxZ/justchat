//go:build !windows

package cli

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/creack/pty"
)

// LoginSessionParams describes inputs for one interactive login session.
type LoginSessionParams struct {
	Executable string   // absolute path, from manifest.Executable
	Argv       []string // manifest.LoginCommand (does NOT include Executable)
	Env        []string // produced by BuildEnv (isolated/shared)
	Cwd        string   // usually plugins/cli/<name>, may be empty
}

// LoginSession owns the PTY-backed child process and its lifecycle.
// All public methods are goroutine-safe.
type LoginSession struct {
	cmd      *exec.Cmd
	ptmx     *os.File
	output   chan []byte
	done     chan struct{}
	waitErr  error
	exitCode int
	mu       sync.Mutex
	closed   bool
	once     sync.Once // guards cleanup: closing ptmx and done
}

// StartLoginSession forks a child process under a PTY and starts draining its output.
// The caller receives ownership of the session; call Wait or Cancel to release resources.
func StartLoginSession(ctx context.Context, p LoginSessionParams) (*LoginSession, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	cmd := exec.Command(p.Executable, p.Argv...)
	if p.Env != nil {
		cmd.Env = p.Env
	}
	if p.Cwd != "" {
		cmd.Dir = p.Cwd
	}

	// We deliberately do not attach cmd.Stdin/Stdout/Stderr — pty.Start wires the
	// PTY slave to all three. Attaching pipes would make exec.CommandContext spawn
	// its own Wait goroutine on cancellation, racing readLoop's cmd.Wait below.
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return nil, err
	}

	s := &LoginSession{
		cmd:    cmd,
		ptmx:   ptmx,
		output: make(chan []byte, 64),
		done:   make(chan struct{}),
	}

	go s.readLoop()
	return s, nil
}

// readLoop drains the PTY master fd and forwards chunks to the output channel.
// It is the sole owner of cleanup: it calls cmd.Wait, closes ptmx, closes output, and signals done.
func (s *LoginSession) readLoop() {
	buf := make([]byte, 4096)
	for {
		n, err := s.ptmx.Read(buf)
		if n > 0 {
			chunk := make([]byte, n)
			copy(chunk, buf[:n])
			s.output <- chunk
		}
		if err != nil {
			// EOF or any read error means the child has closed the PTY master.
			break
		}
	}

	// Collect the exit status before releasing resources.
	waitErr := s.cmd.Wait()
	exitCode := -1
	if s.cmd.ProcessState != nil {
		exitCode = s.cmd.ProcessState.ExitCode()
	}

	s.once.Do(func() {
		s.mu.Lock()
		s.waitErr = waitErr
		s.exitCode = exitCode
		s.closed = true
		s.mu.Unlock()

		close(s.output)
		_ = s.ptmx.Close()
		close(s.done)
	})
}

// Output returns the read-only channel of raw PTY bytes. The channel is closed when the
// child exits or the session is cancelled.
func (s *LoginSession) Output() <-chan []byte { return s.output }

// Write forwards data into the PTY master (the child's stdin).
// Returns "already closed" if the session has finished; if the session closes
// concurrently the write may instead return the underlying syscall error.
func (s *LoginSession) Write(data []byte) error {
	s.mu.Lock()
	closed := s.closed
	s.mu.Unlock()
	if closed {
		return errors.New("cli: login session already closed")
	}
	_, err := s.ptmx.Write(data)
	return err
}

// Resize sends a TIOCSWINSZ ioctl to the PTY master to resize the terminal window.
// Returns "already closed" if the session has finished; if the session closes
// concurrently the resize may instead return the underlying syscall error.
func (s *LoginSession) Resize(rows, cols uint16) error {
	s.mu.Lock()
	closed := s.closed
	s.mu.Unlock()
	if closed {
		return errors.New("cli: login session already closed")
	}
	return pty.Setsize(s.ptmx, &pty.Winsize{Rows: rows, Cols: cols})
}

// Cancel sends SIGTERM to the child process and, if it has not exited within 1.5 seconds,
// sends SIGKILL.
func (s *LoginSession) Cancel() error {
	s.mu.Lock()
	closed := s.closed
	s.mu.Unlock()
	if closed {
		return errors.New("cli: login session already closed")
	}

	if s.cmd.Process == nil {
		return nil
	}
	if err := s.cmd.Process.Signal(syscall.SIGTERM); err != nil {
		if errors.Is(err, os.ErrProcessDone) {
			return errors.New("cli: login session already closed")
		}
		return errors.Join(errors.New("cli: SIGTERM failed"), err)
	}

	go func() {
		timer := time.NewTimer(1500 * time.Millisecond)
		defer timer.Stop()
		select {
		case <-s.done:
			// process exited cleanly after SIGTERM
		case <-timer.C:
			if err := s.cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
				// nothing useful we can do here; log would be ideal but we have no logger
				_ = err
			}
		}
	}()
	return nil
}

// Wait blocks until the session has fully exited and returns the exit code and any error
// from cmd.Wait. Safe to call multiple times.
func (s *LoginSession) Wait() (int, error) {
	<-s.done
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.exitCode, s.waitErr
}
