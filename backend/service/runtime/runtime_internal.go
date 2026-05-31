package runtime

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	goRuntime "runtime"
	"strings"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/dir"
)

// RuntimeState represents the persisted Node runtime install state.
type RuntimeState struct {
	State       string    `json:"state"`
	Version     string    `json:"version"`
	InstallDir  string    `json:"install_dir"`
	NodePath    string    `json:"node_path"`
	NpmPath     string    `json:"npm_path"`
	ErrorMsg    string    `json:"error_msg"`
	InstalledAt time.Time `json:"installed_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

const (
	StateMissing      = "missing"
	StateDownloading  = "downloading"
	StateReady        = "ready"
	StateFailed       = "failed"
	StatePendingLater = "pending_later"
)

// runtimeBaseDir returns the absolute directory holding Node runtime artifacts.
func runtimeBaseDir() (string, error) {
	dataDir, err := dir.GetDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dataDir, NodeSubdir), nil
}

// runtimeStatePath returns the absolute path to the state.json file.
func runtimeStatePath() (string, error) {
	base, err := runtimeBaseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, RuntimeStateFileName), nil
}

// LoadPersistedState reads the Node runtime state.json. Missing file yields a fresh missing state.
// Exposed so other services (e.g. plugin) can resolve the bundled node/npm paths without depending on the runtime Service.
func LoadPersistedState() (RuntimeState, error) {
	path, err := runtimeStatePath()
	if err != nil {
		return RuntimeState{}, err
	}
	bytes, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return RuntimeState{State: StateMissing}, nil
		}
		return RuntimeState{}, err
	}

	var state RuntimeState
	if err := json.Unmarshal(bytes, &state); err != nil {
		return RuntimeState{}, err
	}
	if state.State == "" {
		state.State = StateMissing
	}
	return state, nil
}

// saveState writes the runtime state to state.json, creating the directory tree as needed.
func saveState(state RuntimeState) error {
	base, err := runtimeBaseDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(base, 0o755); err != nil {
		return err
	}

	state.UpdatedAt = time.Now()
	bytes, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	path, err := runtimeStatePath()
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0o644)
}

// platformSlug returns the Node.js dist OS slug for the current platform.
func platformSlug() (string, error) {
	switch goRuntime.GOOS {
	case "darwin":
		return "darwin", nil
	case "linux":
		return "linux", nil
	case "windows":
		return "win", nil
	default:
		return "", fmt.Errorf("unsupported os: %s", goRuntime.GOOS)
	}
}

// archSlug returns the Node.js dist arch slug for the current platform.
func archSlug() (string, error) {
	switch goRuntime.GOARCH {
	case "amd64":
		return "x64", nil
	case "arm64":
		return "arm64", nil
	default:
		return "", fmt.Errorf("unsupported arch: %s", goRuntime.GOARCH)
	}
}

// archiveExt returns the archive extension for the current OS.
func archiveExt() (string, error) {
	switch goRuntime.GOOS {
	case "darwin", "linux":
		return "tar.gz", nil
	case "windows":
		return "zip", nil
	default:
		return "", fmt.Errorf("unsupported os: %s", goRuntime.GOOS)
	}
}

// archiveURL returns the full download URL for the configured Node version.
func archiveURL(version string) (string, error) {
	osSlug, err := platformSlug()
	if err != nil {
		return "", err
	}
	arch, err := archSlug()
	if err != nil {
		return "", err
	}
	ext, err := archiveExt()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s/node-%s-%s-%s.%s", NodeDistBaseURL, version, version, osSlug, arch, ext), nil
}

// sha256SumsURL returns the SHA256SUMS.txt URL for the configured Node version.
func sha256SumsURL(version string) string {
	return fmt.Sprintf("%s/%s/SHASUMS256.txt", NodeDistBaseURL, version)
}

// fetchSha256 downloads SHASUMS256.txt and returns the sha256 for the archive we want.
func fetchSha256(ctx context.Context, sumsURL, archiveName string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sumsURL, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %d", resp.StatusCode)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) == 2 && parts[1] == archiveName {
			return parts[0], nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("checksum not found for %s", archiveName)
}

// downloadArchive streams the archive to dest and returns the SHA256 hex digest.
func downloadArchive(ctx context.Context, url, dest string, report func(received, total int64)) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %d", resp.StatusCode)
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return "", err
	}
	out, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer out.Close()

	hasher := sha256.New()
	total := resp.ContentLength
	var received int64
	buf := make([]byte, 64*1024)
	for {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			if _, err := out.Write(buf[:n]); err != nil {
				return "", err
			}
			if _, err := hasher.Write(buf[:n]); err != nil {
				return "", err
			}
			received += int64(n)
			if report != nil {
				report(received, total)
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return "", readErr
		}
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// extractTarGz extracts the tar.gz archive at src to destDir.
func extractTarGz(ctx context.Context, src, destDir string, report func(received, total int64)) (string, error) {
	f, err := os.Open(src)
	if err != nil {
		return "", err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return "", err
	}
	defer gz.Close()

	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", err
	}

	tr := tar.NewReader(gz)
	var topDir string
	stat, _ := f.Stat()
	total := stat.Size()

	for {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if topDir == "" {
			topDir = strings.SplitN(header.Name, "/", 2)[0]
		}
		target := filepath.Join(destDir, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return "", err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return "", err
			}
			out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(header.Mode))
			if err != nil {
				return "", err
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return "", err
			}
			if err := out.Close(); err != nil {
				return "", err
			}
		case tar.TypeSymlink:
			_ = os.Symlink(header.Linkname, target)
		}
		if report != nil {
			pos, _ := f.Seek(0, io.SeekCurrent)
			report(pos, total)
		}
	}

	return topDir, nil
}

// extractZip extracts a zip archive and returns the top-level directory name.
func extractZip(ctx context.Context, src, destDir string, report func(received, total int64)) (string, error) {
	r, err := zip.OpenReader(src)
	if err != nil {
		return "", err
	}
	defer r.Close()

	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", err
	}

	var topDir string
	total := int64(len(r.File))
	for i, file := range r.File {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		if topDir == "" {
			topDir = strings.SplitN(file.Name, "/", 2)[0]
		}
		target := filepath.Join(destDir, file.Name)
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil {
				return "", err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return "", err
		}
		rc, err := file.Open()
		if err != nil {
			return "", err
		}
		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, file.Mode())
		if err != nil {
			rc.Close()
			return "", err
		}
		if _, err := io.Copy(out, rc); err != nil {
			rc.Close()
			out.Close()
			return "", err
		}
		if err := rc.Close(); err != nil {
			out.Close()
			return "", err
		}
		if err := out.Close(); err != nil {
			return "", err
		}
		if report != nil {
			report(int64(i+1), total)
		}
	}

	return topDir, nil
}
