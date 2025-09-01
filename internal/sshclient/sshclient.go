package sshclient

import (
	"bytes"
	"fmt"
	"sync"

	"ssht/internal/config"
)

var (
	pool     *ConnectionPool
	poolOnce sync.Once
)

type CommandResult struct {
	Output string
	Error  error
}

func ExecuteCommand(host config.HostConfig, command string) CommandResult {
	// Initialize connection pool once
	poolOnce.Do(func() {
		pool = NewConnectionPool()
	})

	// Get connection from pool
	client, err := pool.GetConnection(host)
	if err != nil {
		return CommandResult{Error: fmt.Errorf("failed to get connection: %w", err)}
	}

	// Create session and execute command
	session, err := client.NewSession()
	if err != nil {
		return CommandResult{Error: fmt.Errorf("failed to create session: %w", err)}
	}
	defer session.Close()

	// Capture command output
	var stdoutBuf, stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	if err := session.Run(command); err != nil {
		return CommandResult{
			Error:  fmt.Errorf("command execution failed: %w", err),
			Output: stderrBuf.String(),
		}
	}

	// Return command output
	return CommandResult{
		Output: stdoutBuf.String(),
	}
}
