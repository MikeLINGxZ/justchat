package mcp

import (
	"testing"
	"time"
)

func TestServerConfigToToolSetConfigDefaultsPositiveTimeoutForStdio(t *testing.T) {
	cfg := ServerConfig{
		Transport: "stdio",
		Command:   "uv",
		Args:      []string{"run", "mysql_mcp_server"},
	}

	toolCfg := cfg.toToolSetConfig()
	if toolCfg.Timeout <= 0 {
		t.Fatalf("expected positive timeout for stdio MCP config, got %s", toolCfg.Timeout)
	}
}

func TestServerConfigToToolSetConfigPreservesExplicitTimeout(t *testing.T) {
	cfg := ServerConfig{
		Transport: "stdio",
		Command:   "uv",
		Args:      []string{"run", "mysql_mcp_server"},
		Timeout:   45 * time.Second,
	}

	toolCfg := cfg.toToolSetConfig()
	if toolCfg.Timeout != 45*time.Second {
		t.Fatalf("expected timeout to remain 45s, got %s", toolCfg.Timeout)
	}
}
