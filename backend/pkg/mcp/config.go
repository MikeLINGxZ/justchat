package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	toolmcp "trpc.group/trpc-go/trpc-agent-go/tool/mcp"
)

// ServerConfig is the persisted MCP server connection config used by this app.
type ServerConfig struct {
	Transport   string            `json:"transport"`
	Command     string            `json:"command,omitempty"`
	Args        []string          `json:"args,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	Timeout     time.Duration     `json:"timeout,omitempty"`
	WorkingDir  string            `json:"working_dir,omitempty"`
	ServerURL   string            `json:"server_url,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Description string            `json:"description,omitempty"`
}

const defaultStdioTimeout = 30 * time.Second

type rawServerListConfig struct {
	MCPServers map[string]ServerConfig `json:"mcpServers"`
}

// loadServerConfig reads a JSON MCP config file.
func loadServerConfig(configPath string) (ServerConfig, error) {
	bytes, err := os.ReadFile(configPath)
	if err != nil {
		return ServerConfig{}, err
	}

	cfg := ServerConfig{}
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		return ServerConfig{}, err
	}

	if cfg.Command != "" || cfg.ServerURL != "" || cfg.Transport != "" {
		if cfg.Transport == "" {
			cfg.Transport = "stdio"
		}
		return cfg, nil
	}

	rawList := rawServerListConfig{}
	if err := json.Unmarshal(bytes, &rawList); err != nil {
		return ServerConfig{}, err
	}
	for _, server := range rawList.MCPServers {
		if server.Transport == "" {
			server.Transport = "stdio"
		}
		return server, nil
	}

	return ServerConfig{}, fmt.Errorf("invalid mcp config: no server definition found")
}

// toToolSetConfig converts the local config structure to the trpc-agent-go MCP config.
func (c ServerConfig) toToolSetConfig() toolmcp.ConnectionConfig {
	command := c.Command
	args := append([]string(nil), c.Args...)
	timeout := c.Timeout

	if c.Transport == "stdio" && timeout <= 0 {
		timeout = defaultStdioTimeout
	}

	// trpc-agent-go v0.4.0 does not expose stdio env/working_dir on ConnectionConfig,
	// so flatten env into a shell invocation for local stdio MCP servers.
	if c.Transport == "stdio" && (len(c.Env) > 0 || c.WorkingDir != "") {
		parts := []string{}
		if c.WorkingDir != "" {
			parts = append(parts, "cd "+strconv.Quote(c.WorkingDir), "&&")
		}
		envKeys := make([]string, 0, len(c.Env))
		for key := range c.Env {
			envKeys = append(envKeys, key)
		}
		sort.Strings(envKeys)
		for _, key := range envKeys {
			parts = append(parts, key+"="+strconv.Quote(c.Env[key]))
		}
		parts = append(parts, strconv.Quote(command))
		for _, arg := range args {
			parts = append(parts, strconv.Quote(arg))
		}
		command = "sh"
		args = []string{"-c", strings.Join(parts, " ")}
	}

	return toolmcp.ConnectionConfig{
		Transport:   c.Transport,
		Command:     command,
		Args:        args,
		Timeout:     timeout,
		ServerURL:   c.ServerURL,
		Headers:     c.Headers,
		Description: c.Description,
	}
}
