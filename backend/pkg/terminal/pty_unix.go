//go:build !windows

package terminal

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"sync"
	"syscall"

	"github.com/creack/pty"
)

type unixPTYProcess struct {
	cmd      *exec.Cmd
	ptmx     *os.File
	output   chan []byte
	done     chan struct{}
	waitErr  error
	exitCode int
	once     sync.Once
	mu       sync.Mutex
	closed   bool
}

// startPTYProcess starts the configured command in a Unix pseudo-terminal.
func startPTYProcess(ctx context.Context, params CreateParams) (ptyProcess, int, error) {
	cmd := buildCommand(ctx, params)
	rows, cols := params.Rows, params.Cols
	if rows == 0 {
		rows = 40
	}
	if cols == 0 {
		cols = 120
	}
	ptmx, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: rows, Cols: cols})
	if err != nil {
		return nil, 0, err
	}
	p := &unixPTYProcess{
		cmd:    cmd,
		ptmx:   ptmx,
		output: make(chan []byte, 64),
		done:   make(chan struct{}),
	}
	go p.readLoop()
	go func() {
		<-ctx.Done()
		_ = p.Close()
	}()
	pid := 0
	if cmd.Process != nil {
		pid = cmd.Process.Pid
	}
	return p, pid, nil
}

// readLoop forwards PTY output into the manager channel and records process exit state.
func (p *unixPTYProcess) readLoop() {
	buf := make([]byte, 4096)
	for {
		n, err := p.ptmx.Read(buf)
		if n > 0 {
			chunk := make([]byte, n)
			copy(chunk, buf[:n])
			p.output <- chunk
		}
		if err != nil {
			break
		}
	}
	waitErr := p.cmd.Wait()
	exitCode := -1
	if p.cmd.ProcessState != nil {
		exitCode = p.cmd.ProcessState.ExitCode()
	}
	p.once.Do(func() {
		p.mu.Lock()
		p.waitErr = waitErr
		p.exitCode = exitCode
		p.closed = true
		p.mu.Unlock()
		close(p.output)
		_ = p.ptmx.Close()
		close(p.done)
	})
}

// Output returns the channel that carries raw PTY output bytes.
func (p *unixPTYProcess) Output() <-chan []byte { return p.output }

// Write sends bytes to the PTY master.
func (p *unixPTYProcess) Write(data []byte) error {
	p.mu.Lock()
	closed := p.closed
	p.mu.Unlock()
	if closed {
		return errors.New("terminal: session already closed")
	}
	_, err := p.ptmx.Write(data)
	return err
}

// Resize updates the PTY window size.
func (p *unixPTYProcess) Resize(rows, cols uint16) error {
	p.mu.Lock()
	closed := p.closed
	p.mu.Unlock()
	if closed {
		return errors.New("terminal: session already closed")
	}
	return pty.Setsize(p.ptmx, &pty.Winsize{Rows: rows, Cols: cols})
}

// Close terminates the process group and waits for the read loop to finish.
func (p *unixPTYProcess) Close() error {
	p.mu.Lock()
	closed := p.closed
	p.mu.Unlock()
	if closed {
		return nil
	}
	if p.cmd.Process != nil {
		_ = syscall.Kill(-p.cmd.Process.Pid, syscall.SIGKILL)
		_ = p.cmd.Process.Kill()
	}
	_ = p.ptmx.Close()
	<-p.done
	return nil
}

// Wait waits for process completion and returns the exit code and wait error.
func (p *unixPTYProcess) Wait() (int, error) {
	<-p.done
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.exitCode, p.waitErr
}

// Done returns a channel that closes after the PTY process exits.
func (p *unixPTYProcess) Done() <-chan struct{} { return p.done }
